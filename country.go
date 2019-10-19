package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//two consts, animal(species)URL and countryURL
const animalPage = "http://api.gbif.org/v1/occurrence/search?country="
const countryPage = "https://restcountries.eu/rest/v2/alpha/"

//struct housing all info about a country
type countryInfo struct {
	Name       string
	Alpha2code string
	Flag       string
	Offset     int
	Results    []struct { //2nd struct inside main struct to get species info
		Species    string `json:"species"`
		SpeciesKey int    `json:"speciesKey"`
	} `json:"Results"`
	//ArrayStringTest  []string `json:"species"`
	//ArrayStringTest2 []string `json:"speciesKey"`
}

/*type spfo struct {
	Results []struct {
		Species    string `json:"species"`
		SpeciesKey int    `json:"speciesKey"`
	}
}*/

func getData(urlServer string, w http.ResponseWriter, limit string) {
	parts := strings.Split(urlServer, "/")
	iso := parts[4]
	url := countryPage + iso   //gets correct URL by splicing the PATH
	resp, err := http.Get(url) //checks for response or errors in with http.Get
	if err != nil {            //if error, print 520 code
		fmt.Println("Error with http.Get!")
		fmt.Fprintf(w, "520"+http.StatusText(520))
		return
	}
	defer resp.Body.Close() //close connection
	//spf := &spfo{}
	ci := &countryInfo{}                        //pointer to a struct
	err = json.NewDecoder(resp.Body).Decode(ci) //decoding into pointer
	if err != nil {                             //checking for errors
		fmt.Println("Error with Decoder!")         // if error, print 500 because
		fmt.Fprintf(w, "500"+http.StatusText(500)) //it is an issue with their API
		return
	}
	/*err = json.NewDecoder(resp.Body).Decode(spf)
	if err != nil {
		fmt.Println("Error with Decoder!")
		fmt.Fprintf(w, "500"+http.StatusText(500))
	}*/
	number, err := strconv.Atoi(limit) //converting "number" from limit to an int to get the queried limit
	if err != nil {                    //checking for error, code 500
		fmt.Println("Error with strconv")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}
	if number > 300 { //if queried number is over 300 (max query), then loop and query multiple times
		for i := 0; i < number; i++ {
			newoffset := strconv.Itoa(i)                                    //converting "i" to be offset (querying multiple pages)
			newrl := animalPage + iso + "&offset=" + newoffset + "&limit=1" //new url query
			resp, err = http.Get(newrl)                                     //checking response or error on new url (newrl)
			if err != nil {                                                 //if error, code 520
				fmt.Println("Error with Get Species")
				fmt.Fprintf(w, "520"+http.StatusText(520))
				return
			}
			defer resp.Body.Close()                     //closing connection
			err = json.NewDecoder(resp.Body).Decode(ci) //decoding response to ci
			if err != nil {                             //checking if error, code 500
				fmt.Println("Error with Species decoder")
				fmt.Fprintf(w, "500"+http.StatusText(500))
				return
			}
		}

	} else { //if query limit isn't over 300

		newrl := animalPage + iso + "&limit=" + limit //hard coding limit into API query
		resp, err = http.Get(newrl)                   //response or error on new url
		if err != nil {
			fmt.Println("Error with Get Species")
			fmt.Fprintf(w, "520"+http.StatusText(520))
			return
		}

		defer resp.Body.Close()                     //closing connection
		err = json.NewDecoder(resp.Body).Decode(ci) //decoding hard-coded query into ci
		if err != nil {                             //checking if error
			fmt.Println("Error with Species decoder")
			fmt.Fprintf(w, "500"+http.StatusText(500))
			return
		}

		err = json.NewEncoder(w).Encode(ci)
		if err != nil {
			fmt.Println("Error with Encoder!")
			fmt.Fprintf(w, "520"+http.StatusText(520))
			return
		}
		/*err = json.NewEncoder(w).Encode(spf)
		if err != nil {
			fmt.Println("Error with Encoder!")
			fmt.Fprintf(w, "520"+http.StatusText(520))
			return
		}*/

	}
}
