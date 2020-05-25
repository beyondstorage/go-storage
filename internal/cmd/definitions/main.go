package main

func main() {
	data := parse()
	data.Sort()

	generate(data)
	format(data)
}
