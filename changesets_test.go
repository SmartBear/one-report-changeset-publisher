package changesets

import (
	"encoding/json"
	"github.com/onsi/gomega"
	"github.com/sabhiram/go-gitignore"
	"testing"
)

func TestMakeChangesetNoIgnore(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	_, changeset := MakeChangeset(
		"0ac4dd0d5519bac733f9fcd13792c586317b544d",
		"8bb476618aafc35eafa6beb7f63e286efa3df5d4",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		ignore.CompileIgnoreLines(),
	)
	j, _ := json.MarshalIndent(changeset, "", "  ")

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "0ac4dd0d5519bac733f9fcd13792c586317b544d",
	  "toRev": "8bb476618aafc35eafa6beb7f63e286efa3df5d4",
	  "changes": [
		{
		  "fromPath": "go.mod",
		  "toPath": "go.mod",
		  "lineMappings": [
			[
			  4,
			  -1
			],
			[
			  5,
			  -1
			],
			[
			  6,
			  -1
			],
			[
			  7,
			  -1
			],
			[
			  8,
			  -1
			],
			[
			  9,
			  -1
			],
			[
			  10,
			  -1
			],
			[
			  11,
			  -1
			],
			[
			  12,
			  -1
			]
		  ]
		},
		{
		  "fromPath": "go.sum",
		  "toPath": "go.sum",
		  "lineMappings": [
			[
			  0,
			  -1
			],
			[
			  1,
			  -1
			],
			[
			  2,
			  -1
			],
			[
			  3,
			  -1
			],
			[
			  4,
			  -1
			],
			[
			  5,
			  -1
			],
			[
			  6,
			  -1
			],
			[
			  7,
			  -1
			],
			[
			  8,
			  -1
			],
			[
			  9,
			  -1
			],
			[
			  10,
			  -1
			],
			[
			  11,
			  -1
			],
			[
			  12,
			  -1
			],
			[
			  13,
			  -1
			],
			[
			  14,
			  -1
			],
			[
			  15,
			  -1
			],
			[
			  16,
			  -1
			],
			[
			  17,
			  -1
			],
			[
			  18,
			  -1
			],
			[
			  19,
			  -1
			],
			[
			  20,
			  -1
			],
			[
			  21,
			  -1
			],
			[
			  22,
			  0
			]
		  ]
		},
		{
		  "fromPath": "main.go",
		  "toPath": "main.go",
		  "lineMappings": [
			[
			  4,
			  -1
			],
			[
			  5,
			  -1
			],
			[
			  6,
			  5
			],
			[
			  7,
			  6
			],
			[
			  9,
			  -1
			],
			[
			  10,
			  16
			],
			[
			  11,
			  -1
			],
			[
			  12,
			  -1
			],
			[
			  13,
			  -1
			],
			[
			  14,
			  -1
			],
			[
			  15,
			  -1
			],
			[
			  16,
			  -1
			],
			[
			  17,
			  24
			],
			[
			  18,
			  25
			],
			[
			  -1,
			  4
			],
			[
			  -1,
			  7
			],
			[
			  -1,
			  9
			],
			[
			  -1,
			  10
			],
			[
			  -1,
			  11
			],
			[
			  -1,
			  12
			],
			[
			  -1,
			  13
			],
			[
			  -1,
			  14
			],
			[
			  -1,
			  15
			],
			[
			  -1,
			  17
			],
			[
			  -1,
			  18
			],
			[
			  -1,
			  19
			],
			[
			  -1,
			  20
			],
			[
			  -1,
			  21
			],
			[
			  -1,
			  22
			],
			[
			  -1,
			  23
			]
		  ]
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}

func TestMakeChangesetWithIgnore(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	_, changeset := MakeChangeset(
		"0ac4dd0d5519bac733f9fcd13792c586317b544d",
		"8bb476618aafc35eafa6beb7f63e286efa3df5d4",
		"git@github.com:SmartBear/one-report-changeset-publisher.git",
		ignore.CompileIgnoreLines("*.go", "*.sum"),
	)
	j, _ := json.MarshalIndent(changeset, "", "  ")

	const expected = `{
	  "remote": "git@github.com:SmartBear/one-report-changeset-publisher.git",
	  "fromRev": "0ac4dd0d5519bac733f9fcd13792c586317b544d",
	  "toRev": "8bb476618aafc35eafa6beb7f63e286efa3df5d4",
	  "changes": [
		{
		  "fromPath": "go.mod",
		  "toPath": "go.mod",
		  "lineMappings": [
			[
			  4,
			  -1
			],
			[
			  5,
			  -1
			],
			[
			  6,
			  -1
			],
			[
			  7,
			  -1
			],
			[
			  8,
			  -1
			],
			[
			  9,
			  -1
			],
			[
			  10,
			  -1
			],
			[
			  11,
			  -1
			],
			[
			  12,
			  -1
			]
		  ]
		}
	  ]
	}`

	g.Ω(string(j)).Should(gomega.MatchJSON(expected))
}
