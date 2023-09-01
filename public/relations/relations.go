package relations

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v2"

	uuid "github.com/satori/go.uuid"
	"github.com/wanderer69/MorphologicalSentenceParser/public/natasha"
	"github.com/wanderer69/debug"
)

const (
	RObject          = "объект"
	RAgent           = "агент"
	RLocationIn      = "пункт_оправления"
	RLocationOut     = "пункт_назначения"
	RTool            = "инструмент"
	RAction          = "действие"
	RCoagent         = "соагент"
	RDeliveryVehicle = "средство_доставки"
	RTrajectory      = "траектория"
	RLocation        = "местоположение"
	RRawMaterial     = "сырьё"
	RConsumer        = "потребитель"
	RTime            = "время"
	RWishAction      = "желание"
	RNeedAction      = "необходимость"

	RMemberOrControl = "членство_управление"

	RProperty             = "свойство"
	RPropertyQuery        = "запрос_свойства"
	RObjectQuery          = "запрос_объект"
	RAgentQuery           = "запрос_агент"
	RTimeQuery            = "запрос_время"
	RLocationQuery        = "запрос_местоположение"
	RToolQuery            = "запрос_инструмент"
	RLocationInQuery      = "запрос_пункт_оправления"
	RLocationOutQuery     = "запрос_пункт_назначения"
	RDeliveryVehicleQuery = "запрос_средство_доставки"
	RTrajectoryQuery      = "запрос_траектория"
	RRawMaterialQuery     = "запрос_сырьё"
	RConsumerQuery        = "запрос_потребитель"

	RCoagentQuery = "запрос_соагент"

	// настоятель монастыря - дополнение указывающее на принадлежность настоятеля к монастырю то есть свойство объекту
	RPropertyToObject = "объект_свойство"
	// ваза с цветами - объект с объектом
	RObjectToObject = "объект_объект"

	RCondition   = "условие_логическое"
	RConsequence = "следствие_логическое"

	RConditionTime   = "условие_времени"
	RConsequenceTime = "следствие_времени"
)

type Relation struct {
	Type     string
	ValuePtr string    // предлог
	Value    string    // значение
	Relation *Relation // связь
	WordNum  int
}

/*
описание конструкции
отношение
Часть речи
падеж
управление - от root, от object и т.п. - то есть должны быть определены на этот момент
наличие предлога со <значением>
-> вид отношения с указанием леммы и возможно леммы предлога

Итого, вначале идет описание root, потом объект, потом остальное.
*/

type Action struct {
	Cmd  string
	Args []string
}

type Actions struct {
	Actions []Action
}

func NewAction(cmd string, args ...string) Actions {
	a := Action{Cmd: cmd, Args: args}
	return Actions{Actions: []Action{a}}
}

func (aa Actions) NewAction(cmd string, args ...string) Actions {
	a := Action{Cmd: cmd, Args: args}
	aa.Actions = append(aa.Actions, a)
	return aa
}

type RelationRuleCondition struct {
	Relation     string // имя отношения
	PartOfSpeach string // часть речи
	Case         string // падеж
	Control      string // тип зависимости
	Pretext      string // предлог
	Animated     string // одушевленность
	HaveObject   string // имеет объект
	RootIs       []string
	FromRelation string
	Actions      []Action
	Lemma        string // лемма
	NoPretext    bool   // нет предлога

	DependRelation string // зависимость от отношения
	ID             string
}

func (rrc RelationRuleCondition) Print() string {
	s := []string{}
	if len(rrc.Relation) > 0 {
		s = append(s, fmt.Sprintf("имя_отношения:%v", rrc.Relation))
	}
	if len(rrc.PartOfSpeach) > 0 {
		s = append(s, fmt.Sprintf("часть_речи:%v", rrc.PartOfSpeach))
	}
	if len(rrc.Case) > 0 {
		s = append(s, fmt.Sprintf("падеж:%v", rrc.Case))
	}
	if len(rrc.Control) > 0 {
		s = append(s, fmt.Sprintf("тип_зависимости:%v", rrc.Control))
	}
	if len(rrc.Pretext) > 0 {
		s = append(s, fmt.Sprintf("предлог:%v", rrc.Pretext))
	}
	if len(rrc.Animated) > 0 {
		s = append(s, fmt.Sprintf("одушевленность:%v", rrc.Animated))
	}
	if len(rrc.HaveObject) > 0 {
		s = append(s, fmt.Sprintf("имеет_объект:%v", rrc.HaveObject))
	}
	if len(rrc.RootIs) > 0 {
		ss := ""
		for i := range rrc.RootIs {
			if i == 0 {
				ss = rrc.RootIs[i]
			} else {
				ss = ss + ";" + rrc.RootIs[i]
			}
		}
		s = append(s, fmt.Sprintf("root_is:%v", ss))
	}
	if len(rrc.Lemma) > 0 {
		s = append(s, fmt.Sprintf("лемма:\"%v\"", rrc.Lemma))
	}
	if rrc.NoPretext {
		s = append(s, "нет_предлога:true")
	}
	if len(rrc.DependRelation) > 0 {
		s = append(s, fmt.Sprintf("зависимость_от_отношения:%v", rrc.DependRelation))
	}
	if len(rrc.ID) > 0 {
		s = append(s, fmt.Sprintf("идентификатор:%v", rrc.ID))
	}
	return "\t\tесли(" + strings.Join(s, ", ") + ");"
}

