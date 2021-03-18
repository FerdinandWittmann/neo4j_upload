package main

import (
	"errors"
	"fmt"

	"github.com/FerdinandWittmann/neo4j_extended"
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
	err = r.InsertFields()
	if err != nil {
		return err
	}
	err = r.InsertAttributes()
	if err != nil {
		return err
	}
	err = r.InsertAvailabillity()
	if err != nil {
		return err
	}
	err = r.InsertPrize()
	if err != nil {
		return err
	}
	err = r.InsertAddress()
	if err != nil {
		return err
	}
	return nil
}

//Update checks if fields have been updated and changes their value
//Could just a reinsert, but take that connections stay for collab filtering
//Propably not trivial
func (r *RoomAd) Update() (err error) {
	fmt.Println(r.URL + ": already in DB")
	req := neo4j_extended.NewNeoRequest()
	matchNode, err := addMatchIDCypher("m", "RoomAd", r.ID, req)
	rel, err := matchNode.AddRelation("r", "", nil, 1, req)
	_, err = rel.AddNode("", "", nil, req)
	req.AddMatch(matchNode)
	req.SendNew(neo4j.AccessModeWrite)
	return nil
}

//InsertAddress inserts the roomads address
func (r *RoomAd) InsertAddress() (err error) {
	req := neo4j_extended.NewNeoRequest()
	address := r.Address
	matchNode, err := addMatchIDCypher("m", "RoomAd", r.ID, req)
	if err != nil {
		return err
	}
	matchNode = matchNode.ReuseNode()
	cityNode, err := req.NewNeoNode("c", "City", &[]neo4j_extended.NeoField{{Name: "country", Val: address.Country}, {Name: "city", Val: address.City}})
	if err != nil {
		return err
	}
	err = req.AddMerge(cityNode)
	cityNode = cityNode.ReuseNode()
	if err != nil {
		return err
	}
	postCodeNode, err := req.NewNeoNode("p", "PostCode", &[]neo4j_extended.NeoField{{Name: "value", Val: address.PostCode}})
	if err != nil {
		return err
	}
	inRel, err := postCodeNode.AddRelation("in_c", "IN_CITY", nil, 1, req)
	if err != nil {
		return err
	}
	inRel.NextNode = cityNode
	err = req.AddMerge(postCodeNode)
	if err != nil {
		return err
	}
	postCodeNode = postCodeNode.ReuseNode()
	streetNode, err := req.NewNeoNode("s", "Street", &[]neo4j_extended.NeoField{{Name: "name", Val: address.Streetname}})
	if err != nil {
		return err
	}
	inPostCodeRel, err := streetNode.AddRelation("in_p", "IN_POSTCODE", nil, 1, req)
	if err != nil {
		return err
	}
	inPostCodeRel.NextNode = postCodeNode
	err = req.AddMerge(streetNode)
	if err != nil {
		return err
	}
	streetNode = streetNode.ReuseNode()
	livesInRel, err := matchNode.AddRelation("in_s", "LIVES_IN_STREET", &[]neo4j_extended.NeoField{
		{Name: "number", Val: address.Streetnumber},
		{Name: "long", Val: address.Long},
		{Name: "lat", Val: address.Lat},
	}, 1, req)
	if err != nil {
		return err
	}
	livesInRel.NextNode = streetNode
	err = req.AddCreate(matchNode)
	if err != nil {
		return err
	}
	_, err = req.SendNew(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	return nil
}

//InsertPrize insert prize node and relation to neo4j
func (r *RoomAd) InsertPrize() (err error) {
	req := neo4j_extended.NewNeoRequest()
	prize := r.Prize
	matchNode, err := addMatchIDCypher("m", "RoomAd", r.ID, req)
	if err != nil {
		return err
	}
	matchNode = matchNode.ReuseNode()
	prizeBinMultFloat := prize.Value / 50
	prizeBinInt := int(prizeBinMultFloat)
	prizeBin := prizeBinInt * 50
	prizeBinNode, err := req.NewNeoNode("p", "PrizeBin", &[]neo4j_extended.NeoField{{Name: "value", Val: prizeBin}})
	if err != nil {
		return err
	}
	err = req.AddMatch(prizeBinNode)
	prizeBinNode = prizeBinNode.ReuseNode()
	if err != nil {
		return err
	}
	rentRel, err := matchNode.AddRelation("r", "RENT_PRIZE", &[]neo4j_extended.NeoField{{Name: "prize", Val: prize.Value}}, 1, req)
	if err != nil {
		return err
	}
	rentRel.NextNode = prizeBinNode
	err = req.AddCreate(matchNode)
	if err != nil {
		return err
	}
	_, err = req.SendNew(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	return nil
}

//InsertAvailabillity addsvailable form and to relation to the roomad
func (r *RoomAd) InsertAvailabillity() (err error) {
	//MATCH SECTION
	a := r.Availabillity
	req := neo4j_extended.NewNeoRequest()
	matchNode, err := addMatchIDCypher("m", "RoomAd", r.ID, req)
	if err != nil {
		return err
	}
	matchNode = matchNode.ReuseNode()
	req.SaveReturn(matchNode)
	//BEGIN TO SECTION
	if a.To.Day == 0 && a.To.Month == 0 && a.To.Year == 0 {
		availRel, err := matchNode.AddRelation("a", "AVAILABLE_TO", nil, 1, req)
		if err != nil {
			return err
		}
		req.SaveReturn(availRel)
		if a.To.Text == "" {
			//case unlimited
			//Match the unlimted Node first to make sure always the same ic chosen
			unlNode, err := req.NewNeoNode("u", "Date", &[]neo4j_extended.NeoField{{Name: "name", Val: "unlimited"}})
			if err != nil {
				return err
			}
			req.SaveReturn(unlNode)
			err = req.AddMerge(unlNode)
			if err != nil {
				return err
			}
			unlNode = unlNode.ReuseNode()
			availRel.NextNode = unlNode
			req.AddCreate(matchNode)
		} else {
			//case text
			//new Date Node for every entry
			dateNode, err := availRel.AddNode("u", "Date", &[]neo4j_extended.NeoField{{Name: "text", Val: a.To.Text}}, req)
			if err != nil {
				return err
			}
			req.SaveReturn(dateNode)
			err = req.AddCreate(matchNode)
		}
		if err != nil {
			return nil
		}
	} else {
		//case date
		if a.To.Day != 0 {
			availRel, err := matchNode.AddRelation("a", "AVAILABLE_TO", &[]neo4j_extended.NeoField{{Name: "day", Val: a.To.Day}}, 1, req)
			if err != nil {
				return err
			}
			err = addMonthYear(availRel, a.To, req)
			if err != nil {
				return err
			}
			err = req.AddCreate(matchNode)
		} else {
			return errors.New("@AddAvailabillity: To day is zero")
		}
	}
	//catches all addFunc errors
	if err != nil {
		return err
	}
	//END TO SECTION
	//BEGIN FROM SECTION
	matchNode = matchNode.ReuseNode()
	if a.From.Day != 0 {
		availRel, err := matchNode.AddRelation("a2", "AVAILABLE_FROM", &[]neo4j_extended.NeoField{{Name: "day", Val: a.From.Day}}, 1, req)
		if err != nil {
			return err
		}
		req.SaveReturn(availRel)
		err = addMonthYear(availRel, a.From, req)
		if err != nil {
			return err
		}
	} else {
		return errors.New("@AddAvailabillity: From day is zero")
	}
	err = req.AddCreate(matchNode)
	if err != nil {
		return err
	}
	//END FROM SECTION
	err = req.AddReturns()
	if err != nil {
		return err
	}
	_, err = req.SendNew(neo4j.AccessModeWrite)

	if err != nil {
		return err
	}
	return nil
}

//InsertAttributes takes all Attributes and Inserts them depending on their type
func (r *RoomAd) InsertAttributes() (err error) {
	if len(r.Attributes) == 0 {
		return nil
	}
	for _, attribute := range r.Attributes {
		req := neo4j_extended.NewNeoRequest()
		matchNode, err := addMatchIDCypher("m", "RoomAd", r.ID, req)
		if err != nil {
			return err
		}
		var attributeNode *neo4j_extended.NeoNode
		//Convert Node to be used in later cypher statements
		matchNode = matchNode.ReuseNode()
		//Tyüe Check over actual type instead of entry to ensure correct data format
		//TODO: if attribute without values comes it gets omitted, which case will this be here
		switch val := attribute.Value.(type) {
		case string:
			attribute.Value = val
			hasAttributeRel, err := matchNode.AddRelation("h", "HAS_ATTRIBUTE", nil, 1, req)
			if err != nil {
				return err
			}
			attributeNode, err = hasAttributeRel.AddNode(
				"a",
				"TextAttribute",
				&[]neo4j_extended.NeoField{
					{Name: "value", Val: attribute.Value},
					{Name: "name", Val: attribute.Name}},
				req)
			if err != nil {
				return err
			}

		case float64, int:
			attribute.Value = val
			hasAttributeRel, err := matchNode.AddRelation("h", "HAS_ATTRIBUTE", &[]neo4j_extended.NeoField{{Name: "value", Val: attribute.Value}}, 1, req)
			if err != nil {
				return err
			}
			attributeNode, err = hasAttributeRel.AddNode("a", "ValueAttribute", &[]neo4j_extended.NeoField{{Name: "name", Val: attribute.Name}}, req)
			if err != nil {
				return err
			}
		//case attribute
		default:
			return errors.New("Type not supported: declared-" + attribute.Type)
		}
		//Add the Create Statement for the relation
		err = req.AddMerge(matchNode)
		if err != nil {
			return err
		}
		//add a return statement for attributeNode
		err = req.AddReturn(&[]*neo4j_extended.NeoNode{attributeNode})
		res, err := req.SendNew(neo4j.AccessModeWrite)
		if err != nil {
			return err
		}
		if err = (*res).Err(); err != nil {
			return err
		}
	}
	return nil
}

//InsertDesc inserts the roomad node with its fields
func (r *RoomAd) InsertFields() (err error) {
	//TODO For Update not just create but match
	//New Request
	neoReq := neo4j_extended.NewNeoRequest()
	//CreateNode
	node, err := neoReq.NewNeoNode(
		"r",
		"RoomAd",
		&[]neo4j_extended.NeoField{
			{Name: "ID", Val: r.ID},
			{Name: "description", Val: r.Description},
			{Name: "header", Val: r.Header},
			{Name: "url", Val: r.URL},
			{Name: "crawledAt", Val: r.CrawledAt},
			{Name: "origin", Val: r.Origin},
			{Name: "images", Val: r.Images},
		},
	)
	if err != nil {
		return err
	}
	err = neoReq.AddCreate(node)
	if err != nil {
		return err
	}
	//ReturnNode
	retNode, err := neoReq.NewNeoNode("r", "", nil)
	if err != nil {
		return err
	}
	err = neoReq.AddReturn(&[]*neo4j_extended.NeoNode{retNode})
	//Send
	result, err := neoReq.SendNew(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	if err = (*result).Err(); err != nil {
		return err
	}
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
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
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
