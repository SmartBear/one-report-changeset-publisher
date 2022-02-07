package publisher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"
)

func ExamplePublish() {
	changeset := &Changeset{
		Remote:  "some-remote",
		FromRev: "aaa",
		ToRev:   "bbb",
		Changes: make([]*Change, 0),
	}
	req, err := MakeRequest(changeset, "1CCC7924-051C-496E-8467-D494C1C37B2A", "s3cr3t", "https://host.com")
	printErr(err)
	res := httptest.NewRecorder()

	ChangesetHandler := func(res http.ResponseWriter, req *http.Request) {
		txt, err := httputil.DumpRequest(req, true)
		printErr(err)
		dumpedReq := strings.ReplaceAll(string(txt), "\r\n", "\n")
		fmt.Print(dumpedReq)
	}
	ChangesetHandler(res, req)

	//Output:
	// POST /api/organization/1CCC7924-051C-496E-8467-D494C1C37B2A/changeset HTTP/1.1
	// Host: host.com
	// Content-Type: application/json
	//
	// {
	//   "remote": "some-remote",
	//   "fromRev": "aaa",
	//   "toRev": "bbb",
	//   "changes": []
	// }
}

func printErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
