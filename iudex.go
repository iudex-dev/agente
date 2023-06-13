package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, `{
  "status": "Compilation error",
  "output": "%s"
}`, message)
	os.Exit(1)
}

type CarcerResult struct {
	CpuTime  int64  `json:"cpu_time"`
	RealTime int64  `json:"real_time"`
	Memory   int64  `json:"memory"`
	ExitCode int    `json:"exit_code"`
	Report   string `json:"report"`
}

func processSubmission(filename string) CarcerResult {
	out, err := Compile(filename, "iudex-binary")
	if err != nil {
		exitWithError(out)
	}

	out, err = ExecSandboxed("iudex-binary")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %s\n", err)
	}

	reader := strings.NewReader(out)
	dec := json.NewDecoder(reader)
	result := CarcerResult{}
	dec.Decode(&result)
	return result
}

func main() {
	result := processSubmission("program.cc")
	fmt.Printf("%+v", result)
}
