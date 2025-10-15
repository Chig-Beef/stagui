package stagui

import "fmt"

type Lexer struct {
	source string
	curIndex int
	curChar byte
}

func (lexer *Lexer) skipWhiteSpace() {
	for {
		switch lexer.curChar {
		case ' ', '\t', '\n', '\r':
			lexer.getNextChar()
		default:
			return
		}
	}
}

func (lexer *Lexer) getNextChar() {
	lexer.curIndex++
	if lexer.curIndex >= len(lexer.source) {
		lexer.curChar = 0
		return
	}
	lexer.curChar = lexer.source[lexer.curIndex]
}

func (lexer *Lexer) nextToken() Token {
	out := Token{kind: T_ILLEGAL}
	lexer.skipWhiteSpace()

	switch lexer.curChar {
	case 0:
	// BRACES
	case '{':
		out.kind = T_L_SQUIRLY
		out.data = "{"
	case '}':
		out.kind = T_R_SQUIRLY
		out.data = "}"
	
	case '"': // String
		startIndex := lexer.curIndex
		lexer.getNextChar()
		for lexer.curChar != '"' && lexer.curChar != 0 {
			lexer.getNextChar()
		}

		if lexer.curChar == 0 {
			// END OF FILE??
			panic("Lexer reached end of file inside of string")
		}

		endIndex := lexer.curIndex+1

		str := lexer.source[startIndex:endIndex]
		out.kind = T_STRING
		out.data = str

	default: // Identifiers, and ints
		if charIsNum(lexer.curChar) {
			startIndex := lexer.curIndex
			lexer.getNextChar()
			for charIsNumOrDot(lexer.curChar) {
				lexer.getNextChar()
			}

			endIndex := lexer.curIndex

			str := lexer.source[startIndex:endIndex]
			// TODO: If dot, float
			out.kind = T_INT
			out.data = str

		} else if charIsAlpha(lexer.curChar) {
			startIndex := lexer.curIndex
			lexer.getNextChar()
			for charIsIndentifierValid(lexer.curChar) {
				lexer.getNextChar()
			}

			endIndex := lexer.curIndex

			str := lexer.source[startIndex:endIndex]
			out.kind = getIdentifierTokenCodeFromString(str)
			out.data = str
		} else {
			panic(fmt.Sprintf("Invalid character for start of token %c %d", lexer.curChar, lexer.curChar))
		}
	}

	lexer.getNextChar()
	return out
}
