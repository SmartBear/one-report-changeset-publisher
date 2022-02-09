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
	changeset, err := MakeChangeset(
		"779528fd0fd6648e85fe77b8bf7c1495082e57e8",
		"9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		false,
		nil,
		nil,
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "779528fd0fd6648e85fe77b8bf7c1495082e57e8",
	  "toRev": "9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
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
	changeset, err := MakeChangeset(
		"779528fd0fd6648e85fe77b8bf7c1495082e57e8",
		"9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		false,
		ignore.CompileIgnoreLines("testdata/a.*"),
		nil,
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "779528fd0fd6648e85fe77b8bf7c1495082e57e8",
	  "toRev": "9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
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
	changeset, err := MakeChangeset(
		"779528fd0fd6648e85fe77b8bf7c1495082e57e8",
		"9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		false,
		nil,
		ignore.CompileIgnoreLines("testdata/b.*"),
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "779528fd0fd6648e85fe77b8bf7c1495082e57e8",
	  "toRev": "9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
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
	changeset, err := MakeChangeset(
		"9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
		"6e02e95590db65c905de2f466597a07cd5fd63cd",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		false,
		nil,
		nil,
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "9e2afcc22ab5a68e6ba03ebdde46a3a8b057f16b",
	  "toRev": "6e02e95590db65c905de2f466597a07cd5fd63cd",
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
	changeset, err := MakeChangeset(
		"6e02e95590db65c905de2f466597a07cd5fd63cd",
		"2895a2ce5bb461f56251a9ea3674945f12a3d902",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		false,
		nil,
		nil,
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "6e02e95590db65c905de2f466597a07cd5fd63cd",
	  "toRev": "2895a2ce5bb461f56251a9ea3674945f12a3d902",
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
	changeset, err := MakeChangeset(
		"6e02e95590db65c905de2f466597a07cd5fd63cd",
		"2895a2ce5bb461f56251a9ea3674945f12a3d902",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		true,
		nil,
		nil,
	)
	assert.NoError(t, err)

	j, err := json.MarshalIndent(changeset, "", "  ")
	assert.NoError(t, err)

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "6e02e95590db65c905de2f466597a07cd5fd63cd",
	  "toRev": "2895a2ce5bb461f56251a9ea3674945f12a3d902",
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
