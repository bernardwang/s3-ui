package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
)

func init() {

	// Had to do this because returns svg as text/xml when running on AppEngine: http://goo.gl/hwZSp2
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := mux.NewRouter()
	sr := r.PathPrefix("/api").Subrouter()
	sr.HandleFunc("/photos", Photos)
	r.HandleFunc("/{rest:.*}", handler)
	http.Handle("/", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
	http.ServeFile(w, r, "static/"+r.URL.Path)
}

type Photo struct {
	ID      int    `json:"id"`
	Desc     string `json:"desc"`
	Author string `json:"author"`
	Thumb   string `json:"thumb"`
	Master string   `json:"master"`
}

func Photos(w http.ResponseWriter, r *http.Request) {
	photos := []Photo{}
	// you'd use a real database here
	file, err := ioutil.ReadFile("photos.json")
	if err != nil {
		log.Println("Error reading photos.json:", err)
		panic(err)
	}
	fmt.Printf("file: %s\n", string(file))
	err = json.Unmarshal(file, &photos)
	if err != nil {
		log.Println("Error unmarshalling photos.json:", err)
	}

	bs, err := json.Marshal(photos)
	if err != nil {
		ReturnError(w, err)
		return
	}
	fmt.Fprint(w, string(bs))
}

func ReturnError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, "{\"error\": \"%v\"}", err)
}

func main() {
	http.ListenAndServe("localhost:8080", nil)
}
