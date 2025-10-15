package stagui

import (
	"fmt"
	"os"
)

type Parser struct {
	source []Token
	curToken Token
	curIndex int
}

func (parser *Parser) prevToken() {
	parser.curIndex--
	if parser.curIndex < 0 {
		parser.curToken = Token{kind: T_ILLEGAL}
		return
	}

	parser.curToken = parser.source[parser.curIndex]
}

func (parser *Parser) nextToken() {
	parser.curIndex++
	if parser.curIndex >= len(parser.source) {
		parser.curToken = Token{kind: T_ILLEGAL}
		return
	}

	parser.curToken = parser.source[parser.curIndex]
}

func (parser *Parser) ParseFromFile(filename string) *Node {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return parser.ParseFromSource(string(data))
}

func (parser *Parser) ParseFromSource(source string) *Node {
	// Get all tokens
	l := &Lexer{
		source: source,
	}

	l.curChar = l.source[0]

	t := l.nextToken()
	for t.kind != T_ILLEGAL {
		parser.source = append(parser.source, t)
		t = l.nextToken()
	}

	parser.curToken = parser.source[0]

	// Begin parsing
	return parser.parseProgram()
}

func (parser *Parser) parseProgram() *Node {
	n := Node{Kind: N_PROGRAM}

	for {
		switch parser.curToken.kind {
		case T_COLOR:
			n.addChild(parser.parseColor())
		case T_PAGE:
			n.addChild(parser.parsePage())
		case T_ILLEGAL:
			return &n
		}

		parser.nextToken()
	}
}

// CONSTRUCTS

func (parser *Parser) parseColor() *Node {
	n := Node{Kind: N_COLOR_DEF}

	if parser.curToken.kind != T_COLOR {
		panic("Expected COLOR in color def")
	}
	n.addChild(&Node{Kind: N_COLOR})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Expected STRING in color def")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})
	parser.nextToken()

	n.addChild(parser.parseColorBlock())

	return &n
}

func (parser *Parser) parsePage() *Node {
	n := Node{Kind: N_PAGE_DEF}

	if parser.curToken.kind != T_PAGE {
		panic("Expected PAGE in page def")
	}
	n.addChild(&Node{Kind: N_PAGE})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Expected STRING in page def")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})
	parser.nextToken()

	n.addChild(parser.parsePageBlock())

	return &n
}

func (parser *Parser) parseTextData() *Node {
	n := Node{Kind: N_TEXTDATA_DEF}

	if parser.curToken.kind != T_TEXTDATA {
		panic("Expected TEXTDATA in textdata def")
	}
	n.addChild(&Node{Kind: N_TEXTDATA})
	parser.nextToken()

	n.addChild(parser.parseTextDataBlock())

	return &n
}

func (parser *Parser) parseButtonDef() *Node {
	n := Node{Kind: N_BUTTON_DEF}

	if parser.curToken.kind != T_BUTTON {
		panic("Expected BUTTON in button def")
	}
	n.addChild(&Node{Kind: N_BUTTON})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Expected string in button def")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})
	parser.nextToken()

	n.addChild(parser.parseButtonBlock())

	return &n
}

func (parser *Parser) parsePlainTextDef() *Node {
	n := Node{Kind: N_PLAINTEXT_DEF}

	if parser.curToken.kind != T_PLAINTEXT {
		panic("Expected PLAINTEXT in plaintext def")
	}
	n.addChild(&Node{Kind: N_PLAINTEXT})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Expected string in plaintext def")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})
	parser.nextToken()

	n.addChild(parser.parseButtonBlock())

	return &n
}

// BLOCKS

func (parser *Parser) parseColorBlock() *Node {
	n := Node{Kind: N_COLOR_BLOCK}

	if parser.curToken.kind != T_L_SQUIRLY {
		panic("Expected { in color block")
	}
	parser.nextToken()

	for {
		switch parser.curToken.kind {
		case T_ILLEGAL:
			panic("Unexpected end in color block")

		case T_R_SQUIRLY:
			return &n

		default:
			n.addChild(parser.parseColorBlockStatement())
			parser.nextToken()
		}
	}
}

func (parser *Parser) parsePageBlock() *Node {
	n := Node{Kind: N_PAGE_BLOCK}

	if parser.curToken.kind != T_L_SQUIRLY {
		panic("Expected { in page block")
	}
	parser.nextToken()

	for {
		switch parser.curToken.kind {
		case T_ILLEGAL:
			panic("Unexpected end in page block")

		case T_R_SQUIRLY:
			return &n

		default:
			n.addChild(parser.parsePageBlockStatement())
			parser.nextToken()
		}
	}
}

