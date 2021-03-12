package main

import "github.com/kelvins/geocoder"

//Date -
type EasyDate struct {
	Day   int
	Month int
	Year  int
}

//AvailableUntil can be date unbefristet {-1,-1,-1} or text
type AvailableUntil struct {
	date EasyDate
	text string
}

//WgCh wgzimmer.ch room ad
type WgCh struct {
	AvailableFrom  EasyDate
	AvailableUntil AvailableUntil
	Region         string
	Prize          int
	State          string
	Address        geocoder.Address
	RoomDesc       string
	Images         []string
	WgDesc         string
	InterestedIn   string
}
