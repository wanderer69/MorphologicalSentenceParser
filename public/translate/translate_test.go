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

	rrs, err := rulesLoader.LoadRulesFromFile("../../../SemanticNet/data/rules.script")
	require.NoError(t, err)
	st.SetRRS(rrs)
	tsris, err := st.TranslateText("аббат - это настоятель мужского католического монастыря.", 0)
	require.NoError(t, err)
	fmt.Printf("%#v\r\n", tsris)
	require.NotNil(t, st)
	require.True(t, false)
}
