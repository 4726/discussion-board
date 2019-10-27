package main

func main() {
	cfg, err := ConfigFromJSON("config.json")
	if err != nil {
		panic(err)
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		panic(err)
	}

	panic(api.Run(":14000"))
}
