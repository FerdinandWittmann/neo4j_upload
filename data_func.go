package neo4j_upload

import (
	"errors"

	"github.com/FerdinandWittmann/neo4j_extended"
)

//addMatchIDNode adds to the request the node that matches the roomAd over the ID and returns the node
func addMatchIDCypher(name string, label string, id string, req *neo4j_extended.NeoRequest) (matchNode *neo4j_extended.NeoNode, err error) {
	matchNode, err = req.NewNeoNode(name, label, &[]neo4j_extended.NeoField{{Name: "ID", Val: id}})
	if err != nil {
		return nil, err
	}
	err = req.AddMatch(matchNode)
	if err != nil {
		return nil, err
	}
	return matchNode, nil

}

//addDate adds a date object != to text or unlimited
func addMonthYear(dayRel *neo4j_extended.NeoRelation, date Date, req *neo4j_extended.NeoRequest) (err error) {
	if date.Month != 0 {
		if date.Year != 0 {
			//Check Name that it doesn't exist
			name := req.CheckName("_m")
			monthNode, err := req.NewNeoNode(name, "Month", &[]neo4j_extended.NeoField{{Name: "value", Val: date.Month}, {Name: "year", Val: date.Year}})
			if err != nil {
				return err
			}
			req.AddMerge(monthNode)
			monthNode = monthNode.ReuseNode()
			dayRel.NextNode = monthNode
		} else {
			return errors.New("@AddDate: year is zero")
		}
	} else {
		return errors.New("@AddDate: moth is zero")
	}
	return nil
}