func CreateHash(rrc RelationRuleCondition) string {
	s := ""
	for i := range rrc.RootIs {
		s = s + rrc.RootIs[i]
	}
	s1 := fmt.Sprintf("%v", rrc.NoPretext)

	s2 := rrc.Relation + rrc.PartOfSpeach + rrc.Case + rrc.Control +
		rrc.Pretext + rrc.Animated + rrc.HaveObject + s +
		rrc.FromRelation + rrc.Lemma + s1

	h := sha1.New()
	h.Write([]byte(s2))
	//return string(h.Sum(nil))
	return s2
}

type RelationType struct {
	RelationType   string // отношение
	UsePretext     bool   // заполняем предлог
	ChangeRoot     int
	ChangeObject   int
	ChangeRootBase int

	IsComplex     bool
	IsCondition   bool
	IsComma       bool
	IsConsequence bool
	IsComplicate  bool
	Actions       []Action // список действий
	ID            string
}

func (a Action) Print() string {
	s := []string{}
	s = append(s, a.Cmd)
	s = append(s, a.Args...)
	return "\t\t\tдействие(" + strings.Join(s, ", ") + ");"
}

func (rt RelationType) Print() string {
	s := []string{}
	if len(rt.RelationType) != 0 {
		s = append(s, fmt.Sprintf("RelationType:%v", rt.RelationType))
	}
	if rt.UsePretext {
		s = append(s, "UsePretext:true")
	}
	if rt.ChangeRoot != 0 {
		s = append(s, fmt.Sprintf("ChangeRoot:%v", rt.ChangeRoot))
	}
	if rt.ChangeObject != 0 {
		s = append(s, fmt.Sprintf("ChangeObject:%v", rt.ChangeObject))
	}
	if rt.ChangeRootBase != 0 {
		s = append(s, fmt.Sprintf("ChangeRootBase:%v", rt.ChangeRootBase))
	}
	if rt.IsComplex {
		s = append(s, "IsComplex:true")
	}
	if rt.IsCondition {
		s = append(s, "IsCondition:true")
	}
	if rt.IsComma {
		s = append(s, "IsComma:true")
	}
	if rt.IsConsequence {
		s = append(s, "IsConsequence:true")
	}
	if rt.IsComplicate {
		s = append(s, "IsComplicate:true")
	}
	if len(rt.ID) != 0 {
		s = append(s, fmt.Sprintf("ID:%v", rt.ID))
	}
	ss := []string{}
	for i := range rt.Actions {
		ss = append(ss, rt.Actions[i].Print())
	}
	actions := "\t\t\tпусто;"
	if len(ss) > 0 {
		actions = strings.Join(ss, "\r\n")
	}
	return "\t\tтип(" + strings.Join(s, ", ") + ") {\r\n" + actions + "\r\n\t\t};"
}

type RelationRuleItem struct {
	Type                   string
	RelationRuleConditions []RelationRuleCondition // условие и
	RelationTypes          []RelationType
	ID                     string
}

func (rri *RelationRuleItem) SetID() {
	for i := range rri.RelationRuleConditions {
		rri.RelationTypes[i].ID = uuid.NewV4().String()
	}
	for i := range rri.RelationTypes {
		rri.RelationRuleConditions[i].ID = uuid.NewV4().String()
	}
	rri.ID = uuid.NewV4().String()
}

func (rri RelationRuleItem) Print() string {
	ss := ""
	for i := range rri.RelationRuleConditions {
		rri.RelationTypes[i].ID = uuid.NewV4().String()
		s := rri.RelationRuleConditions[i].Print()
		ss = ss + s + "\r\n"
	}
	for i := range rri.RelationTypes {
		rri.RelationRuleConditions[i].ID = uuid.NewV4().String()
		s := rri.RelationTypes[i].Print()
		ss = ss + s + "\r\n"
	}
	ss = ss + fmt.Sprintf("\t\tидентификатор(%v);\r\n", rri.ID)
	return "\tправило {\r\n" + ss + "\t};"
}

type RelationRules struct {
	Main []RelationRuleItem
}

func (rr RelationRules) Print() string {
	dict := make(map[string][]string)
	for i := range rr.Main {
		rri := rr.Main[i].Print()
		rria := dict[rr.Main[i].Type]
		rria = append(rria, rri+"\r\n")
		dict[rr.Main[i].Type] = rria
	}
	ss := ""
	for k, v := range dict {
		ss = ss + "блок " + k + " {\r\n" + strings.Join(v, "") + "};\r\n"
	}
	return ss
}

/*
Скрипт правила
RelationRuleItem.Type(
	- RelationType, ChangeRoot {
		"SaveVar", "root", "val"
	}
): {
	Relation: "основа"
	PartOfSpeach: "имя_существительное"
	Case: ""
	Control: "zero"
	Pretext: ""
	Animated: ""
	HaveObject: ""
}

Type:
если(
	основа, имя_существительное, zero
) {

}
*/

