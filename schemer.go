package schemer

import (
  "reflect"
  // "fmt"
)

type Schema struct {
  SchemaURI string `json:"$schema,omitempty"`
  Id string `json:"id,omitempty"`
  Title string `json:"title,omitempty"`
  Type string `json:"type,omitempty"`
  Properties map[string]interface{} `json:"properties,omitempty"`
}

func (schema *Schema) Build(val interface{}) error {
  schema.SchemaURI = "http://json-schema.org/schema#"
  return schema.BuildType(reflect.TypeOf(val))
}

func (schema *Schema) BuildType(t reflect.Type) error {
  if t.Kind() == reflect.Ptr {
    t = t.Elem()
  }

  switch t.Kind() {
    case reflect.String:
      schema.Type = "string"
    case reflect.Array, reflect.Slice:
      schema.Type = "array"
    case reflect.Bool:
      schema.Type = "boolean"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
      schema.Type = "integer"
    case reflect.Float32, reflect.Float64:
      schema.Type = "number"
    case reflect.Map:
      schema.Type = "object"
    case reflect.Struct:
      schema.Title = t.Name()
      schema.Type = "object"
      schema.Properties = make(map[string]interface{})
      for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        subSchema := new(Schema)
        subSchema.BuildType(field.Type)
        schema.Properties[field.Tag.Get("json")] = subSchema
      }
    default:
      schema.Type = "null"
  }

  return nil
}