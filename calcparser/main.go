package main

import (
	"errors"
	"fmt"
	"strconv"
)

func main() {
	testStr := "1+(5+5)*(+1+2)+1+1"
	t := "1+1+1"
	fmt.Println(priorityExpressionBracketing(t))
	testStr = addPlusAfterBracket(testStr)
	fmt.Println(testStr)
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
				traverseInterator: func(x float32) func(y float32) float32 {
					return func(y float32) float32 {
						return x + y
					}
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
				traverseInterator: func(x float32) func(y float32) float32 {
					return func(y float32) float32 {
						return x + y
					}
				},
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
					traverseInterator: func(x float32) func(y float32) float32 {
						return func(y float32) float32 {
							if y == 0 {
								y = 1
							}
							return x * y
						}
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
					traverseInterator: func(x float32) func(y float32) float32 {
						return func(y float32) float32 {
							if y == 0 {
								y = 1
							}
							return x * y
						}
					},
				},
			},
		},
	}
	n := Analize(testStr, ops)

	fmt.Println("!!!N:", n)

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
	traverseInterator    func(x float32) func(y float32) float32
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
	}
	i := 0
	for ; i < len(s); i++ {
		for _, op := range oag.ops {
			var cbp int
			if s[i:i+1] == "(" {
				var traverseInter func(x float32) func(y float32) float32
				newS := s[i:]
				cbp, _ = closingBracketPos(newS)
				traverseInter = func(x float32) func(y float32) float32 {
					return func(y float32) float32 {
						return x + y
					}
				}
				if cbp+2 < len(newS) {
					if newS[cbp+1:cbp+2] == "*" || newS[cbp+1:cbp+2] == "/" {
						traverseInter = func(x float32) func(y float32) float32 {
							return func(y float32) float32 {
								if y == 0 {
									y = 1
								}
								return x * y
							}
						}

					}
				}
				left := Analize(newS[i+1:cbp], oag)
				right := Analize(newS[cbp+1:], oag)
				//left.traverseInterator = traverseInter
				//right.traverseInterator = traverseInter
				node.traverseInterator = traverseInter
				node.nodes = append(node.nodes, left)
				node.nodes = append(node.nodes, right)
				fmt.Println("i:", i)
				i += cbp
				fmt.Println("i:", i)
				return node
			}

			if s[i:i+1] == op.name {
				node.traverseInterator = op.traverseInterator
				left := Analize(op.leftStringTransform(s[:i]), oag)
				left.transform = op.leftTransform
				//left.traverseInterator = op.traverseInterator
				node.nodes = append(node.nodes, left)

				fmt.Println(op.rightStringTransform(s[i+1:]))

				right := Analize(op.rightStringTransform((s[i+1:])), oag)
				right.transform = op.rightTransform
				//right.traverseInterator = op.traverseInterator
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
		if s[i:i+1] == "(" {
			newPos, _ := closingBracketPos(s[i:])
			i += newPos
		}
		if s[i:i+1] == "+" {
			b[i] = mp[0]
		}
		if s[i:i+1] == "-" {
			b[i] = mp[1]
		}
	}
	return string(b)
}
func addPlusAfterBracket(s string) string {
	result := make([]byte, 0)
	b := []byte(s)
	search := "(+-"
	for i, v := range b {
		result = append(result, v)
		if v == search[0] &&
			(b[i+1] != search[1] || b[i+1] != search[2]) {
			result = append(result, search[1])
		}
	}
	return string(result)
}
func closingBracketPos(s string) (pos int, err error) {
	var counter int
	if s[:1] != "(" {
		return 0, errors.New("must start with a \"(\"")
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
	return pos, errors.New("no closing bracket found")
}
func priorityExpressionBracketing(s string) string {
	result := make([]byte, 0)
	b := []byte(s)
	search := "(+-)"
	flip := false
	for i, v := range b {
		result = append(result, v)
		if s[i:i+1] == "(" {
			newPos, _ := closingBracketPos(s[i:])
			i += newPos
		}
		if v == search[1] || v == search[2] {
			if !flip {
				result = append(result, search[0])
			}
			if flip {
				result = append(result, search[3])
			}
			flip = !flip
		}

	}
	if flip {
		result = append(result, search[3])
	}
	return string(result)
}
