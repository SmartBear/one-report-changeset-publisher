package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Publish(changeset *Changeset, organizationId string, baseUrl string, username string, password string) (string, error) {
	req, err := MakeRequest(changeset, organizationId, baseUrl, username, password)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode >= 400 {
		txt, err := httputil.DumpResponse(res, true)
		if err != nil {
			return "", err
		}
		cleaned := strings.ReplaceAll(string(txt), "\r\n", "\n")
		return "", fmt.Errorf("HTTP request failed:\n\n%s", cleaned)
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func MakeRequest(changeset *Changeset, organizationId string, baseUrl string, username string, password string) (*http.Request, error) {
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
	req.Header.Set("Content-Type", "application/vnd.smartbear.onereport.changeset.v1+json")
	req.SetBasicAuth(username, password)
	return req, nil
}
