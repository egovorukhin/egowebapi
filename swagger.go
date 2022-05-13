package egowebapi

import (
	"encoding/json"
	"github.com/egovorukhin/egowebapi/security"
	"github.com/mustan989/jsonschema"
)

type Swagger struct {
	ID                  string                 `json:"id,omitempty"`
	Consumes            []string               `json:"consumes,omitempty"`
	Produces            []string               `json:"produces,omitempty"`
	Schemes             []string               `json:"schemes,omitempty"`
	Swagger             string                 `json:"swagger,omitempty"`
	Info                *Info                  `json:"info,omitempty"`
	Host                string                 `json:"host,omitempty"`
	BasePath            string                 `json:"basePath,omitempty"`
	Paths               Paths                  `json:"paths,omitempty"`
	Parameters          map[string]Parameter   `json:"parameters,omitempty"`
	Responses           map[string]Response    `json:"responses,omitempty"`
	SecurityDefinitions SecurityDefinitions    `json:"securityDefinitions,omitempty"`
	Security            Security               `json:"security,omitempty"`
	Tags                []Tag                  `json:"tags,omitempty"`
	ExternalDocs        *ExternalDocs          `json:"externalDocs,omitempty"`
	Definitions         jsonschema.Definitions `json:"definitions,omitempty"`
	models              Models
	//spec.Swagger
}

type Info struct {
	Description    string   `json:"description,omitempty"`
	Title          string   `json:"title,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Tag struct {
	Description  string        `json:"description,omitempty"`
	Name         string        `json:"name,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type Paths map[string]PathItem

type PathItem map[string]Operation

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Security []map[string][]string

type SecurityDefinitions map[string]security.Definition

type Models map[string]interface{}

const (
	InPath     = "path"
	InQuery    = "query"
	InHeader   = "header"
	InBody     = "body"
	InFormData = "formData"

	TypeString  = "string"
	TypeArray   = "array"
	TypeFile    = "file"
	TypeInteger = "integer"
	TypeObject  = "object"
	TypeBoolean = "boolean"

	CollectionFormatMulti = "multi"

	RefDefinitions = "#/definitions/"
)

func (s *Swagger) JSON() ([]byte, error) {
	return json.Marshal(s)
}

// setDefinitions Преобразование моделей в формат JSON Schema
func (s *Swagger) setDefinitions(models ...interface{}) *Swagger {
	for _, model := range models {
		s.setDefinition(model)
	}
	return s
}

// setDefinition Преобразование модели в формат JSON Schema
func (s *Swagger) setDefinition(model interface{}) *Swagger {
	schema := jsonschema.Reflect(model)
	for key, value := range schema.Definitions {
		s.Definitions[key] = value
	}
	return s
}

// SetSchemes устанавливаем схему
func (s *Swagger) SetSchemes(scheme ...string) *Swagger {
	s.Schemes = append(s.Schemes, scheme...)
	return s
}

// setPort добавление порта к хосту
func (s *Swagger) setPort(port string) *Swagger {
	s.Host += port
	return s
}

// setRefDefinitions Проверка модели на существование
func (s *Swagger) setRefDefinitions(ref string) (string, bool) {
	if model, ok := s.models[ref]; ok {
		s.setDefinition(model)
		return RefDefinitions + ref, ok
	}
	return ref, false
}

// Устанавливаем путь с операциями
func (s *Swagger) setPath(path, method string, operation Operation) *Swagger {

	// Настраиваем ссылку на модель в ответах
	for _, response := range operation.Responses {
		if response.Schema == nil {
			continue
		}
		// Пытаемся найти модель в определениях
		var exists bool
		response.Schema.Ref, exists = s.setRefDefinitions(response.Schema.Ref)
		if !exists && response.Schema.Items != nil {
			response.Schema.Items.Ref, _ = s.setRefDefinitions(response.Schema.Items.Ref)
		}
	}

	// Настраиваем ссылку на модель в параметрах
	for _, param := range operation.Parameters {
		switch param.In {
		case InBody:
			if param.Schema == nil {
				break
			}
			if _, ok := s.Definitions[param.Schema.Ref]; ok {
				param.Schema.Ref = RefDefinitions + param.Schema.Ref
			}
			break
		}
	}

	// Проверяем ключ на существование
	if _, ok := s.Paths[path]; !ok {
		s.Paths[path] = PathItem{}
	}

	// Добавляем операцию в список методов
	s.Paths[path][method] = operation
	return s
}

// setSecurityDefinition Устанавливаем необходимые поля для определения авторизации
func (s *Swagger) setSecurityDefinition(authName string, sec security.Definition) *Swagger {
	s.SecurityDefinitions[authName] = sec
	return s
}

// SetInfo Устанавливаем информацию о swagger
func (s *Swagger) SetInfo(host string, info *Info, docs *ExternalDocs) *Swagger {
	s.Info = info
	s.ExternalDocs = docs
	s.Host = host
	return s
}

// SetBasePath Устанавливаем информацию о swagger
func (s *Swagger) SetBasePath(basePath string) *Swagger {
	s.BasePath = basePath
	return s
}

// compareBasePath Сравнение базового пути и пути маршрута.
// Добавляем только те маршруты, которые включают базовый путь
func (s *Swagger) compareBasePath(path string) (bool, int) {
	l := len(s.BasePath)
	if len(path) > l && path[:l] == s.BasePath {
		return true, l
	}
	return false, l
}

// SetModel Добавить модель для определения параметров для swagger
func (s *Swagger) SetModel(name string, model interface{}) *Swagger {
	s.models[name] = model
	return s
}

// SetModels Добавить модель для определения параметров для swagger
func (s *Swagger) SetModels(models Models) *Swagger {
	for key, model := range models {
		s.SetModel(key, model)
	}
	return s
}
