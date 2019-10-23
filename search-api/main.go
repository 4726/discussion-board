package main

func main() {
	app, err := NewRestAPI()
	if err != nil {
		panic(err)
	}

	panic(app.Run(":14000"))
}
