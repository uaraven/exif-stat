package utils

// ClearLine contains ASCII escape code to clear to the end of line
const ClearLine = "\x1b[0K"

// Shorten the string if it's longer than predefined number of character by replacing part of it with ellipsis
func Shorten(text string) string {
	if len(text) > 60 {
		return "…" + text[len(text)-61:]
	}
	return text
}
