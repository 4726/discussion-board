package main

func main() {
	config, err := ConfigFromJSON("config.json")
	if err != nil {
		panic(err)
	}

	api, err := NewRestAPI(config)
	if err != nil {
		panic(err)
	}

	panic(api.Run(":14000"))
}
