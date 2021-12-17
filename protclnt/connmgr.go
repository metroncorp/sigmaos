package protclnt

import (
	"errors"
	"strings"
	"sync"

	db "ulambda/debug"
	"ulambda/lease"
	"ulambda/netclnt"
	np "ulambda/ninep"
)

// XXX duplicate
const (
	Msglen = 64 * 1024
)

type ConnMgr struct {
	mu      sync.Mutex
	session np.Tsession
	seqno   *np.Tseqno
	conns   map[string]*netclnt.NetClnt
}

func makeConnMgr(session np.Tsession, seqno *np.Tseqno) *ConnMgr {
	cm := &ConnMgr{}
	cm.conns = make(map[string]*netclnt.NetClnt)
	cm.session = session
	cm.seqno = seqno
	return cm
}

func (cm *ConnMgr) exit() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for addr, conn := range cm.conns {
		db.DLPrintf("9PCHAN", "exit close connection to %v\n", addr)
		conn.Close()
		delete(cm.conns, addr)
	}
}

// XXX Make array
func (cm *ConnMgr) allocConn(addrs []string) (*netclnt.NetClnt, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Store as concatenation of addresses
	key := strings.Join(addrs, ",")

	var err error
	conn, ok := cm.conns[key]
	if !ok {
		conn, err = netclnt.MkNetClnt(addrs)
		if err == nil {
			cm.conns[key] = conn
		}
	}
	return conn, err
}

func (cm *ConnMgr) lookupConn(addrs []string) (*netclnt.NetClnt, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conn, ok := cm.conns[strings.Join(addrs, ",")]
	return conn, ok
}

func (cm *ConnMgr) makeCall(dst []string, req np.Tmsg) (np.Tmsg, error) {
	conn, err := cm.allocConn(dst)
	if err != nil {
		return nil, err
	}
	reqfc := &np.Fcall{}
	reqfc.Type = req.Type()
	reqfc.Msg = req
	reqfc.Session = cm.session
	reqfc.Seqno = cm.seqno.Next()
	repfc, err := conn.RPC(reqfc)
	if err != nil {
		return nil, err
	}
	return repfc.Msg, nil
}

func (cm *ConnMgr) disconnect(dst []string) bool {
	conn, ok := cm.lookupConn(dst)
	if !ok {
		return false
	}
	conn.Close()
	return true
}

func (cm *ConnMgr) mcast(ch chan error, dst []string, req np.Tmsg) {
	if reply, err := cm.makeCall(dst, req); err != nil {
		ch <- err
	} else {
		if rmsg, ok := reply.(np.Rerror); ok {
			ch <- errors.New(rmsg.Ename)
		} else {
			ch <- nil
		}
	}
}

func (cm *ConnMgr) registerLease(lease *lease.Lease) error {
	ch := make(chan error)
	cm.mu.Lock()
	n := 0
	args := np.Tlease{lease.Fn, lease.Qid}
	for addr, _ := range cm.conns {
		n += 1
		go cm.mcast(ch, strings.Split(addr, ","), args)
	}
	cm.mu.Unlock()
	var err error
	for i := 0; i < n; i++ {
		r := <-ch
		// Ignore EOF, since we cannot talk to that server
		// anymore.  We may try to reconnect and then we will
		// register again.
		if r != nil && r.Error() != "EOF" {
			err = r
		}
	}
	return err
}

// XXX deduplicate
func (cm *ConnMgr) deregisterLease(path []string) error {
	ch := make(chan error)
	cm.mu.Lock()
	n := 0
	args := np.Tunlease{path}
	for addr, _ := range cm.conns {
		n += 1
		go cm.mcast(ch, strings.Split(addr, ","), args)
	}
	cm.mu.Unlock()
	var err error
	for i := 0; i < n; i++ {
		r := <-ch
		if r != nil && r.Error() != "EOF" {
			err = r
		}
	}
	return err
}