func InitRelationRule() *RelationRules {
	rrs := RelationRules{}
	rri := RelationRuleItem{Type: "Root", RelationTypes: []RelationType{{RelationType: RAction, ChangeRoot: 1, Actions: NewAction("SaveVar", "root", "val").Actions}}}
	rr := RelationRuleCondition{Relation: "основа", PartOfSpeach: "глагол_личная_форма", Case: "", Control: "zero", Pretext: "", Animated: "", HaveObject: ""}
	rri.RelationRuleConditions = append(rri.RelationRuleConditions, rr)
	(&rri).SetID()
	rrs.Main = append(rrs.Main, rri)

	rri = RelationRuleItem{Type: "Object", RelationTypes: []RelationType{{RelationType: RObject, ChangeObject: 1, Actions: NewAction("SaveVar", "object", "val").Actions}}}
	rr = RelationRuleCondition{Relation: "объект", PartOfSpeach: "имя_существительное", Case: "винительный_падеж", Control: "root", Pretext: "", Animated: "", HaveObject: ""}
	rri.RelationRuleConditions = append(rri.RelationRuleConditions, rr)
	(&rri).SetID()
	rrs.Main = append(rrs.Main, rri)

	rri = RelationRuleItem{Type: "Case", RelationTypes: []RelationType{{RelationType: RRawMaterial}}}
	rr = RelationRuleCondition{Relation: "номинальный_модификатор", PartOfSpeach: "имя_существительное", Case: "родительный_падеж", Control: "object", Animated: "", Pretext: "из", HaveObject: ""}
	rri.RelationRuleConditions = append(rri.RelationRuleConditions, rr)
	(&rri).SetID()
	rrs.Main = append(rrs.Main, rri)

	rri = RelationRuleItem{Type: "Case", RelationTypes: []RelationType{{RelationType: RAgent}}}
	rr = RelationRuleCondition{Relation: "номинальный_субъект", PartOfSpeach: "имя_существительное", Case: "", Control: "root", Animated: "", Pretext: "", HaveObject: ""}
	rri.RelationRuleConditions = append(rri.RelationRuleConditions, rr)
	(&rri).SetID()
	rrs.Main = append(rrs.Main, rri)

	rri = RelationRuleItem{Type: "Extention", RelationTypes: []RelationType{{RelationType: RProperty}}}
	rr = RelationRuleCondition{Relation: "определение", PartOfSpeach: "имя_прилагательное", Case: "", Control: "by_relation", Pretext: "", Animated: "", HaveObject: ""} // , FromRelation: "RObject"
	rri.RelationRuleConditions = append(rri.RelationRuleConditions, rr)
	(&rri).SetID()
	rrs.Main = append(rrs.Main, rri)
	return &rrs
}

func SaveToYaml(rr *RelationRules, fileName string) error {
	yamlData, err := yaml.Marshal(rr)
	if err != nil {
		//fmt.Printf("Error while Marshaling. %v", err)
		return err
	}
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	return nil
}

func SaveToScript(rr *RelationRules, fileName string) error {
	data := (*rr).Print()
	err := os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	return nil
}

type Pair struct {
	rule      int
	posInRule int
	data      int
}

type VarValue struct {
	rule int
	//tag   string
	value string
}

type ComplexPair struct {
	BeginMain      int
	EndMain        int
	BeginCondition int
	EndCondition   int
	PosCondition   int
	PosComma       int
	PosCosequence  int
	Complicate     int
}

type RelationEnv struct {
	root     int
	object   int
	rootBase int

	pos       int   // положение в списке данных
	dataStack []int // стек положений в списке данных

	posInRRC int    // положение в списке условий правил
	state    int    // состояние 0 - начальный элемент
	stack    []Pair // стек соответствий - сюда пишем сопоставленные

	usedRules []int // список использованных правил

	posRRIA int // положение в списке правил

	Relations []*Relation

	Variables map[string]*VarValue

	CurrentRelation *Relation // отношение для связи с отношением

	RelationsUse map[string]struct{}

	Complex *ComplexPair
}

func NewRelationEnv() *RelationEnv {
	re := RelationEnv{}
	re.Variables = make(map[string]*VarValue)
	re.RelationsUse = make(map[string]struct{})
	return &re
}

func (re *RelationEnv) AddRelationUse(rrc RelationRuleCondition, pos int) {
	key := CreateHash(rrc) + fmt.Sprintf("%v", pos)
	re.RelationsUse[key] = struct{}{}
}

func (re *RelationEnv) CheckRelationUse(rrc RelationRuleCondition, pos int) bool {
	key := CreateHash(rrc) + fmt.Sprintf("%v", pos)
	_, ok := re.RelationsUse[key]
	return ok
}

func Args(args ...any) map[string]interface{} {
	m := make(map[string]interface{})
	if (len(args)/2)*2 != len(args) {
		// чётное?
		return m
	}
	state := 0
	vn := ""
	for i := range args {
		switch state {
		case 0:
			vn = args[i].(string)
			state = 1
		case 1:
			m[vn] = args[i]
			state = 0
		}
	}
	return m
}

func (re *RelationEnv) ExecAction(a Action, val map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	switch a.Cmd {
	case "SaveVar":
		// первый аргумент - имя переменной, ???? второй - имя значения
		if len(a.Args) == 2 {
			vn := a.Args[0]
			vvn := a.Args[1]
			vv := VarValue{rule: re.posRRIA, value: val[vvn].(string)}
			re.Variables[vn] = &vv
		}
	case "CheckVar":
		// первый аргумент - имя переменной, второй - значение
		if len(a.Args) == 2 {
			vn := a.Args[0]
			vvn := a.Args[1]
			vv, ok := re.Variables[vn]
			if ok {
				if vv.value == val[vvn].(string) {
					m[vn] = true
				} else {
					m[vn] = false
				}
			} else {
				m[vn] = false
			}
		}
	}
	return m
}

