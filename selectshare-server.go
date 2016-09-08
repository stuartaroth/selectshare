package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type PostRequest struct {
	Action string `json:"action"`
	Item   string `json:"item"`
}

var (
	files = make(map[string]string)
)

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Error:", err.Error())
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[1:]
	path, pathExists := files[filename]
	if !pathExists {
		errorHandler(w, r, errors.New("File requested does not exist"))
		return
	}

	http.ServeFile(w, r, path)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(files)
	if err != nil {
		errorHandler(w, r, err)
		return
	}

	w.Write(js)
}

func extractFilename(path string) string {
	slashIndex := strings.LastIndex(path, "/")
	if slashIndex == -1 {
		return path
	}

	return path[slashIndex+1:]
}

func addFile(w http.ResponseWriter, r *http.Request, item string) {
	filename := extractFilename(item)
	files[filename] = item

	_, filenameIsPresent := files[filename]
	if !filenameIsPresent {
		fmt.Fprintf(w, `{"error":"`+item+` was not added"}`)
		return
	}

	fmt.Fprintf(w, `{"success":"`+item+` was added"}`)
}

func removeFile(w http.ResponseWriter, r *http.Request, item string) {
	filename := extractFilename(item)
	delete(files, filename)

	_, filenameIsPresent := files[filename]
	if filenameIsPresent {
		fmt.Fprintf(w, `{"error":"`+item+` was not removed"}`)
		return
	}

	fmt.Fprintf(w, `{"success":"`+item+` was removed"}`)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bits, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, r, err)
		return
	}

	var postRequest PostRequest
	err = json.Unmarshal(bits, &postRequest)
	if err != nil {
		errorHandler(w, r, err)
		return
	}

	if postRequest.Action == "list" {
		listFiles(w, r)
	} else if postRequest.Action == "add" {
		addFile(w, r, postRequest.Item)
	} else if postRequest.Action == "remove" {
		removeFile(w, r, postRequest.Item)
	} else {
		errorHandler(w, r, errors.New("PostRequest action was not list, add, or remove"))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getHandler(w, r)
	} else if r.Method == "POST" {
		postHandler(w, r)
	} else {
		errorHandler(w, r, errors.New("Method was not GET or POST"))
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8888", nil)
}
