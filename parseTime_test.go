package common

import (
	"fmt"
	"testing"
)

func TestParseTime(t *testing.T){
	if s := ParseTime("20121232");  s != ""{
		t.Error("error")
	} else {
		fmt.Println(s)
	}

	if s := ParseTime("20121201");  s == "2012-12-01 00:00:00"{
		fmt.Println(s)
	} else {
		t.Error("error")
	}

	if s := ParseTime("二零一七年九月十八日");  s == "2017-09-18 00:00:00"{
		fmt.Println(s)
	} else {
		t.Error("error")
	}

	if s := ParseTime("2012-04-03日16点31分");  s == "2012-04-03 16:31:00"{
		fmt.Println(s)
	} else {
		t.Error("error")
	}

	if s := ParseTime("2012年04月03日");  s == "2012-04-03 00:00:00"{
		fmt.Println(s)
	} else {
		t.Error("error")
	}

	if s := ParseTime("2012年04月03日 15时23分");  s == "2012-04-03 15:23:00"{
		fmt.Println(s)
	} else {
		t.Error("error")
	}
	//
	//fmt.Println(parseTime("二零一七年九月十八日"))
	//fmt.Println(parseTime("2012-04-03日16点31分"))
	//fmt.Println(parseTime("2012年04月03日"))
	//fmt.Println(parseTime("2012年04月03日 15时23分"))
}
