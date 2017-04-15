package proxy

import "fmt"

// Subject returns the communication subject for the given parameters
func Subject(prefix, class, appName, asterisk string) (ret string) {
	ret = fmt.Sprintf("%s%s", prefix, class)
	if appName != "" {
		ret += "." + appName
		if asterisk != "" {
			ret += "." + asterisk
		}
	}
	return
}
