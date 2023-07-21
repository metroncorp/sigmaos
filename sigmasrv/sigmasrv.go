package sigmasrv

import (
	"path"
	"reflect"
	"strings"

	"sigmaos/ctx"
	db "sigmaos/debug"
	"sigmaos/dir"
	"sigmaos/memfs"
	"sigmaos/memfssrv"
	"sigmaos/proc"
	"sigmaos/protdev"
	"sigmaos/sessdevsrv"
	"sigmaos/sesssrv"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
)

//
// Many SigmaOS servers use SigmaSrv to create and run servers.  A server
// typically consists of a MemFS (an in-memory file system accessed
// through sigmap), one or more RPC end points, including an end point
// for leasesrv (to manage leases).  Sigmasrv creates the end-points
// in the memfs. Some servers don't use SigmaSrv and directly interact
// with SessSrv (e.g., ux and knamed/named).
//

type SigmaSrv struct {
	*memfssrv.MemFs
	sti  *protdev.StatInfo
	svc  *svcMap
	lsrv *LeaseSrv
}

// Make a sigmasrv with an memfs, and publish server at fn.
func MakeSigmaSrv(fn string, svci any, uname sp.Tuname) (*SigmaSrv, error) {
	mfs, error := memfssrv.MakeMemFs(fn, uname)
	if error != nil {
		db.DFatalf("MakeSigmaSrv %v err %v\n", fn, error)
	}
	return MakeSigmaSrvMemFs(mfs, svci)
}

func MakeSigmaSrvPublic(fn string, svci any, uname sp.Tuname, public bool) (*SigmaSrv, error) {
	db.DPrintf(db.ALWAYS, "MakeSigmaSrvPublic %T\n", svci)
	if public {
		mfs, error := memfssrv.MakeMemFsPublic(fn, uname)
		if error != nil {
			return nil, error
		}
		return MakeSigmaSrvMemFs(mfs, svci)
	} else {
		return MakeSigmaSrv(fn, svci, uname)
	}
}

// Make a sigmasrv and memfs and publish srv at fn. Note: no lease
// server.
func MakeSigmaSrvNoRPC(fn string, uname sp.Tuname) (*SigmaSrv, error) {
	mfs, err := memfssrv.MakeMemFs(fn, uname)
	if err != nil {
		db.DFatalf("MakeSigmaSrv %v err %v\n", fn, err)
	}

	return newSigmaSrv(mfs), nil
}

func MakeSigmaSrvPort(fn, port string, uname sp.Tuname, svci any) (*SigmaSrv, error) {
	mfs, error := memfssrv.MakeMemFsPort(fn, ":"+port, uname)
	if error != nil {
		db.DFatalf("MakeSigmaSrvPort %v err %v\n", fn, error)
	}
	return MakeSigmaSrvMemFs(mfs, svci)
}

func MakeSigmaSrvClnt(fn string, sc *sigmaclnt.SigmaClnt, uname sp.Tuname, svci any) (*SigmaSrv, error) {
	mfs, error := memfssrv.MakeMemFsPortClnt(fn, ":0", sc)
	if error != nil {
		db.DFatalf("MakeSigmaSrvClnt %v err %v\n", fn, error)
	}
	return makeSigmaSrvRPC(mfs, svci)
}

func MakeSigmaSrvClntNoRPC(fn string, sc *sigmaclnt.SigmaClnt, uname sp.Tuname) (*SigmaSrv, error) {
	mfs, err := memfssrv.MakeMemFsPortClnt(fn, ":0", sc)
	if err != nil {
		db.DFatalf("MakeMemFsPortClnt %v err %v\n", fn, err)
	}
	ssrv := newSigmaSrv(mfs)
	return ssrv, nil
}

