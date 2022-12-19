package fidclnt

import (
	"fmt"

	"sigmaos/fcall"
	"sigmaos/path"
	"sigmaos/protclnt"
	sp "sigmaos/sigmap"
)

//
// Sigma file system API at the level of fids.
//

type FidClnt struct {
	fids *FidMap
	pc   *protclnt.Clnt
	ft   *FenceTable
}

func MakeFidClnt() *FidClnt {
	fidc := &FidClnt{}
	fidc.fids = mkFidMap()
	fidc.pc = protclnt.MakeClnt()
	fidc.ft = MakeFenceTable()
	return fidc
}

func (fidc *FidClnt) String() string {
	str := fmt.Sprintf("Fsclnt fid table %p:\n%v", fidc, fidc.fids)
	return str
}

func (fidc *FidClnt) Len() int {
	return len(fidc.fids.fids)
}

func (fidc *FidClnt) FenceDir(path string, f sp.Tfence) *fcall.Err {
	return fidc.ft.Insert(path, f)
}

func (fidc *FidClnt) ReadSeqNo() sp.Tseqno {
	return fidc.pc.ReadSeqNo()
}

func (fidc *FidClnt) Exit() *fcall.Err {
	return fidc.pc.Exit()
}

func (fidc *FidClnt) allocFid() sp.Tfid {
	return fidc.fids.allocFid()
}

func (fidc *FidClnt) freeFid(np sp.Tfid) {
	// not implemented
}

func (fidc *FidClnt) Free(fid sp.Tfid) {
	fidc.fids.free(fid)
}

func (fidc *FidClnt) Lookup(fid sp.Tfid) *Channel {
	return fidc.fids.lookup(fid)
}

func (fidc *FidClnt) Qid(fid sp.Tfid) *sp.Tqid {
	return fidc.Lookup(fid).Lastqid()
}

func (fidc *FidClnt) Qids(fid sp.Tfid) []*sp.Tqid {
	return fidc.Lookup(fid).qids
}

func (fidc *FidClnt) Path(fid sp.Tfid) path.Path {
	return fidc.Lookup(fid).Path()
}

func (fidc *FidClnt) Insert(fid sp.Tfid, path *Channel) {
	fidc.fids.insert(fid, path)
}

func (fidc *FidClnt) Clunk(fid sp.Tfid) *fcall.Err {
	err := fidc.fids.lookup(fid).pc.Clunk(fid)
	if err != nil {
		return err
	}
	fidc.fids.free(fid)
	return nil
}

func (fidc *FidClnt) Attach(uname string, addrs []string, pn, tree string) (sp.Tfid, *fcall.Err) {
	fid := fidc.allocFid()
	reply, err := fidc.pc.Attach(addrs, uname, fid, path.Split(tree))
	if err != nil {
		fidc.freeFid(fid)
		return sp.NoFid, err
	}
	pc := fidc.pc.MakeProtClnt(addrs)
	fidc.fids.insert(fid, makeChannel(pc, uname, path.Split(pn), []*sp.Tqid{reply.Qid}))
	return fid, nil
}

func (fidc *FidClnt) Detach(fid sp.Tfid) *fcall.Err {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return fcall.MkErr(fcall.TErrUnreachable, "detach")
	}
	if err := ch.pc.Detach(); err != nil {
		return err
	}
	return nil
}

// Walk returns the fid it walked to (which maybe fid) and the
// remaining path left to be walked (which maybe the original path).
func (fidc *FidClnt) Walk(fid sp.Tfid, path []string) (sp.Tfid, []string, *fcall.Err) {
	nfid := fidc.allocFid()
	reply, err := fidc.Lookup(fid).pc.Walk(fid, nfid, path)
	if err != nil {
		fidc.freeFid(nfid)
		return fid, path, err
	}
	channel := fidc.Lookup(fid).Copy()
	channel.AddN(reply.Qids, path)
	fidc.Insert(nfid, channel)
	return nfid, path[len(reply.Qids):], nil
}

// A defensive version of walk because fid is shared among several
// threads (it comes out the mount table) and one thread may free the
// fid while another thread is using it.
func (fidc *FidClnt) Clone(fid sp.Tfid) (sp.Tfid, *fcall.Err) {
	nfid := fidc.allocFid()
	channel := fidc.Lookup(fid)
	if channel == nil {
		return sp.NoFid, fcall.MkErr(fcall.TErrUnreachable, "clone")
	}
	_, err := channel.pc.Walk(fid, nfid, path.Path{})
	if err != nil {
		fidc.freeFid(nfid)
		return fid, err
	}
	channel = channel.Copy()
	fidc.Insert(nfid, channel)
	return nfid, err
}

func (fidc *FidClnt) Create(fid sp.Tfid, name string, perm sp.Tperm, mode sp.Tmode) (sp.Tfid, *fcall.Err) {
	reply, err := fidc.fids.lookup(fid).pc.Create(fid, name, perm, mode)
	if err != nil {
		return sp.NoFid, err
	}
	fidc.fids.lookup(fid).add(name, reply.Qid)
	return fid, nil
}

