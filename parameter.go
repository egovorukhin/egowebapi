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
		p.SetDescription(desc[0])
	}

	return p
}

// NewBodyParam Инициализация параметра in: body
func NewBodyParam(required bool, modelName string, isArray bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InBody,
		Name:     InBody,
		Required: required,
		Schema:   NewSchema(modelName, isArray),
	}
	if desc != nil {
		p.SetDescription(desc[0])
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
		p.SetDescription(desc[0])
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
		p.SetDescription(desc[0])
	}
	return p
}

// NewQueryArrayParam Инициализация параметра in: query с типом массив
func NewQueryArrayParam(name, array string, required bool, desc ...string) *Parameter {

	p := &Parameter{
		In:               InQuery,
		Name:             name,
		Type:             TypeArray,
		CollectionFormat: CollectionFormatMulti,
		Required:         required,
	}
	p.SetItems(array)
	if desc != nil {
		p.SetDescription(desc[0])
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
		p.SetDescription(desc[0])
	}
	return p
}

func ModelToParameters(v interface{}) (p []*Parameter) {

	if v == nil {
		return
	}

	Type := reflect.TypeOf(v)
	if Type.Kind() == reflect.Ptr {
		Type = Type.Elem()
	}
	if Type.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < Type.NumField(); i++ {
		field := Type.Field(i)
		if tag, ok := field.Tag.Lookup(tagEWA); ok {
			var param *Parameter
			for _, tagValue := range strings.Split(tag, ";") {

				inNames := strings.Split(tagValue, ":")
				if (len(inNames) == 0 || len(inNames) < 2) || inNames[1] == "" {
					continue
				}
				t, f := setTypeFormat(reflect.Indirect(reflect.ValueOf(v)).Field(i).Interface())
				inName := strings.ToLower(strings.Trim(inNames[0], " "))
				switch inName {
				case InPath, InHeader, InQuery:

					// Инициализация параметра
					param = NewParameter(inName).SetName(strings.ToLower(field.Name)).SetType(t).SetFormat(f)
					// Если параметр пути, то прописываем свойства *обязательно
					if inName == InPath {
						param.SetRequired(true)
					}

					// Значения
					values := strings.Split(inNames[1], ",")

					// Получаем значения и устанавливаем в параметры
					for _, value := range values {
						items := strings.Split(value, "=")
						if len(items) == 0 {
							continue
						}
						item := strings.ToLower(strings.Trim(items[0], " "))
						switch item {
						case "required":
							param.Required = true
						case "empty":
							param.SetAllowEmptyValue(true)
						default:
							// Проверяем является значение форматом пути, если да, то добавляем в путь параметра
							if len(item) > 0 && item[0] == '/' && inName == InPath {
								param.Path = item
								break
							}
						}
						// Если нет значений после равно, то выходим
						if len(items) < 2 {
							continue
						}
						// Получаем значения после равно
						switch item {
						case "name":
							param.SetName(items[1])
						case "format":
							param.SetFormat(items[1])
						case "type":
							param.SetType(items[1])
						case "array":
							if inName == InQuery {
								param.SetItems(items[1])
							}
						}
					}
					// Добавляем параметры в список
					if len(values) > 0 {
						p = append(p, param)
					}
				}
			}
		}
	}

	return
}

// NewParameter Инициализация параметра
func NewParameter(in string) *Parameter {
	return &Parameter{
		In: in,
	}
}

// SetType Установка типа данных параметра
func (p *Parameter) SetType(t string) *Parameter {
	p.Type = t
	return p
}

// SetName Установка типа параметра
func (p *Parameter) SetName(name string) *Parameter {
	p.Name = name
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

// SetCollectionFormat Установка на отправку пустого параметра
func (p *Parameter) SetCollectionFormat(format string) *Parameter {
	p.CollectionFormat = format
	return p
}

// SetItems Установка на отправку пустого параметра
func (p *Parameter) SetItems(array string) *Parameter {

	var (
		enum        []interface{}
		defaultItem string
	)
	for i, a := range strings.Split(array, "&") {
		a = strings.Trim(a, " ")
		if i == 0 {
			defaultItem = a
		}
		enum = append(enum, a)
	}
	p.SetCollectionFormat(CollectionFormatMulti)
	p.Items = &Items{
		CommonValidations: CommonValidations{
			Enum: enum,
		},
		SimpleSchema: SimpleSchema{
			Type: TypeString,
			//CollectionFormat: CollectionFormatMulti,
			Default: defaultItem,
		},
	}
	return p
}
