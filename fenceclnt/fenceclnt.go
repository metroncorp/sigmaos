package fenceclnt

import (
	"encoding/json"
	"fmt"
	"log"

	"ulambda/fslib"
	np "ulambda/ninep"
	"ulambda/proc"
)

//
// The fence interface for procs.  A proc uses a fence to block its
// requests to servers when its fence has become stale.  A motivating
// use case is a proc that was a primary but has been surplanted by a
// new primary, but the old one doesn't know about the new primary
// (e.g., due to a network partition).  The fence ensures that servers
// will deny requests by the proc after it isn't the primary anymore.
//
// Procs name a fence using a pathname for the file associated with
// the fence.  Internally, a fence object consists of a sequence
// number and a fenceid, which is a tuple containing the pathname at
// the server holding the fence file and a server id.  The server
// increases the seqno associated with the fence pathname when a proc
// creates a new file for that pathname or when a proc modifies the
// file.  The fence interface fences all requests of a proc to servers
// by asking the server hosting the fence file for a fence object and
// asking fsclnt to register that fence object; fsclnt registers the
// fence on all open sessions to servers and will register the fence
// on future sessions that the proc opens when it mounts a new server.
//
// When a proc issues a request on any session, the receiving server
// checks if the registered fence for that session is still valid by
// checking that the seqno is equal or larger than the last seen fence
// for that fenceid on any of the server's sessions. If valid the
// server serves the request, otherwise, it returns a stale-fence
// error.
//
// Procs use the fence interface in two major ways: as write fences
// and read fences.  Write fences are intended for a proc that will
// modify the file associated with the fence and read fences are
// intended for procs that will read the file file.  The read fence
// will alert the reader when the file is modified by a proc that has
// the write fence for the file.  Procs can obtain Write fences in two
// ways: 1) AcquireFenceW(), which exclusively creates an ephemeral
// file for writing or blocks until the server removes the file
// because the proc that created it crashed; 2) OpenFenceFrom(), which
// opens an existing fence file with the new content from the file
// <from>.  A proc obtains a read fence by calling AcquireFenceR(),
// which opens a fence file, potentially blocking until it exists.
// When a writing proc updates a fence file, it is its job to ask for
// a new fence from the server and updates its registered fences.
//
// An intended use case is electing a primary.  Candidate procs invoke
// AcquireFenceW(), and one will succeed and become the primary.  If
// it crashes, one of the other candidates will succeed, and be the
// next primary.  The typical content of the fence file is the
// pathname for the primary.  Backups obtain a read fence for the
// fence file by calling AcquireFenceR().  If there is no primary,
// they will block. Once there is a primary, AcquireFenceR() succeeds
// and they know about the new primary.  If a new primary is elected
// and registers an updated fence, procs holding a read fence will
// observe a stale-fence error. They can then invoke AcquireFenceR()
// again to learn about the new primary.
//
// Another use case is maintaining a configuration file for a service
// (e.g., a mapping from shards to servers).  This file typically
// exists during the life-time of the service. After elected, the
// primary obtains a write fence for this file using OpenFenceFrom(),
// updating the file atomically with the content of file <from> (e.g.,
// a new configuration).  Procs that are clients of the service use
// AcquireFenceR() to obtain a read fence for the file.  Whenever the
// primary updates the fence file (e.g., with a new mapping from
// shards to servers), the server increases the seqno number, and the
// primary asks for a new fence and updates the registered ones.
// Client procs will observe a stale-fence error, and invoke
// AcquireFenceR() again to obtain a new fence and the updates config.
//
// A replicated service can combine these two use cases. One fence
// file to elect a primary and one for the configuration. A proc uses
// the first fence file to become a primary, then do some recovery
// work, including constructing a new configuration to reflect the
// recovered service, which it posts through the second fence file.
//
// Fences are not locks: a fence holder can lose a fence at any time
// (i.e., before the holder releases it).  The read/write usages also
// doesn't correspond to read/write mode of locks: in fact, it is
// common for one proc to have a write fence for a file and another
// proc having a read fence for the same file at the same time.
// Similar to locks, however, it is the application's responsibility
// to use fences correctly.
//

type FenceClnt struct {
	fenceName string // pathname for the fence file
	*fslib.FsLib
	perm    np.Tperm
	mode    np.Tmode
	f       *np.Tfence
	lastSeq np.Tseqno
	paths   []string
}

func MakeFenceClnt(fsl *fslib.FsLib, name string, perm np.Tperm, paths []string) *FenceClnt {
	fc := &FenceClnt{}
	fc.fenceName = name
	fc.FsLib = fsl
	fc.perm = perm
	fc.paths = paths
	return fc
}

func (fc *FenceClnt) IsFenced() *np.Tfence {
	return fc.f
}

func (fc *FenceClnt) Name() string {
	return fc.fenceName
}

func (fc *FenceClnt) Fence() (np.Tfence, error) {
	if fc.f == nil {
		return np.Tfence{}, fmt.Errorf("Fence: not acquired %v\n", fc.fenceName)
	}
	return *fc.f, nil
}