// Removes slice element at index(s) and returns new slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func CheckRelationByRule(rrs *RelationRules, wd []natasha.WordData) ([]*Relation, error) {
	// функция предиката
	predicate := func(mode string, rria []RelationRuleItem, i int, re *RelationEnv) *Relation {
		var r *Relation
		r = nil
		isModeExists := false
		debug.Alias("level1.1").Printf("predicate mode %v wd[i] %v %v pos_in_rrc %v re.pos_rria %v\r\n", mode, wd[i].Text, wd[i].Rel, re.posInRRC, re.posRRIA)
		for {
			state := 0
			ptr := -1
			flagBreak := false

			mm := re.posRRIA
			if rria[mm].Type == mode {
				if mode == "Object" {
					fmt.Printf("%v\r\n", mode)
				}
				isModeExists = true
				rrcl := rria[mm].RelationRuleConditions
				fmt.Printf("predicate re.pos_rria %v re.pos_in_rrc %v\r\n", re.posRRIA, re.posInRRC)

				if re.CheckRelationUse(rrcl[re.posInRRC], i) {
					// fmt.Printf("CheckRelationUse true\r\n")
					if re.posInRRC >= 0 {
						re.state = -3
						state = 15
						break
					} else {
						if len(rria)-1 > re.posRRIA {
							re.posRRIA = re.posRRIA + 1
						} else {
							break
						}
						continue
					}
				}
				//fmt.Printf("rrcl %v re.pos_rria %v\r\n", rrcl, re.pos_rria)
				rr := rrcl[re.posInRRC]

				prev := -1
				prev_ptr := 0
				debug.Alias("level1").Printf("predicate rr %#v\r\n", rr)
				for {
					debug.Alias("level2").Printf("predicate state %v ptr %v\r\n", state, ptr)
					if prev == state {
						prev_ptr = prev_ptr + 1
						if prev_ptr > 5 {
							panic("problem!")
						}
					} else {
						prev_ptr = 0
						prev = state
					}
					switch state {
					case 0:
						if len(rr.Relation) > 0 {
							//fmt.Printf("rr.Relation %v, wd[i].Rel %v, rr.Relation == wd[i].Rel %v\r\n", rr.Relation, wd[i].Rel, rr.Relation == wd[i].Rel)
							if rr.Relation == wd[i].Rel {
								state = 1
							} else {
								state = 20
							}
						} else {
							state = 1
						}
					case 1:
						if len(rr.PartOfSpeach) > 0 {
							if rr.PartOfSpeach == wd[i].Pos {
								state = 2
							} else {
								state = 20
							}
						} else {
							state = 2
						}
					case 2:
						if len(rr.Case) > 0 {
							v, ok := wd[i].Feats["падеж"]
							if ok {
								if v == rr.Case {
									state = 3
								} else {
									state = 20
								}
							} else {
								state = 20
							}
						} else {
							state = 3
						}
					case 3:
						// проверяем управление
						if len(rr.Control) > 0 {
							switch rr.Control {
							case "zero":
								if wd[i].HeadIdN == 0 {
									state = 4
								} else {
									state = 20
								}
							case "self":
								if wd[i].HeadIdN == wd[i].IdN {
									state = 4
								} else {
									state = 20
								}
							case "root":
								if re.root < 0 {
									state = 20
								} else {
									if wd[i].HeadIdN == wd[re.root].IdN {
										state = 4
									} else {
										state = 20
									}
								}
							case "object":
								if re.object < 0 {
									state = 20
								} else {
									if wd[i].HeadIdN == wd[re.object].IdN {
										state = 4
									} else {
										state = 20
									}
								}
							case "root_base":
								if re.rootBase < 0 {
									state = 20
								} else {
									if wd[i].HeadIdN == wd[re.rootBase].IdN {
										state = 4
									} else {
										state = 20
									}
								}
							case "by_relation":
								// цикл по известным отношениям
								flagBreak := true
								for ii := range re.Relations {
									debug.Alias("level3").Printf("re.Relations[i] %v, wd[i].HeadIdN %v\r\n", re.Relations[ii], wd[i].HeadIdN)
									if re.Relations[ii].WordNum == wd[i].HeadIdN {
										state = 4
										re.CurrentRelation = re.Relations[ii]
										flagBreak = false
										break
									}
								}
								if flagBreak {
									state = 20
								}
							case "stack":
								// fmt.Printf("re.stack %#v\r\n", re.stack)
								if len(re.stack) == 0 {
									state = 20
								} else {
									// fmt.Printf("wd[i].HeadIdN %v, wd[re.stack[0].data].IdN %v, wd[i].HeadIdN == wd[re.stack[0].data].IdN %v\r\n", wd[i].HeadIdN, wd[re.stack[0].data].IdN, wd[i].HeadIdN == wd[re.stack[0].data].IdN)
									if wd[i].HeadIdN == wd[re.stack[0].data].IdN {
										state = 4
									} else {
										state = 20
									}
								}
							}
						} else {
							state = 4
						}
					case 4:
						if len(rr.Animated) > 0 {
							// проверяем одушевленность
							dd, ok := wd[i].Feats["одушевлённость"]
							if !ok {
								state = 20
							} else {
								if dd == rr.Animated {
									state = 5
								} else {
									state = 20
								}
							}
						} else {
							state = 5
						}
					case 5:
						if len(rr.HaveObject) > 0 {
							switch rr.HaveObject {
							case "object":
								if re.object < 0 {
									state = 20
								} else {
									if wd[i].HeadIdN == wd[re.object].HeadIdN {
										state = 6
									} else {
										state = 20
									}
								}
							case "is_object":
								if re.object < 0 {
									state = 6
								} else {
									state = 20
								}
							}
						} else {
							state = 6
						}
					case 6:
						if len(rr.RootIs) > 0 {
							flag_nn := true
							for nn := range rr.RootIs {
								switch nn {
								case 0:
									// проверим что root - на самом деле - агент - имя существительное
									if len(rr.RootIs[nn]) > 0 {
										vv, ok := wd[re.root].Feats["падеж"]
										if ok {
											if vv == "именительный_падеж" {
											} else {
												flag_nn = false
											}
										} else {
											flag_nn = false
										}
									}
								case 1:
									if len(rr.RootIs[nn]) > 0 {
										// проверим, что root одушевленный
										vvv, ok := wd[re.root].Feats["одушевлённость"]
										if ok {
											if vvv == "одушевлённое" {

											} else {
												flag_nn = false
											}
										} else {
											flag_nn = false
										}
									}
								}
								if !flag_nn {
									break
								}
							}
							if flag_nn {
								// совпадает
								state = 7
							} else {
								state = 20
							}
						} else {
							state = 7
						}
					case 7:
						if len(rr.Lemma) > 0 {
							// проверяем лемму
							if wd[i].Lemma == rr.Lemma {
								state = 8
							} else {
								state = 20
							}
						} else {
							state = 8
						}

					case 8:
						if rr.NoPretext {
							if len(rr.Pretext) > 0 {
								// надо проверить что может быть предлог rr.Pretext
								flag := true
								for j := i - 1; j >= 0; j-- {
									if wd[j].Rel == "указатель" {
										//fmt.Printf("wd[j].HeadIdN %v, wd[i].IdN %v, wd[j].HeadIdN == wd[i].IdN %v\r\n", wd[j].HeadIdN, wd[i].IdN, wd[j].HeadIdN == wd[i].IdN)
										if wd[j].HeadIdN == wd[i].IdN {
											//fmt.Printf("wd[j] %#v, rr.Pretext %v, wd[j].Text == rr.Pretext %v\r\n", wd[j], rr.Pretext, wd[j].Text == rr.Pretext)
											if wd[j].Lemma == rr.Pretext {
												state = 20
												flag = false
												break
											}
										}
									}
								}
								if flag {
									state = 9
								}
							} else {
								// надо проверить что может быть предлог rr.Pretext
								flag := true
								for j := i - 1; j >= 0; j-- {
									if wd[j].Rel == "указатель" {
										//fmt.Printf("wd[j].HeadIdN %v, wd[i].IdN %v, wd[j].HeadIdN == wd[i].IdN %v\r\n", wd[j].HeadIdN, wd[i].IdN, wd[j].HeadIdN == wd[i].IdN)
										if wd[j].HeadIdN == wd[i].IdN {
											state = 20
											flag = false
											break
										}
									}
								}
								if flag {
									state = 9
								}
							}
						} else {
							state = 9
						}

					case 9:
						//fmt.Printf("rr.DependRelation %v\r\n", rr.DependRelation)
						if len(rr.DependRelation) > 0 {
							flag := true
							// ищем зависимость до
							for j := 0; j < len(wd)-1; j++ {
								//fmt.Printf("rr.DependRelation %v, wd[j].Rel %v\r\n", rr.DependRelation, wd[j].Rel)
								if rr.DependRelation == wd[j].Rel {
									//fmt.Printf("wd[j].HeadIdN %v, wd[i].IdN %v, wd[j].HeadIdN == wd[i].IdN %v\r\n", wd[j].HeadIdN, wd[i].IdN, wd[j].HeadIdN == wd[i].IdN)
									if wd[i].HeadIdN == wd[j].IdN {
										//fmt.Printf("wd[j] %#v, rr.Pretext %v, wd[j].Text == rr.Pretext %v\r\n", wd[j], rr.Pretext, wd[j].Text == rr.Pretext)
										ptr = j
										state = 10
										flag = false
										break
									}
								}
							}
							if flag {
								state = 20
							}
						} else {
							state = 10
						}

					case 10:
						if len(rr.Pretext) > 0 {
							// надо проверить что может быть предлог rr.Pretext
							flag := true
							for j := i - 1; j >= 0; j-- {
								if wd[j].Rel == "указатель" {
									//fmt.Printf("wd[j].HeadIdN %v, wd[i].IdN %v, wd[j].HeadIdN == wd[i].IdN %v\r\n", wd[j].HeadIdN, wd[i].IdN, wd[j].HeadIdN == wd[i].IdN)
									if wd[j].HeadIdN == wd[i].IdN {
										//fmt.Printf("wd[j] %#v, rr.Pretext %v, wd[j].Text == rr.Pretext %v\r\n", wd[j], rr.Pretext, wd[j].Text == rr.Pretext)
										if wd[j].Lemma == rr.Pretext {
											ptr = j
											state = 15
											flag = false
											break
										}
									}
								}
							}
							if flag {
								state = 20
							}
						} else {
							state = 15
						}
					case 15:
						//
						fmt.Printf("predicate len(rrcl)- 1 > re.posInRRC %v\r\n", len(rrcl)-1 > re.posInRRC)
						if len(rrcl)-1 > re.posInRRC {
							// изменяем состояние и выходим
							p := Pair{mm, re.posInRRC, i}
							re.stack = append(re.stack, p)
							if re.state < 0 {
								re.state = 0
							} else {
								re.state = re.state + 1
							}
							re.posInRRC = re.posInRRC + 1
						} else {
							p := Pair{mm, re.posInRRC, i}
							re.stack = append(re.stack, p)
							// получилось сопоставить
							if mode == "Pre" {
								posCondition := -1
								posComma := -1
								posCosequence := -1
								posComplicate := -1
								for j := range rria[mm].RelationTypes {
									if rria[mm].RelationTypes[j].IsComplex {
										if rria[mm].RelationTypes[j].IsCondition {
											posCondition = re.stack[j].data
										}
										if rria[mm].RelationTypes[j].IsComma {
											posComma = re.stack[j].data
										}
										if rria[mm].RelationTypes[j].IsConsequence {
											posCosequence = re.stack[j].data
										}
										if rria[mm].RelationTypes[j].IsComplicate {
											posComplicate = re.stack[j].data
										}
									}
								}
								//debug.Alias("level4.1").Printf("posCondition %v, posComma %v, posCosequence %v\r\n", posCondition, posComma, posCosequence)
								cp := ComplexPair{}
								cp.BeginMain = posCondition + 1
								cp.EndMain = posComma - 1
								cp.BeginCondition = posCosequence + 1
								cp.EndCondition = len(wd) - 1
								cp.Complicate = posComplicate

								cp.PosCondition = posCondition
								cp.PosComma = posComma
								cp.PosCosequence = posCosequence

								debug.Alias("level4.1").Printf("predicate cp %#v\r\n", cp)
								re.Complex = &cp
							}
							for j := range rria[mm].RelationTypes {
								if rria[mm].RelationTypes[j].ChangeRoot > 0 {
									p := rria[mm].RelationTypes[j].ChangeRoot - 1
									re.root = re.stack[p].data
								}
								if rria[mm].RelationTypes[j].ChangeRootBase > 0 {
									p := rria[mm].RelationTypes[j].ChangeRootBase - 1
									re.rootBase = re.stack[p].data
								}
								if rria[mm].RelationTypes[j].ChangeObject > 0 {
									p := rria[mm].RelationTypes[j].ChangeObject - 1
									re.object = re.stack[p].data
								}
								if len(rria[mm].RelationTypes[j].RelationType) > 0 {
									if rria[mm].RelationTypes[j].UsePretext {
										r = &Relation{Type: rria[mm].RelationTypes[j].RelationType, Value: wd[re.stack[j].data].Lemma, ValuePtr: wd[re.stack[j].data].Lemma, WordNum: wd[re.stack[j].data].IdN, Relation: re.CurrentRelation}
										re.Relations = append(re.Relations, r)
									} else {
										r = &Relation{Type: rria[mm].RelationTypes[j].RelationType, Value: wd[re.stack[j].data].Lemma, WordNum: wd[re.stack[j].data].IdN, Relation: re.CurrentRelation}
										re.Relations = append(re.Relations, r)
									}
								}
							}
							for iii := range re.stack {
								rrclNew := rria[re.stack[iii].rule].RelationRuleConditions
								rrNew := rrclNew[re.stack[iii].posInRule]
								re.AddRelationUse(rrNew, re.stack[iii].data)
							}
							debug.Alias("level4").Printf("predicate re.Relations %v\r\n", re.Relations)
							re.state = -2
						}
						flagBreak = true
					case 20:
						if len(re.stack) > 0 {
							state = 15
							if re.posInRRC > 0 {

							} else {
								if len(rria)-1 > re.posRRIA {
									re.posRRIA = re.posRRIA + 1
								}
							}
						}
						re.state = -3
						flagBreak = true
					}
					if flagBreak {
						break
					}
				}
			}
			if state == 15 {
				break
			}
			if len(rria)-1 > re.posRRIA {
				re.posRRIA = re.posRRIA + 1
			} else {
				break
			}
		}
		if !isModeExists {
			re.state = -3
		}
		return r
	}
	// проверяем на рут
	re := NewRelationEnv()

	re.root = -1
	re.rootBase = -1
	re.object = -1

	// текущее слово
	re.pos = 0
	re.state = -1
	re.posInRRC = 0
	re.posRRIA = 0
	// это указатель конечного автомата
	gstate := 0
	// список мод
	modes := []string{"Root", "Pre", "Object", "Case", "Extention"}
	// мода из списка мод
	mode := 0
	// список использованных правил
	re.usedRules = []int{}
	flag_stop := false
	incrementRelations := 0
	for {
		debug.Alias("level0").Printf("gstate %v, mode %v re.pos %v re.pos_rria %v\r\n", gstate, mode, re.pos, re.posRRIA)
		debug.Alias("level0.1").Printf("re %#v\r\n", re)
		fmt.Printf("re %v\r\n", re.RelationsUse)
		switch gstate {
		case 0:
			if re.pos == 0 {
				incrementRelations = len(re.Relations)
			}
			predicate(modes[mode], rrs.Main, re.pos, re)
			debug.Alias("level0").Printf("re.state %v\r\n", re.state)
			switch re.state {
			case -2:
				// выполнилось условие можем либо остановить либо перейти к другим кодам
				// смотрим стек
				if len(re.stack) > 0 {
					// выбираем из стека правило и отправляем его в список
					r := re.stack[0].rule
					re.usedRules = append(re.usedRules, r)
				}
				re.posRRIA = 0
				gstate = 1
			case -3:
				// не получилось. пропускаем.
				if len(re.stack) > 0 {
					if re.posInRRC > 0 {
						// смотрим следующее слово если оно есть
						gstate = 2
					}
				} else {
					re.posRRIA = 0
					gstate = 1
				}
			default:
				if re.state >= 0 {
					// добавляем текущую позицию в стек и начинаем заново
					re.dataStack = append(re.dataStack, re.pos)
					re.pos = 0
				} else {
					flag_stop = true
				}
			}
		case 1:
			if re.Complex != nil {
				flag_stop = true
				break
			}
			re.stack = []Pair{}
			re.posInRRC = 0
			if len(re.dataStack) > 0 {
				// в стеке есть запись, выбираем ее
				pos_n := re.dataStack[0]
				// возвращаем текущее слово
				re.pos = pos_n
			}
			re.dataStack = []int{}
			if re.posRRIA == 0 {
				if len(wd)-1 > re.pos {
					re.pos = re.pos + 1
					re.state = -1
				} else {
					newIncrementRelations := len(re.Relations)
					//fmt.Printf("newIncrementRelations %v, incrementRelations %v, newIncrementRelations == incrementRelations %v\r\n", newIncrementRelations, incrementRelations, newIncrementRelations == incrementRelations)
					if newIncrementRelations == incrementRelations {
						// иначе меняем моду
						if len(modes)-1 > mode {
							mode = mode + 1
							re.pos = 0
							re.state = -1
						} else {
							flag_stop = true
						}
					} else {
						re.pos = 0
						re.state = -1
					}
				}
			}
			gstate = 0
		case 2:
			if len(wd)-1 > re.pos {
				//re.pos_in_rrc = 0
				re.pos = re.pos + 1
				re.state = -1
				gstate = 0
			} else {
				if len(rrs.Main)-1 > re.posRRIA {
					re.posRRIA = re.posRRIA + 1
				}
				gstate = 1
			}
			if re.Complex != nil {
				flag_stop = true
			}
		}
		if flag_stop {
			break
		}
	}
	isPrint := false
	if re.Complex != nil {
		m := -1
		for i := range re.Relations {
			// fmt.Printf("re.Relations[i].WordNum %v, re.root %v\r\n", re.Relations[i].WordNum, re.root)
			if re.Relations[i].WordNum-1 == re.root {
				m = i
				break
			}
		}
		if m >= 0 {
			re.Relations = remove(re.Relations, m)
		}
		var ConditionRel *Relation
		var MainRel *Relation
		for i := range re.Relations {
			//fmt.Printf("re.Relations[i].WordNum %v, re.Complex.BeginCondition %v, re.Complex.BeginMain %v\r\n", re.Relations[i].WordNum, re.Complex.BeginCondition, re.Complex.BeginMain)
			if re.Relations[i].WordNum == re.Complex.BeginCondition {
				ConditionRel = re.Relations[i]
			}
			if re.Relations[i].WordNum == re.Complex.BeginMain {
				MainRel = re.Relations[i]
			}
		}

		wdCondition := wd[re.Complex.BeginCondition:re.Complex.EndCondition]
		for i := range wdCondition {
			wdCondition[i].IdN = wdCondition[i].IdN - re.Complex.BeginCondition
			wdCondition[i].HeadIdN = wdCondition[i].HeadIdN - re.Complex.BeginCondition
			if wdCondition[i].HeadIdN < 0 {
				wdCondition[i].HeadIdN = 0
			}
		}
		if isPrint {
			rows := []Row{}
			for i := range wdCondition {
				rows = TablePrint(rows, i, wdCondition[i].Lemma, wdCondition[i].Rel, wdCondition[i].Pos, wdCondition[i].IdN,
					wdCondition[i].HeadIdN, wdCondition[i].Feats)
			}
			table, err := PrintTable(rows)
			if err != nil {
				return nil, err
			}
			for i := range table {
				fmt.Printf("%v\r\n", table[i])
			}
		}
		wdMain := wd[re.Complex.BeginMain:re.Complex.EndMain]
		for i := range wdMain {
			wdMain[i].IdN = wdMain[i].IdN - re.Complex.BeginMain
			wdMain[i].HeadIdN = wdMain[i].HeadIdN - re.Complex.BeginMain
			if wdMain[i].HeadIdN < 0 {
				wdMain[i].HeadIdN = 0
			}
			if re.Complex.Complicate-1 == i {
				wdMain[i].HeadIdN = 0
				wdMain[i].Rel = "основа"
			}
		}
		if isPrint {
			rows := []Row{}
			for i := range wdMain {
				rows = TablePrint(rows, i, wdMain[i].Lemma, wdMain[i].Rel, wdMain[i].Pos, wdMain[i].IdN,
					wdMain[i].HeadIdN, wdMain[i].Feats)
			}
			table, err := PrintTable(rows)
			if err != nil {
				return nil, err
			}
			for i := range table {
				fmt.Printf("%v\r\n", table[i])
			}
		}
		var relsCondition []*Relation
		relsCondition, err := CheckRelationByRule(rrs, wdCondition)
		if err != nil {
			return nil, err
		}
		if isPrint {
			for i := range relsCondition {
				fmt.Printf("relsCondition %v\r\n", relsCondition[i])
				if relsCondition[i].Relation != nil {
					fmt.Printf("\t %v\r\n", relsCondition[i].Relation)
				}
			}
		}
		var relsMain []*Relation
		relsMain, err = CheckRelationByRule(rrs, wdMain)
		if err != nil {
			return nil, err
		}
		if isPrint {
			for i := range relsMain {
				fmt.Printf("relsMain %v\r\n", relsMain[i])
				if relsMain[i].Relation != nil {
					fmt.Printf("\t %v\r\n", relsMain[i].Relation)
				}
			}
		}
		re.Relations = []*Relation{}
		for i := range relsCondition {
			rc := relsCondition[i]
			rc.Relation = ConditionRel
			re.Relations = append(re.Relations, rc)
		}
		for i := range relsMain {
			rc := relsMain[i]
			rc.Relation = MainRel
			re.Relations = append(re.Relations, rc)
		}
	}
	return re.Relations, nil
}

