package proc

import (
	"encoding/json"
	"fmt"
	"log"
	"path"

	"github.com/thanhpk/randstr"

	db "ulambda/debug"
	"ulambda/fslib"
	np "ulambda/ninep"
)

type Ttype uint32
type Tcore uint32

const (
	T_DEF Ttype = 0
	T_LC  Ttype = 1
	T_BE  Ttype = 2
)

const (
	C_DEF Tcore = 0
)

const (
	// XXX REMOVE BY IMPLEMENTING TRUNC
	WAITFILE_PADDING = 1000
)

const (
	RUNQ          = "name/runq"
	RUNQLC        = "name/runqlc"
	WAITQ         = "name/waitq"
	CLAIMED       = "name/claimed"
	CLAIMED_EPH   = "name/claimed_ephemeral"
	SPAWNED       = "name/spawned"
	RET_STAT      = "name/retstat"
	JOB_SIGNAL    = "job-signal"
	WAIT_LOCK     = "wait-lock."
	CRASH_TIMEOUT = 1
)

type WaitFile struct {
	Started      bool
	RetStatFiles []string
}

//type Proc struct {
//	Pid        string          // SigmaOS PID
//	Program    string          // Program to run
//	WDir       string          // Working directory for the process
//	Args       []string        // Args
//	Env        []string        // Environment variables
//	StartDep   []string        // Start dependencies // XXX Replace somehow?
//	ExitDep    map[string]bool // Exit dependencies// XXX Replace somehow?
//	StartTimer uint32          // Start timer in seconds
//	Type       Ttype           // Type
//	Ncore      Tcore           // Number of cores requested
//}

type ProcCtl struct {
	*fslib.FsLib
}

func MakeProcCtl(fsl *fslib.FsLib) *ProcCtl {
	pctl := &ProcCtl{}
	pctl.FsLib = fsl

	return pctl
}

// ========== SPAWN ==========

func (pctl *ProcCtl) Spawn(p *fslib.Attr) error {
	// Create a file for waiters to watch & wait on
	err := pctl.makeWaitFile(p.Pid)
	if err != nil {
		return err
	}
	pctl.pruneExitDeps(p)
	b, err := json.Marshal(p)
	if err != nil {
		// Unlock the waiter file if unmarshal failed
		pctl.removeWaitFile(p.Pid)
		return err
	}
	err = pctl.MakeFileAtomic(path.Join(WAITQ, p.Pid), 0777, b)
	if err != nil {
		return err
	}
	// Notify localds that a job has become runnable
	pctl.SignalNewJob()
	return nil
}

func (pctl *ProcCtl) makeWaitFile(pid string) error {
	fpath := waitFilePath(pid)
	var wf WaitFile
	wf.Started = false
	b, err := json.Marshal(wf)
	if err != nil {
		log.Printf("Error marshalling waitfile: %v", err)
	}
	// XXX hack around lack of OTRUNC
	for i := 0; i < WAITFILE_PADDING; i++ {
		b = append(b, ' ')
	}
	// Make a writable, versioned file
	err = pctl.MakeFile(fpath, 0777, np.OWRITE, b)
	// Sometimes we get "EOF" on shutdown
	if err != nil && err.Error() != "EOF" {
		return fmt.Errorf("Error on MakeFile MakeWaitFile %v: %v", fpath, err)
	}
	return nil
}

// XXX When we start handling large numbers of lambdas, may be better to stat
// each exit dep individually. For now, this is more efficient (# of RPCs).
// If we know nothing about an exit dep, ignore it by marking it as exited
func (pctl *ProcCtl) pruneExitDeps(p *fslib.Attr) {
	spawned := pctl.getSpawnedLambdas()
	for pid, _ := range p.ExitDep {
		if _, ok := spawned[waitFileName(pid)]; !ok {
			p.ExitDep[pid] = true
		}
	}
}

