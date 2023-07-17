package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const iudexApiVar = "IUDEX_API"

var iudexApi string

func init() {
	err := godotenv.Load()
	if err != nil {
		exitWithError(fmt.Sprintf("Cannot read .env: %s", err))
	}
	value, found := os.LookupEnv(iudexApiVar)
	if !found {
		exitWithError(fmt.Sprintf("Environment variable %s not found", iudexApi))
	}
	iudexApi = value
	fmt.Printf("IUDEX_API=%s\n", iudexApi)
}

func exitWithError(message string) {
	fmt.Fprintf(os.Stderr, `{
  "status": "Compilation error",
  "output": "%s"
}`, message)
	os.Exit(1)
}

func main() {
	agent := Agent{}
	agent.Run()
}
