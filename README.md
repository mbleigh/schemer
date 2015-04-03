# Schemer: Go JSON Schema Builder

[JSON Schema](http://json-schema.org/) provides a standard means of describing
JSON data structures. It has a wide range of uses for documenting and consuming
JSON data, particularly for HTTP APIs.

Schemer allows you to quickly and easily generate JSON Schemas from Go data types
and structures.

### Example

```go
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
```

**Note:** Because Go doesn't allow for the direct passing of types into functions,
it is necessary to pass in a pointer to the type instead. This can be done either
with `new(YourType)` or `(*YourType)(nil)`.

### Progress

- [x] Basic type detection and generation of simple schema structure
- [x] (partial) Tag detection for adding validation and other metadata to fields
- [ ] Configurable handling of sub-resources with `$ref` etc. taken into account
- [ ] Configurable strictness (i.e. all fields required, etc)
- [ ] Schemas can have a definitions registry
- [x] Properly parse the JSON tag data rather than blindly assuming it's the name
- [ ] Create examples that use e.g. [gojsonschema](https://github.com/xeipuuv/gojsonschema) to validate data based on Schemer output
- [ ] Support JSON Hyper Schema