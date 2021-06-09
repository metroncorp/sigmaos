package proxy

import (
	"log"
	"net"
	"os/user"
	"sync"

	db "ulambda/debug"
	"ulambda/fsclnt"
	"ulambda/fslib"
	np "ulambda/ninep"
	"ulambda/npclnt"
	npo "ulambda/npobjsrv"
	"ulambda/npsrv"
)

//
// XXX convert to use npobjsrv
//

const MAXSYMLINK = 20

// The connection from the kernel/client
type NpConn struct {
	mu    sync.Mutex
	conn  net.Conn
	clnt  *npclnt.NpClnt
	uname string
	fids  map[np.Tfid]*npclnt.NpChan // The outgoing channels to servers proxied
	named string
}

func makeNpConn(conn net.Conn, named string) *NpConn {
	npc := &NpConn{}
	npc.conn = conn
	npc.clnt = npclnt.MakeNpClnt()
	npc.fids = make(map[np.Tfid]*npclnt.NpChan)
	npc.named = named
	return npc
}

func (npc *NpConn) npch(fid np.Tfid) *npclnt.NpChan {
	npc.mu.Lock()
	defer npc.mu.Unlock()
	ch, ok := npc.fids[fid]
	if !ok {
		log.Fatal("npch: unknown fid ", fid)
	}
	return ch
}

func (npc *NpConn) addch(fid np.Tfid, ch *npclnt.NpChan) {
	npc.mu.Lock()
	defer npc.mu.Unlock()
	npc.fids[fid] = ch
}

func (npc *NpConn) delch(fid np.Tfid) {
	npc.mu.Lock()
	defer npc.mu.Unlock()
	delete(npc.fids, fid)
}

type Npd struct {
	named string
	st    *npo.SessionTable
}

func MakeNpd() *Npd {
	return &Npd{fslib.Named(), nil}
}

// XXX should/is happen only once for the one mount for :1110
func (npd *Npd) Connect(conn net.Conn) npsrv.NpAPI {
	clnt := makeNpConn(conn, npd.named)
	return clnt
}

func (npd *Npd) SessionTable() *npo.SessionTable {
	if npd.st == nil {
		npd.st = npo.MakeSessionTable()
	}
	return npd.st
}

func (npd *Npd) RegisterSession(sess np.Tsession) {
	if npd.st == nil {
		npd.st = npo.MakeSessionTable()
	}
	npd.st.RegisterSession(sess)
}

func (npc *NpConn) Version(sess np.Tsession, args np.Tversion, rets *np.Rversion) *np.Rerror {
	rets.Msize = args.Msize
	rets.Version = "9P2000"
	return nil
}

func (npc *NpConn) Auth(sess np.Tsession, args np.Tauth, rets *np.Rauth) *np.Rerror {
	return np.ErrUnknownMsg
}

func (npc *NpConn) Attach(sess np.Tsession, args np.Tattach, rets *np.Rattach) *np.Rerror {
	u, err := user.Current()
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	npc.uname = u.Uid

	log.Printf("attach %v\n", args)

	reply, err := npc.clnt.Attach([]string{npc.named}, npc.uname, args.Fid, np.Split(args.Aname))
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	npc.addch(args.Fid, npc.clnt.MakeNpChan([]string{npc.named}))
	rets.Qid = reply.Qid
	return nil
}

func (npc *NpConn) Detach() {
	db.DLPrintf("9POBJ", "Detach\n")
}

// XXX avoid duplication with fsclnt
func (npc *NpConn) autoMount(newfid np.Tfid, target string, path []string) (np.Tqid, error) {
	db.DPrintf("automount %v to %v\n", target, path)
	server, _ := fsclnt.SplitTarget(target)
	reply, err := npc.clnt.Attach([]string{server}, npc.uname, newfid, path)
	if err != nil {
		return np.Tqid{}, err
	}
	npc.addch(newfid, npc.clnt.MakeNpChan([]string{server}))
	return reply.Qid, nil
}

// XXX avoid duplication with fsclnt
func (npc *NpConn) readLink(fid np.Tfid) (string, error) {
	_, err := npc.npch(fid).Open(fid, np.OREAD)
	if err != nil {
		return "", err
	}
	reply, err := npc.npch(fid).Read(fid, 0, 1024)
	if err != nil {
		return "", err
	}
	npc.delch(fid)
	return string(reply.Data), nil
}

