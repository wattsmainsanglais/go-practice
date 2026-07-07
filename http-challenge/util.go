package main

import "strconv"

// parseIntStrict is provided — don't modify it. Used by parseID in
// handlers.go.
func parseIntStrict(s string) (int, error) {
	if s == "" {
		return 0, errBadID
	}
	return strconv.Atoi(s)
}
