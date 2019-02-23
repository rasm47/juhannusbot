package jbot

import (
    "testing"
)

func TeststringSlicesAreEqualNormal(t *testing.T) {
    
    stringSliceA := []string{"a","b","c"}
    stringSliceB := []string{"a","b","c"}
    
    if !stringSlicesAreEqual(stringSliceA,stringSliceB){
        t.Fatalf("Equal string slices did not return true")
    }
}

func TeststringSlicesAreEqualSizeWrong(t *testing.T) {
    
    stringSliceA := []string{"a","b","c"}
    stringSliceB := []string{"a","b"}
    
    if stringSlicesAreEqual(stringSliceA,stringSliceB){
        t.Fatalf("string slices with different sizes returned true")
    }
}

func TeststringSlicesAreEqualContentWrong(t *testing.T) {
    
    stringSliceA := []string{"a","b","c"}
    stringSliceB := []string{"a","b","C"}
    
    if !stringSlicesAreEqual(stringSliceA,stringSliceB){
        t.Fatalf("A difference in string content did not cause the function to return false")
    }
}

func TeststringSlicesAreEqualEmptySlices(t *testing.T) {
    
    stringSliceA := []string{}
    stringSliceB := []string{}
    
    if !stringSlicesAreEqual(stringSliceA,stringSliceB){
        t.Fatalf("Two empty slices were not seen as equal")
    }
}