func (npc *NpConn) Walk(sess np.Tsession, args np.Twalk, rets *np.Rwalk) *np.Rerror {
	path := args.Wnames
	// XXX accumulate qids
	for i := 0; i < MAXSYMLINK; i++ {
		reply, err := npc.npch(args.Fid).Walk(args.Fid, args.NewFid, path)
		if err != nil {
			return np.ErrNotfound
		}
		if len(reply.Qids) == 0 { // clone args.Fid?
			npc.addch(args.NewFid, npc.npch(args.Fid))
			*rets = *reply
			break
		}
		qid := reply.Qids[len(reply.Qids)-1]
		if qid.Type&np.QTSYMLINK == np.QTSYMLINK {
			todo := len(path) - len(reply.Qids)
			db.DPrintf("symlink %v %v\n", todo, path)

			// args.Newfid is fid for symlink
			npc.addch(args.NewFid, npc.npch(args.Fid))

			target, err := npc.readLink(args.NewFid)
			if err != nil {
				return np.ErrUnknownfid
			}
			// XXX assumes symlink is final component of walk
			if fsclnt.IsRemoteTarget(target) {
				qid, err = npc.autoMount(args.NewFid, target, path[todo:])
				if err != nil {
					return np.ErrUnknownfid
				}
				reply.Qids[len(reply.Qids)-1] = qid
				path = path[todo:]
				db.DPrintf("automounted: %v -> %v, %v\n", args.NewFid,
					target, path)
				*rets = *reply
				break
			} else {
				log.Fatal("don't handle")
			}
		} else { // newFid is at same server as args.Fid
			npc.addch(args.NewFid, npc.npch(args.Fid))
			*rets = *reply
			break
		}
	}
	return nil
}

func (npc *NpConn) Open(sess np.Tsession, args np.Topen, rets *np.Ropen) *np.Rerror {
	reply, err := npc.npch(args.Fid).Open(args.Fid, args.Mode)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) WatchV(sess np.Tsession, args np.Twatchv, rets *np.Ropen) *np.Rerror {
	return nil
}

func (npc *NpConn) Create(sess np.Tsession, args np.Tcreate, rets *np.Rcreate) *np.Rerror {
	reply, err := npc.npch(args.Fid).Create(args.Fid, args.Name, args.Perm, args.Mode)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) Clunk(sess np.Tsession, args np.Tclunk, rets *np.Rclunk) *np.Rerror {
	err := npc.npch(args.Fid).Clunk(args.Fid)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	npc.delch(args.Fid)
	return nil
}

func (npc *NpConn) Flush(sess np.Tsession, args np.Tflush, rets *np.Rflush) *np.Rerror {
	return nil
}

func (npc *NpConn) Read(sess np.Tsession, args np.Tread, rets *np.Rread) *np.Rerror {
	reply, err := npc.npch(args.Fid).Read(args.Fid, args.Offset, args.Count)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) ReadV(sess np.Tsession, args np.Treadv, rets *np.Rread) *np.Rerror {
	return nil
}

func (npc *NpConn) Write(sess np.Tsession, args np.Twrite, rets *np.Rwrite) *np.Rerror {
	reply, err := npc.npch(args.Fid).Write(args.Fid, args.Offset, args.Data)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) WriteV(sess np.Tsession, args np.Twritev, rets *np.Rwrite) *np.Rerror {
	return nil
}

func (npc *NpConn) Remove(sess np.Tsession, args np.Tremove, rets *np.Rremove) *np.Rerror {
	err := npc.npch(args.Fid).Remove(args.Fid)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	return nil
}

func (npc *NpConn) Stat(sess np.Tsession, args np.Tstat, rets *np.Rstat) *np.Rerror {
	reply, err := npc.npch(args.Fid).Stat(args.Fid)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) Wstat(sess np.Tsession, args np.Twstat, rets *np.Rwstat) *np.Rerror {
	reply, err := npc.npch(args.Fid).Wstat(args.Fid, &args.Stat)
	if err != nil {
		return &np.Rerror{err.Error()}
	}
	*rets = *reply
	return nil
}

func (npc *NpConn) Renameat(sess np.Tsession, args np.Trenameat, rets *np.Rrenameat) *np.Rerror {
	return np.ErrNotSupported
}
