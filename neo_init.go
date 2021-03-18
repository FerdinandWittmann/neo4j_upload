package main

import (
	"log"

	"github.com/FerdinandWittmann/neo4j_extended"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//neo4jInit helper to initialize nodes in Neo4j
func neo4jInit() {
	createPriceBins()
	//createUnlimitedDateNode()
}

//createUnlimitedDateNode creates the Node all avail to relations that are unlimited point to
func createUnlimitedDateNode() {
	req := neo4j_extended.NewNeoRequest()
	unlNode, err := req.NewNeoNode("_u", "Date", &[]neo4j_extended.NeoField{{Name: "name", Val: "unlimited"}})
	if err != nil {
		printError(createUnlimitedDateNode, err)
		return
	}
	err = req.AddMerge(unlNode)
	if err != nil {
		printError(createUnlimitedDateNode, err)
		return
	}
	_, err = req.SendNew(neo4j.AccessModeWrite)
	if err != nil {
		printError(createUnlimitedDateNode, err)
		return
	}

}

//create Price Bins adds the PrizeBins Container documented
func createPriceBins() {
	// add first node
	seedReq := neo4j_extended.NewNeoRequest()
	node, err := seedReq.NewNeoNode("p", "PrizeBin", &[]neo4j_extended.NeoField{
		{Name: "value", Val: 50},
	})
	err = seedReq.AddCreate(node)
	if err != nil {
		printError(createPriceBins, err)
	}
	err = seedReq.AddReturn(&[]*neo4j_extended.NeoNode{node})
	_, err = seedReq.SendNew(neo4j.AccessModeWrite)
	if err != nil {
		log.Println("@createPriceBins:", err)
		return
	}
	// add and connect nodes
	for i := 50; i < 3000; i += 50 {
		//_, err = req.Send(session)
		req := neo4j_extended.NewNeoRequest()
		root, err := req.NewNeoNode("p", "PrizeBin", &[]neo4j_extended.NeoField{
			{Name: "value", Val: i},
		})
		req.AddMatch(root)
		root, err = req.NewNeoNode("p", "", nil)
		if err != nil {
			printError(createPriceBins, err)
		}
		rel, err := root.AddRelation("h", "HIGHER", nil, 1, req)
		if err != nil {
			printError(createPriceBins, err)
		}
		_, err = rel.AddNode("p2", "PrizeBin", &[]neo4j_extended.NeoField{
			{Name: "value", Val: i + 50},
		}, req)
		if err != nil {
			printError(createPriceBins, err)
		}
		err = req.AddMerge(root)
		if err != nil {
			printError(createPriceBins, err)
		}
		err = req.AddReturn(&[]*neo4j_extended.NeoNode{node})
		_, err = req.SendNew(neo4j.AccessModeWrite)
		if err != nil {
			printError(createPriceBins, err)
		}

		if err != nil {
			log.Println("@createPriceBins:", err)
			return
		}
	}
}
