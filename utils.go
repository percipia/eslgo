package freeswitchesl

import (
	"fmt"
	"strings"
)

func buildVars(vars map[string]string) string {
	var builder strings.Builder
	for key, value := range vars {
		if builder.Len() > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		if strings.ContainsAny(value, " ") {
			builder.WriteString("'")
			builder.WriteString(value)
			builder.WriteString("'")
		} else {
			builder.WriteString(value)
		}
	}
	return fmt.Sprintf("{%s}", builder.String())
}
