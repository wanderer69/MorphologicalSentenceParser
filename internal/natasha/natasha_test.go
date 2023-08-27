package natasha

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNatasha(t *testing.T) {
	n := NewNatasha("../../scripts/python")
	require.NotNil(t, n)
	str := n.ExecParseSentence("лес растёт на холме")
	fmt.Printf("%v", str)
	require.True(t, len(str) > 0)
	wordData, err := n.ParseSentence("лес растёт на холме")
	require.NoError(t, err)
	expectedWordData := []WordData{
		{
			Rel:      "номинальный_субъект",
			Pos:      "имя_существительное",
			Feats:    map[string]string{"Animacy": "Inan", "Case": "Nom", "Gender": "Masc", "Number": "Sing", "одушевлённость": "неодушевлённое", "падеж": "именительный_падеж", "пол": "мужской_род", "число": "единственное_число"},
			Start:    "0",
			Stop:     "3",
			Text:     "лес",
			Lemma:    "лес",
			Id:       "1_1",
			HeadID:   "1_2",
			IdN:      1,
			SidN:     1,
			HeadIdN:  2,
			SheadIdN: 1,
		},
		{
			Rel:      "основа",
			Pos:      "глагол_личная_форма",
			Feats:    map[string]string{"Aspect": "Imp", "Mood": "Ind", "Number": "Sing", "Person": "3", "Tense": "Pres", "VerbForm": "Fin", "Voice": "Act", "атрибут": "спрягаемый", "время": "настоящее_время", "залог": "действительный_залог", "лицо": "третье_лицо", "наклонение": "изъявительное_наклонение", "совершённость": "несовершенный_вид", "число": "единственное_число"},
			Start:    "4",
			Stop:     "10",
			Text:     "растёт",
			Lemma:    "расти",
			Id:       "1_2",
			HeadID:   "1_0",
			IdN:      2,
			SidN:     1,
			HeadIdN:  0,
			SheadIdN: 1,
		},
		{
			Rel:      "указатель",
			Pos:      "предлог",
			Feats:    map[string]string{},
			Start:    "11",
			Stop:     "13",
			Text:     "на",
			Lemma:    "на",
			Id:       "1_3",
			HeadID:   "1_4",
			IdN:      3,
			SidN:     1,
			HeadIdN:  4,
			SheadIdN: 1,
		},
		{
			Rel:      "объект_локальный",
			Pos:      "имя_существительное",
			Feats:    map[string]string{"Animacy": "Inan", "Case": "Loc", "Gender": "Masc", "Number": "Sing", "одушевлённость": "неодушевлённое", "падеж": "предложный_падеж", "пол": "мужской_род", "число": "единственное_число"},
			Start:    "14",
			Stop:     "19",
			Text:     "холме",
			Lemma:    "холм",
			Id:       "1_4",
			HeadID:   "1_2",
			IdN:      4,
			SidN:     1,
			HeadIdN:  2,
			SheadIdN: 1,
		},
	}
	require.Equal(t, expectedWordData, wordData)
}
