package session

import "testing"

var stringSliceEqTests = []struct {
	Left        []string
	Right       []string
	Expected    bool
	Description string
}{
	{[]string{"1"}, []string{"1"}, true, "stringSliceEq('%v','%v') => '%v'. Expected '%v'"},
	{[]string{}, []string{}, true, "stringSliceEq('%v','%v') => '%v'. Expected '%v'"},
	{[]string{"1"}, []string{"2"}, false, "stringSliceEq('%v','%v') => '%v'. Expected '%v'"},
	{[]string{"2", "2"}, []string{"2"}, false, "stringSliceEq('%v','%v') => '%v'. Expected '%v'"},
}

func TestStringSliceEq(t *testing.T) {
	for _, test := range stringSliceEqTests {
		var failed bool

		ok := stringSliceEq(test.Left, test.Right)
		failed = ok != test.Expected
		if failed {
			t.Errorf(test.Description,
				test.Left, test.Right, // input
				ok,            // actual output
				test.Expected, // expected output
			)
		}
	}
}
