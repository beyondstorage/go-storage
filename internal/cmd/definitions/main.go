package main

func main() {
	data := parse()
	data.Handle()
	data.Sort()

	//
	// generate(data)
	format(data)
}
