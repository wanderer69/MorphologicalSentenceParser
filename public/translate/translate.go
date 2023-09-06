package translate

import (
	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/debug"
)

type Translate struct {
	n   *natasha.Natasha
	rrs *relations.RelationRules
}

func NewTranslate() *Translate {
	debug.NewDebug()
	//debug.LoadFromFile("../../cmd/cli/debug.cfg")

	n := natasha.NewNatasha("../../scripts/python")
	rrs := relations.InitRelationRule()
	return &Translate{
		n:   n,
		rrs: rrs,
	}
}

func (t *Translate) SetRRS(rrs *relations.RelationRules) {
	t.rrs = rrs
}

func (t *Translate) Translate(sentence string) ([]*relations.TranslateSentensesResultItem, error) {
	tsris, err := relations.TranslateSentense(t.n, t.rrs, sentence, 0)
	return tsris, err
}
