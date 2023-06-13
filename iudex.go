package main

import (
	"fmt"
	"os"
)

func main() {
	out, err := Compile("program.cc", "program")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation error: %s", err)
	}
	fmt.Printf("OUTPUT:\n%s", out)

	out, err = ExecSandboxed("program")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %s\n", err)
	}
	fmt.Printf("OUTPUT:\n%s", out)
}
