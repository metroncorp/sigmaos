package benchmarks_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ulambda/benchmarks"
	db "ulambda/debug"
	"ulambda/linuxsched"
	"ulambda/mr"
	"ulambda/proc"
	"ulambda/test"
	"ulambda/wc"
)

// Parameters
var MR_APP string
var KV_AUTO string
var NKVD int
var NCLERK int
var CLERK_DURATION string
var CLERK_NCORE proc.Tcore
var KVD_NCORE proc.Tcore
var REALM2 string
var REDIS_ADDR string

// Read & set the proc version.
func init() {
	var nc int
	flag.StringVar(&MR_APP, "mrapp", "mr-wc-wiki1.8G.yml", "Name of mr yaml file.")
	flag.StringVar(&KV_AUTO, "kvauto", "manual", "KV auto-growing/shrinking.")
	flag.IntVar(&NKVD, "nkvd", 1, "Number of kvds.")
	flag.IntVar(&NCLERK, "nclerk", 1, "Number of clerks.")
	flag.StringVar(&CLERK_DURATION, "clerk_dur", "90s", "Clerk duration.")
	flag.IntVar(&nc, "clerk_ncore", 1, "Clerk Ncore")
	CLERK_NCORE = proc.Tcore(nc)
	flag.IntVar(&nc, "kvd_ncore", 2, "KVD Ncore")
	KVD_NCORE = proc.Tcore(nc)
	flag.StringVar(&REALM2, "realm2", "test-realm", "Second realm")
	flag.StringVar(&REDIS_ADDR, "redisaddr", "", "Redis server address")
}

// ========== Common parameters ==========
const (
	OUT_DIR = "name/out_dir"
)

// ========== Nice parameters ==========
const (
	MAT_SIZE        = 2000
	N_TRIALS_NICE   = 10
	CONTENDERS_FRAC = 1.0
)

var MATMUL_NPROCS = linuxsched.NCores
var CONTENDERS_NPROCS = 1

// ========== Micro parameters ==========
const (
	N_TRIALS_MICRO = 1000
	SLEEP_MICRO    = "5000us"
)

var TOTAL_N_CORES_SIGMA_REALM = 0

func TestJsonEncodeTpt(t *testing.T) {
	nruns := 50
	N_KV := 1000000
	kvs := make([]*mr.KeyValue, 0, N_KV)
	for i := 0; i < N_KV; i++ {
		kvs = append(kvs, &mr.KeyValue{"", ""})
	}
	start := time.Now()
	n := 0
	for i := 0; i < nruns; i++ {
		for _, kv := range kvs {
			b, err := json.Marshal(kv)
			assert.Nil(t, err, "Marshal")
			n += len(b)
		}
	}
	mb := 1024.0 * 1024.0
	db.DPrintf(db.ALWAYS, "Marshaling throughput: %v MB/s", float64(n)/mb/time.Since(start).Seconds())
}

func TestWCMapfTpt(t *testing.T) {
	N_WORDS := 1024 * 1024 * 100
	WORD_LEN := 3
	b := make([]byte, 0, N_WORDS*(WORD_LEN+1))
	for i := 0; i < N_WORDS; i++ {
		for j := 0; j < WORD_LEN; j++ {
			b = append(b, 'A')
		}
		b = append(b, ' ')
	}
	s := string(b)
	db.DPrintf(db.ALWAYS, "Input length: %vMB", len(s)/(1024*1024))
	n := 0
	start := time.Now()
	wc.Map("", strings.NewReader(s), func(kv *mr.KeyValue) error {
		b, err := json.Marshal(kv)
		assert.Nil(t, err, "Marshal")
		n += len(b)
		return nil
	})
	n += len(s)
	mb := 1024.0 * 1024.0
	db.DPrintf(db.ALWAYS, "WC Mapping throughput: %v MB/s", float64(n)/mb/time.Since(start).Seconds())
}

// Length of time required to do a simple matrix multiplication.
func TestNiceMatMulBaseline(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_NICE)
	_, ps := makeNProcs(N_TRIALS_NICE, "user/matmul", []string{fmt.Sprintf("%v", MAT_SIZE)}, []string{fmt.Sprintf("GOMAXPROCS=%v", MATMUL_NPROCS)}, 1)
	runOps(ts, ps, runProc, rs)
	printResults(rs)
	ts.Shutdown()
}