func (fc *FenceClnt) registerPaths(fence np.Tfence, paths []string) error {
	for _, p := range paths {
		err := fc.RegisterFence(fence, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fc *FenceClnt) deregisterPaths(fence np.Tfence) error {
	for _, p := range fc.paths {
		err := fc.DeregisterFence(fence, p)
		if err != nil {
			return err
		}
	}
	return nil
}

// XXX register/update may fail because another client has seen a more
// recent seqno, which the server may have not told us about because
// it lost that info due to it crashing. in that case, tell the server
// to use that more recent seqno (if the fence was acquired in write
// mode, which means this client is the current fence holder).
func (fc *FenceClnt) registerFence(mode np.Tmode) error {
	fence, err := fc.MakeFence(fc.fenceName, mode)
	if err != nil {
		log.Printf("%v: MakeFence %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	// log.Printf("%v: MakeFence %v fence %v", proc.GetProgram(), fc.fenceName, fence)
	if fc.lastSeq > fence.Seqno {
		log.Fatalf("%v: FATAL MakeFence bad fence %v last seq %v\n", proc.GetProgram(), fence, fc.lastSeq)
	}
	if fc.f == nil {
		fc.mode = mode
	}
	err = fc.registerPaths(fence, fc.paths)
	if err != nil {
		log.Printf("%v: registerFence %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	fc.lastSeq = fence.Seqno
	fc.f = &fence
	return nil
}

// Acquire a write fence, which may block. Once caller's create
// succeeds, initialize the file with b and register a fence for the
// file.
//
// XXX cleanup on failure XXX create and write atomic
func (fc *FenceClnt) AcquireFenceW(b []byte) error {
	fd, err := fc.Create(fc.fenceName, fc.perm|np.DMTMP, np.OWRITE|np.OWATCH)
	if err != nil {
		log.Printf("%v: Create %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}

	_, err = fc.Write(fd, b)
	if err != nil {
		log.Printf("%v: Write %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	fc.Close(fd)
	return fc.registerFence(np.OWRITE)
}

// Open an existing file as a fence and register the fence.
func (fc *FenceClnt) OpenFenceFrom(from string) error {
	err := fc.Rename(from, fc.fenceName)
	if err != nil {
		log.Printf("%v: Rename %v to %v err %v", proc.GetProgram(), from, fc.fenceName, err)
		return err
	}
	return fc.registerFence(0)
}

// Acquire a read fence, which may block until a writer has created
// the file.  Tell servers to fence our operations.
func (fc *FenceClnt) AcquireFenceR() ([]byte, error) {
	ch := make(chan bool)
	for {
		// log.Printf("%v: file watch %v\n", proc.GetProgram(), fc.fenceName)
		b, err := fc.ReadFileWatch(fc.fenceName, func(string, error) {
			ch <- true
		})
		if err != nil && np.IsErrNotfound(err) {
			// log.Printf("%v: file watch wait %v\n", proc.GetProgram(), fc.fenceName)
			<-ch
		} else if err != nil {
			return nil, err
		} else {
			// log.Printf("%v: file watch return %v\n", proc.GetProgram(), fc.fenceName)
			return b, fc.registerFence(np.OREAD)
		}
	}
}

// Release fence, which deregisters it.
func (fc *FenceClnt) ReleaseFence() error {
	// First deregister fence
	if fc.f == nil {
		log.Fatalf("%v: FATAL ReleaseFence %v\n", proc.GetProgram(), fc.fenceName)
	}
	err := fc.deregisterPaths(*fc.f)
	if err != nil {
		log.Printf("%v: deregister %v err %v\n", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	fc.f = nil
	// Then, remove file so that the next acquirer can run
	if fc.mode == np.OWRITE {
		err := fc.Remove(fc.fenceName)
		if err != nil {
			log.Printf("%v: Remove %v err %v", proc.GetProgram(), fc.fenceName, err)
			return err
		}
	}
	return nil
}

// Remove fence.  The caller better sure there is no client relying on
// server checking this fence.  The caller must have unregistered the
// fence.
func (fc *FenceClnt) RemoveFence() error {
	if fc.f != nil {
		log.Fatalf("%v: FATAL RmFence %v\n", proc.GetProgram(), fc.fenceName)
	}
	err := fc.AcquireFenceW([]byte{})
	if err != nil {
		return err
	}
	err = fc.RmFence(*fc.f, fc.fenceName)
	if err != nil {
		return err
	}
	return fc.ReleaseFence()
}

func (fc *FenceClnt) FencePaths(paths []string) error {
	fence, err := fc.Fence()
	if err != nil {
		return err
	}
	err = fc.registerPaths(fence, paths)
	if err != nil {
		log.Printf("%v: registerFence %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	fc.paths = append(fc.paths, paths...)
	return nil
}

//
// A few writer operations that a fence writer can perform. They will
// increase the fence's seqno, and registerFence will update servers
// to use the new fence.
//

func (fc *FenceClnt) SetFenceFile(b []byte) error {
	_, err := fc.SetFile(fc.fenceName, b)
	if err != nil {
		log.Printf("%v: SetFenceFile %v err %v", proc.GetProgram(), fc.fenceName, err)
		return err
	}
	return fc.registerFence(0)
}

func (fc *FenceClnt) MakeFenceFileFrom(from string) error {
	err := fc.Rename(from, fc.fenceName)
	if err != nil {
		log.Printf("%v: Rename %v to %v err %v", proc.GetProgram(), from, fc.fenceName, err)
		return err
	}
	return fc.registerFence(0)
}

//
// Convenience function
//

func (fc *FenceClnt) AcquireConfig(v interface{}) error {
	//log.Printf("%v: start AcquireConfig %v\n", proc.GetProgram(), fc.Name())
	b, err := fc.AcquireFenceR()
	if err != nil {
		log.Printf("%v: AcquireConfig %v err %v\n", proc.GetProgram(), fc.Name(), err)
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	//log.Printf("%v: AcquireConfig %v %v\n", proc.GetProgram(), fc.Name(), v)
	return nil
}
