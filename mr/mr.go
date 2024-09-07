package mr

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mitchellh/mapstructure"

	db "sigmaos/debug"
	"sigmaos/fslib"
	sp "sigmaos/sigmap"
)

// Map and reduce functions produce and consume KeyValue pairs
type KeyValue struct {
	Key   string
	Value string
}

type EmitT func(key, value string) error

// The mr library calls the reduce function once for each key
// generated by the map tasks, with a list of all the values created
// for that key by any map task.
type ReduceT func(string, []string, EmitT) error

// The mr library calls the map function for each line of input, which
// is passed in as an io.Reader.  The map function outputs its values
// by calling an emit function and passing it a KeyValue.
type MapT func(string, *bufio.Scanner, EmitT) error

// for sorting by key.
type ByKey []*KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

// Use Khash(key) % NReduce to choose the reduce task number for each
// KeyValue emitted by Map.
func Khash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// An input split
type Split struct {
	File   string     `json:"File"`
	Offset sp.Toffset `json:"Offset"`
	Length sp.Tlength `json:"Length"`
}

func (s Split) String() string {
	return fmt.Sprintf("{f %s o %v l %v}", s.File, humanize.Bytes(uint64(s.Offset)), humanize.Bytes(uint64(s.Length)))
}

type Bin []Split

func (b Bin) String() string {
	r := fmt.Sprintf("bins (%d): [ %v, ", len(b), b[0])
	sum := sp.Tlength(b[0].Length)
	for i, s := range b[1:] {
		if s.File == b[i].File {
			r += fmt.Sprintf("_ o %v l %v,", humanize.Bytes(uint64(s.Offset)), humanize.Bytes(uint64(s.Length)))
		} else {
			r += fmt.Sprintf("[ %v, ", s)
		}
		sum += s.Length
	}
	r += fmt.Sprintf("] (sum %v)\n", humanize.Bytes(uint64(sum)))
	return r
}

// Result of mapper or reducer
type Result struct {
	IsM  bool       `json:"IsM"`
	Task string     `json:"Task"`
	In   sp.Tlength `json:"In"`
	Out  sp.Tlength `json:"Out"`
	Ms   int64      `json:"Ms"`
}

func newResult(data interface{}) *Result {
	r := &Result{}
	mapstructure.Decode(data, r)
	return r
}

// Each bin has a slice of splits.  Assign splits of files to a bin
// until the bin is full
func NewBins(fsl *fslib.FsLib, dir string, maxbinsz, splitsz sp.Tlength) ([]Bin, error) {
	bins := make([]Bin, 0)
	binsz := uint64(0)
	bin := Bin{}

	anydir := strings.ReplaceAll(dir, "~local", "~any")
	sts, err := fsl.GetDir(anydir)
	if err != nil {
		return nil, err
	}
	for _, st := range sts {
		for i := uint64(0); ; {
			n := uint64(splitsz)
			if i+n > st.LengthUint64() {
				n = st.LengthUint64() - i
			}
			split := Split{dir + "/" + st.Name, sp.Toffset(i), sp.Tlength(n)}
			bin = append(bin, split)
			binsz += n
			if binsz+uint64(splitsz) > uint64(maxbinsz) { // bin full?
				bins = append(bins, bin)
				bin = Bin{}
				binsz = uint64(0)
			}
			if n < uint64(splitsz) { // next file
				break
			}
			i += n
		}
	}
	if binsz > 0 {
		bins = append(bins, bin)
	}
	db.DPrintf(db.MR, "Bin sizes: %v", bins)
	return bins, nil
}
