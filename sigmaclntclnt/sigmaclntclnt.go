// Package sigmaclntclnt is an RPC-based client for proxing the
// [sigmaos] interface. It sends request to [sigmclntsrv], which has a
// full sigmaclnt to interact with SigmaOS.
package sigmaclntclnt

import (
	"bufio"
	"net"

	db "sigmaos/debug"
	"sigmaos/demux"
	"sigmaos/frame"
	"sigmaos/proc"
	"sigmaos/rpc"
	"sigmaos/rpcclnt"
	"sigmaos/sessp"
	scproto "sigmaos/sigmaclntsrv/proto"
	sp "sigmaos/sigmap"
)

type SigmaClntClnt struct {
	pe           *proc.ProcEnv
	dmx          *demux.DemuxClnt
	rpcc         *rpcclnt.RPCClnt
	seqno        sessp.Tseqno
	conn         net.Conn
	disconnected bool
}

func (scc *SigmaClntClnt) SendReceive(a []byte) ([]byte, error) {
	rep, err := scc.dmx.SendReceive([]frame.Tframe{a})
	return rep[0], err
}

func (scc *SigmaClntClnt) StatsSrv() (*rpc.SigmaRPCStats, error) {
	return nil, nil
}

func (scc *SigmaClntClnt) ReportError(err error) {
	db.DPrintf(db.DEMUXCLNT, "ReportError %v", err)
	go func() {
		scc.close()
	}()
}

func NewSigmaClntClnt(pe *proc.ProcEnv) (*SigmaClntClnt, error) {
	conn, err := net.Dial("unix", sp.SIGMASOCKET)
	if err != nil {
		return nil, err
	}
	scc := &SigmaClntClnt{
		pe:    pe,
		dmx:   nil,
		rpcc:  nil,
		seqno: 0, conn: conn, disconnected: false,
	}
	scc.dmx = demux.NewDemuxClnt(bufio.NewWriterSize(conn, sp.Conf.Conn.MSG_LEN),
		bufio.NewReaderSize(conn, sp.Conf.Conn.MSG_LEN), 1, scc)
	scc.rpcc = rpcclnt.NewRPCClntCh(scc)
	// Initialize the server-side component of sigmaclnt by sending the proc env
	db.DPrintf(db.SIGMACLNTCLNT, "Init sigmaclntclnt for %v", pe.GetPID())
	if err := scc.Init(); err != nil {
		db.DPrintf(db.ERROR, "Error init sigmaclnt: %v", err)
		return nil, err
	}
	return scc, nil
}

// Tell sigmaclntd to shut down
func (scc *SigmaClntClnt) Shutdown() error {
	req := scproto.SigmaNullRequest{}
	rep := scproto.SigmaErrReply{}
	err := scc.rpcErr("SigmaClntSrvAPI.Shutdown", &req, &rep)
	db.DPrintf(db.SIGMACLNTCLNT, "Shutdown %v %v %v", req, rep, err)
	return err
}

// Close the socket connection
func (scc *SigmaClntClnt) close() error {
	if err := scc.conn.Close(); err != nil {
		return err
	}
	return scc.dmx.Close()
}
