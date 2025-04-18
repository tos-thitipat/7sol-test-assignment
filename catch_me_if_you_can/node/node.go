package node

import (
	"errors"
	"fmt"
	"log"
)

type nodeLink struct {
	val  int
	next *nodeLink
	prev *nodeLink
	code *string
}

func InitNodeByInput(input string) (*nodeLink, error) {
	var decodeMapping = map[string]int{
		"L": -1,
		"R": 1,
		"=": 0,
	}
	var headNode *nodeLink
	var currentNode *nodeLink
	var min int = 0

	for _, codeRune := range input {
		code := string(codeRune)
		decodedVal, isContain := decodeMapping[code]
		if !isContain {
			return nil, errors.New("your code is invalid")
		}
		if headNode == nil {
			headNode = &nodeLink{
				val: 0,
			}
			currentNode = headNode
		}
		tmpNode := &nodeLink{
			val:  currentNode.val + decodedVal,
			prev: currentNode,
			code: &code,
		}
		currentNode.next = tmpNode
		currentNode = tmpNode
		if currentNode.val < min {
			min = currentNode.val
		}
	}

	if min < 0 {
		headNode.offsetNodeVal(min)
	}

	return headNode, nil
}

func (node *nodeLink) offsetNodeVal(offsetVal int) {
	if offsetVal == 0 {
		return
	}

	if node != nil {
		node.val -= offsetVal
		node.next.offsetNodeVal(offsetVal)
	}
}

func (node *nodeLink) Decoded() string {
	if node == nil {
		return ""
	}

	for {
		nodeLinkDecoded := node.GetNodeVal()
		optimizedNode := node.optimizeValue()
		optimizedNodeDecoded := optimizedNode.GetNodeVal()
		if nodeLinkDecoded == optimizedNodeDecoded {
			break
		}
	}

	return node.GetNodeVal()
}

func (node *nodeLink) GetNodeVal() string {
	if node == nil {
		return ""
	}
	return fmt.Sprintf("%d", node.val) + node.next.GetNodeVal()
}

func (node *nodeLink) GetNodeCode() string {
	if node == nil {
		return ""
	}

	if node.code == nil {
		return node.next.GetNodeCode()
	}

	return *node.code + node.next.GetNodeCode()
}

func (node *nodeLink) optimizeValue() *nodeLink {
	if node == nil {
		return node
	}

	var lastTmpNode *nodeLink

	if node.prev == nil {
		return node.next.optimizeValue()
	}

	tempNode := node.minimizeVal()
	lastTmpNode = &tempNode
	for lastTmpNode.next != nil && *lastTmpNode.next.code == "=" {
		lastTmpNode = lastTmpNode.next
	}

	if tempNode.isNodeLinkLogicValid() && lastTmpNode.isNodeLinkLogicValid() {
		node.prev.next = &tempNode

		if node.next != nil {
			for node.next != nil && *node.next.code != "R" {
				node = node.next
			}
			if node.next != nil {
				node.next.prev = lastTmpNode
			}
		}
	}

	if node.next != nil {
		node.next.optimizeValue()
	}

	var headNode *nodeLink
	for node.prev != nil {
		node = node.prev
	}
	headNode = node

	if headNode.next != nil && headNode.next.code != nil && *headNode.next.code == "=" {
		headNode.val = headNode.next.val
	}

	return headNode
}

func (node *nodeLink) minimizeVal() nodeLink {

	if node.val == 0 {
		return *node
	}

	tempNode := nodeLink{
		val:  node.val - 1,
		next: node.next,
		prev: node.prev,
		code: node.code,
	}

	if node.next != nil && *node.next.code != "R" {
		newTempNode := node.next.minimizeVal()
		tempNode.next = &newTempNode
		tempNode.next.prev = &tempNode
	}

	return tempNode
}

func (node *nodeLink) isNodeLinkLogicValid() bool {
	if node == nil {
		log.Println("node is nil")
		return true
	}

	if node.code == nil {
		log.Println("code is nil")
		return false
	}

	switch *node.code {
	case "L":
		if node.val < node.prev.val {
			if node.next == nil {
				return true
			} else {
				return node.next.isNodeLinkLogicValid()
			}
		}
	case "R":
		if node.val > node.prev.val {
			if node.next == nil {
				return true
			}
			return node.next.isNodeLinkLogicValid()
		}
	case "=":
		if node.prev.code == nil {
			if node.next == nil {
				return true
			}
			return node.next.isNodeLinkLogicValid()
		}
		if node.val == node.prev.val {
			if node.next == nil {
				return true
			}
			return node.next.isNodeLinkLogicValid()
		}
	default:
		log.Printf("invalid code: %s\n", *node.code)
		return false
	}

	return false
}
