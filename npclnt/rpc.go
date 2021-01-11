package npclnt

import (
	"errors"

	np "ulambda/ninep"
)

type NpClnt struct {
	cm *ChanMgr
}

func MakeNpClnt(debug bool) *NpClnt {
	npc := &NpClnt{}
	npc.cm = makeChanMgr(debug)
	return npc
}

func (npc *NpClnt) callServer(server string, args np.Tmsg) (np.Tmsg, error) {
	reply, err := npc.cm.makeCall(server, args)
	if err != nil {
		return nil, err
	}
	rmsg, ok := reply.(np.Rerror)
	if ok {
		return nil, errors.New(rmsg.Ename)
	}
	return reply, nil
}

// XXX copying msg once too many?

func (npc *NpClnt) Attach(server string, fid np.Tfid, path []string) (*np.Rattach, error) {
	args := np.Tattach{fid, np.NoFid, "fk", ""}
	reply, err := npc.callServer(server, args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rattach)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

type NpChan struct {
	server string
	cm     *ChanMgr
}

func (npc *NpClnt) MakeNpChan(server string) *NpChan {
	npchan := &NpChan{server, npc.cm}
	return npchan
}

func (npch *NpChan) call(args np.Tmsg) (np.Tmsg, error) {
	reply, err := npch.cm.makeCall(npch.server, args)
	if err != nil {
		return nil, err
	}
	rmsg, ok := reply.(np.Rerror)
	if ok {
		return nil, errors.New(rmsg.Ename)
	}
	return reply, nil
}

func (npch *NpChan) Flush(tag np.Ttag) error {
	args := np.Tflush{tag}
	reply, err := npch.call(args)
	if err != nil {
		return err
	}
	_, ok := reply.(np.Rflush)
	if !ok {
		return errors.New("Not correct reply msg")
	}
	return err
}

func (npch *NpChan) Walk(fid np.Tfid, nfid np.Tfid, path []string) (*np.Rwalk, error) {
	args := np.Twalk{fid, nfid, path}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rwalk)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Create(fid np.Tfid, name string, perm np.Tperm, mode np.Tmode) (*np.Rcreate, error) {
	args := np.Tcreate{fid, name, perm, mode}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rcreate)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Remove(fid np.Tfid) error {
	args := np.Tremove{fid}
	reply, err := npch.call(args)
	if err != nil {
		return err
	}
	_, ok := reply.(np.Rremove)
	if !ok {
		return errors.New("Not correct reply msg")
	}
	return err
}

func (npch *NpChan) Clunk(fid np.Tfid) error {
	args := np.Tclunk{fid}
	reply, err := npch.call(args)
	if err != nil {
		return err
	}
	_, ok := reply.(np.Rclunk)
	if !ok {
		return errors.New("Not correct reply msg")
	}
	return err
}

func (npch *NpChan) Open(fid np.Tfid, mode np.Tmode) (*np.Ropen, error) {
	args := np.Topen{fid, mode}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Ropen)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Read(fid np.Tfid, offset np.Toffset, cnt np.Tsize) (*np.Rread, error) {
	args := np.Tread{fid, offset, cnt}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rread)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Write(fid np.Tfid, offset np.Toffset, data []byte) (*np.Rwrite, error) {
	args := np.Twrite{fid, offset, data}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rwrite)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Stat(fid np.Tfid) (*np.Rstat, error) {
	args := np.Tstat{fid}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rstat)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Wstat(fid np.Tfid, st *np.Stat) (*np.Rwstat, error) {
	args := np.Twstat{fid, 0, *st}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rwstat)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}

func (npch *NpChan) Mkpipe(fid np.Tfid, name string, perm np.Tperm) (*np.Rmkpipe, error) {
	args := np.Tmkpipe{fid, name, perm, 0}
	reply, err := npch.call(args)
	if err != nil {
		return nil, err
	}
	msg, ok := reply.(np.Rmkpipe)
	if !ok {
		return nil, errors.New("Not correct reply msg")
	}
	return &msg, err
}
