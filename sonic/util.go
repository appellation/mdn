package sonic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var intArgRegexp = regexp.MustCompile(`^(\w+)\((\d+)\)$`)

func parseIntArg(arg string) (string, int) {
	match := intArgRegexp.FindStringSubmatch(arg)
	if len(match) != 3 {
		return "", 0
	}

	d, _ := strconv.ParseInt(match[2], 10, 0)
	return match[1], int(d)
}

// prepareText splits text and escapes quotes
func prepareText(text string, length int) []string {
	l := len(text)
	if l == 0 {
		return []string{}
	}

	if l > length {
		terms := strings.SplitN(text, " ", 2)
		if len(terms) == 1 {
			terms = []string{text[:l/2], text[l/2:]}
		}
		return append(prepareText(terms[0], length), prepareText(terms[1], length)...)
	}

	return []string{fmt.Sprintf(`"%q"`, text)}
}
