package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/a8m/envsubst"
	"github.com/imdario/mergo"
	"github.com/joho/godotenv"
	concatenate "github.com/paulvollmer/go-concatenate"
	"github.com/tidwall/gjson"
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
		"Services": "test",
		"domain" : "test",
		"namespace" : "kube"
	}`

	finalfile  = "finalfile"
	varfile    = "varfile"
	jsontoyaml = "finalyaml.yml"
)

func main() {
	//concatenate json values and write to yaml
	test()

	// Create each file
	createtemplate(templatefile, templatefilecontents)
	os.Create(varfile)
	os.Create(finalfile)

	// Concatenate all bash vars into one file - the var file
	concatenate.FilesToFile(varfile, 0644, "\n", oldfiles...)
	fmt.Println("Concatenation complete")

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

// func convert(i interface{}) interface{} {
// 	switch x := i.(type) {
// 	case map[interface{}]interface{}:
// 		m2 := map[string]interface{}{}
// 		for k, v := range x {
// 			m2[k.(string)] = convert(v)
// 		}
// 		return m2
// 	case []interface{}:
// 		for i, v := range x {
// 			x[i] = convert(v)
// 		}
// 	}
// 	return i
// }

// make json convert and merge with yaml with https://github.com/imdario/mergo
func test() {

	// Unmarshal JSON from file
	deployJSON := "deploy.json"
	templateJSONfile, err := ioutil.ReadFile(deployJSON)
	if err != nil {
		fmt.Println(deployJSON, "not found - please ensure configs/ipt-envconfig-application contains this file")
		panic(err)
	}
	templateJSONdata := string(templateJSONfile)
	m, ok := gjson.Parse(templateJSONdata).Value().(map[string]interface{})
	if !ok {
		// not a map
	}
	fmt.Println("From file: \n", m)

	// Unmarshal JSON from var raw input
	// fmt.Println("Inputvar:\n", s)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(s), &body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("From variable: \n", body)

	// Merge
	if err := mergo.MapWithOverwrite(&body, m); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Merged: \n", body)

	// Convert to yaml
	// body = convert(body)
	if b, err := yaml.Marshal(body); err != nil {
		panic(err)
	} else {
		fmt.Printf("Output as yaml:\n%s\n", b)
		// Write yaml to file
		f, _ := os.Create(jsontoyaml)
		f.Write([]byte(b))
		f.Close()
	}

}
