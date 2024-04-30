package chunksrv_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"sigmaos/chunk"
	"sigmaos/chunkclnt"
	"sigmaos/chunksrv"
	db "sigmaos/debug"
	sp "sigmaos/sigmap"
	"sigmaos/test"
)

const (
	PROG = "sleeper"
	PATH = "name/ux/~local/bin/user/common/"
)

type Tstate struct {
	*test.Tstate
	ckclnt *chunkclnt.ChunkClnt
	srvs   []string
	bins   *chunkclnt.BinPaths
}

func newTstate(t *testing.T, n int) *Tstate {
	ts := &Tstate{
		bins: chunkclnt.NewBinPaths(),
	}
	s, err := test.NewTstateAll(t)
	assert.Nil(t, err)
	ts.Tstate = s

	err = s.BootNode(n)
	assert.Nil(t, err, "Boot node: %v", err)

	ckclnt := chunkclnt.NewChunkClnt(ts.FsLib)
	ckclnt.UpdateChunkds()
	srvs, err := ckclnt.GetSrvs()
	assert.Nil(t, err)

	ts.srvs = srvs
	ts.ckclnt = ckclnt

	db.DPrintf(db.TEST, "Chunksrvs  %v", ts.srvs)

	for _, srv := range srvs {
		pn := chunksrv.PathHostKernelRealm(srv, sp.ROOTREALM)
		os.Mkdir(pn, 0700)
	}
	return ts
}

func (ts *Tstate) check(srv string, st *sp.Stat) {
	pn := chunksrv.PathHostKernelRealm(srv, sp.ROOTREALM)
	pn = path.Join(pn, PROG)
	fi, err := os.Stat(pn)
	assert.Nil(ts.T, err)
	assert.Equal(ts.T, st.Length, uint64(fi.Size()))
}

func (ts *Tstate) fetch(srv string, paths []string, expect string) {
	pid := ts.ProcEnv().GetPID()

	st, path, err := ts.ckclnt.GetFileStat(srv, PROG, pid, sp.ROOTREALM, paths)
	assert.Nil(ts.T, err)
	assert.Equal(ts.T, expect, path)
	db.DPrintf(db.TEST, "st len %v path %q\n", st.Length, path)

	path, err = ts.ckclnt.FetchBinary(srv, PROG, pid, sp.ROOTREALM, sp.Tsize(st.Length), paths)
	assert.Nil(ts.T, err, "err %v", err)
	assert.Equal(ts.T, expect, path)
	db.DPrintf(db.TEST, "bin %v %q\n", PROG, path)

	ts.bins.SetBinKernelID(PROG, srv)

	ts.check(srv, st)
}

func TestFetchOrigin(t *testing.T) {
	ts := newTstate(t, 0)
	ts.fetch(ts.srvs[0], []string{PATH}, PATH)
	ts.Shutdown()
}

func TestFetchCache(t *testing.T) {
	ts := newTstate(t, 0)

	ts.fetch(ts.srvs[0], []string{PATH}, PATH)
	ts.fetch(ts.srvs[0], []string{PATH}, chunk.ChunkdPath(ts.srvs[0]))

	ts.Shutdown()
}

func TestFetchChunkd(t *testing.T) {
	ts := newTstate(t, 1)

	ts.fetch(ts.srvs[0], []string{PATH}, PATH)

	kid, ok := ts.bins.GetBinKernelID(PROG)
	assert.True(ts.T, ok)

	srv := chunk.ChunkdPath(kid)
	ts.fetch(ts.srvs[0], []string{srv}, srv)

	ts.Shutdown()
}
