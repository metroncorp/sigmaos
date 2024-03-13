package reader

import (
	"io"

	db "sigmaos/debug"
	sp "sigmaos/sigmap"
)

type ReaderI interface {
	Read(sp.Toffset, []byte) (int, error)
	Close() error
}

type Reader struct {
	rdr    ReaderI
	path   string
	off    sp.Toffset
	eof    bool
	fenced bool
}

func (rdr *Reader) Path() string {
	return rdr.path
}

func (rdr *Reader) Nbytes() sp.Tlength {
	return sp.Tlength(rdr.off)
}

func (rdr *Reader) Read(p []byte) (int, error) {
	db.DPrintf(db.ALWAYS, "%p Read[%v] %v bytes", rdr, rdr.Path(), len(p))
	if len(p) == 0 {
		return 0, nil
	}
	if rdr.eof {
		return 0, io.EOF
	}
	n, err := rdr.rdr.Read(rdr.off, p)
	if err != nil {
		db.DPrintf(db.READER_ERR, "Read %v err %v\n", rdr.path, err)
		return 0, err
	}
	db.DPrintf(db.ALWAYS, "%p Read[%v] %v bytes returned %v \nslice %v", rdr, rdr.Path(), len(p), n, p[:n])
	if n == 0 {
		rdr.eof = true
		return 0, io.EOF
	}
	if int(n) < len(p) {
		db.DPrintf(db.READER_ERR, "Read short %v %v %v\n", rdr.path, len(p), n)
		db.DPrintf(db.ERROR, "Read short %v %v %v\n", rdr.path, len(p), n)
		err = io.EOF
		rdr.eof = true
	}
	rdr.off += sp.Toffset(n)
	return int(n), err
}

func (rdr *Reader) GetData() ([]byte, error) {
	// XXX too big?
	db.DPrintf(db.ALWAYS, "%p GetData[%v]", rdr, rdr.Path())
	b := make([]byte, sp.MAXGETSET)
	sz, err := rdr.rdr.Read(0, b)
	return b[:sz], err
}

func (rdr *Reader) Close() error {
	err := rdr.rdr.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rdr *Reader) Unfence() {
	rdr.fenced = false
}

func (rdr *Reader) Reader() ReaderI {
	return rdr.rdr
}

func NewReader(rdr ReaderI, path string) *Reader {
	return &Reader{rdr, path, 0, false, true}
}
