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

	n := natasha.NewNatasha()
	defer n.Close()
	require.NoError(t, n.Init())
	//rrs := InitRelationRule()
	rulesFileName := "/home/user/Go_projects/SemanticNet/data/rules.script"
	//rulesFileName = "/home/user/Go_projects/SemanticNet/data/rules_short_root_это.script"
	rrs, err := LoadRulesFromFile(rulesFileName)
	require.NoError(t, err)

	sentence := "студент собрал дом из деталей."
	//sentence = "руководство это  руководить"
	sentence = "хвост - задняя, конечная часть чего-нибудь движущегося; вообще что-нибудь длинное, движущееся"
	//sentence = "хвост - задняя конечная часть чего-нибудь движущегося; вообще что-нибудь длинное, движущееся"
	//sentence = "хвост - задняя часть чего-нибудь движущегося."
	//sentence = "хвост - задняя часть летательного аппарата."
	//sentence = "случайность - случайное обстоятельство"
	sentence = "сущность - внутреннее содержание предмета, обнаруживщееся во внешних формах его существования"
	sentence = "сущность -  суть"
	sentence = "средство - прием, способ действия для достижения чего-нибудь"
	sentence = "средство - приём или способ действия для достижения чего-нибудь"
	//sentence = "средство - это способ действия для достижения чего-нибудь"
	sentence = "предмет это то, на что направлена мысль, какое-нибудь действие, объект"
	sentence = "предмет это объект на который направлена мысль или какое-нибудь действие"
	//sentence = "объект на который направлена мысль или какое-нибудь действие"
	sentence = "дисциплина это обязательное для всех членов какого-нибудь коллектива подчинение установленному порядку, правилам"
	sentence = "дисциплина - это обязательное подчинение установленному порядку или правилам"
	sentence = "суть это сущность"
	sentence = "плата - диэлектрическая пластина для установки электронных компонентов"                              // !!
	sentence = "команда - автоматически передаваемый сигнал, вызывающий действие какой-нибудь системы или механизма" // !!
	sentence = "случай - это случайное обстоятельство"                                                               // !!
	//sentence = "случай это  случайное событие"
	sentence = "содержание это количество чего-нибудь, находящегося в чем-нибудь другом" // !!
	sentence = "часть -  предмет как составной элемент какого-нибудь целого"             //, организма, механизма"
	sentence = "часть - это составной элемент какого-нибудь целого"                      //, организма, механизма"
	sentence = "часть это элемент в составе какого-нибудь целого"                        // !!                     //, организма, механизма"
	// sentence = "часть - это элемент в составе целого предмета"                           //, организма, механизма"
	// sentence = "часть - элемент в составе целого"                                                     //, организма, механизма"
	sentence = "связь - это совокупность учреждений, обслуживающих техническими средствами общение на расстоянии" // !! //  (почта, телеграф, телефон, радио)
	sentence = "содержание - это единство всех элементов целого, его свойств и связей, существующее и выражаемое в форме и неотделимое от нее"
	sentence = "содержание это единство элементов целого существующее в форме"
	sentence = "содержание - это то, что составляет сущность чего-нибудь"

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
