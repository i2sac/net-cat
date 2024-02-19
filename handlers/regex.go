package handlers

import "regexp"

func IsPort(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func IsIP(s string) bool {
	oct := `([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	return len(s) > 0 && (regexp.MustCompile(`^`+oct+`.`+oct+`.`+oct+`.`+oct+`$`).MatchString(s) || s == "localhost")
	// return s == "localhost"
}
