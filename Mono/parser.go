package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type TestConfiguration struct {
	XMLName        xml.Name        `xml:"scenario"`
	Configurations []Configuration `xml:"configuration"`
}

type Configuration struct {
	XMLName          xml.Name    `xml:"configuration"`
	Configuration_id string      `xml:"configuration_id"`
	Environment      string      `xml:"environment"`
	Machine_type     string      `xml:"machine_type"`
	Testcaseset      TestCaseSet `xml:"test_case_set"`
}

type TestCaseSet struct {
	XMLName              xml.Name   `xml:"test_case_set"`
	PreTestCaseSetScript []string   `xml:"pre_test_case_set_script"`
	Testcase             []TestCase `xml:"test_case"`
}

type TestCase struct {
	XMLName            xml.Name `xml:"test_case"`
	TestCaseId         string   `xml:"test_case_id"`
	MaxDuration        string   `xml:"max_duration"`
	ExecScript         string   `xml:"exec_script"`
	PostExecScript     string   `xml:"post_exec_script"`
	OnErrorScript      string   `xml:"on_error_script"`
	AdditionalDataFile string   `xml:"additional_data_file"`
}

type RequiredVMs = map[string]int

func GetRequiredVMs(testConfiguration TestConfiguration) RequiredVMs {

	requiredVMs := make(RequiredVMs)
	for i := 0; i < len(testConfiguration.Configurations); i++ {
		times, ok := requiredVMs[testConfiguration.Configurations[i].Machine_type]
		if ok {
			requiredVMs[testConfiguration.Configurations[i].Machine_type] = times + 1
		} else {
			requiredVMs[testConfiguration.Configurations[i].Machine_type] = 1
		}
	}
	return requiredVMs
}

func ParseConfigurationXML(filename string) TestConfiguration {

	// Open our xmlFile
	xmlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		os.Exit(1)
	}

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var testConfiguration TestConfiguration
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &testConfiguration)

	return testConfiguration
}
