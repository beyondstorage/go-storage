package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/schollz/progressbar/v2"

	"github.com/Xuanwo/storage/coreutils"
	ps "github.com/Xuanwo/storage/types/pairs"
)

func TestProgress(t *testing.T) {
	// Init a service.
	_, store, err := coreutils.Open("fs:///?work_dir=/tmp")
	if err != nil {
		log.Fatalf("service init failed: %v", err)
	}

	bar := progressbar.New(16 * 1024 * 1024)
	defer bar.Finish()

	r, err := store.Read("test_file", ps.WithReadCallbackFunc(func(b []byte) {
		bar.Add(len(b))
	}))
	if err != nil {
		log.Fatalf("service read failed: %v", err)
	}

	_, err = ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("ioutil read failed: %v", err)
	}
}