func (pctl *ProcCtl) removeWaitFile(pid string) error {
	fpath := waitFilePath(pid)
	err := pctl.Remove(fpath)
	if err != nil {
		log.Printf("Error on RemoveWaitFile  %v: %v", fpath, err)
		return err
	}
	return nil
}

func waitFilePath(pid string) string {
	return path.Join(SPAWNED, waitFileName(pid))
}

func waitFileName(pid string) string {
	return fslib.LockName(WAIT_LOCK + pid)
}

func (pctl *ProcCtl) getSpawnedLambdas() map[string]bool {
	d, err := pctl.ReadDir(SPAWNED)
	if err != nil {
		log.Printf("Error reading spawned dir in pruneExitDeps: %v", err)
	}
	spawned := map[string]bool{}
	for _, l := range d {
		spawned[l.Name] = true
	}
	return spawned
}

// ========== STARTED ==========

/*
 * PairDep-based lambdas are runnable only if they are the producer (whoever
 * claims and runs the producer will also start the consumer, so we disallow
 * unilaterally claiming the consumer for now), and only once all of their
 * consumers have been started. For now we assume that
 * consumers only have one producer, and the roles of producer and consumer
 * are mutually exclusive. We also expect (though not strictly necessary)
 * that producers only have one consumer each. If this is no longer the case,
 * we should handle oversubscription more carefully.
 */
func (pctl *ProcCtl) Started(pid string) error {
	pctl.setWaitFileStarted(pid, true)
	return nil
}

func (pctl *ProcCtl) setWaitFileStarted(pid string, started bool) {
	pctl.LockFile(fslib.LOCKS, waitFilePath(pid))
	defer pctl.UnlockFile(fslib.LOCKS, waitFilePath(pid))

	// Get the current contents of the file & its version
	b1, _, err := pctl.GetFile(waitFilePath(pid))
	if err != nil {
		log.Printf("Error reading when registerring retstat: %v, %v", waitFilePath(pid), err)
		return
	}
	var wf WaitFile
	err = json.Unmarshal(b1, &wf)
	if err != nil {
		log.Fatalf("Error unmarshalling waitfile: %v, %v", string(b1), err)
		return
	}
	wf.Started = started
	b2, err := json.Marshal(wf)
	if err != nil {
		log.Printf("Error marshalling waitfile: %v", err)
		return
	}
	// XXX hack around lack of OTRUNC
	for i := 0; i < WAITFILE_PADDING; i++ {
		b2 = append(b2, ' ')
	}
	_, err = pctl.SetFile(waitFilePath(pid), b2, np.NoV)
	if err != nil {
		log.Printf("Error writing when registerring retstat: %v, %v", waitFilePath(pid), err)
	}
}

// ========== EXITING ==========

func (pctl *ProcCtl) Exiting(pid string, status string) error {
	pctl.WakeupExit(pid)
	err := pctl.Remove(path.Join(CLAIMED, pid))
	if err != nil {
		log.Printf("Error removing claimed in Exiting %v: %v", pid, err)
	}
	err = pctl.Remove(path.Join(CLAIMED_EPH, pid))
	if err != nil {
		log.Printf("Error removing claimed_eph in Exiting %v: %v", pid, err)
	}
	// Write back return statuses
	pctl.writeBackRetStats(pid, status)

	// Release people waiting on this lambda
	return pctl.removeWaitFile(pid)
}

// Write back return statuses
func (pctl *ProcCtl) writeBackRetStats(pid string, status string) {
	pctl.LockFile(fslib.LOCKS, waitFilePath(pid))
	defer pctl.UnlockFile(fslib.LOCKS, waitFilePath(pid))

	b, _, err := pctl.GetFile(waitFilePath(pid))
	if err != nil {
		log.Printf("Error reading waitfile in WriteBackRetStats: %v, %v", waitFilePath(pid), err)
		return
	}
	var wf WaitFile
	err = json.Unmarshal(b, &wf)
	if err != nil {
		log.Printf("Error unmarshalling waitfile: %v, %v, %v", string(b), wf, err)
	}
	for _, p := range wf.RetStatFiles {
		if len(p) > 0 {
			pctl.WriteFile(p, []byte(status))
		}
	}
}

