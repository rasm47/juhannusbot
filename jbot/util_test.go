package jbot

import (
    "testing"
)

func TestReadFileToLinesNormalFile(t *testing.T) {
    actualResult, err := readFileToLines("tests/file_with_lines.txt")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    expectedResult := []string{"line1\n", "line2\n", "line3\n"}
    
    if testEq(actualResult, expectedResult) {
        t.Fatalf("Reading file_with_lines.txt was successful, but it did not generate expected results")
    }    
}

func TestReadFileToLinesNonExistentFile(t *testing.T) {
    _, err := readFileToLines("tests/file_who_does_not_exist.txt")
    if err == nil {
        t.Fatalf("Reading a non-existent file did not cause an error")
    }
}

func TestReadFileToLinesEmptyFile(t *testing.T) {
    _, err := readFileToLines("tests/empty_file.txt")
    if err == nil {
        t.Fatalf("Reading an empty file did not cause an error")
    }
}

func TestFileEmptyTrue(t *testing.T) {
    emptyStringSlice := []string{}
    if !fileEmpty(emptyStringSlice) {
        t.Fatalf("Empty string slice did not register as empty")
    }
}

func TestFileEmptyFalse(t *testing.T) {
    emptyStringSlice := []string{"test","string"}
    if fileEmpty(emptyStringSlice) {
        t.Fatalf("Non-empty string slice did not register as non-empty")
    }
}

// testEq returns true if two slices are fully equal.
func testEq(a []string, b []string) bool {
    if (a == nil) != (b == nil) { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}
