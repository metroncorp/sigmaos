package kv

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ulambda/fslib"
)

type System struct {
	named  *exec.Cmd
	schedd *exec.Cmd
}

func run(name string) (*exec.Cmd, error) {
	cmd := exec.Command(name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, cmd.Start()
}

func Boot(bin string) (*System, error) {
	s := &System{}

	cmd, err := run(bin + "/named")
	if err != nil {
		return nil, err
	}
	s.named = cmd
	s.schedd, err = run(bin + "/schedd")
	if err != nil {
		return nil, err
	}

	time.Sleep(1000 * time.Millisecond)

	return s, nil
}

func (s *System) Shutdown(clnt *fslib.FsLib) {
	err := clnt.WriteFile(fslib.SCHEDDEV, []byte("Exit"))
	if err != nil {
		log.Fatalf("Schedd shutdown %v\n", err)
	}
	err = clnt.WriteFile(fslib.NAMEDEV, []byte("Exit"))
	if err != nil {
		log.Fatalf("Named shutdown %v\n", err)
	}
	s.schedd.Wait()
	s.named.Wait()
}

type Tstate struct {
	t *testing.T
	s *System
	*KvClerk
	ch chan bool
}

func makeTstate(t *testing.T) *Tstate {
	ts := &Tstate{}
	ts.t = t
	ts.ch = make(chan bool)
	bin := "../bin"

	s, err := Boot(bin)
	if err != nil {
		t.Fatalf("Boot %v\n", err)
	}
	ts.s = s

	fsl := fslib.MakeFsLib("boot")

	err = fsl.Mkdir("name/kv", 0777)
	if err != nil {
		t.Fatalf("Mkdir %v\n", err)
	}

	err = fsl.SpawnProgram(bin+"/sharderd", []string{bin})
	if err != nil {
		t.Fatalf("Spawn %v\n", err)

	}

	time.Sleep(1000 * time.Millisecond)

	kc, err := MakeClerk()
	if err != nil {
		t.Fatalf("Make clerk %v\n", err)
	}
	ts.KvClerk = kc

	return ts
}

func (ts *Tstate) cleanup() {
	err := ts.WriteFile(SHARDER+"/dev", []byte("Exit"))
	if err != nil {
		ts.t.Fatalf("Sharder shutdown %v\n", err)
	}
	ts.s.Shutdown(ts.KvClerk.FsLib)
}

func (ts *Tstate) spawnKv() {
	for {
		err := ts.WriteFile(SHARDER+"/dev", []byte("Add"))
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (ts *Tstate) delKv() {
	for {
		err := ts.WriteFile(SHARDER+"/dev", []byte("Del"))
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (ts *Tstate) getKeys() {
	for i := 0; i < 100; i++ {
		k := strconv.Itoa(i)
		v, err := ts.Get(k)
		assert.Nil(ts.t, err, "Get "+k)
		assert.Equal(ts.t, v, k, "Get")
	}
}

// func TestBasic(t *testing.T) {
// 	kc, err := MakeClerk()
// 	assert.Nil(t, err, "MakeClerk")

// 	for i := 0; i < 100; i++ {
// 		err := ts.Put(strconv.Itoa(i), strconv.Itoa(i))
// 		assert.Nil(t, err, "Put")
// 	}

// 	for r := 0; r < NSHARD-1; r++ {
// 		spawnKv(t, kc)
// 		time.Sleep(100 * time.Millisecond)
// 		getKeys(t, kc)
// 	}

// 	for r := NSHARD - 1; r > 0; r-- {
// 		delKv(t, kc)
// 		time.Sleep(100 * time.Millisecond)
// 		getKeys(t, kc)
// 	}
// }

func (ts *Tstate) clerk() {
	for {
		select {
		case <-ts.ch:
			break
		default:
			ts.getKeys()
		}
	}
}

func TestConcur(t *testing.T) {
	ts := makeTstate(t)

	for i := 0; i < 100; i++ {
		err := ts.Put(strconv.Itoa(i), strconv.Itoa(i))
		assert.Nil(t, err, "Put")
	}

	go ts.clerk()

	for r := 0; r < NSHARD-1; r++ {
		ts.spawnKv()
		time.Sleep(1000 * time.Millisecond)
	}

	for r := NSHARD - 1; r > 0; r-- {
		ts.delKv()
		time.Sleep(1000 * time.Millisecond)
	}

	ts.ch <- true

	ts.cleanup()
}
