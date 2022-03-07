package snapshot

import (
	"encoding/json"
	"log"
	"reflect"

	"ulambda/dir"
	"ulambda/fences"
	"ulambda/fs"
	"ulambda/memfs"
	"ulambda/protsrv"
	"ulambda/repl"
	"ulambda/session"
	"ulambda/stats"
	"ulambda/threadmgr"
)

type Snapshot struct {
	fssrv        protsrv.FsServer
	Imap         map[uint64]ObjSnapshot
	Root         uint64
	St           []byte
	Tmt          []byte
	Rft          []byte
	Rc           []byte
	restoreCache map[uint64]fs.FsObj
}

func MakeSnapshot(fssrv protsrv.FsServer) *Snapshot {
	s := &Snapshot{}
	s.fssrv = fssrv
	s.Imap = make(map[uint64]ObjSnapshot)
	s.Root = 0
	s.restoreCache = make(map[uint64]fs.FsObj)
	return s
}

func (s *Snapshot) Snapshot(root fs.FsObj, st *session.SessionTable, tm *threadmgr.ThreadMgrTable, rft *fences.RecentTable, rc *repl.ReplyCache) []byte {
	// Snapshot the FS tree.
	s.Root = s.snapshotFsTree(root)
	// Snapshot the session table.
	s.St = st.Snapshot()
	// Snapshot the thread manager table.
	s.Tmt = tm.Snapshot()
	// Snapshot the recent fence table.
	s.Rft = rft.Snapshot()
	// Snapshot the reply cache.
	s.Rc = rc.Snapshot()
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error marshalling snapshot: %v", err)
	}
	return b
}

func (s *Snapshot) snapshotFsTree(o fs.FsObj) uint64 {
	var stype Tsnapshot
	switch o.(type) {
	case *dir.DirImpl:
		stype = Tdir
	case *memfs.File:
		stype = Tfile
	case *memfs.Symlink:
		stype = Tsymlink
	case *stats.Stats:
		stype = Tstats
	case *Dev:
		stype = Tsnapshotdev
	default:
		log.Fatalf("Unknown FsObj type in snapshot.snapshotFsTree: %v", reflect.TypeOf(o))
	}
	s.Imap[o.Inum()] = MakeObjSnapshot(stype, o.Snapshot(s.snapshotFsTree))
	return o.Inum()
}

func (s *Snapshot) Restore(mkps protsrv.MkProtServer, rps protsrv.RestoreProtServer, fssrv protsrv.FsServer, pfn threadmgr.ProcessFn, oldRc *repl.ReplyCache, b []byte) (fs.FsObj, *session.SessionTable, *threadmgr.ThreadMgrTable, *fences.RecentTable, *repl.ReplyCache) {
	err := json.Unmarshal(b, s)
	if err != nil {
		log.Fatalf("FATAL error unmarshal file in snapshot.Restore: %v", err)
	}
	s.restoreCache[0] = nil
	root := s.RestoreFsTree(s.Root)
	// Restore the thread manager table.
	tmt := threadmgr.Restore(pfn, s.Tmt)
	// Restore the recent fence table.
	rft := fences.RestoreRecentTable(s.Rft)
	// Restore the session table.
	st := session.RestoreTable(mkps, rps, fssrv, rft, tmt, s.St)
	// Restore the reply cache.
	rc := repl.Restore(s.Rc)
	// Merge with the current replyCache, because some ops may have arrived &
	// begun executing since this snapshot was taken, and they expect some state
	// to be in the reply cache.
	rc.Merge(oldRc)
	return root, st, tmt, rft, rc
}

func (s *Snapshot) RestoreFsTree(inum uint64) fs.FsObj {
	if obj, ok := s.restoreCache[inum]; ok {
		return obj
	}
	snap := s.Imap[inum]
	var o fs.FsObj
	switch snap.Type {
	case Tdir:
		// Make a dir with a nil inode so we don't recurse infinitely when trying
		// to set parent pointers.
		d := dir.MakeDir(nil)
		s.restoreCache[inum] = d
		o = dir.Restore(d, s.RestoreFsTree, snap.Data)
	case Tfile:
		o = memfs.RestoreFile(s.RestoreFsTree, snap.Data)
	case Tsymlink:
		o = memfs.RestoreSymlink(s.RestoreFsTree, snap.Data)
	case Tstats:
		o = stats.Restore(s.RestoreFsTree, snap.Data)
	case Tsnapshotdev:
		o = MakeDev(s.fssrv, nil, s.RestoreFsTree(s.Root).(fs.Dir))
	default:
		log.Fatalf("FATAL error unknown type in Snapshot.restore: %v", snap.Type)
		o = nil
	}
	// Store the object in the restore cache.
	s.restoreCache[inum] = o
	return o
}
