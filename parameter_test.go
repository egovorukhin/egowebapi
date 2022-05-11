package egowebapi

import (
	"fmt"
	"testing"
	"time"
)

func TestParameters(t *testing.T) {

	type Person struct {
		Id        int    `json:"id"`
		Firstname string `json:"firstname"`
	}

	param := NewBodyParam(true, NewSchema(Person{}), "Описание")
	fmt.Printf("In Body: %+v\n", param)

	param = NewPathParam("/{id}", "Описание").SetType(TypeInteger)
	fmt.Printf("In Path: %+v\n", param)

	param = NewQueryParam("id", false, "Описание")
	fmt.Printf("In Query: %+v\n", param)

	param = NewHeaderParam("id", false, "Описание")
	fmt.Printf("In Header: %+v\n", param)

	param = NewFormDataParam("file", TypeFile, true, "Описание")
	fmt.Printf("In FormData: %+v\n", param)

	param = NewFormDataParam("id", TypeString, true, "Описание")
	fmt.Printf("In FormData: %+v\n", param)
}

func TestParameter(t *testing.T) {

	param := NewParameter("id").SetType(TypeInteger).SetDescription("Описание").SetRequired(true).SetFormat("int32")
	fmt.Printf("Param created: %+v\n", param)
}

func TestParameter_SetTypeFormat(t *testing.T) {

	param := NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(0)
	fmt.Printf("Param int: %+v\n", param)

	var i64 int64 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i64)
	fmt.Printf("Param int64: %+v\n", param)

	var i32 int32 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i32)
	fmt.Printf("Param int32: %+v\n", param)

	var i16 int16 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i16)
	fmt.Printf("Param int16: %+v\n", param)

	var i8 int8 = 0
	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(i8)
	fmt.Printf("Param int8: %+v\n", param)

	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat("")
	fmt.Printf("Param string: %+v\n", param)

	param = NewParameter("id").SetDescription("Описание").SetRequired(true).SetTypeFormat(true)
	fmt.Printf("Param boolean: %+v\n", param)
}

func TestModelToParameters(t *testing.T) {
	type User struct {
		Id          int       `ewa:"path:name=id;header:name=id,desc=заголовок"`
		Pid         int64     `ewa:"header:name=pid"`
		Firstname   string    `ewa:"query:name=firstname,required"`
		Lastname    string    `ewa:"query:name=lastname, empty"`
		Datetime    time.Time `ewa:"query:name=datetime"`
		Description string    `ewa:"desc"`
	}

	params := ModelToParameters(User{})
	for _, param := range params {
		fmt.Printf("In: %s, ", param.In)
		fmt.Printf("Name: %s, ", param.Name)
		fmt.Printf("Type: %s, ", param.Type)
		fmt.Printf("Format: %s, ", param.Format)
		fmt.Printf("Description: %s, ", param.Description)
		fmt.Printf("Required: %t, ", param.Required)
		fmt.Printf("AllowEmptyValue: %t\n", param.AllowEmptyValue)
	}
}
