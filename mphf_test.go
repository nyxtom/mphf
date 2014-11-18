package mphf

import "testing"

func TestCreate(t *testing.T) {
	dict := map[string]int{
		"hello":     1,
		"world":     4,
		"wonderful": 3,
	}

	table := Create(dict)
	val := table.Lookup("wonderful")
	if val != 3 {
		t.Errorf("Lookup(\"wonderful\") expected 3 but got %v instead", val)
	}
	val2 := table.Lookup("hello")
	if val2 != 1 {
		t.Errorf("Lookup(\"hello\") expected 1 but got %v instead", val)
	}
	val3 := table.Lookup("world")
	if val3 != 4 {
		t.Errorf("Lookup(\"world\") expected 4 but got %v instead", val)
	}
}
