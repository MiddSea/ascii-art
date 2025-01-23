package main

import (
	"testing"
)

func TestMakeAsciiMap(t *testing.T) {
	// Test with valid file
	asciiMap, err := makeAsciiMap("../standard.txt")
	if err != nil {
		t.Errorf("makeAsciiMap failed with valid file: %v", err)
	}
	if len(asciiMap) == 0 {
		t.Error("makeAsciiMap returned empty map for valid file")
	}

	// Test exclamation mark character
	if art, ok := asciiMap['!']; ok {
		if len(art) != 8 {
			t.Errorf("Expected 8 lines for '!', got %d", len(art))
		}
	} else {
		t.Error("'!' character not found in ascii map")
	}

	// Test with invalid file
	_, err = makeAsciiMap("nonexistent.txt")
	if err == nil {
		t.Error("makeAsciiMap should return error for nonexistent file")
	}
}
