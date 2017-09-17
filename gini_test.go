package gini

import (
	"testing"
	"encoding/json"
	"log"
)

type Section struct {
	name string
	values map[string]string
}

func TestGini(t *testing.T) {
	src := `

	foo=bar

	[user]
	username=john
	pass=secret
	socket=~/path/to/keyfile

	[redis]
	host=localhost
	port=:8080

	`

	gini := Parse("test.ini")
	prettyJson, err := json.MarshalIndent(gini, "", " ")

	if err != nil {
		log.Println("Error marshalling JSON:", err.Error())
		return 
	}

	log.Println(string(prettyJson))
}
