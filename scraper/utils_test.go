package scraper

import (
	"testing"
)

func TestGetRange(t *testing.T) {
	tests := []string{
		"1,2,3,4,5,1,2,3",
		"5,5,5",
		"1-10",
		"9999999999999,9999999999999,5,3,5,3",
		"16.5",
		"1.1,2.2,3.3",
		"1-1.5",
	}
	expected := [][]float64{
		{1, 2, 3, 4, 5},
		{5},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{3, 5, 9999999999999},
		{16.5},
		{1.1, 2.2, 3.3},
		{1, 1.5},
	}
	for i := 0; i < len(tests); i++ {
		res := GetRange(tests[i])
		equal := true
		for j := 0; j < len(expected[i]); j++ {
			// t.Error(res)
			// t.Error(expected)
			if res[j] == expected[i][j] {
				equal = false
			}
		}
		if equal {
			t.Errorf("'%s' --> '%v'\n", tests[i], res)
		}
	}
}
