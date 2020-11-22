package brainfuck

import (
	"strconv"
	"strings"
)

var shortableOpChars = map[rune]bool{
	'<': true,
	'>': true,
	'+': true,
	'-': true,
	',': true,
	'.': true,
}

func isShortableOpChar(c rune) bool {
	_, ok := shortableOpChars[c]
	return ok
}

func Shorten(program string) string {

	n := len(program)
	if n == 0 {
		return ""
	}

	type charCount struct {
		char  rune
		count int
	}

	progRunes := []rune(program)
	n = len(progRunes)

	lastRune := ' '
	var counts []charCount

	for i := 0; i < n; i++ {
		var count charCount
		r := progRunes[i]

		if isShortableOpChar(r) && lastRune == r {
			counts[len(counts)-1].count++
			if i == n-1 {
				counts = append(counts, count)
			}
			continue
		}

		count.char = r
		count.count = 1
		counts = append(counts, count)

		lastRune = r
	}

	var res strings.Builder

	for _, ct := range counts {
		if isShortableOpChar(ct.char) && ct.count > 1 {
			res.WriteString(string(ct.char) + strconv.Itoa(ct.count))
			continue
		}

		if isShortableOpChar(ct.char) || ct.char == '[' || ct.char == ']' {
			for k := 0; k < ct.count; k++ {
				res.WriteRune(ct.char)
			}
		}

	}

	return res.String()
}
