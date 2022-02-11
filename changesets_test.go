package publisher

import (
	"encoding/json"
	"github.com/onsi/gomega"
	"github.com/sabhiram/go-gitignore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeChangesetNoExcludeAndIgnore(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	remote := "git@github.com:SmartBear/one-report-changeset-publisher.git"
	changeset, err := MakeChangeset("ad2c70149ccc529ab26588cde2af1312e6aa0c06", "1ae2aabbcdd11948403578a4f2dd32911cc48a00", false, &remote, nil, nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
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
	changeset, err := MakeChangeset("ad2c70149ccc529ab26588cde2af1312e6aa0c06", "1ae2aabbcdd11948403578a4f2dd32911cc48a00", false, &remote, ignore.CompileIgnoreLines("testdata/a.*"), nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
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
	changeset, err := MakeChangeset("ad2c70149ccc529ab26588cde2af1312e6aa0c06", "1ae2aabbcdd11948403578a4f2dd32911cc48a00", false, &remote, nil, ignore.CompileIgnoreLines("testdata/b.*"))
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "ad2c70149ccc529ab26588cde2af1312e6aa0c06",
	  "toRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
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
	changeset, err := MakeChangeset("1ae2aabbcdd11948403578a4f2dd32911cc48a00", "e57bfde5c3591a14c0e199c900174a08b0b94312", false, &remote, nil, nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "1ae2aabbcdd11948403578a4f2dd32911cc48a00",
	  "toRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
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
	changeset, err := MakeChangeset("e57bfde5c3591a14c0e199c900174a08b0b94312", "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3", false, &remote, nil, nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
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
	changeset, err := MakeChangeset("e57bfde5c3591a14c0e199c900174a08b0b94312", "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3", true, &remote, nil, nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
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
	changeset, err := MakeChangeset("e57bfde5c3591a14c0e199c900174a08b0b94312", "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3", true, nil, nil, nil)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "e57bfde5c3591a14c0e199c900174a08b0b94312",
	  "toRev": "082022d1a8bac6a768b0fc9243f3f37ede8c0fc3",
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
