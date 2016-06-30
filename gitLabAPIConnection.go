package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type gitLabAPIConnection struct {
	baseURL string
	token   string
	client  *http.Client
}

type commit struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
	Title   string `json:"title"`
	Author  string `json:"author_name"`
	Email   string `json:"author_email"`
	Message string `json:"message"`
}

type changeDescription struct {
	Type    string
	Scope   string
	Subject string
	Body    string
	Footer  string
}

const apiURL = "/api/v3"

var reCommitMsg = regexp.MustCompile(`(?m:(feat|fix|docs|style|refactor|perf|test|chore)\(([^\(\)]+)\): ([^\n]+)$\n^$\n((?:\n|.)+)^$\n((?:(?:[Rr]efs|[Cc]loses) #\d+\n)+))$`)

func newGitLabAPIConnection(gitLabBaseURL, privateToken string) *gitLabAPIConnection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &gitLabAPIConnection{
		baseURL: gitLabBaseURL,
		token:   privateToken,
		client:  &http.Client{Transport: tr},
	}
}

func (c *gitLabAPIConnection) getRequest(endPoint string) (*http.Request, error) {
	req, err := http.NewRequest("GET", c.baseURL+apiURL+endPoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	return req, nil
}

func (c *gitLabAPIConnection) projectIDFromName(projectName string) (int, error) {
	req, err := c.getRequest("/projects")

	if err != nil {
		return -1, err
	}

	res, err := c.client.Do(req)

	if err != nil {
		return -1, err
	}

	defer res.Body.Close()

	var projectInfo interface{}

	err = json.NewDecoder(res.Body).Decode(&projectInfo)

	if err != nil {
		return -1, err
	}

	projects := projectInfo.([]interface{})

	for _, project := range projects {
		m := project.(map[string]interface{})

		for k, v := range m {
			if k == "name" {
				if v.(string) == projectName {
					return int(m["id"].(float64)), nil
				}
			}
		}
	}

	return -1, errors.New("No project named '" + projectName + "' found")
}

func (c *gitLabAPIConnection) allCommits() ([]commit, error) {
	pid, err := c.projectIDFromName(configInfo.ProjectName)

	if err != nil {
		return nil, err
	}

	var commitInfo []commit

	pageNo := 0

	isInRange := false

TotalLoop:
	for {
		req, err := c.getRequest("/projects/" + strconv.Itoa(pid) + "/repository/commits?page=" + strconv.Itoa(pageNo))

		if err != nil {
			return nil, err
		}

		res, err := c.client.Do(req)

		if res.StatusCode != http.StatusOK || err != nil {
			res.Body.Close()

			return nil, err
		}

		var commits []commit

		err = json.NewDecoder(res.Body).Decode(&commits)

		if err != nil {
			res.Body.Close()

			return nil, err
		}

		if len(commits) == 0 {
			res.Body.Close()

			break
		}

		for _, commit := range commits {
			if commit.ID == configInfo.ToSha || configInfo.ToSha == "HEAD" || isInRange {
				commitInfo = append(commitInfo, commit)
				isInRange = true

				match := reCommitMsg.Match([]byte(commit.Message))

				if match {
					fmt.Printf("Match: %v, %s\n", match, commit.ID)

					matches := reCommitMsg.FindAllStringSubmatch(commit.Message, -1)

					chgDesc := changeDescription{
						Type:    matches[0][1],
						Scope:   matches[0][2],
						Subject: matches[0][3],
						Body:    matches[0][4],
						Footer:  matches[0][5],
					}

					fmt.Printf("Matches: %+v\n", chgDesc)
				}
			}

			if commit.ID == configInfo.FromSha {
				isInRange = false
				break TotalLoop
			}
		}

		pageNo++
		res.Body.Close()
	}

	return commitInfo, nil
}
