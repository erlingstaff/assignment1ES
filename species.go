package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//URL Being modified with user input
const speciesURL = "http://api.gbif.org/v1/species/"

//Struct holding all info about a species
type infoSpecies struct {
	Key            int    `json:"key"`
	Kingdom        string `json:"kingdom"`
	Phylum         string `json:"phylum"`
	Order          string `json:"order"`
	Family         string `json:"family"`
	Genus          string `json:"genus"`
	ScientificName string `json:"scientificName"`
	CanonicalName  string `json:"canonicalName"`
	Year           string `json:"year"`
}

// main function, gets the URL string and the ResponseWriter
// as parameters
func speciesInfo(s string, w http.ResponseWriter) {
	parts := strings.Split(s, "/")
	speciesKey := parts[4]         //Splits the URL string and extrapolates the key
	url := speciesURL + speciesKey //Creates the URL
	resp, err := http.Get(url)     //checks for response / error on http.Get of the URL
	if err != nil {                //If error, print 520 (unknown error)
		fmt.Println("Error with getting URL")
		fmt.Fprintf(w, "520"+http.StatusText(520))
		return
	}
	defer resp.Body.Close() //dereferencing & closing response

	is := &infoSpecies{}                        //species pointer
	err = json.NewDecoder(resp.Body).Decode(is) //checks for decoder error
	if err != nil {
		fmt.Println("Error with Decoder!")
		fmt.Fprintf(w, "500"+http.StatusText(500)) //If error, print 500 (internal server error)
		return
	}

	url += "/name" //appen /name onto the URL for the year

	resp, err = http.Get(url) //check resp/err on new url

	if err != nil {
		fmt.Println("Error getting Year URL")
		fmt.Fprintf(w, "520"+http.StatusText(520))
		return
	}
	defer resp.Body.Close() //dereference
	err = json.NewDecoder(resp.Body).Decode(is)
	if err != nil {
		fmt.Println("Error Decoding Year URL")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}
	err = json.NewEncoder(w).Encode(is) //encodes it to the w parameter
	if err != nil {
		fmt.Println("Error Encoding")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}

}