// ========== WAIT ==========

// Create a file to read return status from, watch wait file, and return
// contents of retstat file.
func (pctl *ProcCtl) Wait(pid string) ([]byte, error) {
	// XXX We can make return statuses optional to save on RPCs if we don't care
	// about them... right now they require a LOT of RPCs.

	// Make a file in which to receive the return status
	fpath := pctl.makeRetStatFile()

	// Communicate the file name to the lambda we're waiting on
	pctl.registerRetStatFile(pid, fpath)

	// Wait on the lambda with a watch
	done := make(chan bool)
	err := pctl.SetRemoveWatch(waitFilePath(pid), func(p string, err error) {
		if err != nil && err.Error() == "EOF" {
			return
		} else if err != nil {
			log.Printf("Error in wait watch: %v", err)
		}
		done <- true
	})
	// if error, don't wait; the lambda may already have exited.
	if err == nil {
		<-done
	}

	// Read the exit status
	b, _, err := pctl.GetFile(fpath)
	if err != nil {
		log.Printf("Error reading retstat file in wait: %v, %v", fpath, err)
		return b, err
	}

	// Clean up our temp file
	err = pctl.Remove(fpath)
	if err != nil {
		log.Printf("Error removing retstat file in wait: %v, %v", fpath, err)
		return b, err
	}
	return b, err
}

// Create a randomly-named ephemeral file to mark into which the return status
// will be written.
func (pctl *ProcCtl) makeRetStatFile() string {
	fname := randstr.Hex(16)
	fpath := path.Join(RET_STAT, fname)
	err := pctl.MakeFile(fpath, 0777|np.DMTMP, np.OWRITE, []byte{})
	if err != nil {
		log.Printf("Error creating return status file: %v, %v", fpath, err)
	}
	return fpath
}

// Register that we want a return status written back
func (pctl *ProcCtl) registerRetStatFile(pid string, fpath string) {
	pctl.LockFile(fslib.LOCKS, waitFilePath(pid))
	defer pctl.UnlockFile(fslib.LOCKS, waitFilePath(pid))

	// Get the current contents of the file & its version
	b1, _, err := pctl.GetFile(waitFilePath(pid))
	if err != nil {
		db.DLPrintf("LOCALD", "Error reading when registerring retstat: %v, %v", waitFilePath(pid), err)
		return
	}
	var wf WaitFile
	err = json.Unmarshal(b1, &wf)
	if err != nil {
		log.Fatalf("Error unmarshalling waitfile: %v, %v", string(b1), err)
		return
	}
	wf.RetStatFiles = append(wf.RetStatFiles, fpath)
	b2, err := json.Marshal(wf)
	if err != nil {
		log.Printf("Error marshalling waitfile: %v", err)
		return
	}
	// XXX hack around lack of OTRUNC
	for i := 0; i < WAITFILE_PADDING; i++ {
		b2 = append(b2, ' ')
	}
	_, err = pctl.SetFile(waitFilePath(pid), b2, np.NoV)
	if err != nil {
		log.Printf("Error writing when registerring retstat: %v, %v", waitFilePath(pid), err)
	}
}

// XXX REMOVE
// ========== SPAWN_NO_OP =========

// Spawn a no-op lambda
func (pctl *ProcCtl) SpawnNoOp(pid string, exitDep []string) error {
	a := &fslib.Attr{}
	a.Pid = pid
	a.Program = fslib.NO_OP_LAMBDA
	exitDepMap := map[string]bool{}
	for _, dep := range exitDep {
		exitDepMap[dep] = false
	}
	a.ExitDep = exitDepMap
	return pctl.Spawn(a)
}
