package script

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	attr "github.com/wanderer69/tools/parser/attributes"
	ns "github.com/wanderer69/tools/parser/new_strings"

	"github.com/wanderer69/MorphologicalSentenceParser/public/relations"
	"github.com/wanderer69/tools/parser/parser"
)

type RelationsParser struct {
	Env *Environment
}

type Environment struct {
	RelationRules *relations.RelationRules

	rri *relations.RelationRuleItem

	conditionTypes        *ConditionTypes
	relationRuleItemTypes *RelationRuleItemTypes
	relationTypeType      *RelationTypeTypes

	actions *relations.Actions

	relationRuleItemType string
	id                   string
}

func NewEnvironment() *Environment {
	ct := NewConditionType()
	rrit := NewRelationRuleItemType()
	rtt := NewRelationTypeType()

	env := Environment{conditionTypes: ct, relationRuleItemTypes: rrit, relationTypeType: rtt}
	env.RelationRules = &relations.RelationRules{}

	return &env
}

func ParseArg(val string) (*attr.Attribute, error) {
	ssl := ns.ParseStringBySignList(val, []string{"?", ":"})
	attrType := attr.AttrTConst
	if len(ssl) > 1 {
		attrType = attr.AttrTArray
	}
	a := attr.NewAttribute(byte(attrType), "", ssl)
	return a, nil
}

func fSymbol(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// символ
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func fString(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// строка
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func fVariable(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s := pi.Items[1].Data
		result = fmt.Sprintf("?%v", s)
		env.CE.State = 1000

		rp := env.Struct.(RelationsParser)

		env.Struct = rp
	}
	return result, nil
}

func fConst(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// символ
		s := pi.Items[0].Data
		result = fmt.Sprintf("%v", s)
		env.CE.State = 1000

		rp := env.Struct.(RelationsParser)

		env.Struct = rp
	}
	return result, nil
}

func fBlockType(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		// список правил
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		typeName := pi.Items[1].Data

		rp := env.Struct.(*RelationsParser)

		if !rp.Env.relationRuleItemTypes.In(typeName) {
			return result, fmt.Errorf("bad block type name %v", typeName)
		}

		rp.Env.relationRuleItemType = typeName
		env.Struct = rp

	case 1:
		body := env.CE.ResultGenerate
		typeName := pi.Items[1].Data

		result = fmt.Sprintf("блок %v {%v};", typeName, body)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)

		env.Struct = rp
	}
	return result, nil
}

func fRule(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		// список правил
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(*RelationsParser)

		rp.Env.rri = &relations.RelationRuleItem{}

		env.Struct = rp

	case 1:
		body := env.CE.ResultGenerate

		result = fmt.Sprintf("правило{%v};", body)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)

		if rp.Env.rri == nil {
			return result, errors.New("empty block type")
		}
		rp.Env.rri.Type = rp.Env.relationRuleItemType
		rp.Env.rri.ID = rp.Env.id
		rp.Env.RelationRules.Main = append(rp.Env.RelationRules.Main, *rp.Env.rri)
		rp.Env.rri = nil

		env.Struct = rp
	}
	return result, nil
}

type ConditionType string

const (
	ConditionTypeRelation     ConditionType = "имя_отношения"
	ConditionTypePartOfSpeach ConditionType = "часть_речи"
	ConditionTypeCase         ConditionType = "падеж"
	ConditionTypeControl      ConditionType = "тип_зависимости"
	ConditionTypePretext      ConditionType = "предлог"
	ConditionTypeAnimated     ConditionType = "одушевленность"
	ConditionTypeHaveObject   ConditionType = "имеет_объект"
	ConditionTypeRootIs       ConditionType = "root_is"
	//ConditionTypeFromRelation   ConditionType = ""
	ConditionTypeLemma          ConditionType = "лемма"
	ConditionTypeNoPretext      ConditionType = "нет_предлога"
	ConditionTypeDependRelation ConditionType = "зависимость_от_отношения"
	ConditionTypeID             ConditionType = "идентификатор"
)

type ConditionTypes struct {
	dict map[string]ConditionType
}

func NewConditionType() *ConditionTypes {
	cts := &ConditionTypes{
		dict: make(map[string]ConditionType),
	}
	cts.dict["имя_отношения"] = ConditionTypeRelation
	cts.dict["часть_речи"] = ConditionTypePartOfSpeach
	cts.dict["падеж"] = ConditionTypeCase
	cts.dict["тип_зависимости"] = ConditionTypeControl
	cts.dict["предлог"] = ConditionTypePretext
	cts.dict["одушевленность"] = ConditionTypeAnimated
	cts.dict["имеет_объект"] = ConditionTypeHaveObject
	cts.dict["root_is"] = ConditionTypeRootIs
	cts.dict["лемма"] = ConditionTypeLemma
	cts.dict["нет_предлога"] = ConditionTypeNoPretext
	cts.dict["зависимость_от_отношения"] = ConditionTypeDependRelation
	cts.dict["идентификатор"] = ConditionTypeID

	return cts
}

