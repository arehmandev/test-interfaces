package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/a8m/envsubst"
	"github.com/joho/godotenv"
	concatenate "github.com/paulvollmer/go-concatenate"
	yaml "gopkg.in/yaml.v2"
)

var (
	oldfiles     = []string{"test.properties", "test2.properties", "test3.properties"}
	templatefile = "templatefile"

	templatefilecontents = `
	I live at $HOME
	Username: $abs
	Password: $test
	My greeting is: $hello
	`

	s = `{
		"Services": "test"
	}`

	finalfile = "finalfile"
	varfile   = "varfile"
)

func main() {

	// Create each file
	createtemplate(templatefile, templatefilecontents)
	os.Create(varfile)
	os.Create(finalfile)

	// Concatenate all bash vars into one file - the var file
	concatenate.FilesToFile(varfile, 0644, "\n", oldfiles...)
	// fmt.Println("Concatenation complete")
	test()

	// Load the variables into memory and interpolate the template file ([]byte output)
	_ = godotenv.Load(varfile)
	str, err := envsubst.ReadFile(templatefile)
	if err != nil {
		log.Fatal(err)
	}

	// Write to te final file
	ioutil.WriteFile(finalfile, str, 0644)

}

func createtemplate(createtemplatepath, createtemplatecontents string) {
	// create the template
	f, _ := os.Create(createtemplatepath)
	f.Write([]byte(createtemplatecontents))
	f.Close()
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

// make json convert and merge with yaml with https://github.com/imdario/mergo
func test() {
	fmt.Printf("Input: %s\n \n", s)
	var body interface{}
	if err := json.Unmarshal([]byte(s), &body); err != nil {
		panic(err)
	}

	body = convert(body)

	if b, err := yaml.Marshal(body); err != nil {
		panic(err)
	} else {
		fmt.Printf("Output:\n%s\n", b)
	}
}
