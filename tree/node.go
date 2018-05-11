package tree

type NodeList []*Node
type Node struct {
	Children NodeList
	Value    interface{}

}

func NewNode(children []*Node, payload interface{}) *Node {
	return &Node {Children: children, Value: payload }
}

func (node *Node) Flatten() NodeList {
	return NodeList{node}.Flatten()
}

func (nodes NodeList) Flatten() NodeList {
	results := NodeList{}
	for _, node := range nodes {
		results = append(results, node.Children.Flatten()...)
		results = append(results, node)
	}
	return results
}
