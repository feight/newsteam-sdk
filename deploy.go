package newsteam

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type Deployment struct {
	Name           string
	ProjectId      string
	Path           string
	DockerfilePath string
	Environment    []string
}

func Deploy(s *Deployment) {

	start := time.Now()

	fmt.Printf("\n> Building %s...\n\n", color.YellowString(s.Name))
	runPrebuild(s)

	fmt.Printf("\n> Creating %s docker image...\n\n", color.YellowString(s.Name))
	runBuild(s)

	fmt.Printf("\n> Pushing %s image to google artifact registry...\n\n", color.YellowString(s.Name))
	pushImage(s)

	fmt.Printf("\n> Deploying %s to %s...\n\n", color.YellowString(s.Name), color.YellowString(s.ProjectId))
	deploy(s)

	took := time.Since(start).Round(time.Millisecond * 100).String()

	fmt.Printf("\n🎉 Successfully deployed %s to %s in %s.\n\n", color.YellowString(s.Name), color.YellowString(s.ProjectId), took)
}

func runPrebuild(s *Deployment) {

	cmd := exec.Command("go", "build", "-o", "./bin/worker.app", s.Path)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, []string{"GOOS=linux", "GOARCH=amd64"}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		panic(errors.Wrap(err, "could not complete build"))
	}

	cmd.Wait()
}

func runBuild(s *Deployment) {

	cmd := exec.Command(
		"docker",
		"build",
		"--platform", "linux/amd64",
		"-t", fmt.Sprintf("eu.gcr.io/%s/%s", s.ProjectId, s.Name),
		s.DockerfilePath,
	)

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		panic(errors.Wrap(err, "could not complete build"))
	}

	cmd.Wait()
}

func pushImage(s *Deployment) {

	cmd := exec.Command(
		"docker",
		"push",
		fmt.Sprintf("eu.gcr.io/%s/%s", s.ProjectId, s.Name),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		panic(errors.Wrap(err, "could not push image"))
	}

	cmd.Wait()
}

func deploy(s *Deployment) {

	cmd := exec.Command(
		"gcloud",
		"run",
		"deploy",
		s.Name,
		"--project", s.ProjectId,
		"--region", "europe-west1",
		"--platform", "managed",
		"--image", fmt.Sprintf("eu.gcr.io/%s/%s", s.ProjectId, s.Name),
		"--allow-unauthenticated",
		"--update-labels", "type=backend")

	cmd.Args = append(cmd.Args, []string{
		"--set-env-vars", strings.Join(s.Environment, ",")}...)

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		panic(errors.Wrap(err, "could not complete deploy"))
	}

	cmd.Wait()
}
