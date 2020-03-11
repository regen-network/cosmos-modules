package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteNaturalKey(t *testing.T) {
	v := Vote{
		Proposal: 1,
		Voter:    []byte{0xff, 0xfe},
	}
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1, 0xff, 0xfe}, v.NaturalKey())
}
