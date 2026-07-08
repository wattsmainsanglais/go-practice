package main

import "testing"

// Throwaway demo — not a real exercise. Delete this file whenever.
// Run it on its own with: go test -run TestNilMapVsMakeMap -v .
func TestNilMapVsMakeMap(t *testing.T) {
	var nilMap map[string]int
	t.Logf("nilMap == nil: %v", nilMap == nil)
	t.Logf("read from nil map, missing key: %v", nilMap["x"]) // reading is fine

	realMap := make(map[string]int)
	t.Logf("realMap == nil: %v", realMap == nil)
	realMap["x"] = 1 // writing is fine
	t.Logf("realMap[\"x\"]: %v", realMap["x"])

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("writing to nilMap panicked, as expected: %v", r)
			}
		}()
		t.Log("now writing to nilMap...")
		nilMap["x"] = 1 // this panics — recovered above so the test doesn't fail
	}()
}
