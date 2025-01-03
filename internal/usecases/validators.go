package usecases

import "regexp"

func validateEmail(email string, rule string) bool {
	pattern := regexp.MustCompile(rule)
	return pattern.MatchString(email)
}

func validatePassword(password string, rules []string) bool {
	for _, rule := range rules {
		matched, _ := regexp.MatchString(rule, password)
		if !matched {
			return false
		}
	}

	return true
}
