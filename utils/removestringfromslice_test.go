package utils

import (
	"fmt"
	"testing"
)

func Test_RemoveStringFromSlice(t *testing.T) {
	originalSlice := []string{"A", "B", "C"}
	removeString := "B"
	expectedOutput := []string{"A", "C"}
	actualOutput := RemoveStringFromSlice(originalSlice, removeString)
	if !compareStringSlice(actualOutput, expectedOutput) {
		t.Errorf("Returned slice [%s] does not match expected output [%s]", actualOutput, expectedOutput)
	}
}

func Test_RemoveStringFromSlice_emptyString(t *testing.T) {
	originalSlice := []string{"A", "B", "C"}
	removeString := ""
	expectedOutput := []string{"A", "B", "C"}
	actualOutput := RemoveStringFromSlice(originalSlice, removeString)
	if !compareStringSlice(actualOutput, expectedOutput) {
		t.Errorf("Returned slice [%s] does not match expected output [%s]", actualOutput, expectedOutput)
	}
}

func Test_RemoveStringFromSlice_emptySlice(t *testing.T) {
	originalSlice := []string{}
	removeString := "WHATEVER"
	expectedOutput := []string{}
	actualOutput := RemoveStringFromSlice(originalSlice, removeString)
	if !compareStringSlice(actualOutput, expectedOutput) {
		t.Errorf("Returned slice [%s] does not match expected output [%s]", actualOutput, expectedOutput)
	}
}

func Test_RemoveStringFromSlice_stringNotFoundInSlice(t *testing.T) {
	originalSlice := []string{"A", "B", "C"}
	removeString := "D"
	expectedOutput := []string{"A", "B", "C"}
	actualOutput := RemoveStringFromSlice(originalSlice, removeString)
	if !compareStringSlice(actualOutput, expectedOutput) {
		t.Errorf("Returned slice [%s] does not match expected output [%s]", actualOutput, expectedOutput)
	}
}

func compareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			fmt.Printf("Checking if [%s] matches [%s]", v, b[i])
			return false
		}
	}

	return true
}
