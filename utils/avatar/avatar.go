package avatar

import "strings"

var colors = []string{
	"#f44336", "#e91e63", "#9c27b0", "#673ab7", "#3f51b5",
	"#2196f3", "#009688", "#4caf50", "#ff9800", "#795548",
	"#607d8b", "#ff5722",
}

func Color(i int) string {
	return colors[i%len(colors)]
}

func Letter(subject string) string {
	r := []rune(subject)
	if len(r) == 0 {
		return "?"
	}
	return strings.ToUpper(string(r[0]))
}
