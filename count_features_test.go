package publisher

import (
	"github.com/libgit2/git2go/v33"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountFeaturesWithLines(t *testing.T) {
	revision := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	oid, err := git.NewOid(revision)
	assert.NoError(t, err)
	commit, err := repo.LookupCommit(oid)
	assert.NoError(t, err)
	tree, err := commit.Tree()
	assert.NoError(t, err)
	loc, files, err := CountFeatures(repo, tree, nil, nil, true)
	assert.NoError(t, err)
	assert.Equal(t, 1297, loc)
	assert.Equal(t, 16, files)
}

func TestCountFeaturesWithoutLines(t *testing.T) {
	revision := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	oid, err := git.NewOid(revision)
	assert.NoError(t, err)
	commit, err := repo.LookupCommit(oid)
	assert.NoError(t, err)
	tree, err := commit.Tree()
	assert.NoError(t, err)
	loc, files, err := CountFeatures(repo, tree, nil, nil, false)
	assert.NoError(t, err)
	assert.Equal(t, -1, loc)
	assert.Equal(t, 16, files)
}
