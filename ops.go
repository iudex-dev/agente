package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func checkDockerImage(imageName string) bool {
	// Check that the docker image exists
	outBytes, err := exec.Command(
		"docker", "image", "list", "-a",
	).CombinedOutput()
	if err != nil {
		return false
	}
	lines := strings.Split(string(outBytes), "\n")
	for _, line := range lines {
		if pos := strings.Index(line, " "); pos != -1 {
			if strings.Compare(line[:pos], imageName) == 0 {
				return true
			}
		}
	}
	return false
}

func Compile(sourceName, binaryName string) (output string, err error) {
	if ok := checkDockerImage("iudex-compiler-clang"); !ok {
		return "", fmt.Errorf("docker image not found for compiler")
	}

	outBytes, err := exec.Command(
		"docker", "run", "--rm", "-v", ".:/iudex",
		"iudex-compiler-clang",
		"clang++", "-static",	"-o",
		fmt.Sprintf("/iudex/%s", binaryName),
		fmt.Sprintf("/iudex/%s", sourceName),
	).CombinedOutput()

	output = string(outBytes)
	return
}

func ExecSandboxed(binaryName string) (output string, err error) {
	outBytes, err := exec.Command(
		"docker", "run",
		"--rm",
		"-v", ".:/iudex",
		"iudex-carcer",
		"carcer", "-o", "output", "-e", "error", fmt.Sprintf("/iudex/%s", binaryName),
	).CombinedOutput()

	output = string(outBytes)
	return
}
