package main

import (
  "fmt"
  "encoding/json"
  "github.com/mbleigh/schemer"
)

type User struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Email string `json:"email"`
  Address Address `json:"address"`
  Balance float32 `json:"balance"`
}

type Address struct {
  City string `json:"city"`
}

func main() {
  schema := new(schemer.Schema)
  schema.Build((*User)(nil))

  out, _ := json.MarshalIndent(schema,"","  ")
  fmt.Printf("%v", string(out[:]))
}