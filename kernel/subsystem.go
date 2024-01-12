package kernel

import (
	"fmt"
	"os/exec"
	"syscall"

	"sigmaos/container"
	db "sigmaos/debug"
	"sigmaos/fslib"
	"sigmaos/kernelsubinfo"
	"sigmaos/port"
	"sigmaos/proc"
	"sigmaos/procclnt"
	sp "sigmaos/sigmap"
)

const NPORT_PER_CONTAINER = 20

// XXX Make interface smaller
type Subsystem interface {
	GetProc() *proc.Proc
	GetHow() proc.Thow
	GetCrashed() bool
	GetContainer() *container.Container
	SetWaited(bool)
	GetWaited() bool
	Evict() error
	Wait() error
	Kill() error
	SetCPUShares(shares int64) error
	GetCPUUtil() (float64, error)
	GetIp(fsl *fslib.FsLib) *sp.Taddr
	AssignToRealm(realm sp.Trealm, ptype proc.Ttype) error
	AllocPort(p sp.Tport) (*port.PortBinding, error)
	Run(how proc.Thow, kernelId string, localIP sp.Thost) error
}

type KernelSubsystem struct {
	*procclnt.ProcClnt
	k         *Kernel
	p         *proc.Proc
	how       proc.Thow
	cmd       *exec.Cmd
	container *container.Container
	waited    bool
	crashed   bool
}

func (ss *KernelSubsystem) GetProc() *proc.Proc {
	return ss.p
}

func (ss *KernelSubsystem) GetContainer() *container.Container {
	return ss.container
}

func (ss *KernelSubsystem) GetHow() proc.Thow {
	return ss.how
}

func (ss *KernelSubsystem) GetCrashed() bool {
	return ss.crashed
}

func (ss *KernelSubsystem) GetWaited() bool {
	return ss.waited
}

func (ss *KernelSubsystem) SetWaited(w bool) {
	ss.waited = w
}

func (ss *KernelSubsystem) String() string {
	s := fmt.Sprintf("subsystem %p: [proc %v how %v]", ss, ss.p, ss.how)
	return s
}

func newSubsystemCmd(pclnt *procclnt.ProcClnt, k *Kernel, p *proc.Proc, how proc.Thow, cmd *exec.Cmd) Subsystem {
	return &KernelSubsystem{pclnt, k, p, how, cmd, nil, false, false}
}

func newSubsystem(pclnt *procclnt.ProcClnt, k *Kernel, p *proc.Proc, how proc.Thow) Subsystem {
	return newSubsystemCmd(pclnt, k, p, how, nil)
}

func (k *Kernel) bootSubsystemWithMcpu(program string, args []string, how proc.Thow, mcpu proc.Tmcpu) (Subsystem, error) {
	pid := sp.GenPid(program)
	p := proc.NewPrivProcPid(pid, program, args, true)
	p.GetProcEnv().SetLocalIP(k.ip)
	p.SetMcpu(mcpu)
	ss := newSubsystem(k.ProcClnt, k, p, how)
	return ss, ss.Run(how, k.Param.KernelId, k.ip)
}

func (k *Kernel) bootSubsystem(program string, args []string, how proc.Thow) (Subsystem, error) {
	return k.bootSubsystemWithMcpu(program, args, how, 0)
}

func (s *KernelSubsystem) Evict() error {
	return s.k.EvictKernelProc(s.p.GetPid(), s.how)
}

func (s *KernelSubsystem) Run(how proc.Thow, kernelId string, localIP sp.Thost) error {
	if how == proc.HLINUX || how == proc.HSCHEDD {
		cmd, err := s.SpawnKernelProc(s.p, s.how, kernelId)
		if err != nil {
			return err
		}
		s.cmd = cmd
	} else {
		if err := s.NewProc(s.p, proc.HDOCKER, kernelId); err != nil {
			return err
		}
		// XXX don't hard code
		h := sp.SIGMAHOME
		s.p.AppendEnv("PATH", h+"/bin/user:"+h+"/bin/user/common:"+h+"/bin/kernel:/usr/sbin:/usr/bin:/bin")
		s.p.FinalizeEnv(localIP, sp.Tpid(proc.NOT_SET))
		var r *port.Range
		up := sp.NO_PORT
		if s.k.Param.Overlays {
			r = &port.Range{FPORT, LPORT}
			up = r.Fport
		}
		c, err := container.StartPContainer(s.p, kernelId, r, up, s.k.Param.GVisor)
		if err != nil {
			return err
		}
		s.container = c
	}
	err := s.WaitStartKernelProc(s.p.GetPid(), how)
	return err
}

