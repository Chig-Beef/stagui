package stagui

type TokenCode byte

const (
	// IDENTIFIERS
	T_COLOR TokenCode = iota
	T_PAGE
	T_BUTTON
	T_PANEL
	T_PLAINTEXT
	T_TEXTDATA
	T_IDENTIFIER

	// COMMANDS
	T_COM_RED
	T_COM_GREEN
	T_COM_BLUE
	T_COM_BGDRAW
	T_COM_BGCOLOR
	T_COM_X
	T_COM_Y
	T_COM_W
	T_COM_H
	T_COM_TEXT
	T_COM_TEXTCOLOR
	T_COM_FONT
	T_COM_FONTSIZE
	T_COM_TEXTGAP
	T_COM_HORALIGN
	T_COM_VERALIGN

	// LITERALS
	T_STRING
	T_INT
	T_ALIGN

	// BRACES
	T_L_SQUIRLY
	T_R_SQUIRLY

	T_ILLEGAL TokenCode = 255
)

func getIdentifierTokenCodeFromString(str string) TokenCode {
	switch str{
	case "COLOR":
		return T_COLOR
	case "PAGE":
		return T_PAGE
	case "BUTTON":
		return T_BUTTON
	case "PANEL":
		return T_PANEL
	case "PLAINTEXT":
		return T_PLAINTEXT
	case "TEXTDATA":
		return T_TEXTDATA
	case "RED":
		return T_COM_RED
	case "GREEN":
		return T_COM_GREEN
	case "BLUE":
		return T_COM_BLUE
	case "BGDRAW":
		return T_COM_BGDRAW
	case "BGCOLOR":
		return T_COM_BGCOLOR
	case "X":
		return T_COM_X
	case "Y":
		return T_COM_Y
	case "W":
		return T_COM_W
	case "H":
		return T_COM_H
	case "TEXT":
		return T_COM_TEXT
	case "TEXTCOLOR":
		return T_COM_TEXTCOLOR
	case "FONT":
		return T_COM_FONT
	case "FONTSIZE":
		return T_COM_FONTSIZE
	case "TEXTGAP":
		return T_COM_TEXTGAP
	case "HORALIGN":
		return T_COM_HORALIGN
	case "VERALIGN":
		return T_COM_VERALIGN
	case "START", "CENTER", "END":
		return T_ALIGN
	default:
		return T_IDENTIFIER
	}
}
