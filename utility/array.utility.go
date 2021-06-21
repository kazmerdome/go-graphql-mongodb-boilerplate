package utility

// Includes ... javascript .includes implementatiom
func Includes(s []string, searchitem string) bool {
	include := false
	for _, item := range s {
		if item == searchitem {
			include = true
		}
	}
	return include
}
