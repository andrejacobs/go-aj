package file

func (e *DefaultDirExcluder) Match(path string) (bool, error) {
	// Exclude the following dirs
	switch path {
	case "/dev":
		return true, nil
	case "/proc":
		return true, nil
	}
	// Don't exclude
	return false, nil
}
