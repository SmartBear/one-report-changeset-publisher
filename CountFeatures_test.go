package publisher

import (
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountLoc(t *testing.T) {
	revision := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)

	loc, files, err := CountFeatures(repo, revision, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1297, loc)
	assert.Equal(t, 16, files)
}
