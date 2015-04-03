package schemer

import (
	"encoding/json"
	"github.com/fatih/structs"
	"reflect"
	"strings"
	"strconv"
)

type Schema struct {
	SchemaURI            string                 `structs:"$schema,omitempty"`
	Id                   string                 `structs:"id,omitempty"`
	Title                string                 `structs:"title,omitempty"`
	Type                 string                 `structs:"type,omitempty"`
	Properties           map[string]interface{} `structs:"properties,omitempty"`
	AdditionalProperties bool                   `structs:"additionalProperties"`
	Required             []string               `structs:"required,omitempty"`
	tagProps             map[string]interface{} `structs:"-"`
}

func NewSchema() *Schema {
	schema := new(Schema)
	schema.SchemaURI = "http://json-schema.org/schema#"
	schema.AdditionalProperties = true
	schema.tagProps = make(map[string]interface{})
	return schema
}

func DetectSchema(val interface{}) *Schema {
	schema := NewSchema()
	schema.ApplyType(reflect.TypeOf(val))

	return schema
}

func (schema *Schema) MarshalJSON() ([]byte, error) {
	m := structs.Map(schema)
	if schema.Type != "object" || schema.AdditionalProperties {
		delete(m, "additionalProperties")
		delete(m, "properties")
	}

	for k, v := range schema.tagProps {
	  m[k] = v
	}

	return json.Marshal(m)
}

func (schema *Schema) ApplyTaggedType(typ reflect.Type, tag string) {
	schema.ApplyType(typ)
	schema.ApplyTag(tag)
}

func (schema *Schema) ApplyTag(tag string) {
  rawProps := strings.Split(tag,",")
  for i := range rawProps {
    parts := strings.Split(rawProps[i], ":")
    if len(parts) == 2 {
      schema.applyTagProp(parts[0], parts[1])
    }
  }
}

func (schema *Schema) applyTagProp(key string, raw string) error {
  var err error

  switch key {
  case "multipleOf", "minimum", "maximum":
    if strings.Index(raw,".") >= 0 {
      schema.tagProps[key], err = strconv.ParseFloat(raw, 32)
    } else {
      schema.tagProps[key], err = strconv.ParseInt(raw, 10, 32)
    }
  default:
    schema.tagProps[key] = raw
  }

  return err
}

func (schema *Schema) ApplyType(t reflect.Type) error {
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
			if fieldName := parseJSONTag(field.Tag.Get("json")); fieldName != "" {
				subSchema := NewSchema()
				subSchema.SchemaURI = ""
				subSchema.ApplyTaggedType(field.Type, field.Tag.Get("schemer"))
				schema.Properties[fieldName] = subSchema
			}
		}
	default:
		schema.Type = "null"
	}

	return nil
}

func parseJSONTag(tag string) string {
	if tag == "-" || tag == "" {
		return ""
	}

	return strings.Split(tag, ",")[0]
}
