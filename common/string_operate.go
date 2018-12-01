package common

import (
	"utils"
	"strings"
)

func LinkInterface(inte ...interface{}) (s string){
	for _, e := range inte{
		s += utils.ParseString(e)
	}
	return
}

func GenPlaceholder(length int) string {
	var statment strings.Builder
	statment.WriteByte('(')
	for i := 0; i < length -1; i++ {
		statment.WriteString("?, ")
	}
	statment.WriteString("?)")
	return statment.String()
}