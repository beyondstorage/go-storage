package main

import (
	"fmt"
)

func main() {
	data := parse()

	for _, v := range data.Services {
		fp := fmt.Sprintf("../services/%s/generated.go", v.Name)
		generateT(serviceT, fp, v)
	}

	// generate(data)
	format(data)
}
