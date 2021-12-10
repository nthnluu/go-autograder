package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func main() {
	res, err := JsonTestRunner()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	file, _ := json.MarshalIndent(res, "", " ")
	_ = ioutil.WriteFile("../results.json", file, 0644)
}
