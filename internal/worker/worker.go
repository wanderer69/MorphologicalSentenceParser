package worker

import (
	"fmt"
	"os"

	"github.com/wanderer69/MorphologicalSentenceParser/internal/script"
	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/debug"
	"github.com/wanderer69/tools/parser/parser"
	"github.com/wanderer69/tools/parser/print"
	process "github.com/wanderer69/tools/worker"
)

type Processor struct {
	proc *process.Process
	n    *natasha.Natasha
	rrs  *relations.RelationRules
}

type payload struct {
	Sentence string
}
type payloadOut struct {
	tsris []*relations.TranslateSentensesResultItem
	// fileRules string
}

func NewProcessor() *Processor {
	debug.NewDebug()

	env := parser.NewEnv()
	buffer := ""
	o := print.NewOutput(func(sfmt string, args ...any) {
		s := fmt.Sprintf(sfmt, args...)
		buffer = buffer + s
	})

	script.MakeRules(env)

	n := natasha.NewNatasha("../../scripts/python")
	//rrs := relations.InitRelationRule()
	rules_file_name := os.Getenv("RULES_FILE_NAME")

	if len(rules_file_name) == 0 {
		rules_file_name = "./rules.script"
	}
	data, err := os.ReadFile(rules_file_name)
	if err != nil {
		panic(fmt.Errorf("failed load rules file name %v: %w", rules_file_name, err))
	}

	rEnv := script.NewEnvironment()
	rp := script.RelationsParser{}
	rp.Env = rEnv
	env.Struct = &rp
	env.Debug = 0

	_, err = env.ParseString(string(data), o)
	if err != nil {
		panic(fmt.Errorf("failed parsing file name %v: %w", rules_file_name, err))
	}

	proc := &Processor{
		n:   n,
		rrs: rp.Env.RelationRules,
	}
	proc.proc = process.NewProcess(proc, procFunc)
	return proc
}

func procFunc(ei interface{}, pli interface{}) (interface{}, error) {
	proc := ei.(*Processor)
	pl := pli.(*payload)

	//defer n.Close()
	tsris, err := relations.TranslateSentense(proc.n, proc.rrs, pl.Sentence, 0)
	if err != nil {
		return nil, err
	}

	plo := &payloadOut{tsris: tsris}
	return plo, nil
}

func (proc *Processor) Send(sentence string) (string, error) {
	return proc.proc.Send(&payload{Sentence: sentence})
}

func (proc *Processor) Check(taskID string) ([]*relations.TranslateSentensesResultItem, error) {
	plo, err := proc.proc.Check(taskID)
	if err != nil {
		return nil, err
	}
	if plo != nil {
		return plo.(*payloadOut).tsris, nil
	}
	return nil, nil
}
