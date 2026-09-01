// Package paths contains helpers for middleware path filtering.
package paths

// IsExcluded reports whether path exactly matches an excluded path.
func IsExcluded(path string, exclude []string) bool {
	for _, p := range exclude {
		if path == p {
			return true
		}
	}
	return false
}
