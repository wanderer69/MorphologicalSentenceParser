package worker

import (
	"github.com/wanderer69/MorphologicalSentenceParser/internal/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/internal/relations"
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
}

func NewProcessor() *Processor {
	n := natasha.NewNatasha("../../scripts/python")
	rrs := relations.InitRelationRule()

	proc := &Processor{
		n:   n,
		rrs: rrs,
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
