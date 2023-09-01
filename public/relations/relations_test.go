package relations

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/debug"

	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
)

func TestRelations(t *testing.T) {
	debug.NewDebug()
	//debug.LoadFromFile("debug.cfg")

	file_in := "../../data/phrases/test1.txt"
	file_out := "../../data/test1.out"

	n := natasha.NewNatasha("../../scripts/python")
	defer n.Close()
	rrs := InitRelationRule()
	err := SaveToYaml(rrs, "rules.yaml")
	require.NoError(t, err)
	err = LoadSentensesNew(n, rrs, file_in, file_out, 0)
	require.NoError(t, err)
}

func TestRelationsStoreToScript(t *testing.T) {
	debug.NewDebug()

	n := natasha.NewNatasha("../../scripts/python")
	defer n.Close()
	rrs := InitRelationRule()
	err := SaveToScript(rrs, "../../data/rules.script")
	require.NoError(t, err)
}

func TestTranslateSentense(t *testing.T) {
	debug.NewDebug()
	debug.LoadFromFile("../../cmd/cli/debug.cfg")

	n := natasha.NewNatasha("../../scripts/python")
	defer n.Close()
	rrs := InitRelationRule()
	tsris, err := TranslateSentense(n, rrs, "студент собрал дом из деталей.", 0)
	require.NoError(t, err)
	tsri := TranslateSentensesResultItem{
		Sentence: "студент собрал дом из деталей.",
		WordsData: []natasha.WordData{
			{
				Rel:      "номинальный_субъект",
				Pos:      "имя_существительное",
				Feats:    map[string]string{"Animacy": "Anim", "Case": "Nom", "Gender": "Masc", "Number": "Sing", "одушевлённость": "одушевлённое", "падеж": "именительный_падеж", "пол": "мужской_род", "число": "единственное_число"},
				Start:    "0",
				Stop:     "7",
				Text:     "студент",
				Lemma:    "студент",
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
				Feats:    map[string]string{"Aspect": "Perf", "Gender": "Masc", "Mood": "Ind", "Number": "Sing", "Tense": "Past", "VerbForm": "Fin", "Voice": "Act", "атрибут": "спрягаемый", "время": "прошедшее_время", "залог": "действительный_залог", "наклонение": "изъявительное_наклонение", "пол": "мужской_род", "совершённость": "совершенный_вид", "число": "единственное_число"},
				Start:    "8",
				Stop:     "14",
				Text:     "собрал",
				Lemma:    "собрать",
				Id:       "1_2",
				HeadID:   "1_0",
				IdN:      2,
				SidN:     1,
				HeadIdN:  0,
				SheadIdN: 1,
			},
			{
				Rel:      "объект",
				Pos:      "имя_существительное",
				Feats:    map[string]string{"Animacy": "Inan", "Case": "Acc", "Gender": "Masc", "Number": "Sing", "одушевлённость": "неодушевлённое", "падеж": "винительный_падеж", "пол": "мужской_род", "число": "единственное_число"},
				Start:    "15",
				Stop:     "18",
				Text:     "дом",
				Lemma:    "дом",
				Id:       "1_3",
				HeadID:   "1_2",
				IdN:      3,
				SidN:     1,
				HeadIdN:  2,
				SheadIdN: 1,
			},
			{
				Rel:      "указатель",
				Pos:      "предлог",
				Feats:    map[string]string{},
				Start:    "19",
				Stop:     "21",
				Text:     "из",
				Lemma:    "из",
				Id:       "1_4",
				HeadID:   "1_5",
				IdN:      4,
				SidN:     1,
				HeadIdN:  5,
				SheadIdN: 1,
			},
			{
				Rel:      "номинальный_модификатор",
				Pos:      "имя_существительное",
				Feats:    map[string]string{"Animacy": "Inan", "Case": "Gen", "Gender": "Fem", "Number": "Plur", "одушевлённость": "неодушевлённое", "падеж": "родительный_падеж", "пол": "женский_род", "число": "множественное_число"},
				Start:    "22",
				Stop:     "29",
				Text:     "деталей",
				Lemma:    "деталь",
				Id:       "1_5",
				HeadID:   "1_3",
				IdN:      5,
				SidN:     1,
				HeadIdN:  3,
				SheadIdN: 1,
			},
			{
				Rel:      "знак_пунктуации",
				Pos:      "знак_пунктуации",
				Feats:    map[string]string{},
				Start:    "29",
				Stop:     "30",
				Text:     ".",
				Lemma:    ".",
				Id:       "1_6",
				HeadID:   "1_2",
				IdN:      6,
				SidN:     1,
				HeadIdN:  2,
				SheadIdN: 1,
			},
		},
		Relations: []*Relation{
			{Type: "действие", ValuePtr: "", Value: "собрать", WordNum: 2},
			{Type: "объект", ValuePtr: "", Value: "дом", WordNum: 3},
			{Type: "агент", ValuePtr: "", Value: "студент", WordNum: 1},
			{Type: "сырьё", ValuePtr: "", Value: "деталь", WordNum: 5},
			//{Type: "объект_свойство", ValuePtr: "", Value: "деталь", WordNum: 5},
		},
	}
	expectTsris := []*TranslateSentensesResultItem{
		&tsri,
	}

	require.Equal(t, *expectTsris[0], *tsris[0])
}
