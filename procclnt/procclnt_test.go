package procclnt_test

import (
	"fmt"
	"log"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	db "ulambda/debug"
	"ulambda/fslib"
	"ulambda/named"
	"ulambda/proc"
	"ulambda/procclnt"
	"ulambda/realm"
)

const (
	SLEEP_MSECS = 2000
)

type Tstate struct {
	*procclnt.ProcClnt
	*fslib.FsLib
	t   *testing.T
	e   *realm.TestEnv
	cfg *realm.RealmConfig
}

func makeTstate(t *testing.T) *Tstate {
	ts := &Tstate{}

	bin := ".."
	e := realm.MakeTestEnv(bin)
	cfg, err := e.Boot()
	if err != nil {
		t.Fatalf("Boot %v\n", err)
	}
	ts.e = e
	ts.cfg = cfg

	ts.FsLib = fslib.MakeFsLibAddr("proc_test", ts.cfg.NamedAddr)
	ts.ProcClnt = procclnt.MakeProcClntInit(ts.FsLib, cfg.NamedAddr)
	ts.t = t
	return ts
}

func (ts *Tstate) procd(t *testing.T) string {
	st, err := ts.ReadDir("name/procd")
	assert.Nil(t, err, "Readdir")
	return st[0].Name
}

func makeTstateNoBoot(t *testing.T, cfg *realm.RealmConfig, e *realm.TestEnv, pid string) *Tstate {
	ts := &Tstate{}
	ts.t = t
	ts.e = e
	ts.cfg = cfg
	ts.FsLib = fslib.MakeFsLibAddr("proc_test", ts.cfg.NamedAddr)
	ts.ProcClnt = procclnt.MakeProcClntInit(ts.FsLib, cfg.NamedAddr)
	return ts
}

func spawnSleeperWithPid(t *testing.T, ts *Tstate, pid string) {
	a := proc.MakeProcPid(pid, "bin/user/sleeper", []string{fmt.Sprintf("%dms", SLEEP_MSECS), "name/out_" + pid})
	err := ts.Spawn(a)
	assert.Nil(t, err, "Spawn")
	db.DLPrintf("SCHEDD", "Spawn %v\n", a)
}

func spawnSleeper(t *testing.T, ts *Tstate) string {
	pid := proc.GenPid()
	spawnSleeperWithPid(t, ts, pid)
	return pid
}

func checkSleeperResult(t *testing.T, ts *Tstate, pid string) bool {
	res := true
	b, err := ts.ReadFile("name/out_" + pid)
	res = assert.Nil(t, err, "ReadFile") && res
	res = assert.Equal(t, string(b), "hello", "Output") && res

	return res
}

func checkSleeperResultFalse(t *testing.T, ts *Tstate, pid string) {
	b, err := ts.ReadFile("name/out_" + pid)
	assert.NotNil(t, err, "ReadFile")
	assert.NotEqual(t, string(b), "hello", "Output")
}

func TestWaitExit(t *testing.T) {
	ts := makeTstate(t)

	start := time.Now()

	pid := spawnSleeper(t, ts)
	status, err := ts.WaitExit(pid)
	assert.Nil(t, err, "WaitExit error")
	assert.Equal(t, "OK", status, "Exit status wrong")

	// cleaned up
	_, err = ts.Stat(proc.PidDir(pid))
	assert.NotNil(t, err, "Stat")

	end := time.Now()

	assert.True(t, end.Sub(start) > SLEEP_MSECS*time.Millisecond)

	checkSleeperResult(t, ts, pid)

	ts.e.Shutdown()
}

func TestWaitExitParentRetStat(t *testing.T) {
	ts := makeTstate(t)

	start := time.Now()

	pid := spawnSleeper(t, ts)
	time.Sleep(2 * SLEEP_MSECS * time.Millisecond)
	status, err := ts.WaitExit(pid)
	assert.Nil(t, err, "WaitExit error")
	assert.Equal(t, "OK", status, "Exit status wrong")

	// cleaned up
	_, err = ts.Stat(proc.PidDir(pid))
	assert.NotNil(t, err, "Stat")

	end := time.Now()

	assert.True(t, end.Sub(start) > SLEEP_MSECS*time.Millisecond)

	checkSleeperResult(t, ts, pid)

	ts.e.Shutdown()
}

func TestWaitStart(t *testing.T) {
	ts := makeTstate(t)

	start := time.Now()

	pid := spawnSleeper(t, ts)
	err := ts.WaitStart(pid)
	assert.Nil(t, err, "WaitStart error")

	end := time.Now()

	assert.True(t, end.Sub(start) < SLEEP_MSECS*time.Millisecond, "WaitStart waited too long")

	// Check if proc exists
	sts, err := ts.ReadDir(path.Join("name/procd", ts.procd(t), named.PROCD_RUNNING))
	assert.Nil(t, err, "Readdir")

	// skip ctl entry
	i := 0
	if sts[i].Name == "ctl" {
		i = 1
	}
	assert.Equal(t, pid, sts[i].Name, "pid")

	// Make sure the proc hasn't finished yet...
	checkSleeperResultFalse(t, ts, pid)

	ts.WaitExit(pid)

	checkSleeperResult(t, ts, pid)

	ts.e.Shutdown()
}

