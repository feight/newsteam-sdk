package main

import "github.com/feight/newsteam-sdk"

func main() {
	newsteam.Deploy(&newsteam.Deployment{
		Name:           "importers",
		ProjectId:      "newsteam-stage",
		Path:           "./cmd/worker",
		DockerfilePath: "./bin",
		Environment:    []string{},
	})
}
