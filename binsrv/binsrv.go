// The binsrv package serves sigmaos binaries to the linux kernel.  It
// fetches the binary using the sigmaos pathname and caches them
// locally.  This allow support demand paging: the kernel can start
// running the binary before the complete binary has been downloaded.
//
// binsrv is based on go-fuse's loopback.
package binsrv

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	db "sigmaos/debug"
	"sigmaos/fslib"
	"sigmaos/proc"
	"sigmaos/rpcclnt"
	"sigmaos/sigmaclnt"
	sp "sigmaos/sigmap"
	"sigmaos/sigmarpcchan"
	"sigmaos/uprocclnt"
)

const (
	// binfsd mounts itself here:
	BINFSMNT = "/mnt/binfs/"

	// The directory /tmp/sigmaos-bin/realms/<realm> in the host file
	// system is mounted here by uprocd:
	BINCACHE = sp.SIGMAHOME + "/bin/user/"

	DEBUG = false
)

func BinPath(program string) string {
	return BINFSMNT + program
}

func binCachePath(program string) string {
	return BINCACHE + program
}

type binFsRoot struct {
	Path     string // the directory that holds cached binaries
	bincache *bincache
}

func (r *binFsRoot) newNode(parent *fs.Inode, name string, sz sp.Tsize) fs.InodeEmbedder {
	n := &binFsNode{
		RootData: r,
		name:     name,
		sz:       sz,
	}
	return n
}

type binFsNode struct {
	fs.Inode

	RootData *binFsRoot
	name     string
	sz       sp.Tsize
}

func (n *binFsNode) String() string {
	return fmt.Sprintf("{N %q}", n.path())
}

func newBinRoot(kernelId string, sc *sigmaclnt.SigmaClnt, updc *uprocclnt.UprocdClnt) (fs.InodeEmbedder, error) {
	var st syscall.Stat_t
	err := syscall.Stat(BINCACHE, &st)
	if err != nil {
		return nil, err
	}
	root := &binFsRoot{
		bincache: newBinCache(kernelId, sc, updc),
	}
	return root.newNode(nil, "", 0), nil
}

func RunBinFS(kernelId, uprocdpid string) error {
	pe := proc.GetProcEnv()

	proc.SetSigmaDebugPid("binfsd-" + uprocdpid)

	db.DPrintf(db.BINSRV, "MkDir %q", BINFSMNT)

	if err := os.MkdirAll(BINFSMNT, 0750); err != nil {
		return err
	}

	db.DPrintf(db.BINSRV, "%s", db.LsDir(BINCACHE))

	sc, err := sigmaclnt.NewSigmaClnt(pe)
	if err != nil {
		return err
	}

	// if sts, err := sc.GetDir(sp.CHUNKD); err == nil {
	// 	db.DPrintf(db.ALWAYS, "chunksrvs %v", sp.Names(sts))
	// } else {
	// 	db.DPrintf(db.ALWAYS, "chunksrvs err %v", err)
	// }

	pn := path.Join(sp.SCHEDD, kernelId, sp.UPROCDREL, uprocdpid)
	ch, err := sigmarpcchan.NewSigmaRPCCh([]*fslib.FsLib{sc.FsLib}, pn)
	if err != nil {
		db.DPrintf(db.ERROR, "rpcclnt err %v", err)
		return err
	}
	rc := rpcclnt.NewRPCClnt(ch)
	updc := uprocclnt.NewUprocdClnt(sp.Tpid(uprocdpid), rc)

	loopbackRoot, err := newBinRoot(kernelId, sc, updc)
	if err != nil {
		return err
	}
	sec := time.Second
	opts := &fs.Options{
		AttrTimeout:  &sec,
		EntryTimeout: &sec,

		NullPermissions: true, // Leave file permissions on "000" files as-is

		MountOptions: fuse.MountOptions{
			Debug:  DEBUG,
			FsName: BINCACHE, // First column in "df -T": original dir
			Name:   "binfs",  // Second column in "df -T" will be shown as "fuse." + Name
		},
	}
	opts.MountOptions.Options = append(opts.MountOptions.Options, "ro")

	server, err := fs.Mount("/mnt/binfs", loopbackRoot, opts)
	if err != nil {
		return err
	}

	// Tell ExecBinSrv we are running
	if _, err := io.WriteString(os.Stdout, "r"); err != nil {
		return err
	}
	go func() {
		buf := make([]byte, 1)
		if _, err := io.ReadFull(os.Stdin, buf); err != nil {
			db.DFatalf("read pipe err %v\n", err)
		}
		db.DPrintf(db.BINSRV, "exiting\n")
		server.Unmount()
		os.Exit(0)
	}()

	server.Wait()
	db.DPrintf(db.ALWAYS, "Wait returned\n")
	return nil
}

type BinSrvCmd struct {
	cmd *exec.Cmd
	out io.WriteCloser
}

func ExecBinSrv(kernelId, uprocdpid string) (*BinSrvCmd, error) {
	cmd := exec.Command("binfsd", kernelId, uprocdpid)
	// cmd.Env = p.GetEnv()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		db.DPrintf(db.BINSRV, "Error start %v %v", cmd, err)
		return nil, err
	}
	buf := make([]byte, 1)
	if _, err := io.ReadFull(stdout, buf); err != nil {
		db.DPrintf(db.BINSRV, "read pipe err %v\n", err)
		return nil, err
	}

	return &BinSrvCmd{cmd: cmd, out: stdin}, nil
}

func (bsc *BinSrvCmd) Shutdown() error {
	if _, err := io.WriteString(bsc.out, "e"); err != nil {
		return err
	}
	return nil
}
