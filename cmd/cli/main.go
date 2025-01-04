package main

import (
	"flag"
	"fmt"

	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/debug"
)

func main() {
	debug.NewDebug()
	debug.LoadFromFile("debug.cfg")

	var file_in string
	flag.StringVar(&file_in, "file_in", "", "input phrase file")

	var file_out string
	flag.StringVar(&file_out, "file_out", "", "output parsing file")

	flag.Parse()

	n := natasha.NewNatasha()
	fmt.Printf("natasha loaded\r\n")
	defer n.Close()
	rrs := relations.InitRelationRule()
	fmt.Printf("rules!\r\n")
	if true {
		err := relations.LoadSentensesNew(n, rrs, file_in, file_out, 0)
		if err != nil {
			fmt.Printf("error %v\r\n", err)
			return
		}
	}
}
