package binsrv

import (
	"fmt"
	"sync"
	"time"

	"sigmaos/chunksrv"
	db "sigmaos/debug"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
	"sigmaos/uprocclnt"
)

const (
	SEEK_DATA = 3
	SEEK_HOLE = 4
)

type downloader struct {
	mu       sync.Mutex
	pn       string
	kernelId string
	sc       *sigmaclnt.SigmaClnt
	sz       sp.Tsize
	pid      uint32
	updc     *uprocclnt.UprocdClnt
	upds     uprocclnt.UprocSrv
	err      error
	tot      time.Duration
}

func newDownloader(pn string, sc *sigmaclnt.SigmaClnt, updc *uprocclnt.UprocdClnt, upds uprocclnt.UprocSrv, kernelId string, sz sp.Tsize, pid uint32) *downloader {
	dl := &downloader{
		pn:       pn,
		sc:       sc,
		updc:     updc,
		upds:     upds,
		sz:       sz,
		kernelId: kernelId,
		pid:      pid,
	}
	return dl
}

func (dl *downloader) String() string {
	return fmt.Sprintf("{pn %q sz %d ckclnt %v}", dl.pn, dl.sz, dl.updc)
}

// Fetch chunk through uprocd, which will fill in the realm and
// write it a local file, which binsrv can read.
func (dl *downloader) fetchChunk(ck int) (int64, error) {
	db.DPrintf(db.BINSRV, "fetchChunk invoke %q ck %d\n", dl.pn, ck)
	sz := sp.Tsize(0)
	var err error
	if dl.upds != nil {
		db.DPrintf(db.BINSRV, "Fetch direct %q %d\n", dl.pn, ck)
		sz, err = dl.upds.Fetch(int(dl.pid), ck, dl.pn, dl.sz)
	} else {
		sz, err = dl.updc.Fetch(dl.pn, ck, dl.sz, dl.pid)
		if err != nil {
			db.DPrintf(db.BINSRV, "fetchChunk %q fetch %d err %v\n", dl.pn, ck, err)
			return 0, err
		}
	}
	db.DPrintf(db.BINSRV, "fetchChunk done %q ck %d sz %d\n", dl.pn, ck, sz)
	return int64(sz), nil
}

func min(n int64, l int) int {
	if int64(l) < n {
		return l
	} else {
		return int(n)
	}
}

func (dl *downloader) read(off int64, nbyte int) (int, error) {
	i := chunksrv.Index(off)
	o := off - chunksrv.Ckoff(i)
	j := chunksrv.Index(off+int64(nbyte)) + 1
	n := int64(0)
	db.DPrintf(db.BINSRV, "read %d %d: chunks [%d,%d)", off, nbyte, i, j)
	for c := i; c < j; c++ {
		pn := binCachePath(dl.pn)
		sz, ok := chunksrv.IsPresent(pn, c, dl.sz)
		if !ok {
			db.DPrintf(db.BINSRV, "read %d %d: chunk %v not present, need to fetch", off, nbyte, c)
			s := time.Now()
			sz0, err := dl.fetchChunk(c)
			if err != nil {
				return 0, err
			} else {
				sz = sz0
			}
			d := time.Since(s)
			dl.tot += d
			db.DPrintf(db.SPAWN_LAT, "[%v] fetchChunk %d dur %v tot %v", dl.pn, c, d, dl.tot)
		}
		n += sz
		db.DPrintf(db.BINSRV, "read %q ck %d sz %d", pn, c, sz)
	}
	db.DPrintf(db.BINSRV, "read done %d %d: chunks [%d,%d)", off, nbyte, i, j)
	return min(n-o, nbyte), nil
}
