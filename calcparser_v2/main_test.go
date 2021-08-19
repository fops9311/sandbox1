package main

import (
	"fmt"
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
	got, err := SplitExpr("2+3", "+")
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if got[0] != "2" {
		t.Fatal("Want 2 got ", got[0])
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\")[0]=", got[0])
	}

	got, err = SplitExpr("2+3", "+")
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if got[1] != "3" {
		t.Fatal("Want 3 got ", got[1])
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\")[1]=", got[1])
	}

	got, err = SplitExpr("2+3", "+")
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if len(got) != 2 {
		t.Fatal("Want 2 got ", len(got))
	} else {
		t.Log("SplitExpr(\"2+3\", \"+\") len=", len(got))
	}

	got, err = SplitExpr("2+3-1", "+")
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if got[1] != "3-1" {
		t.Fatal("Want 3-1 got ", got[1])
	} else {
		t.Log("SplitExpr(\"2+3-1\", \"+\")[1]=", got[1])
	}
}
func TestNode(t *testing.T) {
	OAG := DefaultOAG()

	testString := "1+2"
	root, err := Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err := root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 3 {
		t.Fatal("1+2 should be 3 but result is ", result)
	}

	testString = "2*2"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 4 {
		t.Fatal("2*2 should be 4 but result is ", result)
	} else {
		t.Log("2*2 result is ", result)
	}

	testString = "(2*2)"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 4 {
		t.Fatal("(2*2) should be 4 but result is ", result)
	} else {
		t.Log("(2*2) result is ", result)
	}

	testString = "(2/2)"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 1 {
		t.Fatal("(2/2) should be 1 but result is ", result)
	} else {
		t.Log("(2/2) result is ", result)
	}

	testString = "(4/2)"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 2 {
		t.Fatal("(4/2) should be 2 but result is ", result)
	} else {
		t.Log("(4/2) result is ", result)
	}

	testString = "(1/2)"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if result != 0.5 {
		t.Fatal("(1/2) should be 0.5 but result is ", result)
	} else {
		t.Log("(1/2) result is ", result)
	}

	testString = "(1/0)"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if fmt.Sprint(result) != "+Inf" {
		t.Fatal("(1/0) should be +Inf but result is ", result)
	} else {
		t.Log("(1/0) result is ", result)
	}

	testString = "(1/0)+5"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if fmt.Sprint(result) != "+Inf" {
		t.Fatal("(1/0)+5 should be +Inf but result is ", result)
	} else {
		t.Log("(1/0)+5 result is ", result)
	}

	testString = "(1/0)*0"
	root, err = Analize(testString, OAG)
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	result, err = root.Traverse()
	if err != nil {
		t.Fatal("Got err ", err.Error())
	}
	if fmt.Sprint(result) != "NaN" {
		t.Fatal("(1/0)*0 should be NaN but result is ", result)
	} else {
		t.Log("(1/0)*0 result is ", result)
	}
}
