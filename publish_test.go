package publisher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
)

func ExamplePublish() {
	changeset := &MetaChangeset{
		Remote:     "some-remote",
		ParentShas: []string{"aaa"},
		Sha:        "bbb",
		Changes:    make([]Change, 0),
		Loc:        9876,
		Files:      31,
	}
	req, err := MakeRequest(changeset, "1CCC7924-051C-496E-8467-D494C1C37B2A", "https://host.com", "anyone", "secret")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	res := httptest.NewRecorder()

	ChangesetHandler := func(res http.ResponseWriter, req *http.Request) {
		txt, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dumpedReq := strings.ReplaceAll(string(txt), "\r\n", "\n")
		fmt.Print(dumpedReq)
	}
	ChangesetHandler(res, req)

	//Output:
	// POST /api/organization/1CCC7924-051C-496E-8467-D494C1C37B2A/changeset HTTP/1.1
	// Host: host.com
	// Authorization: Basic YW55b25lOnNlY3JldA==
	// Content-Type: application/vnd.smartbear.onereport.changeset.v1+json
	//
	// {
	//   "remote": "some-remote",
	//   "parentShas": [
	//     "aaa"
	//   ],
	//   "sha": "bbb",
	//   "changes": [],
	//   "loc": 9876,
	//   "files": 31
	// }
}
