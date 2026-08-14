package avatar

import (
	"hash/fnv"
	"strings"
)

var colors = []string{
	"#f44336", "#e91e63", "#9c27b0", "#673ab7", "#3f51b5",
	"#2196f3", "#009688", "#4caf50", "#ff9800", "#795548",
	"#607d8b", "#ff5722",
}

// Color picks a palette entry deterministically from an id string so the
// colour is stable regardless of where the email sits in the list.
func Color(id string) string {
	h := fnv.New32a()
	h.Write([]byte(id))
	return colors[h.Sum32()%uint32(len(colors))]
}

func Letter(subject string) string {
	r := []rune(subject)
	if len(r) == 0 {
		return "?"
	}
	return strings.ToUpper(string(r[0]))
}
