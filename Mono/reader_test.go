package main

import (
	"testing"
)

func TestExecute(t *testing.T) {

	var s string = "testSets/only-configs.test_scenario.xml"

	execute(&s)
}
