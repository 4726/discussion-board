package main

func main() {
	api, err := NewRestAPI("todo")
	if err != nil {
		panic(err)
	}

	panic(api.Run(":14000"))
}
