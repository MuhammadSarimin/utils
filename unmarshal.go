package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Location  *struct {
		CamelCase string `json:"camelCase"`
		City      string `json:"city"`
		State     string `json:"state"`
	} `json:"location"`
	IsStudent bool `json:"isStudent"`
	Salary    int  `json:"salary"`
	Contacts  []struct {
		TypeValue string `json:"typeValue"`
		Value     string `json:"value"`
		Contacts  []struct {
			TypeValue string `json:"typeValue"`
			Value     string `json:"value"`
		} `json:"contacts"`
	} `json:"contacts"`
	Maping map[string]interface{} `json:"maping"`
}

func TestUnmarshal() {

	jsonData := []byte(`{
		"firstName": "John",
		"lastname": "John",
		"age": 30,
		"location": {
			"camelcase": "CamelCase",
			"city": "New York",
			"state": "NY"
			},
		"is_student": true,
		"salary": 50000,
		"contacts": [
			{"type_value": "email", "value": "john@example.com"},
			{"typeValue": "phone", "value": "123-456-7890"},
			{
				"typevalue": "okeoke", 
				"value": "123-456-7890",
				"contacts": [
					{"type_value": "email", "value": "john@example.com"},
					{"typeValue": "phone", "value": "123-456-7890"},
					{"typevalue": "okeoke", "value": "123-456-7890"}
				]
			}
		],
		"maping": {"key": "value"}
	}`)

	var person Person
	err := Unmarshal(jsonData, &person)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("")
	fmt.Println("======================================")

	fmt.Println("FirstName:", person.FirstName)
	fmt.Println("LastName:", person.LastName)
	fmt.Println("Age:", person.Age)
	fmt.Println("Location:", person.Location)
	fmt.Println("IsStudent:", person.IsStudent)
	fmt.Println("Salary:", person.Salary)
	fmt.Println("Contacts:")
	for _, contact := range person.Contacts {
		fmt.Printf("  Type: %s, Value: %s\n", contact.TypeValue, contact.Value)
		for _, subContact := range contact.Contacts {
			fmt.Printf("sub    Type: %s, Value: %s\n", subContact.TypeValue, subContact.Value)
		}
	}
	fmt.Println("Maping:", person.Maping)

}

func Unmarshal(b []byte, v interface{}) error {

	var temp map[string]interface{}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	temp = filter(temp, v)

	b, err := json.Marshal(temp)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func filter(m map[string]interface{}, v interface{}) map[string]interface{} {

	strct := reflect.TypeOf(v).Elem()
	filtered := make(map[string]interface{})

	for i := 0; i < strct.NumField(); i++ {

		field := strct.Field(i)
		jsonTag := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if !(jsonTag == "" || jsonTag == "-") {
			if value, ok := m[jsonTag]; ok {

				switch field.Type.Kind() {

				case reflect.Struct:
					filtered[jsonTag] = filter(value.(map[string]interface{}), reflect.New(field.Type).Interface())

				case reflect.Slice:
					sliceValue := make([]map[string]interface{}, 0)
					for _, item := range value.([]interface{}) {
						if itemMap, ok := item.(map[string]interface{}); ok {

							elemType := field.Type.Elem()
							filteredItem := filter(itemMap, reflect.New(elemType).Interface())
							sliceValue = append(sliceValue, filteredItem)

						}
					}
					filtered[jsonTag] = sliceValue

				case reflect.Pointer:
					filtered[jsonTag] = filter(value.(map[string]interface{}), reflect.New(field.Type.Elem()).Interface())

				default:
					filtered[jsonTag] = value

				}

			}
		}

	}

	return filtered
}
