package translate

import (
	"fmt"

	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/debug"
)

type Translate struct {
	isHaveInit bool
	n          *natasha.Natasha
	rrs        *relations.RelationRules
}

func NewTranslate() *Translate {
	debug.NewDebug()
	//debug.LoadFromFile("../../cmd/cli/debug.cfg")

	n := natasha.NewNatasha()
	rrs := relations.InitRelationRule()
	return &Translate{
		n:   n,
		rrs: rrs,
	}
}

func (t *Translate) Init() error {
	err := t.n.Init()
	if err != nil {
		return err
	}
	t.isHaveInit = true
	return nil
}

func (t *Translate) SetRRS(rrs *relations.RelationRules) {
	t.rrs = rrs
}

func (t *Translate) TranslateText(sentence string, debug int) ([]*relations.TranslateSentensesResultItem, error) {
	if !t.isHaveInit {
		return nil, fmt.Errorf("not initialized")
	}
	tsris, err := relations.TranslateText(t.n, t.rrs, sentence, 0)
	return tsris, err
}

func (t *Translate) TranslateSentence(sentence string, debug int) (*relations.TranslateSentensesResultItem, error) {
	if !t.isHaveInit {
		return nil, fmt.Errorf("not initialized")
	}
	tsri, err := relations.TranslateSentence(t.n, t.rrs, sentence, 0)
	return tsri, err
}
