package translate

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/debug"
)

func TestTranslateSentense(t *testing.T) {
	debug.NewDebug()
	debug.LoadFromFile("../../cmd/cli/debug.cfg")

	st := NewTranslate()

	//st.SetRRS()
	require.NotNil(t, st)
}
