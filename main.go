package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var seconds int //global value to store starttime of the server in UNIX time

func countryHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json") //makes content-type application/json
	s := r.URL.Path
	keys, ok := r.URL.Query()["limit"] //checks the query key limit
	if !ok || len(keys[0]) < 1 {       //if error or incorrect amount of queries
		//print error code 400, bad request
		fmt.Println("URL Parameter limit missing")
		fmt.Fprintf(w, "400"+"Bad Request")
		return
	}
	key := keys[0]     //gets limit query key
	getData(s, w, key) //parameters = url path, ResponseWriter and limit query key

}

func speciesHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	s := r.URL.Path
	speciesInfo(s, w) //parameters = url path, ResponseWriter
}

func diagHandler(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	diagnosTics(w, seconds) //parameters = ResponseWriter and UNIX time of server start
}

func main() {
	serverstart := int(time.Now().Unix()) //logging unix time of server start as
	//a global variable, used as parameter above
	seconds = serverstart
	http.HandleFunc("/conservation/v1/country/", countryHandler) //3 different handlers
	http.HandleFunc("/conservation/v1/species/", speciesHandler)
	http.HandleFunc("/conservation/v1/diag/", diagHandler)
	log.Fatal(http.ListenAndServe(":5067", nil)) //server on port 5067
}
