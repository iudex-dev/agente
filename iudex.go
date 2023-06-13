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

func main() {
	out, err := Compile("program.cc", "program")
	if err != nil {
		exitWithError(out)	
	}

	out, err = ExecSandboxed("program")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %s\n", err)
	}

	reader := strings.NewReader(out)
	dec := json.NewDecoder(reader)
	result := map[string]interface{}{}
	dec.Decode(&result)
	fmt.Printf("%+v", result)
}
