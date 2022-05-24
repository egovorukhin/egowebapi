package egowebapi

import (
	"time"
)

type Operation struct {
	Description  string               `json:"description,omitempty"`
	Consumes     []string             `json:"consumes,omitempty"`
	Produces     []string             `json:"produces,omitempty"`
	Schemes      []string             `json:"schemes,omitempty"`
	Tags         []string             `json:"tags,omitempty"`
	Summary      string               `json:"summary,omitempty"`
	ExternalDocs *ExternalDocs        `json:"externalDocs,omitempty"`
	ID           string               `json:"operationId,omitempty"`
	Deprecated   bool                 `json:"deprecated,omitempty"`
	Security     Security             `json:"security,omitempty"`
	Parameters   []*Parameter         `json:"parameters,omitempty"`
	Responses    map[string]*Response `json:"responses,omitempty"`
}

type Schema struct {
	Ref   string `json:"$ref,omitempty"`
	Type  string `json:"type,omitempty"`
	Items *Items `json:"items,omitempty"`
	/*Discriminator string        `json:"discriminator,omitempty"`
	ReadOnly      bool          `json:"readOnly,omitempty"`
	XML           *XMLObject    `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocs `json:"externalDocs,omitempty"`
	Example       interface{}   `json:"example,omitempty"`*/
}

/*type XMLObject struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}*/

type Response struct {
	Description string                 `json:"description"`
	Schema      *Schema                `json:"schema,omitempty"`
	Headers     Headers                `json:"headers,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
}

type Header struct {
	Description string `json:"description,omitempty"`
	CommonValidations
	SimpleSchema
}

type Headers map[string]Header

type CommonValidations struct {
	Maximum          *float64      `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64      `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64        `json:"maxLength,omitempty"`
	MinLength        *int64        `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         *int64        `json:"maxItems,omitempty"`
	MinItems         *int64        `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       *float64      `json:"multipleOf,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
}

type SimpleSchema struct {
	Type             string      `json:"type,omitempty"`
	Nullable         bool        `json:"nullable,omitempty"`
	Format           string      `json:"format,omitempty"`
	Items            *Items      `json:"items,omitempty"`
	CollectionFormat string      `json:"collectionFormat,omitempty"`
	Default          interface{} `json:"default,omitempty"`
	Example          interface{} `json:"example,omitempty"`
}

type Items struct {
	Ref string `json:"$ref,omitempty"`
	CommonValidations
	SimpleSchema
}

// NewSchema Инициализация схемы для параметров
func NewSchema(modelName string, isArray bool) *Schema {
	if len(modelName) == 0 {
		return nil
	}
	if isArray {
		return &Schema{
			Type: TypeArray,
			Items: &Items{
				Ref: modelName,
			},
		}
	}

	return &Schema{
		Ref: modelName,
	}
}

// getPathParams Извлекаем пути из параметров
func (o Operation) getPathParams() (params string) {
	for _, param := range o.Parameters {
		if param.In != InPath {
			continue
		}
		params += param.Path
	}
	return
}

// getParams Извлекаем пути из параметров
func (o Operation) getParams(excludes ...string) (params []*Parameter) {
	for _, param := range o.Parameters {
		isTrue := true
		for _, ex := range excludes {
			if param.In == InPath && param.Path == ex {
				isTrue = false
				break
			}
		}
		if isTrue {
			params = append(params, param)
		}
	}
	return
}

// addTag Добавить tag в список
func (o *Operation) addTag(tag string) {
	o.Tags = append(o.Tags, tag)
}

// NewResponse Инициализация ответа
func NewResponse(schema *Schema, desc ...string) Response {
	r := Response{
		Schema:  schema,
		Headers: map[string]Header{},
	}
	if desc != nil {
		r.Description = desc[0]
	}
	return r
}

// AddHeader Добавить заголовок в ответ
func (r Response) AddHeader(name string, header Header) Response {
	r.Headers[name] = header
	return r
}

// NewHeader Инициализация заголовка
func NewHeader(t interface{}, nullable bool, desc ...string) Header {

	h := Header{
		SimpleSchema: SimpleSchema{
			Nullable: nullable,
		},
	}
	if desc != nil {
		h.Description = desc[0]
	}

	h.Type, h.Format = setTypeFormat(t)

	return h
}

// setTypeFormat Получение типа и формата на основе интерфейса
func setTypeFormat(t interface{}) (Type string, Format string) {

	switch t.(type) {
	case string, *string:
		Type = TypeString
		break
	case int, *int:
		Type = TypeInteger
		break
	case int8, *int8:
		Type = TypeInteger
		Format = "int8"
		break
	case int16, *int16:
		Type = TypeInteger
		Format = "int16"
		break
	case int32, *int32:
		Type = TypeInteger
		Format = "int32"
		break
	case int64, *int64:
		Type = TypeInteger
		Format = "int64"
		break
	case time.Time, *time.Time:
		Type = TypeString
		Format = "date-time"
		break
	case bool, *bool:
		Type = TypeBoolean
		break
	}
	return
}
