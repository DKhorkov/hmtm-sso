package usecases

import "regexp"

func validateValueByRule(value, rule string) bool {
	pattern := regexp.MustCompile(rule)
	return pattern.MatchString(value)
}

func validateValueByRules(value string, rules []string) bool {
	for _, rule := range rules {
		matched, _ := regexp.MatchString(rule, value)
		if !matched {
			return false
		}
	}

	return true
}
