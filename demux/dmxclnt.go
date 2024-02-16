// The demux package multiplexes calls over a transport (e.g., TCP
// connection, unix socket, etc.), and matches responses with requests
// using the call's tag.
package demux

import (
	"io"
	"sync"

	db "sigmaos/debug"
	"sigmaos/serr"
	"sigmaos/sessp"
)

type DemuxClntI interface {
	ReportError(err error)
}

type WriteCallF func(io.Writer, CallI) *serr.Err

type DemuxClnt struct {
	out     io.Writer
	in      io.Reader
	callmap *callMap
	clnti   DemuxClntI
	rf      ReadCallF
	wf      WriteCallF
	mu      sync.Mutex
}

type reply struct {
	rep CallI
	err *serr.Err
}

func NewDemuxClnt(out io.Writer, in io.Reader, rf ReadCallF, wf WriteCallF, clnti DemuxClntI) *DemuxClnt {
	dmx := &DemuxClnt{
		out:     out,
		in:      in,
		callmap: newCallMap(),
		clnti:   clnti,
		rf:      rf,
		wf:      wf,
	}
	go dmx.reader()
	return dmx
}

func (dmx *DemuxClnt) reply(tag sessp.Ttag, rep CallI, err *serr.Err) {
	if ch, ok := dmx.callmap.remove(tag); ok {
		ch <- reply{rep, err}
	} else {
		db.DFatalf("reply remove missing %v\n", tag)
	}
}

func (dmx *DemuxClnt) reader() {
	for {
		c, err := dmx.rf(dmx.in)
		if err != nil {
			db.DPrintf(db.DEMUXCLNT, "reader rf err %v\n", err)
			dmx.callmap.close()
			break
		}
		dmx.reply(c.Tag(), c, nil)
	}
	for _, t := range dmx.callmap.outstanding() {
		db.DPrintf(db.DEMUXCLNT, "reader fail %v\n", t)
		dmx.reply(t, nil, serr.NewErr(serr.TErrUnreachable, "reader"))
	}

}

func (dmx *DemuxClnt) SendReceive(req CallI) (CallI, *serr.Err) {
	ch := make(chan reply)
	if err := dmx.callmap.put(req.Tag(), ch); err != nil {
		db.DPrintf(db.DEMUXCLNT, "SendReceive: enqueue req %v err %v\n", req, err)
		return nil, err
	}
	dmx.mu.Lock()
	err := dmx.wf(dmx.out, req)
	dmx.mu.Unlock()
	if err != nil {
		db.DPrintf(db.DEMUXCLNT, "wf req %v error %v\n", req, err)
		return nil, err
	}
	rep := <-ch
	return rep.rep, rep.err
}

func (dmx *DemuxClnt) Close() error {
	return dmx.callmap.close()
}

func (dmx *DemuxClnt) IsClosed() bool {
	return dmx.callmap.isClosed()
}