// Start a bunch of spinning procs to contend with one matmul task, and then
// see how long the matmul task took.
func TestNiceMatMulWithSpinners(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_NICE)
	makeOutDir(ts)
	nContenders := int(float64(linuxsched.NCores) / CONTENDERS_FRAC)
	// Make some spinning procs to take up nContenders cores.
	psSpin, _ := makeNProcs(nContenders, "user/spinner", []string{OUT_DIR}, []string{fmt.Sprintf("GOMAXPROCS=%v", CONTENDERS_NPROCS)}, 0)
	// Burst-spawn BE procs
	spawnBurstProcs(ts, psSpin)
	// Wait for the procs to start
	waitStartProcs(ts, psSpin)
	// Make the LC proc.
	_, ps := makeNProcs(N_TRIALS_NICE, "user/matmul", []string{fmt.Sprintf("%v", MAT_SIZE)}, []string{fmt.Sprintf("GOMAXPROCS=%v", MATMUL_NPROCS)}, 1)
	// Spawn the LC procs
	runOps(ts, ps, runProc, rs)
	printResults(rs)
	evictProcs(ts, psSpin)
	rmOutDir(ts)
	ts.Shutdown()
}

// Invert the nice relationship. Make spinners high-priority, and make matul
// low priority. This is intended to verify that changing priorities does
// actually affect application throughput for procs which have their priority
// lowered, and by how much.
func TestNiceMatMulWithSpinnersLCNiced(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_NICE)
	makeOutDir(ts)
	nContenders := int(float64(linuxsched.NCores) / CONTENDERS_FRAC)
	// Make some spinning procs to take up nContenders cores. (AS LC)
	psSpin, _ := makeNProcs(nContenders, "user/spinner", []string{OUT_DIR}, []string{fmt.Sprintf("GOMAXPROCS=%v", CONTENDERS_NPROCS)}, 1)
	// Burst-spawn spinning procs
	spawnBurstProcs(ts, psSpin)
	// Wait for the procs to start
	waitStartProcs(ts, psSpin)
	// Make the matmul procs.
	_, ps := makeNProcs(N_TRIALS_NICE, "user/matmul", []string{fmt.Sprintf("%v", MAT_SIZE)}, []string{fmt.Sprintf("GOMAXPROCS=%v", MATMUL_NPROCS)}, 0)
	// Spawn the matmul procs
	runOps(ts, ps, runProc, rs)
	printResults(rs)
	evictProcs(ts, psSpin)
	rmOutDir(ts)
	ts.Shutdown()
}

// Test how long it takes to init a semaphore.
func TestMicroInitSemaphore(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_MICRO)
	makeOutDir(ts)
	_, is := makeNSemaphores(ts, N_TRIALS_MICRO)
	runOps(ts, is, initSemaphore, rs)
	printResults(rs)
	rmOutDir(ts)
	ts.Shutdown()
}

// Test how long it takes to up a semaphore.
func TestMicroUpSemaphore(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_MICRO)
	makeOutDir(ts)
	_, is := makeNSemaphores(ts, N_TRIALS_MICRO)
	// Init semaphores first.
	for _, i := range is {
		initSemaphore(ts, time.Now(), i)
	}
	runOps(ts, is, upSemaphore, rs)
	printResults(rs)
	rmOutDir(ts)
	ts.Shutdown()
}

// Test how long it takes to down a semaphore.
func TestMicroDownSemaphore(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_MICRO)
	makeOutDir(ts)
	_, is := makeNSemaphores(ts, N_TRIALS_MICRO)
	// Init semaphores first.
	for _, i := range is {
		initSemaphore(ts, time.Now(), i)
		upSemaphore(ts, time.Now(), i)
	}
	runOps(ts, is, downSemaphore, rs)
	printResults(rs)
	rmOutDir(ts)
	ts.Shutdown()
}

// Test how long it takes to Spawn, run, and WaitExit a 5ms proc.
func TestMicroSpawnWaitExit5msSleeper(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(N_TRIALS_MICRO)
	makeOutDir(ts)
	_, ps := makeNProcs(N_TRIALS_MICRO, "user/sleeper", []string{SLEEP_MICRO, OUT_DIR}, []string{}, 1)
	runOps(ts, ps, runProc, rs)
	printResults(rs)
	rmOutDir(ts)
	ts.Shutdown()
}

func TestAppMR(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(1)
	jobs, apps := makeNMRJobs(ts, 1, MR_APP)
	// XXX Clean this up/hide this somehow.
	go func() {
		for _, j := range jobs {
			// Wait until ready
			<-j.ready
			// Ack to allow the job to proceed.
			j.ready <- true
		}
	}()
	p := monitorCoresAssigned(ts)
	defer p.Done()
	runOps(ts, apps, runMR, rs)
	printResults(rs)
	ts.Shutdown()
}

