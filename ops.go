package iudex

import (
	"fmt"
	"os/exec"
)

func Compile(sourcePath, binaryName string) (output string, err error) {
	print("Compiling... ")
	defer println("ok")

	outBytes, err := exec.Command(
		"docker", "run",
		"--rm",
		"-v", ".:/iudex",
		"iudex-compile-clang",
		"clang++", sourcePath, "-o", 
		fmt.Sprintf("/iudex/%s", binaryName), 
		"-static",
	).CombinedOutput()

	output = string(outBytes)
	return
}

func ExecSandboxed(binaryName string) (output string, err error) {
	print("Running... ")
	defer println("ok")

	outBytes, err := exec.Command(
		"docker", "run",
		"--rm",
		"-v", ".:/iudex",
		"iudex-carcer",
		"carcer", fmt.Sprintf("/iudex/%s", binaryName),
	).CombinedOutput()

	output = string(outBytes)
	return
}
