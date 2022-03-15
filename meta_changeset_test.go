package publisher

import (
	"encoding/json"
	"github.com/libgit2/git2go/v33"
	"github.com/onsi/gomega"
	"github.com/sabhiram/go-gitignore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMakeMetaChangesetNoExcludeAndIgnore(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-metaChangeset-publisher.git"
	oldShas := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	sha := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	metaChangeset, err := MakeMetaChangeset(&oldShas, &sha, true, &remote, repo, nil, nil, true)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(metaChangeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-metaChangeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["ad2c70149ccc529ab26588cde2af1312e6aa0c06"],
	  "sha": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
      "loc": 823,
      "files": 7,
	  "changes": [
		{
		  "oldPath": "",
		  "newPath": "testdata/a.txt",
		  "lineMappings": [
			[-1,0],
			[-1,1],
			[-1,2],
			[-1,3]
		  ]
		},
		{
		  "oldPath": "",
		  "newPath": "testdata/b.txt",
		  "lineMappings": [
			[-1,0],
			[-1,1],
			[-1,2],
			[-1,3]
		  ]
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithExclude(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-metaChangeset-publisher.git"
	oldShas := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	sha := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	metaChangeset, err := MakeMetaChangeset(&oldShas, &sha, true, &remote, repo, ignore.CompileIgnoreLines("testdata/a.*"), nil, true)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(metaChangeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-metaChangeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["ad2c70149ccc529ab26588cde2af1312e6aa0c06"],
	  "sha": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
      "loc": 820,
      "files": 6,
	  "changes": [
		{
		  "oldPath": "",
		  "newPath": "testdata/b.txt",
		  "lineMappings": [
			[-1,0],
			[-1,1],
			[-1,2],
			[-1,3]
		  ]
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithInclude(t *testing.T) {
	//g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-metaChangeset-publisher.git"
	oldShas := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	sha := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	metaChangeset, err := MakeMetaChangeset(&oldShas, &sha, true, &remote, repo, nil, ignore.CompileIgnoreLines("testdata/b.*"), false)
	assert.NoError(t, err)

	expected := &MetaChangeset{
		Remote:   "git@github.com:SmartBear/one-report-metaChangeset-publisher.git",
		UnixTime: 1644410531,
		OldShas:  []string{"ad2c70149ccc529ab26588cde2af1312e6aa0c06"},
		Sha:      "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
		Loc:      -1,
		Files:    1,
		Changes: []Change{
			{
				OldPath:      "",
				NewPath:      "testdata/b.txt",
				LineMappings: [][]int{},
			},
		},
	}
	assert.Equal(t, expected, metaChangeset)

	//j, err := json.MarshalIndent(metaChangeset, "", "  ")
	//assert.NoError(t, err)
	//
	//
	//
	//const expected = `{
	//  "remote": "git@github.com:SmartBear/one-report-metaChangeset-publisher.git",
    //  "unixTime": 1644410531,
	//  "oldShas": ["ad2c70149ccc529ab26588cde2af1312e6aa0c06"],
	//  "sha": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
    //  "loc": -1,
    //  "files": 1,
	//  "changes": [
	//	{
	//	  "oldPath": "",
	//	  "newPath": "testdata/b.txt",
	//	  "lineMappings": []
	//	}
	//  ]
	//}`
	//
	//g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithDeleteAndModification(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	oldShas := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	sha := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	changeset, err := MakeMetaChangeset(&oldShas, &sha, true, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["1ae2aabbcdd11948403578a4f2dd32911cc48a00"],
	  "sha": "e57bfde5c3591a14c0e199c900174a08b0b94312",
      "loc": -1,
      "files": 6,
	  "changes": [
		{
		  "oldPath": "testdata/a.txt",
		  "newPath": "",
		  "lineMappings": []
		},
		{
		  "oldPath": "testdata/b.txt",
		  "newPath": "testdata/b.txt",
		  "lineMappings": []
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithMovedFile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-metaChangeset-publisher.git"
	oldShas := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	sha := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	metaChangeset, err := MakeMetaChangeset(&oldShas, &sha, true, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(metaChangeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-metaChangeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["e57bfde5c3591a14c0e199c900174a08b0b94312"],
	  "sha": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": 6,
	  "changes": [
		{
		  "oldPath": "testdata/b.txt",
		  "newPath": "testdata/c.txt",
		  "lineMappings": []
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithHashedPaths(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	oldShas := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	sha := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	changeset, err := MakeMetaChangeset(&oldShas, &sha, false, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["e57bfde5c3591a14c0e199c900174a08b0b94312"],
	  "sha": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": 6,
	  "changes": [
		{
		  "oldPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
		  "newPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
		  "lineMappings": []
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeMetaChangesetWithoutRemote(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	oldShas := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	sha := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.OpenRepository(".")
	assert.NoError(t, err)
	metaChangeset, err := MakeMetaChangeset(&oldShas, &sha, false, nil, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(metaChangeset, "", "  ")
	assert.NoError(t, err)

	expected := `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
      "unixTime": 1644410531,
	  "oldShas": ["e57bfde5c3591a14c0e199c900174a08b0b94312"],
	  "sha": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": 6,
	  "changes": [
		{
		  "oldPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
		  "newPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
		  "lineMappings": []
		}
	  ]
	}`

	if os.Getenv("CI") != "" {
		// Different remote on GitHub Actions
		expected = `{
		  "remote": "https://github.com/SmartBear/one-report-changeset-publisher",
          "unixTime": 1644410531,
		  "oldShas": ["e57bfde5c3591a14c0e199c900174a08b0b94312"],
		  "sha": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
		  "loc": -1,
		  "files": 6,
		  "changes": [
			{
			  "oldPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
			  "newPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
			  "lineMappings": []
			}
		  ]
		}`
	}

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}
