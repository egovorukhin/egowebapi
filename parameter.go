package egowebapi

import (
	"reflect"
	"regexp"
	"strings"
)

type Parameter struct {
	Path             string  `json:"-"`
	Type             string  `json:"type,omitempty"`
	Format           string  `json:"format,omitempty"`
	Description      string  `json:"description,omitempty"`
	Name             string  `json:"name,omitempty"`
	In               string  `json:"in,omitempty"`
	Required         bool    `json:"required,omitempty"`
	Schema           *Schema `json:"schema,omitempty"`
	CollectionFormat string  `json:"collectionFormat,omitempty"`
	AllowEmptyValue  bool    `json:"allowEmptyValue,omitempty"`
	Items            *Items  `json:"items,omitempty"`
}

const (
	tagEWA = "ewa"
)

// NewPathParam Инициализация параметра in: path
func NewPathParam(path string, desc ...string) *Parameter {

	// Извлекаем параметр из пути
	matches := regexp.MustCompile(`{(\w+)}`).FindStringSubmatch(path)
	if len(matches) < 2 {
		return nil
	}

	p := &Parameter{
		Path:     path,
		Name:     matches[1],
		In:       InPath,
		Required: true,
		Type:     TypeString,
	}
	if desc != nil {
		p.Description = desc[0]
	}

	return p
}

// NewBodyParam Инициализация параметра in: body
func NewBodyParam(required bool, schema *Schema, desc ...string) *Parameter {
	p := &Parameter{
		In:       InBody,
		Name:     InBody,
		Required: required,
		Schema:   schema,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewHeaderParam Инициализация параметра in: header
func NewHeaderParam(name string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InHeader,
		Name:     name,
		Type:     TypeString,
		Required: required,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewQueryParam Инициализация параметра in: query
func NewQueryParam(name string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InQuery,
		Name:     name,
		Type:     TypeString,
		Required: required,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewQueryArrayParam Инициализация параметра in: query с типом массив
func NewQueryArrayParam(name, array string, required bool, desc ...string) *Parameter {
	var (
		enum        []interface{}
		defaultItem string
	)
	for i, a := range strings.Split(array, ",") {
		a = strings.Trim(a, " ")
		if i == 0 {
			defaultItem = a
		}
		enum = append(enum, a)
	}
	p := &Parameter{
		In:               InQuery,
		Name:             name,
		Type:             TypeArray,
		CollectionFormat: CollectionFormatMulti,
		Required:         required,
		Items: &Items{
			CommonValidations: CommonValidations{
				Enum: enum,
			},
			SimpleSchema: SimpleSchema{
				Type: TypeString,
				//CollectionFormat: CollectionFormatMulti,
				Default: defaultItem,
			},
		},
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewFormDataParam Инициализация параметра in: formData
func NewFormDataParam(name, t string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		Name:     name,
		In:       InFormData,
		Required: required,
		Type:     t,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

func ModelToParameters(v interface{}) (p []*Parameter) {

	if v == nil {
		return nil
	}

	Type := reflect.TypeOf(v)
	if Type.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < Type.NumField(); i++ {
		field := Type.Field(i)
		if tag, ok := field.Tag.Lookup(tagEWA); ok {
			var param *Parameter
			for _, tagValue := range strings.Split(tag, ";") {

				names := strings.Split(tagValue, ":")
				if (len(names) == 0 || len(names) < 2) || names[1] == "" {
					continue
				}
				t, f := setTypeFormat(reflect.Indirect(reflect.ValueOf(v)).Field(i).Interface())
				name := strings.ToLower(strings.Trim(names[0], " "))
				switch name {
				case InHeader:
					param = &Parameter{
						Type:   t,
						Format: f,
						In:     name,
					}
				case InPath:
					param = &Parameter{
						Type:     t,
						Format:   f,
						In:       name,
						Required: true,
					}
				case InQuery:
					param = &Parameter{
						Type:   t,
						Format: f,
						In:     name,
					}
				default:
					continue
				}

				for _, value := range strings.Split(names[1], ",") {
					items := strings.Split(value, "=")
					if len(items) == 0 {
						continue
					}
					item := strings.ToLower(strings.Trim(items[0], " "))
					switch item {
					case "required":
						param.Required = true
					case "empty":
						param.AllowEmptyValue = true
					}
					if len(items) < 2 {
						continue
					}
					switch item {
					case "name":
						param.Name = items[1]
					case "format":
						param.Format = items[1]
					case "type":
						param.Type = items[1]
					case "desc":
						param.Description = items[1]
					default:
						continue
					}
				}

				p = append(p, param)
			}
		}
	}

	return
}

// NewParameter Инициализация параметра
func NewParameter(name string) *Parameter {
	return &Parameter{
		Name: name,
	}
}

// SetType Установка типа данных параметра
func (p *Parameter) SetType(t string) *Parameter {
	p.Type = t
	return p
}

// SetIn Установка типа параметра
func (p *Parameter) SetIn(i string) *Parameter {
	p.In = i
	return p
}

// SetFormat Установка формата данных параметра
func (p *Parameter) SetFormat(format string) *Parameter {
	p.Format = format
	return p
}

// SetRequired Установка флага обязательности параметра
func (p *Parameter) SetRequired(required bool) *Parameter {
	p.Required = required
	return p
}

// SetDescription Установка описания параметра
func (p *Parameter) SetDescription(desc string) *Parameter {
	p.Description = desc
	return p
}

// SetSchema Установка схемы параметра
func (p *Parameter) SetSchema(schema *Schema) *Parameter {
	p.Schema = schema
	return p
}

// SetTypeFormat Установка типа и формата параметра
func (p *Parameter) SetTypeFormat(t interface{}) *Parameter {
	p.Type, p.Format = setTypeFormat(t)
	return p
}

// SetAllowEmptyValue Установка на отправку пустого параметра
func (p *Parameter) SetAllowEmptyValue(allowEmptyValue bool) *Parameter {
	p.AllowEmptyValue = allowEmptyValue
	return p
}
