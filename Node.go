package stagui

type Node struct {
	Kind NodeCode
	Data string
	Children []*Node
}

func (node *Node) addChild(child *Node) {
	node.Children = append(node.Children, child)
}
