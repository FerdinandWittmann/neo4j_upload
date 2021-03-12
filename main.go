package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FerdinandWittmann/coli_crawler/neo4j_extended"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

//debubg enables prints
const debug bool = false

//NDriver one open connection
var NDriver *neo4j.Driver

func main() {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("wittmafe", "130345", ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	NDriver = (&driver)
	if err != nil {
		log.Println("@Neo4j:", err)
		return // handle error
	}
}

func processData() {
	path := "/home/workerferd/coli_data/crawler"
	iterateOverFolder(path)
}

func iterateOverFolder(path string) {
	err := filepath.Walk("/home/workerferd/coli_data/crawler",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path)
			if strings.HasSuffix(path, ".json") {
				processJSONFile(path)
			}
			return nil
		})
	if err != nil {
		printError(iterateOverFolder, err)
	}
}

func processJSONFile(path string) {
	// Load JSON Data
	file, err := ioutil.ReadFile(path)
	if err != nil {
		printError(processJSONFile, err)
		return
	}
	data := []RoomAd{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		printError(processJSONFile, err)
	}
	for _, roomAd := range data {
		processRoomAd(roomAd, 1)
	}
}

func processRoomAd(roomAd RoomAd, hashCount int) {
	//Compute Hash over URL
	hashFunction := sha1.New()
	hashFunction.Write([]byte(roomAd.URL))
	URLHash := hashFunction.Sum(nil)
	URLHashHex := hex.EncodeToString(URLHash)
	// hash until unique found
	for i := 0; i < hashCount-1; i++ {
		hashFunction.Write([]byte(URLHashHex))
		URLHash = hashFunction.Sum(nil)
		URLHashHex = hex.EncodeToString(URLHash)
	}

	fmt.Println(roomAd.URL)
	//open neo4j session
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	//Check if exisit in neo4j
	cypherText := "MATCH (r:RoomAd {id: \"" + URLHashHex + "\" }) RETURN r"
	res, err := SendSimple(session, cypherText)
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	//Check if we have a hash duplicate
	if res != nil {
		cypherText = "MATCH (r:RoomAd {id: \"" + URLHashHex + "\", url: \"" + roomAd.URL + "\"}) RETURN r"
		if err != nil {
			printError(processRoomAd, err)
			return
		}
		//increase hashCount rec-loop back
		if res != nil {
			processRoomAd(roomAd, hashCount+1)
		}
		//update the stored RoomAd
		//Update
	}

	session.Close()

	roomAd.ID = URLHashHex
	err = roomAd.Address.Fill()
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	fmt.Printf("%+v\n", roomAd)
	//add Obj to Neo4j
	err = roomAd.Insert()
	if err != nil {
		printError(processRoomAd, err)
		return
	}
}

func printError(i interface{}, err error) {
	log.Println("@" + runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name() + ":" + err.Error())
}
