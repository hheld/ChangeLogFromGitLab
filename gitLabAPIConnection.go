package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
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
}

const apiURL = "/api/v3"

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
			if commit.ID == configInfo.ToSha || isInRange {
				commitInfo = append(commitInfo, commit)
				isInRange = true
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
