// The netclnt package establishes a TCP connection to a server.
package clnt

import (
	"net"
	"sync"

	db "sigmaos/debug"
	dialproxyclnt "sigmaos/dialproxy/clnt"
	"sigmaos/proc"
	"sigmaos/serr"
	sp "sigmaos/sigmap"
)

type NetClnt struct {
	mu     sync.Mutex
	pe     *proc.ProcEnv
	npc    *dialproxyclnt.DialProxyClnt
	conn   net.Conn
	ep     *sp.Tendpoint
	addr   *sp.Taddr
	closed bool
	realm  sp.Trealm
}

func NewNetClnt(pe *proc.ProcEnv, npc *dialproxyclnt.DialProxyClnt, ep *sp.Tendpoint) (*NetClnt, *serr.Err) {
	db.DPrintf(db.NETCLNT, "NewNetClnt to %v\n", ep)
	nc := &NetClnt{
		pe:  pe,
		npc: npc,
		ep:  ep,
	}
	if err := nc.connect(ep); err != nil {
		db.DPrintf(db.NETCLNT_ERR, "NewNetClnt connect %v err %v\n", ep, err)
		return nil, err
	}
	return nc, nil
}

func (nc *NetClnt) Conn() net.Conn {
	return nc.conn
}

func (nc *NetClnt) Dst() string {
	return nc.conn.RemoteAddr().String()
}

func (nc *NetClnt) Src() string {
	return nc.conn.LocalAddr().String()
}

func (nc *NetClnt) Close() error {
	return nc.conn.Close()
}

func (nc *NetClnt) connect(ep *sp.Tendpoint) *serr.Err {
	db.DPrintf(db.NETCLNT, "NetClnt connect to any of %v, starting w. %v\n", ep, ep.Addrs()[0])
	for i, addr := range ep.Addrs() {
		if i > 0 {
			ep.Addr = append(ep.Addr[1:], ep.Addr[0])
		}
		c, err := nc.npc.Dial(ep)
		db.DPrintf(db.NETCLNT, "Dial %v addr.Addr %v\n", addr.IPPort(), err)
		if err != nil {
			continue
		}
		nc.conn = c
		nc.addr = addr
		db.DPrintf(db.NETCLNT, "NetClnt connected %v -> %v\n", c.LocalAddr(), nc.addr)
		return nil
	}
	db.DPrintf(db.NETCLNT_ERR, "NetClnt unable to connect to any of %v\n", ep)
	return serr.NewErr(serr.TErrUnreachable, "no connection")
}
