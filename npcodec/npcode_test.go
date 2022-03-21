package npcodec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	np "ulambda/ninep"
)

func TestPutfile(t *testing.T) {
	b := []byte("hello")
	fence := np.Tfence1{np.Tfenceid1{36, 2}, 7}
	msg := np.Tputfile{1, np.OWRITE, 0777, 101, []string{"f"}, b}
	fcall := np.MakeFcall(msg, 13, nil, fence)
	frame, error := marshal(fcall)
	assert.Nil(t, error)
	fcall1 := &np.Fcall{}
	error = unmarshal(frame, fcall1)
	assert.Nil(t, error)
	assert.Equal(t, fcall1, fcall, "fcall")
}
