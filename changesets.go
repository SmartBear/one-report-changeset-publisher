package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SmartBear/lhdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/sabhiram/go-gitignore"
	"net/http"
	"net/url"
)

type Changeset struct {
	Remote  string  `json:"remote"`
	FromRev string  `json:"fromRev"`
	ToRev   string    `json:"toRev"`
	Changes []*Change `json:"changes"`
}

type Change struct {
	FromPath string  `json:"fromPath"`
	ToPath       string  `json:"toPath"`
	LineMappings [][]int `json:"lineMappings"`
}

func MakeChangeset(fromRev string, toRev string, remote string, gitIgnore *ignore.GitIgnore) (error, *Changeset) {
	contextSize := 4

	r, err := git.PlainOpen(".")
	if err != nil {
		return err, nil
	}

	leftTree, err := getTree(r, fromRev)
	if err != nil {
		return err, nil
	}

	rightTree, err := getTree(r, toRev)
	if err != nil {
		return err, nil
	}

	gitChanges, err := leftTree.Diff(rightTree)
	if err != nil {
		return err, nil
	}

	changes := make([]*Change, 0)

	for _, gitChange := range gitChanges {
		action, err := gitChange.Action()
		if err != nil {
			return err, nil
		}
		switch action {
		case merkletrie.Insert:
			// TODO
		case merkletrie.Delete:
			// TODO
		case merkletrie.Modify:
			if gitIgnore.MatchesPath(gitChange.To.Name) {
				continue
			}
			leftFile, err := leftTree.File(gitChange.From.Name)
			if err != nil {
				return err, nil
			}
			leftBinary, err := leftFile.IsBinary()
			if err != nil {
				return err, nil
			}

			rightFile, err := rightTree.File(gitChange.To.Name)
			if err != nil {
				return err, nil
			}
			rightBinary, err := rightFile.IsBinary()
			if err != nil {
				return err, nil
			}

			if !leftBinary && !rightBinary {
				leftContents, err := leftFile.Contents()
				if err != nil {
					return err, nil
				}

				rightContents, err := rightFile.Contents()
				if err != nil {
					return err, nil
				}

				mapping, err := lhdiff.Lhdiff(leftContents, rightContents, contextSize, false)
				if err != nil {
					return err, nil
				}

				change := &Change{
					FromPath:     gitChange.From.Name,
					ToPath:       gitChange.To.Name,
					LineMappings: mapping,
				}
				changes = append(changes, change)
			}
		default:
			panic(fmt.Sprintf("unsupported action: %d", action))
		}
	}

	changeset := &Changeset{
		Remote:  remote,
		FromRev: fromRev,
		ToRev:   toRev,
		Changes: changes,
	}
	return err, changeset
}

func getTree(r *git.Repository, revision string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(*h)
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func Publish(changeset *Changeset, organizationId string, password string, url string) error {
	req, err := MakeRequest(changeset, organizationId, password, url)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("expected 201, got %d", res.StatusCode)
	}
	return nil
}

func MakeRequest(changeset *Changeset, organizationId string, password string, baseUrl string) (*http.Request, error) {
	body, err := json.MarshalIndent(changeset, "", "  ")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	u.Path = "/api/organization/" + url.PathEscape(organizationId) + "/changeset"
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