func (fidc *FidClnt) Open(fid sp.Tfid, mode sp.Tmode) (*sp.Tqid, *fcall.Err) {
	reply, err := fidc.fids.lookup(fid).pc.Open(fid, mode)
	if err != nil {
		return nil, err
	}
	return reply.Qid, nil
}

func (fidc *FidClnt) Watch(fid sp.Tfid) *fcall.Err {
	return fidc.fids.lookup(fid).pc.Watch(fid)
}

func (fidc *FidClnt) Wstat(fid sp.Tfid, st *sp.Stat) *fcall.Err {
	f := fidc.ft.Lookup(fidc.fids.lookup(fid).Path())
	_, err := fidc.fids.lookup(fid).pc.WstatF(fid, st, f)
	return err
}

func (fidc *FidClnt) Renameat(fid sp.Tfid, o string, fid1 sp.Tfid, n string) *fcall.Err {
	f := fidc.ft.Lookup(fidc.fids.lookup(fid).Path())
	if fidc.fids.lookup(fid).pc != fidc.fids.lookup(fid1).pc {
		return fcall.MkErr(fcall.TErrInval, "paths at different servers")
	}
	_, err := fidc.fids.lookup(fid).pc.Renameat(fid, o, fid1, n, f)
	return err
}

func (fidc *FidClnt) Remove(fid sp.Tfid) *fcall.Err {
	f := fidc.ft.Lookup(fidc.fids.lookup(fid).Path())
	return fidc.fids.lookup(fid).pc.RemoveF(fid, f)
}

func (fidc *FidClnt) RemoveFile(fid sp.Tfid, wnames []string, resolve bool) *fcall.Err {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return fcall.MkErr(fcall.TErrUnreachable, "getfile")
	}
	f := fidc.ft.Lookup(ch.Path().AppendPath(wnames))
	return ch.pc.RemoveFile(fid, wnames, resolve, f)
}

func (fidc *FidClnt) Stat(fid sp.Tfid) (*sp.Stat, *fcall.Err) {
	reply, err := fidc.fids.lookup(fid).pc.Stat(fid)
	if err != nil {
		return nil, err
	}
	return reply.Stat, nil
}

func (fidc *FidClnt) ReadV(fid sp.Tfid, off sp.Toffset, cnt sp.Tsize, v sp.TQversion) ([]byte, *fcall.Err) {
	f := fidc.ft.Lookup(fidc.fids.lookup(fid).Path())
	data, err := fidc.fids.lookup(fid).pc.ReadVF(fid, off, cnt, f, v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Unfenced read
func (fidc *FidClnt) ReadVU(fid sp.Tfid, off sp.Toffset, cnt sp.Tsize, v sp.TQversion) ([]byte, *fcall.Err) {
	data, err := fidc.fids.lookup(fid).pc.ReadVF(fid, off, cnt, sp.MakeFenceNull(), v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fidc *FidClnt) WriteV(fid sp.Tfid, off sp.Toffset, data []byte, v sp.TQversion) (sp.Tsize, *fcall.Err) {
	f := fidc.ft.Lookup(fidc.fids.lookup(fid).Path())
	reply, err := fidc.fids.lookup(fid).pc.WriteVF(fid, off, f, v, data)
	if err != nil {
		return 0, err
	}
	return reply.Tcount(), nil
}

func (fidc *FidClnt) WriteRead(fid sp.Tfid, data []byte) ([]byte, *fcall.Err) {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return nil, fcall.MkErr(fcall.TErrUnreachable, "WriteRead")
	}
	data, err := fidc.fids.lookup(fid).pc.WriteRead(fid, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fidc *FidClnt) GetFile(fid sp.Tfid, path []string, mode sp.Tmode, off sp.Toffset, cnt sp.Tsize, resolve bool) ([]byte, *fcall.Err) {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return nil, fcall.MkErr(fcall.TErrUnreachable, "getfile")
	}
	f := fidc.ft.Lookup(ch.Path().AppendPath(path))
	data, err := ch.pc.GetFile(fid, path, mode, off, cnt, resolve, f)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (fidc *FidClnt) SetFile(fid sp.Tfid, path []string, mode sp.Tmode, off sp.Toffset, data []byte, resolve bool) (sp.Tsize, *fcall.Err) {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return 0, fcall.MkErr(fcall.TErrUnreachable, "getfile")
	}
	f := fidc.ft.Lookup(ch.Path().AppendPath(path))
	reply, err := ch.pc.SetFile(fid, path, mode, off, resolve, f, data)
	if err != nil {
		return 0, err
	}
	return reply.Tcount(), nil
}

func (fidc *FidClnt) PutFile(fid sp.Tfid, path []string, mode sp.Tmode, perm sp.Tperm, off sp.Toffset, data []byte) (sp.Tsize, *fcall.Err) {
	ch := fidc.fids.lookup(fid)
	if ch == nil {
		return 0, fcall.MkErr(fcall.TErrUnreachable, "putfile")
	}
	f := fidc.ft.Lookup(ch.Path().AppendPath(path))
	reply, err := ch.pc.PutFile(fid, path, mode, perm, off, f, data)
	if err != nil {
		return 0, err
	}
	return reply.Tcount(), nil
}
