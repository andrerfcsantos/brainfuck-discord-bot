package brainfuck

import "strings"

var G [256][256]string

// Taken from https://codegolf.stackexchange.com/a/5440
func init() {
	// initial state for G[x][y]: go from x to y using +s or -s.
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			delta := y - x

			if delta > 128 {
				delta -= 256
			}

			if delta < -128 {
				delta += 256
			}

			if delta >= 0 {
				G[x][y] = strings.Repeat("+", delta)
			} else {
				G[x][y] = strings.Repeat("-", -delta)
			}
		}
	}

	// keep applying rules until we can't find any more shortenings
	iter := true
	for iter {
		iter = false

		// multiplication by n/d
		for x := 0; x < 256; x++ {
			for n := 1; n < 40; n++ {
				for d := 1; d < 40; d++ {
					j := x
					y := 0
					for i := 0; i < 256; i++ {
						if j == 0 {
							break
						}
						j = (j - d + 256) & 255
						y = (y + n) & 255
					}
					if j == 0 {
						s := "[" + strings.Repeat("-", d) + ">" + strings.Repeat("+", n) + "<]>"
						if len(s) < len(G[x][y]) {
							G[x][y] = s
							iter = true
						}
					}

					j = x
					y = 0
					for i := 0; i < 256; i++ {
						if j == 0 {
							break
						}
						j = (j + d) & 255
						y = (y - n + 256) & 255
					}
					if j == 0 {
						s := "[" + strings.Repeat("+", d) + ">" + strings.Repeat("-", n) + "<]>"
						if len(s) < len(G[x][y]) {
							G[x][y] = s
							iter = true
						}
					}
				}
			}
		}

		// combine number schemes
		for x := 0; x < 256; x++ {
			for y := 0; y < 256; y++ {
				for z := 0; z < 256; z++ {
					if len(G[x][z])+len(G[z][y]) < len(G[x][y]) {
						G[x][y] = G[x][z] + G[z][y]
						iter = true
					}
				}
			}
		}
	}
}

// Encode creates a Brainfuck program that outputs the given string
func Encode(s string) string {
	var res strings.Builder
	lastc := rune(0)

	for _, c := range s {
		a := G[lastc][c]
		b := G[0][c]
		if len(a) <= len(b) {
			res.WriteString(a)
		} else {
			res.WriteString(">" + b)
		}
		res.WriteString(".")
		lastc = c
	}
	return res.String()
}
