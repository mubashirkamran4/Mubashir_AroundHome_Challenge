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

//// Social struct which contains a
//// list of links
//type experienced_in struct {
//	Wood int `json:"wood"`
//	Tiles int `json:"tiles"`
//	Carpet int `json:"carpet"`
//}

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
	// Open our jsonFile
	jsonFile, err := os.Open("./partners.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened partners.json")
	// defer the closing of our jsonFile so that we can parse it later on

	// read our opened jsonFile as a byte array.
	byteValue,_ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var partners Partners

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &partners)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get_partner/{name}", getOnePartner).Methods("GET")
	router.HandleFunc("/match_partners", returnPartners).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	//for i := 0; i < len(partners.Partners); i++ {
	//	fmt.Println("Partner # %s Details: ", i)
	//	fmt.Println("Partner Name: " + partners.Partners[i].Name)
	//	fmt.Println("Partner latitude: %f" , partners.Partners[i].Latitude)
	//	fmt.Println("Partner longitude: %f" , partners.Partners[i].Longitude)
	//	fmt.Println("Partner operating_radius_latitude: %f" , partners.Partners[i].Operating_radius_latitude)
	//	fmt.Println("Partner operating_radius_longitude: %f" , partners.Partners[i].Operating_radius_longitude)
	//	fmt.Println("Partner Experience in wood: " + strconv.Itoa(partners.Partners[i].Experienced_in.Wood))
	//	fmt.Println("Partner Experience in tiles: " + strconv.Itoa(partners.Partners[i].Experienced_in.Tiles))
	//	fmt.Println("Partner Experince in carpet: " + strconv.Itoa(partners.Partners[i].Experienced_in.Carpet))
	//	fmt.Println("Partner rating: " + strconv.Itoa(partners.Partners[i].Rating))
	//}
}