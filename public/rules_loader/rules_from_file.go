package rulesloader

import (
	"fmt"
	"os"

	"github.com/wanderer69/MorphologicalSentenceParser/internal/script"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"

	"github.com/wanderer69/debug"
	"github.com/wanderer69/tools/parser/parser"
	"github.com/wanderer69/tools/parser/print"
)

func LoadRulesFromFile(fileName string) (*relations.RelationRules, error) {
	debug.NewDebug()

	env := parser.NewEnv()
	buffer := ""
	o := print.NewOutput(func(sfmt string, args ...any) {
		s := fmt.Sprintf(sfmt, args...)
		buffer = buffer + s
	})

	script.MakeRules(env)

	if len(fileName) == 0 {
		fileName = "./rules.script"
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	rEnv := script.NewEnvironment()
	rp := script.RelationsParser{}
	rp.Env = rEnv
	env.Struct = &rp
	env.Debug = 0

	_, err = env.ParseString(string(data), o)
	if err != nil {
		return nil, err
	}
	return rp.Env.RelationRules, nil
}