func (cts *ConditionTypes) In(ctv string) bool {
	_, ok := cts.dict[ctv]
	return ok
}

func (cts *ConditionTypes) Set(rrc *relations.RelationRuleCondition, ctv string, value string) error {
	ct, ok := cts.dict[ctv]
	if !ok {
		return fmt.Errorf("bad condition type - %v", ctv)
	}
	switch ct {
	case ConditionTypeRelation:
		rrc.Relation = value
	case ConditionTypePartOfSpeach:
		rrc.PartOfSpeach = value
	case ConditionTypeCase:
		rrc.Case = value
	case ConditionTypeControl:
		rrc.Control = value
	case ConditionTypePretext:
		rrc.Pretext = value
	case ConditionTypeAnimated:
		rrc.Animated = value
	case ConditionTypeHaveObject:
		rrc.HaveObject = value
	case ConditionTypeRootIs:
		ssl := strings.Split(value, ";")
		rrc.RootIs = ssl
	case ConditionTypeLemma:
		rrc.Lemma = value[1 : len(value)-1]
	case ConditionTypeNoPretext:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("bad condition type - %v: %w", ctv, err)
		}
		rrc.NoPretext = b
	case ConditionTypeDependRelation:
		rrc.DependRelation = value
	case ConditionTypeID:
		rrc.ID = value
	}
	return nil
}

type RelationRuleItemType string

const (
	RelationRuleItemTypeRoot      RelationRuleItemType = "Root"
	RelationRuleItemTypeObject    RelationRuleItemType = "Object"
	RelationRuleItemTypeCase      RelationRuleItemType = "Case"
	RelationRuleItemTypeExtention RelationRuleItemType = "Extention"
	RelationRuleItemTypePre       RelationRuleItemType = "Pre"
)

type RelationRuleItemTypes struct {
	dict map[string]RelationRuleItemType
}

func NewRelationRuleItemType() *RelationRuleItemTypes {
	rrits := &RelationRuleItemTypes{
		dict: make(map[string]RelationRuleItemType),
	}
	rrits.dict["Root"] = RelationRuleItemTypeRoot
	rrits.dict["Object"] = RelationRuleItemTypeObject
	rrits.dict["Case"] = RelationRuleItemTypeCase
	rrits.dict["Extention"] = RelationRuleItemTypeExtention
	rrits.dict["Pre"] = RelationRuleItemTypePre

	return rrits
}

func (rrits *RelationRuleItemTypes) In(rritv string) bool {
	_, ok := rrits.dict[rritv]
	return ok
}

type RelationTypeType string

const (
	RelationTypeRelationType   RelationTypeType = "RelationType"
	RelationTypeUsePretext     RelationTypeType = "UsePretext"
	RelationTypeChangeRoot     RelationTypeType = "ChangeRoot"
	RelationTypeChangeObject   RelationTypeType = "ChangeObject"
	RelationTypeChangeRootBase RelationTypeType = "ChangeRootBase"
	RelationTypeIsComplex      RelationTypeType = "IsComplex"
	RelationTypeIsCondition    RelationTypeType = "IsCondition"
	RelationTypeIsComma        RelationTypeType = "IsComma"
	RelationTypeIsConsequence  RelationTypeType = "IsConsequence"
	RelationTypeIsComplicate   RelationTypeType = "IsComplicate"
	RelationTypeID             RelationTypeType = "ID"
)

type RelationTypeTypes struct {
	dict map[string]RelationTypeType
}

func NewRelationTypeType() *RelationTypeTypes {
	rtts := &RelationTypeTypes{
		dict: make(map[string]RelationTypeType),
	}
	rtts.dict["RelationType"] = RelationTypeRelationType
	rtts.dict["UsePretext"] = RelationTypeUsePretext
	rtts.dict["ChangeRoot"] = RelationTypeChangeRoot
	rtts.dict["ChangeObject"] = RelationTypeChangeObject
	rtts.dict["ChangeRootBase"] = RelationTypeChangeRootBase
	rtts.dict["IsComplex"] = RelationTypeIsComplex
	rtts.dict["IsCondition"] = RelationTypeIsCondition
	rtts.dict["IsComma"] = RelationTypeIsComma
	rtts.dict["IsConsequence"] = RelationTypeIsConsequence
	rtts.dict["IsComplicate"] = RelationTypeIsComplicate
	rtts.dict["ID"] = RelationTypeID
	return rtts
}

