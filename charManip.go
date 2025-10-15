package stagui

func charIsNum(c byte) bool {
	return c >= '0' && c <= '9'
}

func charIsNumOrDot(c byte) bool {
	return charIsNum(c) || c == '.'
}

func charIsLowerCaseAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func charIsUpperCaseAlpha(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func charIsAlpha(c byte) bool {
	return charIsLowerCaseAlpha(c) || charIsUpperCaseAlpha(c)
}

func charIsIndentifierValid(c byte) bool {
	return charIsAlpha(c) ||
		charIsNum(c) ||
		c == '_'
}

func charIsWhiteSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
