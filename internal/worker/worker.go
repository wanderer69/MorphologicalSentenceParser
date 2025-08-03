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
	isHaveInit bool
	proc       *process.Process
	n          *natasha.Natasha
	rrs        *relations.RelationRules
	env        *parser.Env
}

type payload struct {
	ClientID string
	Sentence string
}
type payloadOut struct {
	clientID string
	tsris    []*relations.TranslateSentensesResultItem
}

func NewProcessor() *Processor {
	debug.NewDebug()

	env := parser.NewEnv()
	script.MakeRules(env)

	n := natasha.NewNatasha()

	rEnv := script.NewEnvironment()
	rp := script.RelationsParser{}
	rp.Env = rEnv
	env.Struct = &rp
	env.Debug = 0

	proc := &Processor{
		n:   n,
		env: env,
		rrs: rp.Env.RelationRules,
	}
	return proc
}

func (p *Processor) Init() error {
	buffer := ""
	o := print.NewOutput(func(sfmt string, args ...any) {
		s := fmt.Sprintf(sfmt, args...)
		buffer = buffer + s
	})

	script.MakeRules(p.env)

	err := p.n.Init()
	if err != nil {
		return err
	}
	rules_file_name := os.Getenv("RULES_FILE_NAME")

	if len(rules_file_name) == 0 {
		rules_file_name = "./rules.script"
	}
	data, err := os.ReadFile(rules_file_name)
	if err != nil {
		return fmt.Errorf("failed load rules file name %v: %w", rules_file_name, err)
	}

	_, err = p.env.ParseString(string(data), o)
	if err != nil {
		return fmt.Errorf("failed parsing file name %v: %w", rules_file_name, err)
	}

	p.proc = process.NewProcess(p, procFunc)
	p.proc.Run()
	p.isHaveInit = true
	return nil
}

func procFunc(ei interface{}, pli interface{}) (interface{}, error) {
	proc := ei.(*Processor)
	pl := pli.(*payload)

	tsris, err := relations.TranslateText(proc.n, proc.rrs, pl.Sentence, 0)
	if err != nil {
		return nil, err
	}

	plo := &payloadOut{clientID: pl.ClientID, tsris: tsris}
	return plo, nil
}

func (proc *Processor) Send(sentence string) (string, error) {
	if !proc.isHaveInit {
		return "", fmt.Errorf("not initialized")
	}
	return proc.proc.Send(&payload{Sentence: sentence})
}

func (proc *Processor) Check(taskID string) (string, string, []*relations.TranslateSentensesResultItem, error) {
	if !proc.isHaveInit {
		return "", "", nil, fmt.Errorf("not initialized")
	}
	ploi, _, err := proc.proc.Check(taskID)
	if err != nil {
		return "", "error", nil, err
	}
	if ploi != nil {
		plo := ploi.(*payloadOut)
		return plo.clientID, "Ok", plo.tsris, nil
	}
	return "", "", nil, nil
}
