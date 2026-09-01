package paths

func IsExcluded(path string, exclude []string) bool {
	for _, p := range exclude {
		if path == p {
			return true
		}
	}
	return false
}