func LoadSentensesNew(n *natasha.Natasha, rrs *RelationRules, fileNameIn string, fileNameOut string, debug int) error {
	bs, err := os.ReadFile(fileNameIn)
	if err != nil {
		return err
	}
	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		bs = bs[3:]
	}
	str := string(bs)
	lines_list := strings.Split(str, "\n")

	var f *os.File
	if len(fileNameOut) > 0 {
		f, err = os.OpenFile(fileNameOut, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(lines_list); i++ {
		ln := lines_list[i]
		line := strings.TrimSpace(ln)
		line = strings.ToLower(line)
		if len(line) > 0 {
			ch, _ := utf8.DecodeRune([]byte(line))
			if ch == '#' {
				// комментарий пропускаем
			} else {
				res, err := n.ParseSentence(line)
				if err != nil {
					return err
				}
				if len(fileNameOut) > 0 {
					fmt.Fprintf(f, "> %s\r\n", line)
					rows := []Row{}
					for i := range res {
						rows = TablePrint(rows, i, res[i].Lemma, res[i].Rel, res[i].Pos, res[i].IdN, res[i].HeadIdN, res[i].Feats)
					}
					table, err := PrintTable(rows)
					if err != nil {
						return err
					}
					for i := range table {
						fmt.Fprintf(f, "%v\r\n", table[i])
					}
					var rels []*Relation
					rels, err = CheckRelationByRule(rrs, res)
					if err != nil {
						return err
					}

					for i := range rels {
						fmt.Fprintf(f, "rels %v\r\n", rels[i])
						if rels[i].Relation != nil {
							fmt.Fprintf(f, "\t %v\r\n", rels[i].Relation)
						}
					}
				} else {
					for i := range res {
						fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\r\n", i, res[i].Lemma, res[i].Rel, res[i].Pos, res[i].IdN, res[i].HeadIdN, res[i].Feats)
					}
					var rels []*Relation
					rels, err = CheckRelationByRule(rrs, res)
					if err != nil {
						return err
					}
					for i := range rels {
						fmt.Printf("rels %v\r\n", rels[i])
						if rels[i].Relation != nil {
							fmt.Printf("\t %v\r\n", rels[i].Relation)
						}
					}
				}
			}
		}
	}

	if len(fileNameOut) > 0 {
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

type TranslateSentensesResultItem struct {
	Sentence  string
	WordsData []natasha.WordData
	Relations []*Relation
}

func TranslateSentense(n *natasha.Natasha, rrs *RelationRules, str_in string, debug int) ([]*TranslateSentensesResultItem, error) {
	lines_list := strings.Split(str_in, "\n")

	tsris := []*TranslateSentensesResultItem{}
	for i := 0; i < len(lines_list); i++ {
		ln := lines_list[i]
		line := strings.TrimSpace(ln)
		line = strings.ToLower(line)
		if len(line) > 0 {
			ch, _ := utf8.DecodeRune([]byte(line))
			if ch == '#' {
				// комментарий пропускаем
			} else {
				res, err := n.ParseSentence(line)
				if err != nil {
					return nil, err
				}
				var rels []*Relation
				rels, err = CheckRelationByRule(rrs, res)
				if err != nil {
					return nil, err
				}
				for i := range rels {
					fmt.Printf("rels %#v\r\n", rels[i])
					if rels[i].Relation != nil {
						fmt.Printf("\t %#v\r\n", rels[i].Relation)
					}
				}
				tsri := TranslateSentensesResultItem{
					Sentence:  line,
					WordsData: res,
					Relations: rels,
				}
				tsris = append(tsris, &tsri)
			}
		}
	}
	return tsris, nil
}

func TablePrint(rows []Row, args ...interface{}) []Row {
	row := Row{}
	for i := range args {
		s := fmt.Sprintf("%v", args[i])
		row = append(row, s)
	}
	rows = append(rows, row)
	return rows
}

type Row []string

func PrintTable(table []Row) ([]string, error) {
	y := len(table)
	if y == 0 {
		return nil, nil
	}
	// проверяем что длина всех строк одинаковая
	x := len(table[0])
	for i := 1; i < y; i++ {
		if len(table[i]) != x {
			return nil, fmt.Errorf("table has rows with differents length")
		}
	}
	// по каждому столбцу находим максимум
	lenRows := make([]int, x)
	for i := range lenRows {
		lenRows[i] = -1
	}
	for i := 0; i < y; i++ {
		for j := 0; j < len(table[i]); j++ {
			st := strings.Trim(table[i][j], " ")
			stLen := utf8.RuneCountInString(st)
			if stLen > lenRows[j] {
				lenRows[j] = stLen
			}
		}
	}
	str := []string{}
	for i := 0; i < y; i++ {
		line := ""
		for j := 0; j < x; j++ {
			if j == 0 {
				line = table[i][j]
				stLen := utf8.RuneCountInString(table[i][j])
				t := lenRows[j] - stLen
				space := ""
				for k := 0; k < t; k++ {
					space = space + " "
				}
				line = line + space
				line = line + "  "
			} else {
				line = line + table[i][j]
				stLen := utf8.RuneCountInString(table[i][j])
				t := lenRows[j] - stLen
				if j < x-1 {
					space := ""
					for k := 0; k < t; k++ {
						space = space + " "
					}
					line = line + space
					line = line + "  "
				}
			}
		}
		str = append(str, line)
	}
	return str, nil
}
