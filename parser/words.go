package parser

func Lines(text string) []string {
	lines := []string{}

	var s int
	var n int
	var ch rune
	for n, ch = range text {
		switch ch {
		case '\n':
			if n > s {
				lines = append(lines, text[s:n])
			}
			lines = append(lines, string(ch))
			s = n + 1
		default:
			break
		}
	}

	if len(text) > s {
		lines = append(lines, text[s:n])
	}

	return lines
}

func Words(text string) []string {
	tokens := []string{}

	var s int
	var n int
	var ch rune
	for n, ch = range text {
		switch ch {
		case '(', ')', '[', ']', ',', '\n':
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			tokens = append(tokens, string(ch))
			s = n + 1
		case ' ':
			if n > s {
				tokens = append(tokens, text[s:n])
			}
			s = n + 1
		default:
			break
		}
	}

	if len(text) > s {
		tokens = append(tokens, text[s:len(text)])
	}

	return tokens
}
