package stanza

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/MorphologicalSentenceParser/public/entity"
)

func TestNewStanza(t *testing.T) {
	n := NewStanza()
	require.NoError(t, n.Init())
	str, err := n.ExecParseSentence("лес растёт на холме")
	require.NoError(t, err)
	fmt.Printf("%v", str)
	require.True(t, len(str) > 0)
	wordData, err := n.ParseSentence("лес растёт на холме")
	require.NoError(t, err)
	expectedWordData := []entity.WordData{
		{
			Rel:   "номинальный_субъект",
			Pos:   "имя_существительное",
			Feats: map[string]string{"одушевлённость": "неодушевлённое", "падеж": "именительный_падеж", "пол": "мужской_род", "число": "единственное_число"},
			//			Start:    "0",
			//			Stop:     "3",
			Text:  "лес",
			Lemma: "лес",
			//			Id:       "1_1",
			//			HeadID:   "1_2",
			IdN: 1,
			//			SidN:     1,
			HeadIdN: 2,
			//			SheadIdN: 1,
			StartN: 0,
			StopN:  3,
		},
		{
			Rel:   "основа",
			Pos:   "глагол_личная_форма",
			Feats: map[string]string{"атрибут": "спрягаемый", "время": "настоящее_время", "залог": "действительный_залог", "лицо": "третье_лицо", "наклонение": "изъявительное_наклонение", "совершённость": "несовершенный_вид", "число": "единственное_число"},
			//			Start:    "4",
			//			Stop:     "10",
			Text:  "растёт",
			Lemma: "расти",
			//			Id:       "1_2",
			//			HeadID:   "1_0",
			IdN: 2,
			//			SidN:     1,
			HeadIdN: 0,
			//			SheadIdN: 1,
			StartN: 4,
			StopN:  10,
		},
		{
			Rel: "указатель",
			Pos: "предлог",
			//			Feats: map[string]string{},
			//			Start:    "11",
			//			Stop:     "13",
			Text:  "на",
			Lemma: "на",
			//			Id:       "1_3",
			//			HeadID:   "1_4",
			IdN: 3,
			//			SidN:     1,
			HeadIdN: 4,
			//			SheadIdN: 1,
			StartN: 11,
			StopN:  13,
		},
		{
			Rel:   "объект_локальный",
			Pos:   "имя_существительное",
			Feats: map[string]string{"одушевлённость": "неодушевлённое", "падеж": "предложный_падеж", "пол": "мужской_род", "число": "единственное_число"},
			//			Start:    "14",
			//			Stop:     "19",
			Text:  "холме",
			Lemma: "холм",
			//			Id:       "1_4",
			//			HeadID:   "1_2",
			IdN: 4,
			//			SidN:     1,
			HeadIdN: 2,
			//			SheadIdN: 1,
			StartN: 14,
			StopN:  19,
		},
	}
	require.Equal(t, expectedWordData, wordData)
}
