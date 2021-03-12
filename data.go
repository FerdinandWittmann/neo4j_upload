package main

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//RoomAd represents a RoomAd in general representation
type RoomAd struct {
	ID            string        `json:"id,omitempty"`
	Address       Address       `json:"address"`
	Prize         Prize         `json:"prize"`
	Header        string        `json:"header,omitempty"`
	Description   string        `json:"description"`
	Attributes    []Attribute   `json:"attributes"`
	Availabillity Availabillity `json:"availabillity"`
	URL           string        `json:"url,omitempty"`
	Origin        string        `json:"origin,omitempty"`
	CrawledAt     float64       `json:"crawledat,omitemtpy"`
	Images        []string      `json:"images,omitempty"`
}

//Insert adds RoomAd to neo4j db instance
func (r *RoomAd) Insert() (err error) {
	session, err := CreateSession(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()
	err = r.InsertDesc(session)
	if err != nil {
		return err
	}
	return nil
}

//InsertDesc inserts the roomad node with its fields
func (r RoomAd) InsertDesc(session neo4j.Session) (err error) {
	neoReq := NewNeoRequest()
	err = neoReq.AddCreate(NeoVariable{"r", "RoomAd"}, []NeoParam{
		{"id", r.ID},
		{"description", r.Description},
		{"header", r.Header},
		{"url", r.URL},
		{"crawledAt", r.CrawledAt},
		{"origin", r.Origin},
		{"images", r.Images},
	})
	if err != nil {
		return err
	}
	result, err := neoReq.send(session)
	if err != nil {
		return err
	}
	if err = result.Err(); err != nil {
		return err
	}
	record := result.Record()
	fmt.Printf("+v\n", record)

	return nil
}

//Address represents a address
type Address struct {
	Country      string  `json:"country"`
	City         string  `json:"city,omitempty"`
	State        string  `json:"state,omitempty"`
	PostCode     int     `json:"postcode"`
	Streetname   string  `json:"streetname"`
	Streetnumber int     `json:"streetnumber,omitempty"`
	Long         float32 `json:"long,omitempty"`
	Lat          float32 `json:"lat,omitempty"`
}

//Fill Fills all fields hat are nil in Address
func (a *Address) Fill() (err error) {
	err = FillAddress(a)
	return err
}

//Attribute General Attribute object type can be [name-string, name-value, name]
type Attribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

//Availabillity represents duration the room is free To None is equal to unlimited
type Availabillity struct {
	From Date `json:"avFrom"`
	To   Date `json:"avTo,omitempty"`
}

//Prize represents the monthly rent
type Prize struct {
	Currency string `json:"currency"`
	Value    int    `json:"value"`
}

// Date represents a date in dd.mm.yyyy optional text if date couldn t be parsed
type Date struct {
	Day   int    `json:"day"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
	Text  string `json:"text,omitempty"`
}
