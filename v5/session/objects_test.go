package session

import "testing"

var addTests = []struct {
	Items       []string
	Expected    []string
	Return      []bool
	Description string
}{
	// Single adding
	{[]string{"1"}, []string{"1"}, []bool{true}, "Add('%v') => '%v'; '%v'. Expected '%v'; '%v'."},

	// Duplicate
	{[]string{"1", "1"}, []string{"1"}, []bool{true, false}, "Add('%v') => '%v'; '%v'. Expected '%v'; '%v'."},

	// Sorting
	{[]string{"2", "1"}, []string{"1", "2"}, []bool{true, true}, "Add('%v') => '%v'; '%v'. Expected '%v'; '%v'."},
}

func TestObjectsAdd(t *testing.T) {

	for _, test := range addTests {
		var failed bool
		var objects Objects
		var returns []bool
		for idx, item := range test.Items {
			returns = append(returns, objects.Add(item))
			if returns[idx] != test.Return[idx] {
				failed = true
			}
		}

		failed = failed || !stringSliceEq(objects.ids, test.Expected)

		if failed {
			t.Errorf(test.Description,
				test.Items,           // input
				returns, objects.ids, //actual output
				test.Return, test.Expected, // expected output
			)
		}
	}
}

var removeTests = []struct {
	Initial     []string
	Items       []string
	Return      []bool
	Expected    []string
	Description string
}{
	// single remove
	{[]string{"1", "2", "3"}, []string{"2"}, []bool{true}, []string{"1", "3"}, "Remove('%v') => '%v'; '%v'. Expected '%v'; '%v'."},

	// duplicate remove
	{[]string{"1", "2", "3"}, []string{"2", "2"}, []bool{true, false}, []string{"1", "3"}, "Remove('%v') => '%v'; '%v'. Expected '%v'; '%v'."},

	// missing remove
	{[]string{}, []string{"2"}, []bool{false}, []string{}, "Remove('%v') => '%v'; '%v'. Expected '%v'; '%v'."},

	// remove at end
	{[]string{"1", "2", "3"}, []string{"3"}, []bool{true}, []string{"1", "2"}, "Remove('%v') => '%v'; '%v'. Expected '%v'; '%v'."},
}

func TestObjectsRemove(t *testing.T) {

	for _, test := range removeTests {

		var objects Objects
		objects.ids = test.Initial

		var failed bool
		var returns []bool
		for idx, item := range test.Items {
			returns = append(returns, objects.Remove(item))
			if returns[idx] != test.Return[idx] {
				failed = true
			}
		}

		failed = failed || !stringSliceEq(objects.ids, test.Expected)

		if failed {
			t.Errorf(test.Description,
				test.Items,           // input
				returns, objects.ids, //actual output
				test.Return, test.Expected, // expected output
			)
		}
	}
}

var clearTests = []struct {
	Initial     []string
	Description string
}{
	{[]string{"1", "2", "3"}, "'%v'.Clear() => '%v'. Expected empty list."},
	{[]string{}, "'%v'.Clear() => '%v'. Expected empty list."},
}

func TestObjectsClear(t *testing.T) {
	for _, test := range clearTests {

		var failed bool

		var objects Objects
		objects.ids = test.Initial

		objects.Clear()

		failed = len(objects.ids) != 0

		if failed {
			t.Errorf(test.Description,
				test.Initial, // input
				objects.ids,  //actual output
			)
		}
	}
}
