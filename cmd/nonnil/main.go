// Command nonnil runs the nonnil analyzer standalone or as a `go vet` tool:
//
//	go install github.com/go-vet-analyzers/nonnil/cmd/nonnil@latest
//	go vet -vettool=$(which nonnil) ./...
package main

import (
	"github.com/go-vet-analyzers/nonnil"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(nonnil.Analyzer) }
