package main

import (
	"fmt"
	"strconv"
)

func main() {
	testStr := "+1-4-6+9-4-"
	n := Analize(testStr)
	fmt.Println(n)

	fmt.Println(n.Traverse())

}

type Node struct {
	value float32
	sign  float32
	err   error
	nodes []*Node
}

func Analize(s string) *Node {
	node := &Node{value: 0}
	for i := 0; i < len(s); i++ {
		switch s[i : i+1] {
		case "+":
			left := Analize(s[:i])
			left.sign = 1.0
			node.nodes = append(node.nodes, left)
			right := Analize(s[i+1:])
			right.sign = (1.0)
			node.nodes = append(node.nodes, right)
			return node
		case "-":
			left := Analize(s[:i])
			left.sign = 1.0
			node.nodes = append(node.nodes, left)
			right := Analize(s[i+1:])
			right.sign = (-1.0)
			fmt.Println(right.value)
			node.nodes = append(node.nodes, right)
			return node
		}
	}
	if s == "" {
		return node
	}
	parseFloat, err := strconv.ParseFloat(s, 32)
	if err != nil {
		fmt.Println(err)
	}
	node.value = float32(parseFloat)
	return node
}
func (n *Node) Traverse() (result float32) {
	if len(n.nodes) == 0 {
		return n.value * n.sign
	}
	for _, pN := range n.nodes {
		result += pN.Traverse() * pN.sign
	}
	n.value = result * n.sign
	fmt.Println(result)
	return
}
func ReverceSign(s string) string {
	b := []byte(s)
	mp := "-+"
	for i := 0; i < len(b); i++ {
		if s[i:i+1] == "+" {
			b[i] = mp[0]
		}
		if s[i:i+1] == "-" {
			b[i] = mp[1]
		}
	}
	return string(b)
}
