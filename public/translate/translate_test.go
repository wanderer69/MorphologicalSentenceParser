package translate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/debug"

	rulesLoader "github.com/wanderer69/MorphologicalSentenceParser/public/rules_loader"
)

func TestTranslateSentense(t *testing.T) {
	debug.NewDebug()
	debug.LoadFromFile("../../cmd/cli/debug.cfg")

	st := NewTranslate()
	st.Init()
	require.NotNil(t, st)

	rrs, err := rulesLoader.LoadRulesFromFile("../../../SemanticNet/data/rules.script")
	require.NoError(t, err)
	st.SetRRS(rrs)
	tsris, err := st.TranslateText("аббат - это настоятель мужского католического монастыря.", 0)
	require.NoError(t, err)
	for i := range tsris {
		fmt.Printf("%v %v %v %v %v\r\n", tsris[i].ObjectPos, tsris[i].RootBasePos, tsris[i].RootPos, tsris[i].Sentence, tsris[i].WordsData)
		for j := range tsris[i].Relations {
			fmt.Printf("%v %v %v %v", tsris[i].Relations[j].Type, tsris[i].Relations[j].Value, tsris[i].Relations[j].ValuePtr,
				tsris[i].Relations[j].WordNum)
			if tsris[i].Relations[j].Relation != nil {
				fmt.Printf(" %v %v %v", tsris[i].Relations[j].Relation.Type, tsris[i].Relations[j].Relation.Value, tsris[i].Relations[j].Relation.ValuePtr)
			}
			fmt.Printf("\r\n")
		}
	}
	//require.True(t, false)
}
