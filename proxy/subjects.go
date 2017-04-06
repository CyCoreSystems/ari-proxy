package proxy

import "fmt"

// GetSubject returns the subject for a 'get' request.
func GetSubject(prefix string, application string, asteriskID string) string {
	if application == "" && asteriskID == "" {
		return fmt.Sprintf("%sget", prefix)
	} else if asteriskID == "" {
		return fmt.Sprintf("%sget.%s", prefix, application)
	}

	return fmt.Sprintf("%sget.%s.%s", prefix, application, asteriskID)
}

// CommandSubject returns the subject for a 'command' request.
func CommandSubject(prefix string, application string, asteriskID string) string {
	if application == "" && asteriskID == "" {
		return fmt.Sprintf("%scommand", prefix)
	} else if asteriskID == "" {
		return fmt.Sprintf("%scommand.%s", prefix, application)
	}

	return fmt.Sprintf("%scommand.%s.%s", prefix, application, asteriskID)
}

// CreateSubject returns the subject for an 'create' request.
func CreateSubject(prefix string, application string, asteriskID string) string {
	if application == "" && asteriskID == "" {
		return fmt.Sprintf("%screate", prefix)
	} else if asteriskID == "" {
		return fmt.Sprintf("%screate.%s", prefix, application)
	}

	return fmt.Sprintf("%screate.%s.%s", prefix, application, asteriskID)
}
