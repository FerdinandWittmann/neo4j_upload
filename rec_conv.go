package main

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func recToRoomAd(rec *neo4j.Record) (err error) {
	//var r RoomAd
	roomAdNode := (*rec).Values()[0].(neo4j.Node)
	var r RoomAd
	err = mapstructure.Decode(roomAdNode.Props(), &r)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", r)
	return nil

}
