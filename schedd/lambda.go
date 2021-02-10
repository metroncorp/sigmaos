package schedd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	db "ulambda/debug"
	"ulambda/fslib"
)

type Lambda struct {
	mu         sync.Mutex
	cond       *sync.Cond
	condWait   *sync.Cond
	sd         *Sched
	pid        string
	status     string
	exitStatus string
	program    string
	args       []string
	env        []string
	consDep    map[string]bool // if true, consumer has finished
	prodDep    map[string]bool // if true, producer is running
	exitDep    map[string]bool
}

func makeLambda(sd *Sched, a string) (*Lambda, error) {
	l := &Lambda{}
	l.sd = sd
	l.cond = sync.NewCond(&l.mu)
	l.condWait = sync.NewCond(&l.mu)
	l.consDep = make(map[string]bool)
	l.prodDep = make(map[string]bool)
	l.exitDep = make(map[string]bool)
	var attr fslib.Attr
	err := json.Unmarshal([]byte(a), &attr)
	if err != nil {
		return nil, err
	}
	l.pid = attr.Pid
	l.program = attr.Program
	l.args = attr.Args
	l.env = attr.Env
	for _, p := range attr.PairDep {
		if l.pid != p.Producer {
			c, ok := sd.ls[p.Producer]
			if ok {
				l.prodDep[p.Producer] = c.isRunnning()
			} else {
				l.prodDep[p.Producer] = false
			}
		}
		if l.pid != p.Consumer {
			l.consDep[p.Consumer] = false
		}
	}
	for _, p := range attr.ExitDep {
		l.exitDep[p] = false
	}
	if l.runnableConsumer() {
		l.status = "Runnable"
	} else {
		l.status = "Waiting"
	}
	return l, nil
}

func (l *Lambda) String() string {
	str := fmt.Sprintf("λ pid %v st %v %v args %v env %v cons %v prod %v exit %v",
		l.pid, l.status, l.program, l.args, l.env, l.consDep, l.prodDep,
		l.exitDep)
	return str
}

func (l *Lambda) changeStatus(new string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.status = new
	return nil
}

// XXX Might want to optimize this.
func (l *Lambda) swapExitDependency(depSwaps map[string]string) {
	// Assuming len(depSwaps) << len(l.exitDeps)
	for from, to := range depSwaps {
		// Check if present & false (hasn't exited yet)
		if val, ok := l.exitDep[from]; ok && !val {
			l.exitDep[to] = false
			l.exitDep[from] = true
		}
	}
}

// XXX if remote, keep-alive?
func (l *Lambda) wait(cmd *exec.Cmd) {
	err := cmd.Wait()
	if err != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		log.Printf("Lambda %v finished with error: %v", l, err)
	}
}

// XXX if had remote machines, this would be run on the remote machine
// maybe we should have machines register with ulambd; have a
// directory with machines?
func (l *Lambda) run() error {
	db.DPrintf("Run %v\n", l)
	err := l.changeStatus("Started")
	if err != nil {
		return err
	}
	args := append([]string{l.pid}, l.args...)
	env := append(os.Environ(), l.env...)
	cmd := exec.Command(l.program, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	go l.wait(cmd)
	return nil
}

func (l *Lambda) startConsDep() {
	for p, _ := range l.consDep {
		c := l.sd.findLambda(p)
		if c != nil {
			c.markProducer(l.pid)
		}
	}
}

func (l *Lambda) startExitDep(pid string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, ok := l.exitDep[pid]
	if ok {
		l.exitDep[pid] = true
		for _, b := range l.exitDep {
			if !b {
				return
			}
		}
		l.status = "Runnable"
	}
}

func (l *Lambda) stopProducers() {
	for p, _ := range l.prodDep {
		c := l.sd.findLambda(p)
		if c != nil {
			c.markConsumer(l.pid)
		}
	}
}

func (l *Lambda) markProducer(pid string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prodDep[pid] = true
}

func (l *Lambda) markExit(pid string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.exitDep[pid] = true
}

func (l *Lambda) markConsumer(pid string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.consDep[pid] = true
	l.cond.Signal()
}

// caller holds lock
func (l *Lambda) runnableConsumer() bool {
	if len(l.exitDep) != 0 {
		return false
	}
	run := true
	for _, b := range l.prodDep {
		if !b {
			return false
		}
	}
	return run
}

func (l *Lambda) runnableWaitingConsumer() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.status != "Waiting" {
		return false
	}
	return l.runnableConsumer()
}

// Wait for consumers that depend on me to exit too.
func (l *Lambda) waitExit() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for !l.exitable() {
		l.cond.Wait()
	}
}

// Caller hold locks
func (l *Lambda) exitable() bool {
	exit := true
	for _, b := range l.consDep {
		if !b {
			return false
		}
	}
	return exit
}

func (l *Lambda) isRunnable() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.status == "Runnable"
}

func (l *Lambda) isRunnning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.status == "Running"
}

// A caller wants to Wait for l to exit
func (l *Lambda) waitFor() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	db.DPrintf("Wait for %v\n", l)
	if l.status != "Exiting" {
		l.condWait.Wait()
	}
	return l.exitStatus
}

// l is exiting; wakeup waiters who are waiting for me to exit
func (l *Lambda) wakeupWaiter(exitStatus string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.exitStatus = exitStatus
	db.DPrintf("Wakeup waiters for %v\n", l)
	l.condWait.Broadcast()
}
