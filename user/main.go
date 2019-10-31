package main

func main() {
	cfg, err := ConfigFromJSON("config.json")
	if err != nil {
		log.WithFields(appFields).Fatal(err)
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		log.WithFields(appFields).Fatal(err)
	}

	err = api.Run(":14000")

	log.WithFields(appFields).Fatal(err)
}
