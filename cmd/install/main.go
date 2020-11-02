package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	f, err := os.Open("go.mod")
	if err != nil {
		log.Fatalf("open go.mod: %v", err)
	}
	defer f.Close()

	var version string

	buf := bufio.NewScanner(f)

	for buf.Scan() {
		content := buf.Text()

		// github.com/aos-dev/go-storage/v2 v2.0.0-alpha.1.0.20201102100531-59890d25fd75
		if !strings.Contains(content, "github.com/aos-dev/go-storage/v2") {
			continue
		}

		vs := strings.Split(content, "-")
		if len(vs) < 1 {
			log.Fatalf("invalid version: %s", content)
		}

		version = vs[len(vs)-1]
	}

	// TODO: we need to handle replace here.
	cmd := exec.Command("go", "get", fmt.Sprintf("github.com/aos-dev/go-storage/cmd/definitions@%s", version))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
