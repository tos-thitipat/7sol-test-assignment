package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type treeNode struct {
	left  *treeNode
	right *treeNode
	val   int
}

func main() {
	input, err := getInputFromFile("hard.json")
	if err != nil {
		log.Println("Error getting input from file:", err)
		return
	}
	nodes := initNode(input)
	connectedNodes := connectNodes(nodes)
	mostValuablePath := findMostValuablePath(connectedNodes)
	fmt.Println("Most Valuable Path:", mostValuablePath)
	var result int
	for _, val := range mostValuablePath {
		result += val
	}
	fmt.Println("Result:", result)

}

func getInputFromFile(filePath string) ([][]int, error) {
	var result [][]int
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func initNode(input [][]int) [][]treeNode {
	var result [][]treeNode
	for _, row := range input {
		nodeRow := []treeNode{}
		for _, val := range row {
			newNode := treeNode{
				val: val,
			}
			nodeRow = append(nodeRow, newNode)
		}
		result = append(result, nodeRow)
	}
	return result
}

func connectNodes(nodes [][]treeNode) *treeNode {
	var headNode *treeNode
	var currentNode *treeNode

	for rowIdx, row := range nodes {
		for nodeIdx := range row {
			if rowIdx == 0 {
				headNode = &row[rowIdx]
				currentNode = headNode
			} else {
				currentNode = &nodes[rowIdx][nodeIdx]
			}
			if rowIdx < len(nodes)-1 {
				currentNode.left = &nodes[rowIdx+1][nodeIdx]
				if nodeIdx < len(nodes[rowIdx+1])-1 {
					currentNode.right = &nodes[rowIdx+1][nodeIdx+1]
				} else {
					currentNode.right = nil
				}
			}
		}
	}

	return headNode
}

func findMostValuablePath(root *treeNode) []int {
	var maxSum int
	var maxPath []int
	var currentPath []int

	var findPath func(node *treeNode, currentSum int)
	findPath = func(node *treeNode, currentSum int) {
		if node == nil {
			return
		}

		currentSum += node.val
		currentPath = append(currentPath, node.val)

		if node.left == nil && node.right == nil {
			if currentSum > maxSum {
				maxSum = currentSum
				maxPath = make([]int, len(currentPath))
				copy(maxPath, currentPath)
			}
		} else {
			findPath(node.left, currentSum)
			findPath(node.right, currentSum)
		}
		currentPath = currentPath[:len(currentPath)-1]
	}

	findPath(root, 0)
	return maxPath
}
