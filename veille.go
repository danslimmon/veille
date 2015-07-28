package main

func main() {
	var err error

	err = LoadConfig()
	if err != nil {
		panic(err)
	}
}
