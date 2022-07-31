package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"strings"
)

type Partners struct {
	Partners []Partner `json:"partners"`
}

type Partner struct {
	Name   string `json:"name"`
	Experienced_in string `json:"experienced_in"`
	Latitude   float32 `json:"latitude"`
	Longitude    float32    `json:"longitude"`
	Operating_radius_latitude float32 `json:"operating_radius_latitude"`
	Operating_radius_longitude float32 `json:"operating_radius_longitude"`
	Rating int `json:"rating"`
}

type crequest struct {
	Material string `json:"material"`
	Clatitude   float32 `json:"clatitude"`
	Clongitude   float32 `json:"clongitude"`
	SquareMetres  string `json:"squarmetres"`
	PhoneNumber string `json:"phonenumber"`
}

// Abs returns the absolute value of x.
func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

//function to return specific partner based on Partner's Name
func getOnePartner(w http.ResponseWriter, r *http.Request) {
	partnerName := mux.Vars(r)["name"]
	// Open our jsonFile
	jsonFile, err := os.Open("./partners.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened partners.json")
	byteValue,_ := ioutil.ReadAll(jsonFile)

	var partners Partners

	json.Unmarshal(byteValue, &partners)

	for i := 0; i < len(partners.Partners); i++ {
		if partnerName == partners.Partners[i].Name{
			json.NewEncoder(w).Encode(partners.Partners[i])
		}
	}
}

// function to handle the matching partners request based on customer's request
func returnPartners(w http.ResponseWriter, r *http.Request) {

	// Open our jsonFile
	jsonFile, err := os.Open("./partners.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened partners.json")
	byteValue,_ := ioutil.ReadAll(jsonFile)

	var partners Partners
	json.Unmarshal(byteValue, &partners)

	var matched_partners []Partner

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter correct data to get matched with partners")
	}
	var customerRequest crequest

	json.Unmarshal(reqBody, &customerRequest)
	fmt.Println("customer latitude")
	fmt.Println(customerRequest)

	for i := 0; i < len(partners.Partners); i++ {
		experiences := strings.Split(partners.Partners[i].Experienced_in, ",")
		for _, exp := range experiences {
			if exp == customerRequest.Material {
					var customerLat = customerRequest.Clatitude
					var customerLon = customerRequest.Clongitude
					var partnerLat = partners.Partners[i].Latitude
					var partnerLon = partners.Partners[i].Longitude
					if (Abs(customerLat - partnerLat) <= partners.Partners[i].Operating_radius_latitude) && (Abs(customerLon - partnerLon) <= partners.Partners[i].Operating_radius_longitude){
						matched_partners = append(matched_partners, partners.Partners[i])
					}
			}
		}
	}

	//sort partners based on rating
	for i := 0; i < len(matched_partners); i++ {
		for i := 0; i < len(matched_partners); i++ {
			if i == len(matched_partners) - 1{
				continue
			}
			if matched_partners[i].Rating < matched_partners[i+1].Rating {
				temp_partner := matched_partners[i]
				matched_partners[i] = matched_partners[i+1]
				matched_partners[i+1] = temp_partner
			}
		}
	}

	//sort partners based on radius
	for i := 0; i < len(matched_partners); i++ {
		for i := 0; i < len(matched_partners); i++ {
			if i == len(matched_partners) - 1{
				continue
			}
			var customerLat = customerRequest.Clatitude
			var customerLon = customerRequest.Clongitude
			if (Abs(customerLat - matched_partners[i].Latitude) > Abs(customerLat - matched_partners[i+1].Latitude)) &&
				(Abs(customerLon - matched_partners[i].Longitude) > Abs(customerLon - matched_partners[i+1].Longitude)){
				temp_partner := matched_partners[i]
				matched_partners[i] = matched_partners[i+1]
				matched_partners[i+1] = temp_partner
			}
		}
	}

	json.NewEncoder(w).Encode(matched_partners)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get_partner/{name}", getOnePartner).Methods("GET")
	router.HandleFunc("/match_partners", returnPartners).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}