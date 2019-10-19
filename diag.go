package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//diagnosticsStruct
type diagStruct struct {
	Gbif          int
	Restcountries int
	Version       string
	Uptime        int
}

//three constants, test-URLS and version (to change version need to go into code, so i figured
//it would be good practise to keep it as a constant)
const exampleURLGBIF = "http://api.gbif.org/v1/occurrence/search?country=no"
const exampleURLRESTCOUNTRIES = "https://restcountries.eu/rest/v2/name/no"
const versjon = "v1"

func diagnosTics(w http.ResponseWriter, seconds int) {
	resp, err := http.Get(exampleURLGBIF) //checking response/error from first URL
	ds := &diagStruct{}
	if err != nil {
		fmt.Println("Something wrong with Get exampleURLGBIF")
		fmt.Fprintf(w, "520"+http.StatusText(520)) //520 Server error code
		return
	}
	statusCode1 := resp.StatusCode
	fmt.Println("HTTP Response Status GBIF: ", resp.StatusCode, http.StatusText(resp.StatusCode))
	defer resp.Body.Close() //closing response.Body

	resp, err = http.Get(exampleURLRESTCOUNTRIES) //changing test-URL to restcountries
	if err != nil {                               //checking if error or new URL
		fmt.Println("Something wrong with Get exampleURLRESTCOUNTRIES")
		fmt.Fprintf(w, "520"+http.StatusText(520))
		return
	}
	statusCode2 := resp.StatusCode
	fmt.Println("HTTP Response Status Rest: ", resp.StatusCode, http.StatusText(resp.StatusCode))
	defer resp.Body.Close()
	revisedSeconds := int(time.Now().Unix()) - seconds //getting uptime from unix-parameter
	ds.Gbif = statusCode1                              //setting ds (diagnostics) struct to correct values manually
	ds.Restcountries = statusCode2
	ds.Uptime = revisedSeconds
	ds.Version = versjon
	err = json.NewEncoder(w).Encode(ds) //encoding and essentially printing the struct
	if err != nil {                     //checking for issues
		fmt.Println("Error with Encoder!")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}

}
