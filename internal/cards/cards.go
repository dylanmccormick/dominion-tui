package cards

import "fmt"

type Action string
type Card struct {
	Name string
	Image []byte //TODO
	Actions []Action
	Cost int
	Type string
}

var CardDict = map[string]Card{
	"cellar": {
		Name: "cellar",
		Image: []byte{},
		Actions: []Action{},
		Cost: 2,
		Type: "action",
	},
	"market": {
		Name: "market",
		Image: []byte{},
		Actions: []Action{},
		Cost: 5,
		Type: "action",
	},
	"merchant": {
		Name: "market",
		Image: []byte{},
		Actions: []Action{},
		Cost: 3,
		Type: "action",
	},
}


func (c Card) String() string{
	s := fmt.Sprintf(`
	Name: %s,
	Cost: %d,
	Actions: %s
	Type: %s
	`, c.Name, c.Cost, c.Actions, c.Type)

	return s
}



