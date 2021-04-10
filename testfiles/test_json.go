package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type EmailData struct {
	Domain string `json:"domain"`
	APIKey string `json:"emailAPIkey"`
	Sender string `json:"sender"`
}

func main() {
	jsonFile, err := ioutil.ReadFile("email.json")
	if err != nil {
		fmt.Println(err)
	}
	var emailData EmailData
	_ = json.Unmarshal([]byte(jsonFile), &emailData)
	fmt.Println("Domain: ", emailData.Domain)
	fmt.Println("APIKey: ", emailData.APIKey)
	fmt.Println("Sender: ", emailData.Sender)
}
