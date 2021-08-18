package main

import (
	"errors"
	"fmt"
)

func main() {
	testString := "1+2*2+(3-4)"
	exprBoundrySigns := "+-()"

	pos := 0
	result := make([]string, 0)
	//for i, v := range []byte(testString) {
	for i := 0; i < len(testString); i++ {
		fmt.Println(i)

		if testString[i] == exprBoundrySigns[2] {
			newPos, err := closingBracketPos(testString[i:])
			if err != nil {
				panic(err)
			}
			i += newPos
			fmt.Println(newPos)
			continue
		}
		if testString[i] == exprBoundrySigns[0] || testString[i] == exprBoundrySigns[1] {
			result = append(result, testString[pos:i])
			pos = i + 1
		}
	}
	result = append(result, testString[pos:])
	fmt.Println(result, " l:", len(result))
}

type Expr struct {
	Interator ExpressionInterator
	content   []*Expr
}
type ExpressionInterator struct {
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