// Makes a sigmasrv with an memfs, rpc server, and LeaseSrv RPC
// service.
func MakeSigmaSrvMemFs(mfs *memfssrv.MemFs, svci any) (*SigmaSrv, error) {
	ssrv, err := makeSigmaSrvRPC(mfs, svci)
	if err != nil {
		return nil, err
	}
	if err := ssrv.NewLeaseSrv(); err != nil {
		return nil, err
	}
	return ssrv, nil
}

func newSigmaSrv(mfs *memfssrv.MemFs) *SigmaSrv {
	ssrv := &SigmaSrv{MemFs: mfs, svc: newSvcMap()}
	return ssrv
}

// Make a sigmasrv with an RPC server
func makeSigmaSrvRPC(mfs *memfssrv.MemFs, svci any) (*SigmaSrv, error) {
	ssrv := newSigmaSrv(mfs)
	return ssrv, ssrv.makeRPCSrv(svci)
}

// Create the rpc server directory in memfs and register the RPC
// service svci to the RPC server.
func (ssrv *SigmaSrv) makeRPCSrv(svci any) error {
	db.DPrintf(db.SIGMASRV, "makeRPCSrv: %v\n", svci)
	if _, err := ssrv.Create(protdev.RPC, sp.DMDIR|0777, sp.ORDWR, sp.NoLeaseId); err != nil {
		return err
	}
	if err := ssrv.registerRPCSrv(svci); err != nil {
		return err
	}
	return nil
}

func MakeSigmaSrvSess(sesssrv *sesssrv.SessSrv, uname sp.Tuname) *SigmaSrv {
	mfs := memfssrv.MakeMemFsSrv(uname, "", sesssrv)
	return newSigmaSrv(mfs)
}

// Mount the rpc directory in sessrv and create the RPC service in
// it. This function is useful for SigmaSrv that don't have an MemFs
// (e.g., knamed/named).
func (ssrv *SigmaSrv) MountRPCSrv(svci any) error {
	d := dir.MkRootDir(ctx.MkCtxNull(), memfs.MakeInode, nil)
	ssrv.MemFs.SessSrv.Mount(protdev.RPC, d.(*dir.DirImpl))
	if err := ssrv.registerRPCSrv(svci); err != nil {
		return err
	}
	return nil
}

// Make the rpc server
func (ssrv *SigmaSrv) registerRPCSrv(svci any) error {
	ssrv.svc.NewRPCService(svci)
	rd := mkRpcDev(ssrv)

	if err := sessdevsrv.MkSessDev(ssrv.MemFs, path.Join(protdev.RPC, protdev.RPC), rd.mkRpcSession, nil); err != nil {
		return err
	}
	if si, err := makeStatsDev(ssrv.MemFs, protdev.RPC); err != nil {
		return err
	} else {
		ssrv.sti = si
	}
	return nil
}

// Assumes RPCSrv has been created and create a LeaseSrv service.
func (ssrv *SigmaSrv) NewLeaseSrv() error {
	lsrv := newLeaseSrv(ssrv.MemFs)
	ssrv.svc.NewRPCService(lsrv)
	return nil
}

func (ssrv *SigmaSrv) QueueLen() int64 {
	return ssrv.MemFs.QueueLen()
}

func structName(svci any) string {
	typ := reflect.TypeOf(svci)
	name := typ.String()
	dot := strings.LastIndex(name, ".")
	return name[dot+1:]
}

func (ssrv *SigmaSrv) RunServer() error {
	db.DPrintf(db.SIGMASRV, "Run %v\n", proc.GetProgram())
	ssrv.MemFs.Serve()
	if ssrv.lsrv != nil {
		ssrv.lsrv.Stop()
	}
	ssrv.Exit(proc.MakeStatus(proc.StatusEvicted))
	return nil
}

func (ssrv *SigmaSrv) Exit(status *proc.Status) error {
	db.DPrintf(db.SIGMASRV, "Run %v\n", proc.GetProgram())
	if ssrv.lsrv != nil {
		ssrv.lsrv.Stop()
	}
	return ssrv.MemFs.Exit(proc.MakeStatus(proc.StatusEvicted))
}
