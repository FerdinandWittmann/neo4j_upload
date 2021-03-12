package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//NominatimServer entrypoint for geocoding requests
const nominatimServer = "http://192.168.1.159:7070"

//NomJSONResponse Noimnatim JsonResponse Obj
type NomJSONResponse struct {
	PlaceID     int         `json:"place_id"`
	Licence     interface{} `json:"licence"`
	BoundingBox interface{} `json:"boundingbox"`
	Lat         float32     `json:"lat,string"`
	Long        float32     `json:"lon,string"`
	DisplayName string      `json:"display_name"`
	Class       string      `json:"class"`
	Type        string      `json:"type"`
	Importance  float32     `json:"importance"`
	Address     NomAddress  `json:"address"`
}

//NomAddress additional Info (addressdetails=1)
type NomAddress struct {
	Streetnumber  int    `json:"house_number,string"`
	Streetname    string `json:"road"`
	City          string `json:"city,town,village,muinipality"`
	Town          string `json:"town"`
	Village       string `json:"village"`
	Municipality  string `json:"municipality"`
	State         string `json:"state"`
	Region        string `json:"region"`
	StateDistrict string `json:"state_district"`
	County        string `json:"county"`
	PostCode      int    `json:"postcode,string"`
	Country       string `json:"country"`
}

//FillAddress takes a Address obj and fills the missing data by requesting it from Nominatim
func FillAddress(a *Address) error {
	queryString := "/search?q="
	//Build string for nominatim req
	if a.Streetname != "" {
		street := ""
		if a.Streetnumber != 0 {
			street = a.Streetname + " " + strconv.Itoa(a.Streetnumber)
		} else {
			street = a.Streetname
		}
		queryString += street
	}
	if a.City != "" {
		queryString += " " + a.City
	}
	if a.State != "" {
		queryString += " " + a.State
	}
	if a.PostCode != 0 {
		queryString += " " + strconv.Itoa(a.PostCode)
	}
	if a.Country != "" {
		queryString += " " + a.Country
	}
	// space -> +
	queryString = strings.Replace(queryString, " ", "+", -1)
	//addressdetails adds address to response
	queryString += "&format=json&limit=1&addressdetails=1"
	fmt.Println(nominatimServer + queryString)
	//Get Request
	resp, err := http.Get(nominatimServer + queryString)
	if err != nil {
		return err
	}

	var respJSON []NomJSONResponse
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(respBody, &respJSON)
	//check if zero arr
	//TODO: ignore streetname if empty result
	if len(respJSON) == 0 {
		return errors.New("@http/nominatim: no results found")
	}
	nomAddress := respJSON[0].Address
	//find most acurate city representation
	if obj := nomAddress.Municipality; obj != "" {
		a.City = obj
	}
	if obj := nomAddress.Village; obj != "" {
		a.City = obj
	}
	if obj := nomAddress.City; obj != "" {
		a.City = obj
	}
	if obj := nomAddress.Town; obj != "" {
		a.City = obj
	}
	//find most accurate State representation
	if obj := nomAddress.Region; obj != "" {
		a.State = obj
	}
	if obj := nomAddress.StateDistrict; obj != "" {
		a.State = obj
	}
	if obj := nomAddress.State; obj != "" {
		a.State = obj
	}
	if obj := nomAddress.County; obj != "" {
		a.State = obj
	}
	//adapt streetname and streetnumber/post code and to have consistent database entries
	a.Streetname = nomAddress.Streetname
	a.Streetnumber = nomAddress.Streetnumber
	a.PostCode = nomAddress.PostCode
	a.Long = respJSON[0].Long
	a.Lat = respJSON[0].Lat
	//fmt.Printf("%+v\n", a)
	return nil
}
