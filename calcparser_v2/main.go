package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	OAG := OperationAnalizerGroup{
		ops: []OperationAnalizer{
			{
				name: "+",
				traverseInterator: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						return x + y
					}
				},
			},
			{
				name: "-",
				traverseInterator: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						if first {
							return y + x
						}
						return -x + y
					}
				},
			},
			{
				name: "*",
				traverseInterator: func(x float32, first bool) func(y float32) float32 {
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
				traverseInterator: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						if first {
							return x
						}
						return y / x
					}
				},
			},
		},
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		Expr := c.QueryParam("Expr")
		str := fmt.Sprint((Analize(Expr, OAG).Traverse()))
		fmt.Println(Expr)
		fmt.Println(str)
		return c.String(http.StatusOK, str)
	})
	e.Logger.Fatal(e.Start(":8000"))

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

type Node struct {
	value             float32
	substring         string
	traverseInterator func(x float32, first bool) func(y float32) float32
	err               error
	nodes             []*Node
}
type OperationAnalizer struct {
	name              string
	traverseInterator func(x float32, first bool) func(y float32) float32
}
type OperationAnalizerGroup struct {
	ops []OperationAnalizer
}

func Analize(s string, oag OperationAnalizerGroup) *Node {
	var splitResult []string
	n := &Node{
		substring: s,
	}
	for _, ops := range oag.ops {
		splitResult = SplitExpr(s, ops.name)
		if len(splitResult) > 1 {
			n.traverseInterator = ops.traverseInterator
			for _, sr := range splitResult {
				n.nodes = append(n.nodes, Analize(sr, oag))
			}
			return n
		}
	}
	return n
}

func (n *Node) Traverse() (result float32) {
	if len(n.nodes) == 0 {
		fl64, _ := strconv.ParseFloat(n.substring, 32)
		n.value = float32(fl64)
		return n.value
	}
	for i, pN := range n.nodes {
		if n.traverseInterator != nil {
			result = n.traverseInterator(pN.Traverse(), i == 0)(result)
		}
	}
	n.value = result
	return result
}

//SplitExpr takes string and sub of +/-*
func SplitExpr(s string, sub string) []string {
	var pos int
	exprBoundrySigns := "("
	result := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == exprBoundrySigns[0] {
			newPos, err := closingBracketPos(s[i:])
			if err != nil {
				panic(err)
			}
			i += newPos
			continue
		}
		if s[i] == sub[0] {
			result = append(result, RemBrackets(s[pos:i]))
			pos = i + 1
		}
	}
	result = append(result, RemBrackets(s[pos:]))
	return result
}
func RemBrackets(s string) string {
	if s == "" {
		return s
	}
	if s[0:1] == "(" && s[len(s)-1:] == ")" {
		return s[1 : len(s)-1]
	}
	return s
}
