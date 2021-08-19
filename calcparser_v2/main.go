package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	OAG := DefaultOAG()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		Expr := c.QueryParam("Expr")
		node, err := Analize(Expr, OAG)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		str, err := node.Traverse()
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.String(http.StatusOK, fmt.Sprint(str))
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
	value     float32
	substring string
	Aggregate func(x float32, first bool) func(y float32) float32
	nodes     []*Node
}
type OperationAnalizer struct {
	op        string
	Aggregate func(x float32, first bool) func(y float32) float32
}
type OperationAnalizerGroup struct {
	ops []OperationAnalizer
}

//Analize creates graph of expression strings and aggregate functions
func Analize(s string, oag OperationAnalizerGroup) (res *Node, err error) {
	s = strings.TrimSpace(s)
	s = RemBrackets(s)
	var splitResult []string
	n := &Node{
		substring: s,
	}
	for _, ops := range oag.ops {
		splitResult, err = SplitExpr(s, ops.op)
		if err != nil {
			return n, err
		}
		if len(splitResult) > 1 {
			n.Aggregate = ops.Aggregate
			for _, sr := range splitResult {
				newN, err := Analize(sr, oag)
				if err != nil {
					return n, err
				}
				n.nodes = append(n.nodes, newN)
			}
			return n, nil
		}
	}
	return n, nil
}

//Traverse calculates result with Aggregate func
func (n *Node) Traverse() (result float32, err error) {
	if len(n.nodes) == 0 {
		if n.substring == "" {
			return 0, nil
		}
		fl64, err := strconv.ParseFloat(n.substring, 32)
		if err != nil {
			return 0, errors.New("parsing error use +-*/() and [0-9] only")
		}
		n.value = float32(fl64)
		return n.value, err
	}
	for i, pN := range n.nodes {
		if n.Aggregate != nil {
			nodeResult, err := pN.Traverse()
			if err != nil {
				return result, err
			}
			result = n.Aggregate(nodeResult, i == 0)(result)
		}
	}
	n.value = result
	return result, nil
}

//SplitExpr takes string and sub of +/-* returns array of substrings and error
func SplitExpr(s string, sub string) ([]string, error) {
	var pos int

	exprBoundrySigns := "("
	result := make([]string, 0)
	for i := 0; i < len(s); i++ {
		if s[i] == exprBoundrySigns[0] {
			newPos, err := closingBracketPos(s[i:])
			if err != nil {
				return []string{}, err
			}
			i += newPos
			continue
		}
		if s[i] == sub[0] {
			result = append(result, (s[pos:i]))
			pos = i + 1
		}
	}
	result = append(result, (s[pos:]))
	return result, nil
}

//RemBrackets returns string with no () if they wrap original string
func RemBrackets(s string) string {
	if s == "" {
		return s
	}
	if s[0:1] == "(" && s[len(s)-1:] == ")" {
		return s[1 : len(s)-1]
	}
	return s
}
func DefaultOAG() OperationAnalizerGroup {
	return OperationAnalizerGroup{
		ops: []OperationAnalizer{
			{
				op: "+",
				Aggregate: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						return x + y
					}
				},
			},
			{
				op: "-",
				Aggregate: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						if first {
							return y + x
						}
						return -x + y
					}
				},
			},
			{
				op: "*",
				Aggregate: func(x float32, first bool) func(y float32) float32 {
					return func(y float32) float32 {
						if y == 0 {
							y = 1
						}
						return x * y
					}
				},
			},
			{
				op: "/",
				Aggregate: func(x float32, first bool) func(y float32) float32 {
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
}
