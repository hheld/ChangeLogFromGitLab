package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	GitLabURL   string
	GitLabToken string
	ProjectName string
	FromSha     string
	ToSha       string
}

const configFileLocation = "config"

var (
	configInfo config
	apiConn    *gitLabAPIConnection
)

func init() {
	if _, err := os.Stat(configFileLocation); os.IsNotExist(err) {
		f, _ := os.Create(configFileLocation)

		json.NewEncoder(f).Encode(configInfo)

		log.Println("There is no config file yet, fill the generated empty one and restart!")

		os.Exit(1)
	}

	f, err := os.Open(configFileLocation)

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(f).Decode(&configInfo)

	if err != nil {
		panic(err)
	}

	apiConn = newGitLabAPIConnection(configInfo.GitLabURL, configInfo.GitLabToken)
}