func (ss *KernelSubsystem) AssignToRealm(realm sp.Trealm, ptype proc.Ttype) error {
	return ss.container.AssignToRealm(realm, ptype)
}

func (ss *KernelSubsystem) SetCPUShares(shares int64) error {
	return ss.container.SetCPUShares(shares)
}

func (ss *KernelSubsystem) GetCPUUtil() (float64, error) {
	return ss.container.GetCPUUtil()
}

func (ss *KernelSubsystem) AllocPort(p sp.Tport) (*port.PortBinding, error) {
	if p == sp.NO_PORT {
		return ss.container.AllocPort()
	} else {
		return ss.container.AllocPortOne(p)
	}
}

func (ss *KernelSubsystem) GetIp(fsl *fslib.FsLib) *sp.Taddr {
	return kernelsubinfo.GetSubsystemInfo(fsl, sp.KPIDS, ss.p.GetPid().String()).Addr
}

// Send SIGTERM to a system.
func (s *KernelSubsystem) Terminate() error {
	db.DPrintf(db.KERNEL, "Terminate %v %v\n", s.cmd.Process.Pid, s.cmd)
	if s.how != proc.HLINUX {
		db.DFatalf("Tried to terminate a kernel subsystem spawned through procd: %v", s.p)
	}
	return syscall.Kill(s.cmd.Process.Pid, syscall.SIGTERM)
}

// Kill a subsystem, either by sending SIGKILL or Evicting it.
func (s *KernelSubsystem) Kill() error {
	s.crashed = true
	db.DPrintf(db.KERNEL, "Kill %v\n", s)
	if s.p.GetProgram() == "knamed" {
		return stopKNamed(s.cmd)
	}
	if s.how == proc.HSCHEDD || s.how == proc.HDOCKER {
		db.DPrintf(db.ALWAYS, "Killing a kernel subsystem spawned through %v: %v", s.p, s.how)
		err := s.EvictKernelProc(s.p.GetPid(), s.how)
		if err != nil {
			db.DPrintf(db.ALWAYS, "Error killing procd-spawned kernel proc: %v err %v", s.p.GetPid(), err)
		}
		return err
	}
	db.DPrintf(db.ALWAYS, "kill %v\n", s.cmd.Process.Pid)
	return syscall.Kill(s.cmd.Process.Pid, syscall.SIGKILL)
}

func (s *KernelSubsystem) Wait() error {
	db.DPrintf(db.KERNEL, "Wait subsystem for %v", s)
	defer db.DPrintf(db.KERNEL, "Wait subsystem done for %v", s)

	if !s.GetCrashed() {
		s.SetWaited(true)
		// Only wait if this proc has not been waited for already, since calling
		// WaitExit twice leads to an error.
		status, err := s.WaitExitKernelProc(s.p.GetPid(), s.how)
		if err != nil || !status.IsStatusEvicted() {
			db.DPrintf(db.ALWAYS, "shutdown susbystem [%v] exit with status %v err %v", s.p.GetPid(), status, err)
			return err
		}
	}

	if s.how == proc.HSCHEDD || s.how == proc.HDOCKER {
		// Do nothing (already waited)
		return nil
	} else {
		db.DPrintf(db.KERNEL, "Wait subsystem via cmd %v", s)
		if err := s.cmd.Wait(); err != nil {
			return err
		}
	}
	if s.container != nil {
		db.DPrintf(db.KERNEL, "Container shutdown %v", s)
		return s.container.Shutdown()
	}
	return nil
}
