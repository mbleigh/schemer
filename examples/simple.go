package main

import (
	"encoding/json"
	"fmt"
	"github.com/mbleigh/schemer"
)

type User struct {
	Id      int     `json:"id" schemer:"minimum:1000,maximum:2000"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
	Balance float32 `json:"-"`
	priv    string  `json:"-"`
	anon    string  `json:",omitempty"`
}

type Address struct {
	City string `json:"city"`
}

func main() {
	schema := schemer.DetectSchema((*User)(nil))
	schema.AdditionalProperties = false

	out, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Printf("%v", string(out[:]))
}
