package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func check(path string) {
	fi, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalf("check %s: %v", path, err)
	}

	for _, v := range fi {
		if v.IsDir() && !strings.HasPrefix(v.Name(), ".") {
			err = os.Chdir(v.Name())
			if err != nil {
				log.Fatal(err)
			}

			check(v.Name())

			err = os.Chdir("..")
			if err != nil {
				log.Fatal(err)
			}
		}
		if !strings.HasSuffix(v.Name(), ".mod") {
			continue
		}

		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("go", "mod", "verify")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	check(".")
}
