package main

func main() {
	data := parse()
	data.Handle()
	data.Sort()

	generateT(operationT, "../generated.go", data)
	// generate(data)
	format(data)
}
