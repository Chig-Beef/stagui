package stagui

import (
	"fmt"
	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UI struct {
	vis    VisualHandler
	Pages  map[string]*Page
	colors map[string]color.RGBA
}

func (ui *UI) Init(vis VisualHandler) {
	ui.Pages = map[string]*Page{}
	ui.colors = map[string]color.RGBA{}
	ui.vis = vis
	ui.loadUI()
}

func (ui *UI) loadUI() {
	parser := Parser{}
	node := parser.ParseFromFile("main.stagui")
	ui.loadProgram(node)
}

func (ui *UI) loadProgram(node *Node) {
	// fmt.Println("load program", node)
	for _, child := range node.Children {
		switch child.Kind {
		case N_PAGE_DEF:
			ui.loadPageDef(child)
		case N_COLOR_DEF:
			ui.loadColorDef(child)
		default:
			panic("Invalid definition in program")
		}
	}
}

// DEFS

func (ui *UI) loadPageDef(node *Node) {
	page := Page{}

	// Get name without quotes
	name := node.Children[1].Data
	name = name[1:len(name)-1]

	ui.loadPageBlock(node.Children[2], &page)

	page.Name = name
	ui.Pages[name] = &page
}

func (ui *UI) loadButtonDef(node *Node, page *Page) {
	button := Button{}

	// Get name without quotes
	name := node.Children[1].Data
	name = name[1:len(name)-1]
	button.Name = name

	ui.loadButtonBlock(node.Children[2], &button)

	page.Buttons = append(page.Buttons, &button)
}

func (ui *UI) loadPanelDef(node *Node, page *Page) {
	panel := Panel{}

	// Get name without quotes
	name := node.Children[1].Data
	name = name[1:len(name)-1]
	panel.Name = name

	ui.loadPanelBlock(node.Children[2], &panel)

	page.Panels = append(page.Panels, &panel)
}

func (ui *UI) loadPlainTextDef(node *Node, page *Page) {
	plainText := PlainText{}

	// Get name without quotes
	name := node.Children[1].Data
	name = name[1:len(name)-1]
	plainText.Name = name

	ui.loadPlainTextBlock(node.Children[2], &plainText)

	page.PlainTexts = append(page.PlainTexts, &plainText)
}

func (ui *UI) loadTextDataDef(node *Node) TextData {
	textData := TextData{}
	ui.loadTextDataBlock(node.Children[1], &textData)
	return textData
}

func (ui *UI) loadColorDef(node *Node) {
	clr := color.RGBA{}

	// Get name without quotes
	name := node.Children[1].Data
	name = name[1:len(name)-1]

	ui.loadColorBlock(node.Children[2], &clr)

	ui.colors[name] = clr
}

// BLOCKS

func (ui *UI) loadPageBlock(node *Node, page *Page) {
	for _, command := range node.Children {
		ui.loadPageBlockStatement(command, page)
	}
}

func (ui *UI) loadButtonBlock(node *Node, button *Button) {
	for _, command := range node.Children {
		ui.loadButtonBlockStatement(command, button)
	}
}

func (ui *UI) loadPanelBlock(node *Node, panel *Panel) {
	for _, command := range node.Children {
		ui.loadPanelBlockStatement(command, panel)
	}
}

func (ui *UI) loadPlainTextBlock(node *Node, plainText *PlainText) {
	for _, command := range node.Children {
		ui.loadPlainTextBlockStatement(command, plainText)
	}
}

func (ui *UI) loadTextDataBlock(node *Node, textData *TextData) {
	for _, command := range node.Children {
		ui.loadTextDataBlockStatement(command, textData)
	}
}

func (ui *UI) loadColorBlock(node *Node, clr *color.RGBA) {
	for _, command := range node.Children {
		ui.loadColorBlockStatement(command, clr)
	}
}

// STATEMENTS

func (ui *UI) loadPageBlockStatement(node *Node, page *Page) {
	switch node.Kind {
	case N_STAT_BGDRAW:
		ui.loadBgDrawCommand(node, page)
	case N_STAT_BGCOLOR:
		page.BgColor = ui.loadBgColorCommand(node)
	case N_BUTTON_DEF:
		ui.loadButtonDef(node, page)
	case N_PLAINTEXT_DEF:
		ui.loadPlainTextDef(node, page)
	default:
		panic(fmt.Sprintf("Invalid start to page block statement %s (%s)", node.Data, node.Kind.String()))
	}
}

func (ui *UI) loadButtonBlockStatement(node *Node, button *Button) {
	switch node.Kind {
	case N_STAT_BGCOLOR:
		button.Bg.Color = ui.loadBgColorCommand(node)
	case N_STAT_X:
		button.X = ui.loadXCommand(node)
	case N_STAT_Y:
		button.Y = ui.loadYCommand(node)
	case N_STAT_W:
		button.W = ui.loadWCommand(node)
	case N_STAT_H:
		button.H = ui.loadHCommand(node)
	case N_TEXTDATA_DEF:
		button.Text = ui.loadTextDataDef(node)
	default:
		panic(fmt.Sprintf("Invalid start to button statement %s (%s)", node.Data, node.Kind.String()))
	}
}

func (ui *UI) loadPanelBlockStatement(node *Node, panel *Panel) {
	switch node.Kind {
	case N_STAT_BGCOLOR:
		panel.Bg.Color = ui.loadBgColorCommand(node)
	case N_STAT_X:
		panel.X = ui.loadXCommand(node)
	case N_STAT_Y:
		panel.Y = ui.loadYCommand(node)
	case N_STAT_W:
		panel.W = ui.loadWCommand(node)
	case N_STAT_H:
		panel.H = ui.loadHCommand(node)
	default:
		panic(fmt.Sprintf("Invalid start to panel statement %s (%s)", node.Data, node.Kind.String()))
	}
}

func (ui *UI) loadPlainTextBlockStatement(node *Node, plainText *PlainText) {
	switch node.Kind {
	case N_STAT_X:
		plainText.X = ui.loadXCommand(node)
	case N_STAT_Y:
		plainText.Y = ui.loadYCommand(node)
	case N_TEXTDATA_DEF:
		plainText.Text = ui.loadTextDataDef(node)
	default:
		panic(fmt.Sprintf("Invalid start to plaintext statement %s (%s)", node.Data, node.Kind.String()))
	}
}

func (ui *UI) loadTextDataBlockStatement(node *Node, textData *TextData) {
	switch node.Kind {
	case N_STAT_TEXT:
		ui.loadTextCommand(node, textData)
	case N_STAT_TEXTCOLOR:
		ui.loadTextColorCommand(node, textData)
	case N_STAT_FONTSIZE:
		ui.loadFontSizeCommand(node, textData)
	case N_STAT_FONT:
		ui.loadFontCommand(node, textData)
	case N_STAT_TEXTGAP:
		ui.loadTextGapCommand(node, textData)
	case N_STAT_HORALIGN:
		ui.loadHorAlignCommand(node, textData)
	case N_STAT_VERALIGN:
		ui.loadVerAlignCommand(node, textData)
	default:
		panic(fmt.Sprintf("Invalid start to textdata statement %s (%s)", node.Data, node.Kind.String()))
	}
}

func (ui *UI) loadColorBlockStatement(node *Node, clr *color.RGBA) {
	switch node.Kind {
	case N_STAT_RED:
		ui.loadRedCommand(node, clr)
	case N_STAT_GREEN:
		ui.loadGreenCommand(node, clr)
	case N_STAT_BLUE:
		ui.loadBlueCommand(node, clr)
	default:
		panic("Invalid start to color block statement")
	}
}

// COMMANDS

func (ui *UI) loadBgDrawCommand(node *Node, page *Page) {
	page.BgDraw = true
}

func (ui *UI) loadBgColorCommand(node *Node) color.RGBA {
	colorName := node.Children[1].Data
	colorName = colorName[1:len(colorName)-1]
	clr := ui.colors[colorName]
	return clr
}

func (ui *UI) loadHorAlignCommand(node *Node, textData *TextData) {
	strAlign := node.Children[1].Data
	align := text.AlignStart
	switch strAlign{
	case "START":
		align = text.AlignStart
	case "CENTER":
		align = text.AlignCenter
	case "END":
		align = text.AlignEnd
	}
	textData.HorAlign = align
}

func (ui *UI) loadVerAlignCommand(node *Node, textData *TextData) {
	strAlign := node.Children[1].Data
	align := text.AlignStart
	switch strAlign{
	case "START":
		align = text.AlignStart
	case "CENTER":
		align = text.AlignCenter
	case "END":
		align = text.AlignEnd
	}
	textData.VerAlign = align
}

func (ui *UI) loadTextGapCommand(node *Node, textData *TextData) {
	textGap, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	textData.Size = textGap
}

func (ui *UI) loadTextColorCommand(node *Node, textData *TextData) {
	colorName := node.Children[1].Data
	colorName = colorName[1:len(colorName)-1]
	clr := ui.colors[colorName]
	textData.Color = clr
}

func (ui *UI) loadTextCommand(node *Node, textData *TextData) {
	text := node.Children[1].Data
	text = text[1:len(text)-1]
	textData.Text = []string{text}
}

func (ui *UI) loadFontCommand(node *Node, textData *TextData) {
	fontName := node.Children[1].Data
	fontName = fontName[1:len(fontName)-1]
	textData.Font = ui.vis.GetFont(fontName)
}

func (ui *UI) loadFontSizeCommand(node *Node, textData *TextData) {
	fontSize, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	textData.Size = fontSize
}

func (ui *UI) loadXCommand(node *Node) float64 {
	x, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	return x
}

func (ui *UI) loadYCommand(node *Node) float64 {
	y, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	return y
}

func (ui *UI) loadWCommand(node *Node) float64 {
	w, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	return w
}

func (ui *UI) loadHCommand(node *Node) float64 {
	h, err := strconv.ParseFloat(node.Children[1].Data, 64)
	if err != nil {
		panic(err)
	}
	return h
}

func (ui *UI) loadRedCommand(node *Node, clr *color.RGBA) {
	red, err := strconv.Atoi(node.Children[1].Data)
	if err != nil {
		panic(err)
	}
	clr.R = byte(red)
}

func (ui *UI) loadGreenCommand(node *Node, clr *color.RGBA) {
	green, err := strconv.Atoi(node.Children[1].Data)
	if err != nil {
		panic(err)
	}
	clr.G = byte(green)
}

func (ui *UI) loadBlueCommand(node *Node, clr *color.RGBA) {
	blue, err := strconv.Atoi(node.Children[1].Data)
	if err != nil {
		panic(err)
	}
	clr.B = byte(blue)
}
