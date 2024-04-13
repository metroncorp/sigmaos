package procqsrv

import (
	"fmt"
	"sync"
	"time"

	"sigmaos/chunk"
	"sigmaos/chunksrv"
	db "sigmaos/debug"
	"sigmaos/proc"
	sp "sigmaos/sigmap"
)

const (
	DEF_Q_SZ = 10
)

type Qitem struct {
	p     *proc.Proc
	kidch chan string
	enqTS time.Time
}

func newQitem(p *proc.Proc, kidch chan string) *Qitem {
	return &Qitem{
		p:     p,
		kidch: kidch,
		enqTS: time.Now(),
	}
}

type Queue struct {
	sync.Mutex
	procs []*Qitem
	pmap  map[sp.Tpid]*proc.Proc
}

func newQueue() *Queue {
	return &Queue{
		procs: make([]*Qitem, 0, DEF_Q_SZ),
		pmap:  make(map[sp.Tpid]*proc.Proc, 0),
	}
}

func (q *Queue) Enqueue(p *proc.Proc, kidch chan string) {
	q.Lock()
	defer q.Unlock()

	q.pmap[p.GetPid()] = p
	qi := newQitem(p, kidch)
	q.procs = append(q.procs, qi)
}

func isEligible(p *proc.Proc, mem proc.Tmem, kernelID string) bool {
	if p.GetMem() > mem {
		return false
	}
	if p.HasNoKernelPref() {
		return true
	}
	return p.HasKernelPref(kernelID)
}

func (q *Queue) Dequeue(mem proc.Tmem, kernelID string) (*proc.Proc, chan string, time.Time, bool) {
	q.Lock()
	defer q.Unlock()

	for i := 0; i < len(q.procs); i++ {
		if isEligible(q.procs[i].p, mem, kernelID) {
			// Save the proc we want to return
			qi := q.procs[i]
			// Delete the i-th proc from the queue
			copy(q.procs[i:], q.procs[i+1:])
			q.procs = q.procs[:len(q.procs)-1]
			delete(q.pmap, qi.p.GetPid())
			return qi.p, qi.kidch, qi.enqTS, true
		}
	}
	return nil, nil, time.UnixMicro(0), false
}

func (q *Queue) updateSigmaPath(prog, kernelId string) {
	for _, qi := range q.procs {
		if !chunksrv.IsChunkSrvPath(qi.p.GetSigmaPath()[0]) &&
			qi.p.GetProgram() == prog {
			db.DPrintf(db.TEST, "PrependSigmaPath: %v GetBinKernelId %v %v\n", qi.p.GetProgram(), qi.p.GetPid(), kernelId)
			qi.p.PrependSigmaPath(chunk.ChunkdPath(kernelId))
		}
	}
}

func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()

	return len(q.procs)
}

func (q *Queue) String() string {
	q.Lock()
	defer q.Unlock()

	return fmt.Sprintf("{ procs:%v }", q.procs)
}
