package stats_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"ulambda/kernel"
	"ulambda/stats"
)

type Tstate struct {
	t *testing.T
	*kernel.System
}

func makeTstate(t *testing.T) *Tstate {
	ts := &Tstate{}
	ts.t = t
	ts.System = kernel.MakeSystemNamed("statstest", "..")
	return ts
}

func TestStatsd(t *testing.T) {
	ts := makeTstate(t)

	st := stats.StatInfo{}
	err := ts.ReadFileJson("name/statsd", &st)
	assert.Nil(t, err, "statsd")
	assert.Equal(t, stats.Tcounter(0), st.Nread, "Nread")
	for i := 0; i < 1000; i++ {
		_, err := ts.ReadFile("name/statsd")
		assert.Nil(t, err, "statsd")
	}
	err = ts.ReadFileJson("name/statsd", &st)
	assert.Nil(t, err, "statsd")
	assert.Equal(t, st.Nopen, stats.Tcounter(1000), "statsd")

	err = ts.ReadFileJson("name/statsd", &st)
	assert.Nil(t, err, "statsd")

	for i := 0; i < 10; i++ {
		log.Printf("util %v load %v\n", st.Util, st.Load)
		time.Sleep(1000 * time.Millisecond)
		//assert.Equal(t, st.Nopen, stats.Tcounter(1000), "statsd")
	}

	ts.Shutdown()
}