// Should exit immediately
func TestWaitNonexistentProc(t *testing.T) {
	ts := makeTstate(t)

	ch := make(chan bool)

	pid := proc.GenPid()
	go func() {
		ts.WaitExit(pid)
		ch <- true
	}()

	done := <-ch
	assert.True(t, done, "Nonexistent proc")

	close(ch)

	ts.e.Shutdown()
}

func TestEarlyExit1(t *testing.T) {
	ts := makeTstate(t)

	pid1 := proc.GenPid()
	a := proc.MakeProc("bin/user/parentexit", []string{fmt.Sprintf("%dms", SLEEP_MSECS), pid1})
	err := ts.Spawn(a)
	assert.Nil(t, err, "Spawn")

	// Wait for parent to finish
	status, err := ts.WaitExit(a.Pid)
	assert.Nil(t, err, "WaitExit")
	assert.Equal(t, "OK", status, "WaitExit")

	// Child should be still running
	_, err = ts.Stat(proc.PidDir(pid1))
	assert.Nil(t, err, "Stat")

	time.Sleep(2 * SLEEP_MSECS * time.Millisecond)

	// Child should have exited
	b, err := ts.ReadFile("name/out_" + pid1)
	assert.Nil(t, err, "ReadFile")
	assert.Equal(t, string(b), "hello", "Output")

	// .. and cleaned up
	_, err = ts.Stat(proc.PidDir(pid1))
	assert.NotNil(t, err, "Stat")

	ts.e.Shutdown()
}

func TestEarlyExitN(t *testing.T) {
	ts := makeTstate(t)
	nProcs := 500
	var done sync.WaitGroup
	done.Add(nProcs)

	for i := 0; i < nProcs; i++ {
		go func() {
			pid1 := proc.GenPid()
			a := proc.MakeProc("bin/user/parentexit", []string{fmt.Sprintf("%dms", 0), pid1})
			err := ts.Spawn(a)
			assert.Nil(t, err, "Spawn")

			// Wait for parent to finish
			status, err := ts.WaitExit(a.Pid)
			assert.Nil(t, err, "WaitExit")
			assert.Equal(t, "OK", status, "WaitExit")

			time.Sleep(2 * SLEEP_MSECS * time.Millisecond)

			// Child should have exited
			b, err := ts.ReadFile("name/out_" + pid1)
			assert.Nil(t, err, "ReadFile")
			assert.Equal(t, string(b), "hello", "Output")

			// .. and cleaned up
			_, err = ts.Stat(proc.PidDir(pid1))
			assert.NotNil(t, err, "Stat")
			done.Done()
		}()
	}
	done.Wait()

	log.Printf("DONE\n")

	ts.e.Shutdown()
}

// Spawn a bunch of procs concurrently, then wait for all of them & check
// their result
func TestConcurrentProcs(t *testing.T) {
	ts := makeTstate(t)

	nProcs := 8
	pids := map[string]int{}

	// Make a bunch of fslibs to avoid concurrency issues
	tses := []*Tstate{}

	var barrier sync.WaitGroup
	barrier.Add(nProcs)
	var started sync.WaitGroup
	started.Add(nProcs)
	var done sync.WaitGroup
	done.Add(nProcs)

	for i := 0; i < nProcs; i++ {
		pid := proc.GenPid()
		_, alreadySpawned := pids[pid]
		for alreadySpawned {
			pid = proc.GenPid()
			_, alreadySpawned = pids[pid]
		}
		pids[pid] = i
		newts := makeTstateNoBoot(t, ts.cfg, ts.e, pid)
		tses = append(tses, newts)
		go func(pid string, started *sync.WaitGroup, i int) {
			barrier.Done()
			barrier.Wait()
			spawnSleeperWithPid(t, tses[i], pid)
			started.Done()
		}(pid, &started, i)
	}

	started.Wait()

	for pid, i := range pids {
		_ = i
		go func(pid string, done *sync.WaitGroup, i int) {
			defer done.Done()
			ts.WaitExit(pid)
			checkSleeperResult(t, tses[i], pid)
			_, err := ts.Stat(proc.PidDir(pid))
			assert.NotNil(t, err, "Stat")
		}(pid, &done, i)
	}

	done.Wait()

	ts.e.Shutdown()
}

func (ts *Tstate) evict(pid string) {
	time.Sleep(SLEEP_MSECS / 2 * time.Millisecond)
	err := ts.Evict(pid)
	assert.Nil(ts.t, err, "evict")
}

func TestEvict(t *testing.T) {
	ts := makeTstate(t)

	start := time.Now()
	pid := spawnSleeper(t, ts)

	go ts.evict(pid)

	status, err := ts.WaitExit(pid)
	assert.Nil(t, err, "WaitExit")
	assert.Equal(t, "EVICTED", status, "WaitExit status")
	end := time.Now()

	assert.True(t, end.Sub(start) < SLEEP_MSECS*time.Millisecond, "Didn't evict early enough.")
	assert.True(t, end.Sub(start) > (SLEEP_MSECS/2)*time.Millisecond, "Evicted too early")

	// Make sure the proc didn't finish
	checkSleeperResultFalse(t, ts, pid)

	ts.e.Shutdown()
}