func (rtts *RelationTypeTypes) In(rttv string) bool {
	_, ok := rtts.dict[rttv]
	return ok
}

func (rtts *RelationTypeTypes) Set(rt *relations.RelationType, rttv string, value string) error {
	rtt, ok := rtts.dict[rttv]
	if !ok {
		return fmt.Errorf("bad relation type - %v", rttv)
	}
	switch rtt {
	case RelationTypeRelationType:
		rt.RelationType = value
	case RelationTypeUsePretext:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}

		rt.UsePretext = b
	case RelationTypeChangeRoot:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
		}
		rt.ChangeRoot = int(n)
	case RelationTypeChangeObject:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
		}
		rt.ChangeObject = int(n)
	case RelationTypeChangeRootBase:
		n, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
		}
		rt.ChangeRootBase = int(n)
	case RelationTypeIsComplex:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}
		rt.IsComplex = b
	case RelationTypeIsCondition:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}
		rt.IsCondition = b
	case RelationTypeIsComma:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}
		rt.IsComma = b
	case RelationTypeIsConsequence:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}
		rt.IsConsequence = b
	case RelationTypeIsComplicate:
		b, err := strconv.ParseBool(value)
		if err != nil {
			if value != "!" {
				return fmt.Errorf("bad %v type - %v: %w", rttv, value, err)
			}
			b = true
		}
		rt.IsComplicate = b
	}
	return nil
}

func fCondition(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		//body := env.CE.ResultGenerate
		data := pi.Items[1].Data

		result = fmt.Sprintf("если(%v);", data)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)
		// Relation: "основа", PartOfSpeach: "глагол_личная_форма", Case: "", Control: "zero", Pretext: "", Animated: "", HaveObject: ""

		rrc := &relations.RelationRuleCondition{}
		b1 := strings.Trim(data, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}

				t, str, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					fmt.Printf("const %v %v\r\n", a, str)
					if len(arg) != 7 {
						fmt.Printf("error - expected 7 args")
						return result, errors.New("error - expected 7 args")
					}

				case attr.AttrTArray:
					//fmt.Printf("array %v %v\r\n", a, array)
					if len(array) != 3 {
						fmt.Printf("error - argument want have <field name>:<value> format, got %v", array)
						return result, fmt.Errorf("error - argument want have <field name>:<value> format, got %v", array)
					}
					if array[1] != ":" {
						fmt.Printf("error - expected :")
						return result, errors.New("error - expected colon")
					}
					if !rp.Env.conditionTypes.In(array[0]) {
						fmt.Printf("error - bad field name %v", array[0])
						return result, fmt.Errorf("error - bad field name %v", array[0])
					}
					err := rp.Env.conditionTypes.Set(rrc, array[0], array[2])
					if err != nil {
						return result, err
					}
				}
			}
		}
		rp.Env.rri.RelationRuleConditions = append(rp.Env.rri.RelationRuleConditions, *rrc)
		//r := Frames{rp.Operators}
		//rp.Env.AddFrames(r)
		env.Struct = rp
	}
	return result, nil
}

func fType(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// это просто секция по сути
		// список правил
		env.CE.PiCnt = 0
		env.CE.NextState = 1
		env.CE.State = 100

		rp := env.Struct.(*RelationsParser)

		env.Struct = rp

	case 1:
		body := env.CE.ResultGenerate
		data := pi.Items[1].Data

		result = fmt.Sprintf("тип(%v) {%v};", data, body)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)
		rt := &relations.RelationType{}
		b1 := strings.Trim(data, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}

				t, str, array := attr.GetAttribute(a)
				switch t {
				case attr.AttrTConst:
					fmt.Printf("const %v %v\r\n", a, str)
					if len(arg) != 7 {
						fmt.Printf("error - expected 10 args")
						return result, errors.New("error - expected 10 args")
					}

				case attr.AttrTArray:
					//fmt.Printf("array %v %v\r\n", a, array)
					if len(array) != 3 {
						fmt.Printf("error - argument want have <field name>:<value> format, got %v", array)
						return result, fmt.Errorf("error - argument want have <field name>:<value> format, got %v", array)
					}
					if array[1] != ":" {
						fmt.Printf("error - expected :")
						return result, errors.New("error - expected colon")
					}
					if !rp.Env.relationTypeType.In(array[0]) {
						fmt.Printf("error - bad field name %v", array[0])
						return result, fmt.Errorf("error - bad field name %v", array[0])
					}
					err := rp.Env.relationTypeType.Set(rt, array[0], array[2])
					if err != nil {
						return result, err
					}
				}
			}
		}
		if rp.Env.actions != nil {
			acts := *rp.Env.actions
			rt.Actions = acts.Actions
			rp.Env.actions = nil
		}
		rp.Env.rri.RelationTypes = append(rp.Env.rri.RelationTypes, *rt)

		env.Struct = rp
	}
	return result, nil
}

