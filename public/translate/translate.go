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

func (t *Translate) TranslateText(sentence string, debug int) ([]*relations.TranslateSentensesResultItem, error) {
	tsris, err := relations.TranslateText(t.n, t.rrs, sentence, 0)
	return tsris, err
}

func (t *Translate) TranslateSentence(sentence string, debug int) (*relations.TranslateSentensesResultItem, error) {
	tsri, err := relations.TranslateSentence(t.n, t.rrs, sentence, 0)
	return tsri, err
}
