package parser

import (
	"fmt"
	"strings"
)

var traceLevel = 0

func padIdent() string {
	return strings.Repeat("\t", traceLevel-1)
}

func tracePrint(fnName string) {
	fmt.Printf("%s%s\n", padIdent(), fnName)
}

func incIdent() {
	traceLevel += 1
}

func decIdent() {
	traceLevel -= 1
}

func trace(message string) string {
	incIdent()
	tracePrint("BEGIN " + message)
	return message
}

func untrace(message string) {
	tracePrint("END " + message)
	decIdent()
}
