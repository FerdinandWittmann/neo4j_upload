package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/FerdinandWittmann/neo4j_extended"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//debubg enables prints
const debug bool = false

func main() {
	//Driver Setup
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("wittmafe", "130345", ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		log.Println("@Neo4j:", err)
		return // handle error
	}
	neo4j_extended.SetDriver(&driver)
	//Create Wireframe for Neo4j Data Represntation
	neo4jInit()
	//ProcessJSON Data
	processData()

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
	session, err := neo4j_extended.CreateSession(neo4j.AccessModeRead)
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	if URLHashHex == "aac31256a89869bc91e4553c53c892d566300625" {
		fmt.Println("WTF")
	}
	//Check if exisit in neo4j
	cypherText := "MATCH (r:RoomAd {ID: \"" + URLHashHex + "\" }) RETURN r"
	res, err := neo4j_extended.SendSimple(*session, cypherText)
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	//Check if we have a hash duplicate
	if res != nil {
		cypherText = "MATCH (r:RoomAd {ID: \"" + URLHashHex + "\", url: \"" + roomAd.URL + "\"}) RETURN r"
		res, err := neo4j_extended.SendSimple(*session, cypherText)
		if err != nil {
			printError(processRoomAd, err)
			return
		}
		//increase hashCount rec-loop back
		if res == nil {
			processRoomAd(roomAd, hashCount+1)
		}
		roomAd.ID = URLHashHex
		err = roomAd.Update()
		if err != nil {
			printError(processRoomAd, err)
			return
		}
		return
		//update the stored RoomAd
		//Update
	}

	(*session).Close()

	roomAd.ID = URLHashHex
	err = roomAd.Address.Fill()
	if err != nil {
		printError(processRoomAd, err)
		return
	}
	//fmt.Printf("%+v\n", roomAd)
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
