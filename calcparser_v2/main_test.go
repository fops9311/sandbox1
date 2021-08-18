package main

import (
	"testing"
)

func TestRemBrackets(t *testing.T) {
	result := RemBrackets("(fdfdf)")
	if result != "fdfdf" {
		t.Fatal("Want fdfdf got ", result)
	} else {
		t.Log("RemBrackets(\"(fdfdf)\")=", result)
	}

	result = RemBrackets("(fdfdf)gf")
	if result != "(fdfdf)gf" {
		t.Fatal("Want (fdfdf)gf got ", result)
	} else {
		t.Log("RemBrackets(\"(fdfdf)gf\")=", result)
	}

}
func TestSpiltExpr(t *testing.T) {
	got := SplitExpr("2+3", "+")[0]
	if got != "2" {
		t.Fatal("Want 2 got ", got)
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\")[0]=", got)
	}

	got = SplitExpr("2+3", "+")[1]
	if got != "3" {
		t.Fatal("Want 3 got ", got)
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\")[1]=", got)
	}

	len := len(SplitExpr("2+3", "+"))
	if len != 2 {
		t.Fatal("Want 2 got ", len)
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\") len=", len)
	}

	got = SplitExpr("2+3-1", "+")[1]
	if got != "3-1" {
		t.Fatal("Want 3-1 got ", got)
	} else {
		t.Log("SplitExpr(\"2+3-1\", \"+\")[1]=", got)
	}
}
func TestNode(t *testing.T) {
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
	testString := "1+1"
	root := Analize(testString, OAG)
	result := root.Traverse()
	if result != 2 {
		t.Fatal("1+1 should be 2 but resul is ", result)
	}
}