func fAction(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		data := pi.Items[1].Data

		result = fmt.Sprintf("действие(%v);", data)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)

		b1 := strings.Trim(data, " ")
		sl := ns.ParseStringBySignList(b1, []string{","})
		cmd := ""
		args := []string{}
		for i := range sl {
			arg := strings.Trim(sl[i], " ,")
			if len(arg) > 0 {
				if i == 0 {
					cmd = arg
				} else {
					args = append(args, arg)
				}
			}
		}
		actions := relations.NewAction(cmd, args...)
		rp.Env.actions = &actions

		env.Struct = rp
	}
	return result, nil
}

func fEmpty(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		result = "пусто;"
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)

		env.Struct = rp
	}
	return result, nil
}

func fId(pi parser.ParseItem, env *parser.Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		data := pi.Items[1].Data

		result = fmt.Sprintf("идентификатор(%v);", data)
		env.CE.State = 1000

		rp := env.Struct.(*RelationsParser)

		rp.Env.id = data
		env.Struct = rp
	}
	return result, nil
}

/*
	блок <type> {
		правило {
			тип() {
				действие();
				пусто();
			};
			если();
		};
	};
*/
func MakeRules(env *parser.Env) {
	if true {
		defer func() {
			r := recover()
			if r != nil {
				env.Output.Print("%v\r\n", r)
				return
			}
		}()
	}

	//<symbols, == тип>  <symbols, > <{, > - блок правил
	gr := parser.MakeRule("блок типа", env)
	gr.AddItemToRule("symbols", "", 1, "блок", "", []string{}, env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"правило"}, env)
	gr.AddRuleHandler(fBlockType, env)

	//<symbols, == правило> <{, > - определение правила
	gr = parser.MakeRule("правило", env)
	gr.AddItemToRule("symbols", "", 1, "правило", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"если", "тип", "идентификатор"}, env)
	gr.AddRuleHandler(fRule, env)

	//<symbols, == если> <(, > - определение условия
	gr = parser.MakeRule("если", env)
	gr.AddItemToRule("symbols", "", 1, "если", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fCondition, env)

	//<symbols, == тип> <(, >  <{, > - определение типа
	gr = parser.MakeRule("тип", env)
	gr.AddItemToRule("symbols", "", 1, "тип", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"действие", "пусто"}, env)
	gr.AddRuleHandler(fType, env)

	//<symbols, == действие> <(, > - определение действия
	gr = parser.MakeRule("действие", env)
	gr.AddItemToRule("symbols", "", 1, "действие", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env) // , "список"
	gr.AddRuleHandler(fAction, env)

	//<symbols, == идентификатор> <(, > - определение идентификатор
	gr = parser.MakeRule("идентификатор", env)
	gr.AddItemToRule("symbols", "", 1, "идентификатор", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{}, env) // , "список"
	gr.AddRuleHandler(fId, env)

	//<symbols, == пусто> <(, > - определение пусто
	gr = parser.MakeRule("пусто", env)
	gr.AddItemToRule("symbols", "", 1, "пусто", "", []string{}, env)
	gr.AddRuleHandler(fEmpty, env)

	// среднеуровневые элементы
	// список в определении тринара или шаблона
	// <symbols, > - просто символ
	gr = parser.MakeRule("символ", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fSymbol, env)

	// <string, > - просто строка
	gr = parser.MakeRule("строка", env)
	gr.AddItemToRule("string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fString, env)

	// ?<variable name>
	// <symbols, == ?> - переменная
	gr = parser.MakeRule("переменная", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fVariable, env)

	// <string, > - просто строка
	gr = parser.MakeRule("константа", env)
	gr.AddItemToRule("symbols|string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(fConst, env)

	// ?<variable name>:<attribute name>
	// <symbols, == ?> <symbols, > <symbols, == :> <symbols, >- атрибут переменной
	gr = parser.MakeRule("атрибут переменной", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	// ?<variable name>:<attribute name>
	// <symbols, > <symbols, == :> <symbols, >- атрибут
	gr = parser.MakeRule("атрибут", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	high_level_array := []string{"блок типа"}

	expr_array := []string{"атрибут переменной", "строка", "символ"}

	env.SetHLAEnv(high_level_array)
	env.SetEAEnv(expr_array)
	env.SetBGRAEnv()
}