func runKVTest(t *testing.T, nReplicas int) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(1)
	setNCoresSigmaRealm(ts)
	nclerks := []int{NCLERK}
	db.DPrintf(db.ALWAYS, "Running with %v clerks", NCLERK)
	jobs, ji := makeNKVJobs(ts, 1, NKVD, nReplicas, nclerks, nil, CLERK_DURATION, KVD_NCORE, CLERK_NCORE, KV_AUTO, REDIS_ADDR)
	// XXX Clean this up/hide this somehow.
	go func() {
		for _, j := range jobs {
			// Wait until ready
			<-j.ready
			// Ack to allow the job to proceed.
			j.ready <- true
		}
	}()
	p := monitorCoresAssigned(ts)
	runOps(ts, ji, runKV, rs)
	defer p.Done()
	printResults(rs)
	ts.Shutdown()
}

func TestAppKVUnrepl(t *testing.T) {
	runKVTest(t, 0)
}

func TestAppKVRepl(t *testing.T) {
	runKVTest(t, 3)
}

// Burst a bunch of spinning procs, and see how long it takes for all of them
// to start.
//
// XXX Maybe we should do a version with procs that don't spin & consume so
// much CPU?
//
// XXX We should probably try this one both warm and cold.
func TestRealmBurst(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(1)
	makeOutDir(ts)
	// Find the total number of cores available for spinners across all machines.
	// We need to get this in order to find out how many spinners to start.
	setNCoresSigmaRealm(ts)
	db.DPrintf(db.ALWAYS, "Bursting %v spinning procs", TOTAL_N_CORES_SIGMA_REALM)
	ps, _ := makeNProcs(TOTAL_N_CORES_SIGMA_REALM, "user/spinner", []string{OUT_DIR}, []string{}, 1)
	p := monitorCoresAssigned(ts)
	runOps(ts, []interface{}{ps}, spawnBurstWaitStartProcs, rs)
	defer p.Done()
	printResults(rs)
	evictProcs(ts, ps)
	rmOutDir(ts)
	ts.Shutdown()
}

func TestLambdaBurst(t *testing.T) {
	ts := test.MakeTstateAll(t)
	rs := benchmarks.MakeRawResults(1)
	makeOutDir(ts)
	// Find the total number of cores available for spinners across all machines.
	// We need to get this in order to find out how many spinners to start.
	setNCoresSigmaRealm(ts)
	db.DPrintf(db.ALWAYS, "Invoking %v lambdas", 1) //TOTAL_N_CORES_SIGMA_REALM)
	ss, is := makeNSemaphores(ts, 1)                //TOTAL_N_CORES_SIGMA_REALM)
	// Init semaphores first.
	for _, i := range is {
		initSemaphore(ts, time.Now(), i)
	}
	runOps(ts, []interface{}{ss}, invokeWaitStartLambdas, rs)
	printResults(rs)
	rmOutDir(ts)
	ts.Shutdown()
}

// Start a realm with a long-running BE mr job. Then, start a realm with a kv
// job. In phases, ramp the kv job's CPU utilization up and down, and watch the
// realm-level software balance resource requests across realms.
func TestRealmBalance(t *testing.T) {
	done := make(chan bool)
	// Find the total number of cores available for spinners across all machines.
	ts := test.MakeTstateAll(t)
	setNCoresSigmaRealm(ts)
	// Structures for mr
	ts1 := test.MakeTstateRealm(t, ts.RealmId())
	rs1 := benchmarks.MakeRawResults(1)
	// Structure for kv
	ts2 := test.MakeTstateRealm(t, REALM2)
	rs2 := benchmarks.MakeRawResults(1)
	// Prep MR job
	mrjobs, mrapps := makeNMRJobs(ts1, 1, MR_APP)
	// Prep KV job
	nclerks := []int{NCLERK}
	// TODO move phases to new clerk type.
	// phases := parseDurations(ts2, []string{"5s", "5s", "5s", "5s", "5s"})
	kvjobs, ji := makeNKVJobs(ts2, 1, NKVD, 0, nclerks, nil, CLERK_DURATION, KVD_NCORE, CLERK_NCORE, KV_AUTO, REDIS_ADDR)
	p1 := monitorCoresAssigned(ts1)
	defer p1.Done()
	p2 := monitorCoresAssigned(ts2)
	defer p2.Done()
	// Run KV job
	go func() {
		runOps(ts2, ji, runKV, rs2)
		done <- true
	}()
	// Wait for KV jobs to set up.
	<-kvjobs[0].ready
	// Run MR job
	go func() {
		runOps(ts1, mrapps, runMR, rs1)
		done <- true
	}()
	// Wait for MR jobs to set up.
	<-mrjobs[0].ready
	// Kick off MR jobs.
	mrjobs[0].ready <- true
	// Sleep for a bit
	time.Sleep(70 * time.Second)
	// Kick off KV jobs
	kvjobs[0].ready <- true
	// Wait for both jobs to finish.
	<-done
	<-done
	printResults(rs1)
	printResults(rs2)
	ts1.Shutdown()
	ts2.Shutdown()
}
