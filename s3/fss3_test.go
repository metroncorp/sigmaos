package fss3

import (
	"bufio"
	"context"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/stretchr/testify/assert"

	"sigmaos/fslib"
	sp "sigmaos/sigmap"
	"sigmaos/test"
)

var ROOT = []string{"a", "gutenberg", "wiki", "img"}

func TestCompile(t *testing.T) {
}

func TestOne(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	dirents, err := ts.GetDir(sp.S3)
	assert.Nil(t, err, "GetDir")

	assert.Equal(t, 1, len(dirents))

	ts.Shutdown()
}

func TestReadOff(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	rdr, err := ts.OpenReader(path.Join(sp.S3, "~local/9ps3/gutenberg/pg-being_ernest.txt"))
	assert.Nil(t, err, "Error ReadOff %v", err)
	rdr.Lseek(1 << 10)
	brdr := bufio.NewReaderSize(rdr.Reader, 1<<16)
	scanner := bufio.NewScanner(brdr)
	l := sp.Tlength(1 << 10)
	n := 0
	for scanner.Scan() {
		line := scanner.Text()
		n += len(line) + 1 // 1 for newline
		if sp.Tlength(n) > l {
			break
		}
	}
	assert.Equal(t, 1072, n)

	ts.Shutdown()
}

func s3Name(ts *test.Tstate) string {
	sts, err := ts.GetDir(sp.S3)
	assert.Nil(ts.T, err, sp.S3)
	assert.Equal(ts.T, 1, len(sts))
	name := path.Join(sp.S3, sts[0].Name)
	return name
}

func TestSymlinkFile(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	dn := s3Name(ts)
	fn := path.Join(dn, "9ps3", "gutenberg/pg-being_ernest.txt")

	_, err := ts.GetFile(fn)
	assert.Nil(t, err, "GetFile")

	fn = dn + "/9ps3" + "//gutenberg/pg-being_ernest.txt"
	_, err = ts.GetFile(fn)
	assert.Nil(t, err, "GetFile")

	ts.Shutdown()
}

func TestSymlinkDir(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	dn := s3Name(ts)

	_, err := ts.GetFile(dn)
	assert.Nil(t, err, "GetFile")

	dirents, err := ts.GetDir(dn + "/" + "9ps3")
	assert.Nil(t, err, "GetDir")

	assert.True(t, fslib.Present(dirents, ROOT))

	ts.Shutdown()
}

func TestReadSplit(t *testing.T) {
	const SPLITSZ = 64 * sp.MBYTE

	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	rdr, err := ts.OpenReader(path.Join(sp.S3, "~local/9ps3/wiki/enwiki-latest-pages-articles-multistream.xml"))
	assert.Nil(t, err)
	err = rdr.Lseek(SPLITSZ)
	assert.Nil(t, err)
	brdr := bufio.NewReaderSize(rdr.Reader, sp.BUFSZ)
	b := make([]byte, SPLITSZ)
	n, err := brdr.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, SPLITSZ, n)
	assert.Equal(t, "s released", string(b[0:10]))

	ts.Shutdown()
}

const NOBJ = 100

func put(clnt *s3.Client, i int, wg *sync.WaitGroup) {
	prefix := "s3test/" + strconv.Itoa(i) + "/"
	for j := 0; j < NOBJ; j++ {
		key := prefix + strconv.Itoa(j)
		input := &s3.PutObjectInput{
			Bucket: aws.String("9ps3"),
			Key:    &key,
		}
		_, err := clnt.PutObject(context.TODO(), input)
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func cleanup(cfg aws.Config) {
	maxKeys := 0
	clnt := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String("9ps3"),
		Prefix: aws.String("s3test/"),
	}
	p := s3.NewListObjectsV2Paginator(clnt, params,
		func(o *s3.ListObjectsV2PaginatorOptions) {
			if v := int32(maxKeys); v != 0 {
				o.Limit = v
			}
		})
	for p.HasMorePages() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return
		}
		wg := &sync.WaitGroup{}
		wg.Add(len(page.Contents))
		for _, obj := range page.Contents {
			input := &s3.DeleteObjectInput{
				Bucket: aws.String("9ps3"),
				Key:    obj.Key,
			}
			go func() {
				defer wg.Done()
				_, err = clnt.DeleteObject(context.TODO(), input)
				if err != nil {
					panic(err)
				}
			}()
		}
		wg.Wait()
	}
}

// Run: go test -v sigmaos/s3 -bench=. -benchtime=1x -run PutObj
func BenchmarkPutObj(b *testing.B) {
	const N = 200

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("me-mit"))
	if err != nil {
		panic(err)
	}

	clnt := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	wg := &sync.WaitGroup{}
	wg.Add(N)

	start := time.Now()
	for i := 0; i < N; i++ {
		go put(clnt, i, wg)
	}
	wg.Wait()
	ms := time.Since(start).Milliseconds()
	s := float64(ms) / 1000
	n := N * NOBJ

	log.Printf("%d took %vms (%.1f file/s)", n, ms, float64(n)/s)

	cleanup(cfg)
}

func TestTwo(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	// Make a second one
	ts.BootFss3d()

	time.Sleep(100 * time.Millisecond)

	dirents, err := ts.GetDir(sp.S3)
	assert.Nil(t, err, "GetDir")

	assert.Equal(t, 2, len(dirents))

	ts.Shutdown()
}

func TestUnionSimple(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	// Make a second one
	ts.BootFss3d()

	dirents, err := ts.GetDir(path.Join(sp.S3, "~local/9ps3/"))
	assert.Nil(t, err, "GetDir: %v", err)

	assert.True(t, fslib.Present(dirents, ROOT), "%v not in %v", ROOT, dirents)

	ts.Shutdown()
}

func TestUnionDir(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	// Make a second one
	ts.BootFss3d()

	dirents, err := ts.GetDir(path.Join(sp.S3, "~local/9ps3/gutenberg"))
	assert.Nil(t, err, "GetDir")

	assert.Equal(t, 8, len(dirents))

	ts.Shutdown()
}

func TestUnionFile(t *testing.T) {
	ts, err1 := test.NewTstateAll(t)
	if !assert.Nil(t, err1, "Error New Tstate: %v", err1) {
		return
	}

	// Make a second one
	ts.BootFss3d()

	file, err := os.ReadFile("../input/pg-being_ernest.txt")
	assert.Nil(t, err, "ReadFile")

	name := path.Join(sp.S3, "~local/9ps3/gutenberg/pg-being_ernest.txt")
	st, err := ts.Stat(name)
	assert.Nil(t, err, "Stat")

	fd, err := ts.Open(name, sp.OREAD)
	if assert.Nil(ts.T, err, "Error Open: %v", err) {
		n := len(file)
		for {
			b := make([]byte, 8192)
			n, err := ts.Read(fd, b)
			if n == 0 {
				break
			}
			if !assert.Nil(ts.T, err, "Error Read: %v", err) {
				break
			}
			b = b[:n]
			for i := 0; i < int(n); i++ {
				assert.Equal(t, file[i], b[i])
			}
			file = file[len(b):]
		}
		assert.Equal(ts.T, int(st.Length), n)
	}

	ts.Shutdown()
}
