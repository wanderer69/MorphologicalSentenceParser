package script

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wanderer69/MorphologicalSentenceParser/internal/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/internal/relations"
	"github.com/wanderer69/debug"
	"github.com/wanderer69/tools/parser/parser"
	"github.com/wanderer69/tools/parser/print"
)

func TestScriptTranslate(t *testing.T) {
	debug.NewDebug()

	env := parser.NewEnv()
	buffer := ""
	o := print.NewOutput(func(sfmt string, args ...any) {
		s := fmt.Sprintf(sfmt, args...)
		buffer = buffer + s
	})

	MakeRules(env)

	t.Run("bad type name", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок пример {
				правило {
					если(параметр1;данные1, параметр2:данные2);
					тип(параметр1;данные1, параметр2:данные2) {
						действие(SaveVar, root, val);					
					};					
				};
			};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate блок типа: bad block type name пример")
	})

	t.Run("bad arg 1", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(параметр1;данные1, параметр2:данные2);
				тип(параметр1;данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate если: error - expected 7 args")
		//require.Contains(t, res, "тип пример {правило{если();действие();};}; ")
	})

	t.Run("bad arg 2", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(параметр1?данные1, параметр2:данные2);
				тип(параметр1;данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate если: error - expected colon")
	})

	t.Run("bad arg 3", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(?параметр1?данные1, параметр2:данные2);
				тип(параметр1;данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};				
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate если: error - argument want have <field name>:<value> format, got [? параметр1 ? данные1]")
		//require.Contains(t, res, "тип пример {правило{если();действие();};}; ")
	})

	t.Run("bad arg 4", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(параметр1:данные1, параметр2:данные2);
				тип(параметр1;данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate если: error - bad field name параметр1")
		//require.Contains(t, res, "тип пример {правило{если();действие();};}; ")
	})

	t.Run("bad arg 5", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(имя_отношения : данные1, часть_речи:данные2);
				тип(параметр1;данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};				
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate тип: error - expected 10 args")
	})

	t.Run("bad arg 6", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(имя_отношения : данные1, часть_речи:данные2);
				тип(параметр1:данные1, параметр2:данные2) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate тип: error - bad field name параметр1")
	})

	t.Run("bad arg 7", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(имя_отношения : данные1, часть_речи:данные2);
				тип(RelationType: RAction, ChangeRoot: n) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		_, err := env.ParseString(data, o)
		require.Error(t, err)
		require.EqualError(t, err, "error when translate тип: bad ChangeRoot type - n: strconv.ParseInt: parsing \"n\": invalid syntax")
	})

	t.Run("success", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		data := `блок Root {
			правило {
				если(имя_отношения : данные1, часть_речи:данные2);
				тип(RelationType: RAction, ChangeRoot: 2) {
					действие(SaveVar, root, val);					
				};
			};
		};	
		`
		res, err := env.ParseString(data, o)
		require.NoError(t, err)
		require.Contains(t, res, "блок Root {правило{если(имя_отношения : данные1, часть_речи:данные2);тип(RelationType: RAction, ChangeRoot: 2) {действие(SaveVar, root, val);};};};")
		expectedRelationRules := &relations.RelationRules{
			Main: []relations.RelationRuleItem{
				{
					Type: "Root",
					RelationRuleConditions: []relations.RelationRuleCondition{
						{
							Relation:     "данные1",
							PartOfSpeach: "данные2",
							Actions:      []relations.Action(nil),
						},
					},
					RelationTypes: []relations.RelationType{
						{
							RelationType: "RAction",
							ChangeRoot:   2,
							Actions: []relations.Action{
								{
									Cmd:  "SaveVar",
									Args: []string{"root", "val"},
								},
							},
						},
					},
				},
			},
		}
		require.Equal(t, expectedRelationRules, rp.Env.RelationRules)
	})

}

func TestScriptTranslateRule(t *testing.T) {
	debug.NewDebug()

	env := parser.NewEnv()
	buffer := ""
	o := print.NewOutput(func(sfmt string, args ...any) {
		s := fmt.Sprintf(sfmt, args...)
		buffer = buffer + s
	})

	MakeRules(env)

	//data, err := os.ReadFile("../../data/rules.script")
	//require.NoError(t, err)

	n := natasha.NewNatasha("")
	defer n.Close()
	rrs := relations.InitRelationRule()
	data := rrs.Print()
	err := os.WriteFile("../../data/rules.script", []byte(data), 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

	t.Run("success", func(t *testing.T) {
		rEnv := NewEnvironment()
		rp := RelationsParser{}
		rp.Env = rEnv
		env.Struct = &rp
		env.Debug = 0

		_, err := env.ParseString(string(data), o)
		require.NoError(t, err)

		//require.Contains(t, res, "блок Rule {правило{если(имя_отношения : данные1, часть_речи:данные2);тип(RelationType: RAction, ChangeRoot: 2) {действие(SaveVar, root, val);};};};")
		require.Equal(t, len(rrs.Main), len(rp.Env.RelationRules.Main))
		dict := make(map[string]*relations.RelationRuleItem)
		dictCond := make(map[string]*relations.RelationRuleCondition)
		for i := range rrs.Main {
			dict[rp.Env.RelationRules.Main[i].ID] = &rp.Env.RelationRules.Main[i]
			//fmt.Printf("rp.Env.RelationRules.Main[i].ID %v\r\n", rp.Env.RelationRules.Main[i].ID)
			rri := rp.Env.RelationRules.Main[i]
			for j := range rri.RelationRuleConditions {
				//fmt.Printf("\trri.RelationRuleConditions[j].ID %v\r\n", rri.RelationRuleConditions[j].ID)
				dictCond[rri.RelationRuleConditions[j].ID] = &rri.RelationRuleConditions[j]
			}
		}
		for i := range rrs.Main {
			//fmt.Printf("-> %v\r\n", rrs.Main[i].ID)
			rri, ok := dict[rrs.Main[i].ID]
			require.True(t, ok)
			require.Equal(t, len(rrs.Main[i].RelationRuleConditions), len(rri.RelationRuleConditions))
			for j := range rrs.Main[i].RelationRuleConditions {
				//fmt.Printf("--> %v\r\n", rrs.Main[i].RelationRuleConditions[j].ID)
				rrs.Main[i].RelationRuleConditions[j].ID = ""
				rri.RelationRuleConditions[j].ID = ""
			}
			require.ElementsMatch(t, rrs.Main[i].RelationRuleConditions, rri.RelationRuleConditions)
			for j := range rrs.Main[i].RelationTypes {
				//fmt.Printf("--> %v\r\n", rrs.Main[i].RelationTypes[j].ID)
				rrs.Main[i].RelationTypes[j].ID = ""
				rri.RelationTypes[j].ID = ""
			}
			require.ElementsMatch(t, rrs.Main[i].RelationTypes, rri.RelationTypes)
		}
	})
}
