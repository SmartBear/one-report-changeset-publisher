package publisher

import (
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/onsi/gomega"
	"github.com/sabhiram/go-gitignore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMakeChangesets(t *testing.T) {
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	r1 := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	r2 := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	r3 := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	revisions := []string{
		r1,
		r2,
		r3,
	}
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)

	changesets, err := MakeChangesets(revisions, false, &remote, repo, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(changesets))
	assert.Equal(t, r1, changesets[0].FromRev)
	assert.Equal(t, r2, changesets[0].ToRev)
	assert.Equal(t, -1, changesets[0].Loc)
	assert.Equal(t, r2, changesets[1].FromRev)
	assert.Equal(t, r3, changesets[1].ToRev)
	assert.Equal(t, 821, changesets[1].Loc)
	assert.Equal(t, 6, changesets[1].Files)
}

func TestMakeChangesetNoExcludeAndIgnore(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	toRev := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, false, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "",
		  "toPath": "testdata/a.txt",
		  "lineMappings": [
			[-1,0],
			[-1,1],
			[-1,2],
			[-1,3]
		  ]
		},
		{
		  "fromPath": "",
		  "toPath": "testdata/b.txt",
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

func TestMakeChangesetWithExclude(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	toRev := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, false, &remote, repo, ignore.CompileIgnoreLines("testdata/a.*"), nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "",
		  "toPath": "testdata/b.txt",
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

func TestMakeChangesetWithInclude(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "ad2c70149ccc529ab26588cde2af1312e6aa0c06"
	toRev := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, false, &remote, repo, nil, ignore.CompileIgnoreLines("testdata/b.*"), false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "",
		  "toPath": "testdata/b.txt",
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

func TestMakeChangesetWithDeleteAndModification(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "1ae2aabbcdd11948403578a4f2dd32911cc48a00"
	toRev := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, false, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
	  "toRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "testdata/a.txt",
		  "toPath": "",
		  "lineMappings": [
			[0,-1],
			[1,-1],
			[2,-1],
			[3,-1]
		  ]
		},
		{
		  "fromPath": "testdata/b.txt",
		  "toPath": "testdata/b.txt",
		  "lineMappings": [
			[3,4],
			[-1,3]
		  ]
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeChangesetWithMovedFile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	toRev := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, false, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "testdata/b.txt",
		  "toPath": "testdata/c.txt",
		  "lineMappings": []
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeChangesetWithHashedPaths(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	fromRev := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	toRev := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, true, &remote, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
		  "toPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
		  "lineMappings": []
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeChangesetWithoutRemote(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	fromRev := "e57bfde5c3591a14c0e199c900174a08b0b94312"
	toRev := "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3"
	repo, err := git.PlainOpen(".")
	assert.NoError(t, err)
	changeset, err := MakeChangeset(&fromRev, &toRev, true, nil, repo, nil, nil, false)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	expected := `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
      "loc": -1,
      "files": -1,
	  "changes": [
		{
		  "fromPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
		  "toPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
		  "lineMappings": []
		}
	  ]
	}`

	if os.Getenv("CI") != "" {
		// Different remote on GitHub Actions
		expected = `{
		  "remote": "https://github.com/SmartBear/one-report-changeset-publisher",
		  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
		  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
		  "loc": -1,
		  "files": -1,
		  "changes": [
			{
			  "fromPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
			  "toPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
			  "lineMappings": []
			}
		  ]
		}`
	}

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}
