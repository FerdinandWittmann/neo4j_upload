package main

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//neo4jInit helper to initialize nodes in Neo4j
func neo4jInit() {
	createPriceBins()
}

func createPriceBins() {
	session, err := CreateSession(neo4j.AccessModeWrite)
	defer session.Close()
	if err != nil {
		log.Println("@createPriceBins:", err)
		return
	}
	seedReq := NeoRequest{
		multiCypher: []string{
			"CREATE (p:PrizeBin {amount: 50})",
			"RETURN p",
		},
		params: map[string]interface{}{},
	}
	_, err = seedReq.send(session)
	if err != nil {
		log.Println("@createPriceBins:", err)
		return
	}
	for i := 50; i < 3000; i += 50 {
		req := NeoRequest{
			multiCypher: []string{
				"MATCH (p:PrizeBin {amount: $amount})",
				"MERGE (p)-[u:NEXT_BIN]->(p2:PrizeBin {amount: $amount2})",
				"RETURN p,p2",
			},
			params: map[string]interface{}{
				"amount":  i,
				"amount2": i + 50,
			},
		}
		_, err = req.send(session)
		if err != nil {
			log.Println("@createPriceBins:", err)
			return
		}
	}
}
