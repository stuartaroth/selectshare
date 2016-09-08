package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	helpMessage  = "\nselectshare\n\nYou must provide either list, add, or remove. \n\nAdd should be followed by the path of the file. \n\nRemove should be followed by the path or name of the file."
	errorMessage = "Oops, an error occurred."
	httpClient   = &http.Client{}
)

type PostRequest struct {
	Action string `json:"action"`
	Item   string `json:"item"`
}

func makeRequest(action string, item string) ([]byte, error) {
	postRequest := PostRequest{action, item}

	js, err := json.Marshal(postRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://localhost:8888", bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bits, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bits, nil
}

func listRequest() {
	bits, err := makeRequest("list", "")
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	var files map[string]string
	err = json.Unmarshal(bits, &files)
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	fmt.Println("List count:", len(files))
	for key, value := range files {
		fmt.Println(key+":", value)
	}
}

func addRequest(item string) {
	bits, err := makeRequest("add", item)
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	var files map[string]string
	err = json.Unmarshal(bits, &files)
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	for key, value := range files {
		fmt.Println(key+":", value)
	}
}

func removeRequest(item string) {
	bits, err := makeRequest("remove", item)
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	var files map[string]string
	err = json.Unmarshal(bits, &files)
	if err != nil {
		fmt.Println(errorMessage)
		return
	}

	for key, value := range files {
		fmt.Println(key+":", value)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		return
	}

	action := os.Args[1]

	if action == "list" {
		listRequest()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println(helpMessage)
		return
	}

	item := os.Args[2]

	if action == "add" {
		addRequest(item)
	} else if action == "remove" {
		removeRequest(item)
	} else {
		fmt.Println(helpMessage)
	}
}