func (parser *Parser) parseButtonBlock() *Node {
	n := Node{Kind: N_BUTTON_BLOCK}

	if parser.curToken.kind != T_L_SQUIRLY {
		panic("Expected { in button block")
	}
	parser.nextToken()

	for {
		switch parser.curToken.kind {
		case T_ILLEGAL:
			panic("Unexpected end in button block")

		case T_R_SQUIRLY:
			return &n

		default:
			n.addChild(parser.parseButtonBlockStatement())
			parser.nextToken()
		}
	}
}

func (parser *Parser) parsePlainTextBlock() *Node {
	n := Node{Kind: N_PLAINTEXT_BLOCK}

	if parser.curToken.kind != T_L_SQUIRLY {
		panic("Expected { in plaintext block")
	}
	parser.nextToken()

	for {
		switch parser.curToken.kind {
		case T_ILLEGAL:
			panic("Unexpected end in plaintext block")

		case T_R_SQUIRLY:
			return &n

		default:
			n.addChild(parser.parsePlainTextBlockStatement())
			parser.nextToken()
		}
	}
}

func (parser *Parser) parseTextDataBlock() *Node {
	n := Node{Kind: N_TEXTDATA_BLOCK}

	if parser.curToken.kind != T_L_SQUIRLY {
		panic("Expected { in textdata block")
	}
	parser.nextToken()

	for {
		switch parser.curToken.kind {
		case T_ILLEGAL:
			panic("Unexpected end in textdata block")

		case T_R_SQUIRLY:
			return &n

		default:
			n.addChild(parser.parseTextDataBlockStatement())
			parser.nextToken()
		}
	}
}

// STATEMENT

func (parser *Parser) parseColorBlockStatement() *Node {
	switch parser.curToken.kind {
	case T_COM_RED:
		return parser.parseRedCommand()
	case T_COM_GREEN:
		return parser.parseGreenCommand()
	case T_COM_BLUE:
		return parser.parseBlueCommand()
	default:
		panic("Invalid start to button block statement")
	}
}


func (parser *Parser) parsePageBlockStatement() *Node {
	switch parser.curToken.kind {
	case T_COM_BGDRAW:
		return parser.parseBgDrawCommand()
	case T_COM_BGCOLOR:
		return parser.parseBgColorCommand()
	case T_BUTTON:
		return parser.parseButtonDef()
	case T_PLAINTEXT:
		return parser.parsePlainTextDef()
	default:
		panic("Invalid start to page block statement")
	}
}

func (parser *Parser) parseButtonBlockStatement() *Node {
	switch parser.curToken.kind {
	case T_COM_BGCOLOR:
		return parser.parseBgColorCommand()
	case T_COM_X:
		return parser.parseXCommand()
	case T_COM_Y:
		return parser.parseYCommand()
	case T_COM_W:
		return parser.parseWCommand()
	case T_COM_H:
		return parser.parseHCommand()
	case T_TEXTDATA:
		return parser.parseTextData()
	default:
		panic("Invalid start to block statement")
	}
}

func (parser *Parser) parsePlainTextBlockStatement() *Node {
	switch parser.curToken.kind {
	case T_COM_X:
		return parser.parseXCommand()
	case T_COM_Y:
		return parser.parseYCommand()
	case T_TEXTDATA:
		return parser.parseTextData()
	default:
		panic("Invalid start to plaintext statement")
	}
}

func (parser *Parser) parseTextDataBlockStatement() *Node {
	switch parser.curToken.kind {
	case T_COM_TEXT:
		return parser.parseTextCommand()
	case T_COM_TEXTCOLOR:
		return parser.parseTextColorCommand()
	case T_COM_FONTSIZE:
		return parser.parseFontSizeCommand()
	case T_COM_FONT:
		return parser.parseFontCommand()
	case T_COM_TEXTGAP:
		return parser.parseTextGapCommand()
	case T_COM_HORALIGN:
		return parser.parseHorAlignCommand()
	case T_COM_VERALIGN:
		return parser.parseVerAlignCommand()
	default:
		panic("Invalid start to text data statement")
	}
}

// COMMANDS

