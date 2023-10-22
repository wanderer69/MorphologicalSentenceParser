package rulesloader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/debug"
)

func TestTranslateSentenseExt(t *testing.T) {
	debug.NewDebug()
	debug.LoadFromFile("../../cmd/cli/debug.cfg")

	n := natasha.NewNatasha("../../scripts/python")
	defer n.Close()
	//rrs := InitRelationRule()
	rrs, err := LoadRulesFromFile("/home/user/Go_projects/SemanticNet/data/rules.script")
	require.NoError(t, err)

	sentence := "студент собрал дом из деталей."
	//sentence = "руководство это  руководить"
	sentence = "хвост - это задняя, конечная часть чего-нибудь движущегося; вообще что-нибудь длинное, движущееся"
	sentence = "хвост - это задняя конечная часть чего-нибудь движущегося; вообще что-нибудь длинное, движущееся"
	sentence = "хвост - это задняя часть чего-нибудь движущегося."
	//sentence = "хвост - это задняя часть летательного аппарата."

	tsri, err := relations.TranslateSentence(n, rrs, sentence, 0)
	require.NoError(t, err)
	fmt.Printf("sentence %v\r\nroot %v obj %v rootBase %v\r\n", tsri.Sentence, tsri.RootPos, tsri.ObjectPos, tsri.RootBasePos)
	for j := range tsri.WordsData {
		fmt.Printf("id %v\tidN %v\theadId %v\theadIdN %v\trel %v lemma %v POS %v case %v sidN %v \r\n",
			tsri.WordsData[j].Id, tsri.WordsData[j].IdN, tsri.WordsData[j].HeadID, tsri.WordsData[j].HeadIdN,
			tsri.WordsData[j].Rel, tsri.WordsData[j].Lemma, tsri.WordsData[j].Pos, tsri.WordsData[j].Feats["падеж"], tsri.WordsData[j].SidN,
		)
	}
	for j := range tsri.Relations {
		fmt.Printf("%#v\r\n", tsri.Relations[j])
	}

	require.True(t, false)
}
