package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	testStr := "(2*2*2+1)"
	ops := OperationAnalizerGroup{
		ops: []OperationAnalizer{
			{
				name: "+",
				leftTransform: func(x float32) float32 {
					return x
				},
				rightTransform: func(x float32) float32 {
					return x
				},
				leftStringTransform: func(s string) string {
					return s
				},
				rightStringTransform: func(s string) string {
					return s
				},
			},
			{
				name: "-",
				leftTransform: func(x float32) float32 {
					return x
				},
				rightTransform: func(x float32) float32 {
					return -x
				},
				leftStringTransform: func(s string) string {
					return s
				},
				rightStringTransform: ReverceSign,
			},
		},
		next: &OperationAnalizerGroup{
			ops: []OperationAnalizer{
				{
					name: "*",
					leftTransform: func(x float32) float32 {
						return x
					},
					rightTransform: func(x float32) float32 {
						return x
					},
					leftStringTransform: func(s string) string {
						return s
					},
					rightStringTransform: func(s string) string {
						return s
					},
				},
				{
					name: "/",
					leftTransform: func(x float32) float32 {
						return x
					},
					rightTransform: func(x float32) float32 {
						return 1 / x
					},
					leftStringTransform: func(s string) string {
						return s
					},
					rightStringTransform: func(s string) string {
						return s
					},
				},
			},
		},
	}
	n := Analize(testStr, ops)
	fmt.Println(n)

	fmt.Println(n.Traverse())
	fmt.Println(n.value)

}

type Node struct {
	value             float32
	transform         func(x float32) float32
	traverseInterator func(x float32) func(y float32) float32
	err               error
	nodes             []*Node
}
type OperationAnalizer struct {
	name                 string
	leftTransform        func(x float32) float32
	rightTransform       func(x float32) float32
	leftStringTransform  func(s string) string
	rightStringTransform func(s string) string
	next                 *OperationAnalizer
}
type OperationAnalizerGroup struct {
	ops  []OperationAnalizer
	next *OperationAnalizerGroup
}

func Analize(s string, oag OperationAnalizerGroup) *Node {
	node := &Node{
		value: 0,
		transform: func(x float32) float32 {
			return x
		},
		traverseInterator: func(x float32) func(y float32) float32 {
			return func(y float32) float32 {
				return x + y
			}
		},
	}
	i := 0
	for ; i < len(s); i++ {
		for _, op := range oag.ops {
			if s[i:i+1] == "(" {
				newS := s[i:]
				cbp, _ := closingBracketPos(newS)
				node.nodes = append(node.nodes, Analize(newS[i+1:cbp], oag))
				fmt.Println("i:", i)
				i += cbp
				fmt.Println("i:", i)
			}

			if s[i:i+1] == op.name {

				left := Analize(op.leftStringTransform(s[:i]), oag)
				left.transform = op.leftTransform
				node.nodes = append(node.nodes, left)

				fmt.Println(op.rightStringTransform(s[i+1:]))

				right := Analize(op.rightStringTransform((s[i+1:])), oag)
				right.transform = op.rightTransform
				node.nodes = append(node.nodes, right)

				return node
			}
		}
	}
	if node.transform == nil {
		node.transform = func(x float32) float32 {
			return x
		}
	}
	if s == "" {
		return node
	}
	if oag.next != nil {
		return Analize(s, *oag.next)
	}
	parseFloat, err := strconv.ParseFloat(s, 32)
	if err != nil {
		fmt.Println(err)
	}
	node.value = float32(parseFloat)
	return node
}
func (n *Node) Traverse() (result float32) {
	fmt.Println(n)

	if len(n.nodes) == 0 {
		return n.transform(n.value)
	}
	for _, pN := range n.nodes {
		result = n.traverseInterator(n.transform(pN.Traverse()))(result)
	}
	n.value = n.transform(result)
	//fmt.Println(result)
	return result
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
func closingBracketPos(s string) (pos int, err error) {
	var counter int
	if s[:1] != "(" {
		return 0, errors.New("Must start with a \"(\"")
	}
	for i := 0; i < len(s); i++ {
		if s[i:i+1] == "(" {
			counter++
		}
		if s[i:i+1] == ")" {
			counter--
		}
		if counter == 0 {
			return i, nil
		}
	}
	return pos, errors.New("No closing bracket found")
}