func (parser *Parser) parseRedCommand() *Node {
	n := Node{Kind: N_STAT_RED}

	if parser.curToken.kind != T_COM_RED {
		panic("Invalid red command")
	}
	n.addChild(&Node{Kind: N_COM_RED, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid red amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseGreenCommand() *Node {
	n := Node{Kind: N_STAT_GREEN}

	if parser.curToken.kind != T_COM_GREEN {
		panic("Invalid green command")
	}
	n.addChild(&Node{Kind: N_COM_GREEN, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid green amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseBlueCommand() *Node {
	n := Node{Kind: N_STAT_BLUE}

	if parser.curToken.kind != T_COM_BLUE {
		panic("Invalid blue command")
	}
	n.addChild(&Node{Kind: N_COM_BLUE, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid blue amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseBgDrawCommand() *Node {
	n := Node{Kind: N_STAT_BGDRAW}

	if parser.curToken.kind != T_COM_BGDRAW {
		panic("Invalid bgdraw command")
	}
	n.addChild(&Node{Kind: N_COM_BGDRAW, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseBgColorCommand() *Node {
	n := Node{Kind: N_STAT_BGCOLOR}

	if parser.curToken.kind != T_COM_BGCOLOR {
		panic("Invalid bgcolor command")
	}
	n.addChild(&Node{Kind: N_COM_BGCOLOR, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Invalid bgcolor name")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseXCommand() *Node {
	n := Node{Kind: N_STAT_X}

	if parser.curToken.kind != T_COM_X {
		panic("Invalid x command")
	}
	n.addChild(&Node{Kind: N_COM_X, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid x amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseYCommand() *Node {
	n := Node{Kind: N_STAT_Y}

	if parser.curToken.kind != T_COM_Y {
		panic("Invalid y command")
	}
	n.addChild(&Node{Kind: N_COM_Y, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid y amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseWCommand() *Node {
	n := Node{Kind: N_STAT_W}

	if parser.curToken.kind != T_COM_W {
		panic("Invalid w command")
	}
	n.addChild(&Node{Kind: N_COM_W, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid w amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseHCommand() *Node {
	n := Node{Kind: N_STAT_H}

	if parser.curToken.kind != T_COM_H {
		panic("Invalid h command")
	}
	n.addChild(&Node{Kind: N_COM_H, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid h amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseTextColorCommand() *Node {
	n := Node{Kind: N_STAT_TEXTCOLOR}

	if parser.curToken.kind != T_COM_TEXTCOLOR {
		panic("Invalid textcolor command")
	}
	n.addChild(&Node{Kind: N_COM_TEXTCOLOR, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Invalid textcolor name")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseTextCommand() *Node {
	n := Node{Kind: N_STAT_TEXT}

	if parser.curToken.kind != T_COM_TEXT {
		panic("Invalid text command")
	}
	n.addChild(&Node{Kind: N_COM_TEXT, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic("Invalid text")
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseFontSizeCommand() *Node {
	n := Node{Kind: N_STAT_FONTSIZE}

	if parser.curToken.kind != T_COM_FONTSIZE {
		panic("Invalid fontsize command")
	}
	n.addChild(&Node{Kind: N_COM_FONTSIZE, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid fontsize amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseFontCommand() *Node {
	n := Node{Kind: N_STAT_FONT}

	if parser.curToken.kind != T_COM_FONT {
		panic("Invalid font command")
	}
	n.addChild(&Node{Kind: N_COM_FONT, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_STRING {
		panic(fmt.Sprintf("Invalid font name %s (%d)", parser.curToken.data, parser.curToken.kind))
	}
	n.addChild(&Node{Kind: N_STRING, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseTextGapCommand() *Node {
	n := Node{Kind: N_STAT_TEXTGAP}

	if parser.curToken.kind != T_COM_TEXTGAP {
		panic("Invalid textgap command")
	}
	n.addChild(&Node{Kind: N_COM_TEXTGAP, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_INT {
		panic("Invalid textgap amount")
	}
	n.addChild(&Node{Kind: N_INT, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseHorAlignCommand() *Node {
	n := Node{Kind: N_STAT_HORALIGN}

	if parser.curToken.kind != T_COM_HORALIGN {
		panic("Invalid horalign command")
	}
	n.addChild(&Node{Kind: N_COM_HORALIGN, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_ALIGN {
		panic("Invalid alignment")
	}
	n.addChild(&Node{Kind: N_ALIGN, Data: parser.curToken.data})

	return &n
}

func (parser *Parser) parseVerAlignCommand() *Node {
	n := Node{Kind: N_STAT_VERALIGN}

	if parser.curToken.kind != T_COM_VERALIGN {
		panic("Invalid veralign command")
	}
	n.addChild(&Node{Kind: N_COM_VERALIGN, Data: parser.curToken.data})
	parser.nextToken()

	if parser.curToken.kind != T_ALIGN {
		panic("Invalid alignment")
	}
	n.addChild(&Node{Kind: N_ALIGN, Data: parser.curToken.data})

	return &n
}
