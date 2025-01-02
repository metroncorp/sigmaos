package fid

import (
	"fmt"
	"sync"

	"sigmaos/api/fs"
	db "sigmaos/debug"
	"sigmaos/serr"
	sessp "sigmaos/session/proto"
	sp "sigmaos/sigmap"
)

// Several fids may name the same Pobj. For example, each session's
// fid 0 refers to the root of the file system.
type Pobj struct {
	name   string   // name for obj
	obj    fs.FsObj // the obj in the backing file system
	parent fs.Dir   // parent dir of obj
	ctx    fs.CtxI  // the context of the attached sesssion
}

func NewPobj(name string, o fs.FsObj, dir fs.Dir, ctx fs.CtxI) *Pobj {
	return &Pobj{name: name, parent: dir, obj: o, ctx: ctx}
}

func (po *Pobj) String() string {
	return fmt.Sprintf("{name '%v'(p %d) o %v parent %v ctx %v}", po.name, po.Path(), po.obj, po.parent, po.ctx)
}

func (po *Pobj) Name() string {
	return po.name
}

func (po *Pobj) Path() sp.Tpath {
	return po.obj.Path()
}

func (po *Pobj) Ctx() fs.CtxI {
	return po.ctx
}

func (po *Pobj) SetName(name string) {
	po.name = name
}

func (po *Pobj) Obj() fs.FsObj {
	return po.obj
}

func (po *Pobj) SetObj(o fs.FsObj) {
	po.obj = o
}

func (po *Pobj) Parent() fs.Dir {
	return po.parent
}

type Fid struct {
	mu     sync.Mutex
	po     *Pobj
	isOpen bool     // has Create/Open() been called for po.Obj?
	m      sp.Tmode // mode for Create/Open()
	qid    sp.Tqid  // the qid of obj at the time of invoking NewFidPath
	cursor int      // for directories
}

func NewFid(pobj *Pobj, m sp.Tmode, qid sp.Tqid) *Fid {
	return &Fid{isOpen: false, po: pobj, m: m, qid: qid}
}

func (f *Fid) String() string {
	return fmt.Sprintf("{po %v o? %v %v v %v}", f.po, f.isOpen, f.m, f.qid)
}

func (f *Fid) Mode() sp.Tmode {
	return f.m
}

func (f *Fid) SetMode(m sp.Tmode) {
	f.isOpen = true
	f.m = m
}

func (f *Fid) Pobj() *Pobj {
	return f.po
}

func (f *Fid) Parent() fs.Dir {
	return f.po.parent
}

func (f *Fid) IsOpen() bool {
	return f.isOpen
}

func (f *Fid) Qid() *sp.Tqid {
	return &f.qid
}

func (f *Fid) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.isOpen = false
}

func (f *Fid) Write(off sp.Toffset, b []byte, fence sp.Tfence) (sp.Tsize, *serr.Err) {
	o := f.Pobj().Obj()
	var err *serr.Err
	sz := sp.Tsize(0)

	switch i := o.(type) {
	case fs.File:
		sz, err = i.Write(f.Pobj().Ctx(), off, b, fence)
	default:
		db.DFatalf("Write: obj type %T isn't Dir or File\n", o)
	}
	return sz, err
}

func (f *Fid) WriteRead(req sessp.IoVec) (sessp.IoVec, *serr.Err) {
	o := f.Pobj().Obj()
	var err *serr.Err
	var iov sessp.IoVec
	switch i := o.(type) {
	case fs.RPC:
		iov, err = i.WriteRead(f.Pobj().Ctx(), req)
	default:
		db.DFatalf("Write: obj type %T isn't RPC\n", o)
	}
	return iov, err
}

func (f *Fid) readDir(o fs.FsObj, off sp.Toffset, count sp.Tsize) ([]byte, *serr.Err) {
	d := o.(fs.Dir)
	dirents, err := d.ReadDir(f.Pobj().Ctx(), f.cursor, count)
	if err != nil {
		return nil, err
	}
	b, n, err := fs.MarshalDir(count, dirents)
	if err != nil {
		return nil, err
	}
	f.cursor += n
	return b, nil
}

func (f *Fid) Read(off sp.Toffset, count sp.Tsize, fence sp.Tfence) ([]byte, *serr.Err) {
	po := f.Pobj()
	switch i := po.Obj().(type) {
	case fs.Dir:
		return f.readDir(po.Obj(), off, count)
	case fs.File:
		b, err := i.Read(po.Ctx(), off, count, fence)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		db.DFatalf("Read: obj %v type %T isn't Dir or File\n", po.Obj(), po.Obj())
		return nil, nil
	}
}
