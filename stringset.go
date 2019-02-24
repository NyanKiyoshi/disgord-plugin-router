// Package drouter was taken from docker's distribution repo, under Apache-2.0 license:
// https://github.com/docker/distribution/blob/0d3efadf015/registry/auth/token/stringset.go
// Thank you for your work! <3
package drouter

// StringSet is a useful type for looking up strings.
type StringSet map[string]struct{}

// NewStringSet creates a new StringSet with the given strings.
func NewStringSet(keys ...string) StringSet {
	ss := make(StringSet, len(keys))
	ss.Add(keys...)
	return ss
}

// Add inserts the given keys into this StringSet.
func (ss StringSet) Add(keys ...string) {
	for _, key := range keys {
		ss[key] = struct{}{}
	}
}

// Contains returns whether the given key is in this StringSet.
func (ss StringSet) Contains(key string) bool {
	_, ok := ss[key]
	return ok
}

// Keys returns a slice of all keys in this StringSet.
func (ss StringSet) Keys() []string {
	keys := make([]string, 0, len(ss))

	for key := range ss {
		keys = append(keys, key)
	}

	return keys
}
