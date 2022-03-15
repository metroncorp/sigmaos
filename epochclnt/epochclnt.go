package epochclnt

import (
	db "ulambda/debug"
	"ulambda/fslib"
	np "ulambda/ninep"
)

type EpochClnt struct {
	*fslib.FsLib
	path string
	perm np.Tperm
}

func MakeEpochClnt(fsl *fslib.FsLib, path string, perm np.Tperm) *EpochClnt {
	ec := &EpochClnt{}
	ec.FsLib = fsl
	ec.path = path
	ec.perm = perm
	return ec
}

func (ec *EpochClnt) Name() string {
	return ec.path
}

// XXX should use writeV
func (ec *EpochClnt) AdvanceEpoch() (string, error) {
	fd, err := ec.CreateOpen(ec.path, ec.perm, np.ORDWR)
	if err != nil {
		db.DLPrintf("EPOCHCLNT_ERR", "CreateOpen %v err %v", ec.path, err)
	}
	defer ec.Close(fd)
	b, err := ec.Read(fd, 100)
	if err != nil {
		db.DLPrintf("EPOCHCLNT_ERR", "Read %v err %v", ec.path, err)
		return "", err
	}
	n := np.Tepoch(0)
	if len(b) > 0 {
		n, err = np.String2Epoch(string(b))
		if err != nil {
			db.DLPrintf("EPOCHCLNT_ERR", "String2Epoch %v err %v", string(b), err)
			return "", err
		}
	}
	n += 1
	if err := ec.Seek(fd, 0); err != nil {
		db.DLPrintf("EPOCHCLNT_ERR", "Seek %v err %v", fd, err)
		return "", err
	}

	db.DLPrintf("EPOCHCLNT", "AdvanceEpoch %v %v", ec.path, n)

	epoch := n.String()
	_, err = ec.Write(fd, []byte(epoch))
	if err != nil {
		db.DLPrintf("EPOCHCLNT_ERR", "Write %v err %v", ec.path, err)
		return "", err
	}
	return epoch, nil
}
