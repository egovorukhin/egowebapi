package egowebapi

import (
	"path"
	"reflect"
	"regexp"
	"strings"
)

type IGet interface {
	Get(route *Route)
}

type IPost interface {
	Post(route *Route)
}
type IPut interface {
	Put(route *Route)
}

type IDelete interface {
	Delete(route *Route)
}

type IOptions interface {
	Options(route *Route)
}

type IPatch interface {
	Patch(route *Route)
}

type IHead interface {
	Head(route *Route)
}

type ITrace interface {
	Trace(route *Route)
}

type IConnect interface {
	Connect(route *Route)
}

type Controller struct {
	Interface interface{}
	IsShow    bool
	Name      string
	Path      string
	Suffix    []Suffix
	PathTree  []string
	FileTree  []string
	Tag       Tag
	Model     interface{}
}

// SetName Устанавливаем имя контроллера
func (c *Controller) SetName(name string) *Controller {
	c.Name = name
	c.Tag.Name = name
	return c
}

// SetDocs Устанавливаем имя контроллера
func (c *Controller) SetDocs(desc, url string) *Controller {
	c.Tag.ExternalDocs = &ExternalDocs{
		Description: desc,
		URL:         url,
	}
	return c
}

// SetPath Устанавливаем путь контроллера
func (c *Controller) SetPath(path string) *Controller {
	c.Path = path
	return c
}

// SetDescription Устанавливаем описание контроллера
func (c *Controller) SetDescription(desc string) *Controller {
	c.Tag.Description = desc
	return c
}

// SetSuffix Устанавливаем суффикс пути контроллера
func (c *Controller) SetSuffix(suffix ...Suffix) *Controller {
	c.Suffix = append(c.Suffix, suffix...)
	return c
}

// NotShow Установка флага отображения контроллера в swagger
func (c *Controller) NotShow() *Controller {
	c.IsShow = false
	return c
}

// initialize инициализация контролера
func (c *Controller) initialize(basePath string) {

	//Извлекаем имя и путь до "controllers"
	var t reflect.Type
	value := reflect.ValueOf(c.Interface)
	if value.Type().Kind() == reflect.Ptr {
		t = reflect.Indirect(value).Type()
	} else {
		t = value.Type()
	}

	pkgPath := strings.Replace(
		regexp.MustCompile(`controllers(.*)$`).FindString(t.PkgPath()),
		"controllers",
		"",
		-1,
	)

	// Путь указанный вручную
	if c.Path != "" {
		pkgPath = c.Path
	}

	// Формирование дерева путей
	c.FileTree = strings.Split(pkgPath, "/")
	c.PathTree = c.FileTree
	// Вставляем суффиксы по индексу пути
	for _, item := range c.Suffix {
		if regexp.MustCompile(`{\w+}`).MatchString(item.Value) {
			item.isParam = true
		}
		c.PathTree = insert(c.FileTree, item.Index, item.Value)
	}
	pkgPath = strings.Join(c.PathTree, "/")

	// Имя контроллера указанное в ручную
	if c.Name == "" {
		c.Name = t.Name()
	}
	c.Name = strings.ToLower(c.Name)

	// Формирование имени для тэга контроллера
	var p string
	name := c.Name
	if pkgPath != "" && pkgPath != "/" && pkgPath[:len(basePath)] == basePath {
		index := len(pkgPath)
		loc := regexp.MustCompile(`{\w+}`).FindStringIndex(pkgPath)
		if loc != nil {
			index = loc[1]
			name = ""
		}
		p = strings.ToLower(pkgPath[len(basePath):index])
	}

	c.Tag.Name = path.Join(p, name)
	if c.Path != "" {
		c.Path = pkgPath
	} else {
		c.Path = path.Join(pkgPath, c.Name)
	}
}

func (c *Controller) SetModel(model interface{}) *Controller {
	c.Model = model
	return c
}

func insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	} else if len(a) < index {
		return a
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
