package shardsvcmgr

import (
	"strconv"
	"sync"

	"sigmaos/fslib"
	"sigmaos/proc"
	"sigmaos/procclnt"
)

// XXX shard should be a proc or a group

type ShardMgr struct {
	sync.Mutex
	*fslib.FsLib
	*procclnt.ProcClnt
	bin    string
	job    string
	shards []proc.Tpid
	nshard int
	pn     string
}

func (sm *ShardMgr) addShard(i int) error {
	p := proc.MakeProc(sm.bin, []string{sm.job, strconv.Itoa(i)})
	if err := sm.Spawn(p); err != nil {
		return err
	}
	if err := sm.WaitStart(p.Pid); err != nil {
		return err
	}
	sm.shards = append(sm.shards, p.Pid)
	return nil
}

func MkShardMgr(fsl *fslib.FsLib, pclnt *procclnt.ProcClnt, n int, job, bin, pn string) (*ShardMgr, error) {
	sm := &ShardMgr{FsLib: fsl, ProcClnt: pclnt, bin: bin, job: job, shards: make([]proc.Tpid, 0), nshard: n, pn: pn}
	for i := 0; i < n; i++ {
		if err := sm.addShard(i); err != nil {
			return nil, err
		}
	}
	return sm, nil
}

func (sm *ShardMgr) AddShard() error {
	sm.Lock()
	defer sm.Unlock()

	n := len(sm.shards)
	return sm.addShard(n)
}

func Shard(i int) string {
	return strconv.Itoa(i)
}

func (sm *ShardMgr) Nshard() int {
	return len(sm.shards)
}

func (sm *ShardMgr) Path() string {
	return sm.pn
}

func (sm *ShardMgr) Server(i int) string {
	return sm.pn + Shard(i)
}

func (sm *ShardMgr) Stop() error {
	for _, pid := range sm.shards {
		if err := sm.Evict(pid); err != nil {
			return err
		}
		if _, err := sm.WaitExit(pid); err != nil {
			return err
		}
	}
	return nil
}
