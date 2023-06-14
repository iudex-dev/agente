package main

import (
	"encoding/json"
	"fmt"
	"os"
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
		"clang++", "-static", "-o",
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

type CarcerResult struct {
	CpuTime  int64  `json:"cpu_time"`
	RealTime int64  `json:"real_time"`
	Memory   int64  `json:"memory"`
	ExitCode int    `json:"exit_code"`
	Report   string `json:"report"`
}

func ProcessSubmission(code string, result *CarcerResult) (err error) {
	file, err := os.CreateTemp(".", "iudex-*.cc")
	if err != nil {
		return fmt.Errorf("cannot process submission: %s", err)
	}
	_, err = file.WriteString(code)
	if err != nil {
		return fmt.Errorf("cannot write code into '%s': %s", file.Name(), err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("cannot close file '%s': %s", file.Name(), err)
	}
	out, err := Compile(file.Name(), "iudex-binary")
	if err != nil {
		exitWithError(out)
	}
	out, err = ExecSandboxed("iudex-binary")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %s\n", err)
	}

	reader := strings.NewReader(out)
	dec := json.NewDecoder(reader)
	dec.Decode(&result)
	return nil
}
