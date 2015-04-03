package schemer

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"reflect"
	"strconv"
	"strings"
)

type Schema struct {
	SchemaURI            string                 `structs:"$schema,omitempty"`
	Id                   string                 `structs:"id,omitempty"`
	Title                string                 `structs:"title,omitempty"`
	Type                 string                 `structs:"type,omitempty"`
	Format               string                 `structs:"format,omitempty"`
	Required             []string               `structs:"required,omitempty"`
	AdditionalProperties bool                   `structs:"additionalProperties"`
	AdditionalItems      bool                   `structs:"additionalItems"`
	Properties           map[string]*Schema     `structs:"-"`
	Items                *Schema                `structs:"-"`
	Not                  *Schema                `structs:"-"`
	AnyOf                []*Schema              `structs:"-"`
	AllOf                []*Schema              `structs:"-"`
	OneOf                []*Schema              `structs:"-"`
	tagProps             map[string]interface{} `structs:"-"`
	customProps          map[string]interface{} `structs:"-"`
}

func RootSchema() *Schema {
	schema := NewSchema()
	schema.SchemaURI = "http://json-schema.org/schema#"
	return schema
}

func NewSchema() *Schema {
	schema := new(Schema)
	schema.AdditionalProperties = true
	schema.AdditionalItems = true
	schema.tagProps = make(map[string]interface{})
	schema.customProps = make(map[string]interface{})
	return schema
}

func DetectSchema(val interface{}) *Schema {
	schema := RootSchema()
	schema.ApplyType(reflect.TypeOf(val))

	return schema
}

func (schema *Schema) MarshalJSON() ([]byte, error) {
	fmt.Println(schema.Type)
	m := structs.Map(schema)

	if schema.Type != "object" || schema.AdditionalProperties {
		delete(m, "additionalProperties")
	}
	if schema.Type != "array" || schema.AdditionalItems {
		delete(m, "additionalItems")
	}
	if schema.Properties != nil {
		m["properties"] = schema.Properties
	}
	if schema.Items != nil {
		m["items"] = schema.Items
	}
	if schema.Not != nil {
		m["not"] = schema.Not
	}
	if schema.AnyOf != nil {
		m["anyOf"] = schema.AnyOf
	}
	if schema.AllOf != nil {
		m["allOf"] = schema.AllOf
	}
	if schema.OneOf != nil {
		m["oneOf"] = schema.OneOf
	}

	for k, v := range schema.tagProps {
		m[k] = v
	}
	for k, v := range schema.customProps {
		m[k] = v
	}

	return json.Marshal(m)
}

func (schema *Schema) ApplyTaggedType(typ reflect.Type, tag string) {
	schema.ApplyType(typ)
	schema.ApplyTag(tag)
}

func (schema *Schema) ApplyTag(tag string) {
	rawProps := strings.Split(tag, ",")
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
		if strings.Index(raw, ".") >= 0 {
			schema.tagProps[key], err = strconv.ParseFloat(raw, 32)
		} else {
			schema.tagProps[key], err = strconv.ParseInt(raw, 10, 32)
		}
	default:
		schema.tagProps[key] = raw
	}

	return err
}

func (schema *Schema) SetProp(prop string, val interface{}) {
	schema.customProps[prop] = val
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

		subSchema := NewSchema()
		subSchema.SchemaURI = ""
		subSchema.ApplyType(t.Elem())

		schema.Items = subSchema
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
		schema.Properties = make(map[string]*Schema)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if fieldName := parseJSONTagName(field.Tag.Get("json")); fieldName != "" {
				subSchema := NewSchema()
				subSchema.ApplyTaggedType(field.Type, field.Tag.Get("schemer"))
				schema.Properties[fieldName] = subSchema
			}
		}
	default:
		schema.Type = "null"
	}

	return nil
}

func parseJSONTagName(tag string) string {
	if tag == "-" || tag == "" {
		return ""
	}

	return strings.Split(tag, ",")[0]
}
