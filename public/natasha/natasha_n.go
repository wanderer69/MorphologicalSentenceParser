package natasha

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ocorp "github.com/wanderer69/OpCorpora"
)

type Natasha struct {
	pyCmd *PyCmd
}

// UniversalDependecies
var internalUniversalDependecies map[string]string

func InitUD() {
	internalUniversalDependecies = make(map[string]string)
	internalUniversalDependecies["amod"] = "определение"
	internalUniversalDependecies["nsubj"] = "номинальный_субъект" // подлежащее
	internalUniversalDependecies["obj"] = "объект"                // дополнение
	internalUniversalDependecies["punct"] = "знак_пунктуации"
	internalUniversalDependecies["root"] = "основа"          // может быть глагол
	internalUniversalDependecies["obl"] = "объект_локальный" //дополнение
	internalUniversalDependecies["case"] = "указатель"       // предлог
	internalUniversalDependecies["fixed"] = "фиксированный"
	internalUniversalDependecies["det"] = "определитель" // детерминатив
	internalUniversalDependecies["parataxis"] = "сочиненный"
	internalUniversalDependecies["iobj"] = "косвенное_дополнение"
	internalUniversalDependecies["nmod"] = "номинальный_модификатор"
	internalUniversalDependecies["advmod"] = "наречный_модификатор"
	internalUniversalDependecies["conj"] = "соединение"
	internalUniversalDependecies["nsubj:pass"] = "пассивный_номинальный_субъект"
	internalUniversalDependecies["aux:pass"] = "пассивный_вспомогательный"
	internalUniversalDependecies["aux"] = "вспомогательный"
	internalUniversalDependecies["xcomp"] = "открытое_клаузальное_дополнение"
	internalUniversalDependecies["cop"] = "связка"
	internalUniversalDependecies["advcl"] = "модификатор_придаточного_предложения"
	internalUniversalDependecies["mark"] = "маркер"
	internalUniversalDependecies["expl"] = "ругательство"
	internalUniversalDependecies["cc"] = "координирующее_соединение"
	internalUniversalDependecies["acl:relcl"] = "модификатор_относительного_предложения"
	internalUniversalDependecies["csubj"] = "клаузальный_субъект"
	internalUniversalDependecies["acl"] = "клаузальный_модификатор_существительного" // _придаточного_предложения
	internalUniversalDependecies["nummod:gov"] = "числовой_модификатор_регулирующий_падеж_существительного"
	internalUniversalDependecies["ccomp"] = "клаузальное_дополнение"
}

func Tag2strByUniversalDependecies(tag string) string {
	return internalUniversalDependecies[tag]
}

func NewNatasha(path string) *Natasha {
	n := Natasha{}

	n.pyCmd = &PyCmd{}
	n.pyCmd.Py_cmd_Init(path)
	//fmt.Printf("Python init!\r\n")
	ModulePy := "natasha1"
	// fmt.Printf("ModulePy %v\r\n", ModulePy)
	// инициализация модуля Python3 если он есть
	//        n.G_pc.Py_cmd_Import(ModulePy)

	if true {
		if ModulePy != "" {
			//		go n.G_pc.Py_cmd_Import(ModulePy)
			n.pyCmd.Py_cmd_Import(ModulePy)
		}
		//fmt.Printf("Module init!\r\n")
		//n.G_pc.Py_cmd_Wait()
		//fmt.Printf("Module inited!\r\n")
	}
	ocorp.Global_Dict_Init()
	InitUD()

	return &n
}

func (n *Natasha) ExecParseSentence(sentence string) string {
	sf := []interface{}{}
	sf = append(sf, sentence)
	res := n.pyCmd.Py_cmd_Call( /*n.pyCmd.Module,*/ "parse_sentence", sf)

	return res
}

func (n *Natasha) Close() {
	n.pyCmd.Py_cmd_Close()
}

type WordData struct {
	Rel      string
	Pos      string
	Feats    map[string]string
	Start    string
	Stop     string
	Text     string
	Lemma    string
	Id       string
	HeadID   string
	IdN      int
	SidN     int
	HeadIdN  int
	SheadIdN int
}

func (n *Natasha) ParseSentence(str string) ([]WordData, error) {
	res := n.ExecParseSentence(str)
	res = res[3 : len(res)-2]
	res_, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		fmt.Printf("error %v\r\n", err)
		return nil, err
	}
	res = string(res_)
	//fmt.Printf("res %v\r\n", res)

	var s []map[string]string
	wds := []WordData{}
	err = json.Unmarshal([]byte(res), &s)
	if err != nil {
		fmt.Printf("error %v\r\n", err)
		return nil, err
	}
	for i := range s {
		wd := WordData{}
		for k, v := range s[i] {
			switch k {
			case "rel":
				vn := Tag2strByUniversalDependecies(v)
				if len(vn) > 0 {
					wd.Rel = vn
				} else {
					wd.Rel = v
				}
			case "pos":
				_, vn := ocorp.Tag2str_attr_int(v)
				if len(vn) > 0 {
					wd.Pos = vn
				} else {
					wd.Pos = v
				}
			case "feats":
				//				var sr map[string]string
				var s map[string]string
				err := json.Unmarshal([]byte(v), &s)
				if err != nil {
					fmt.Printf("error %v\r\n", err)
					return nil, err
				}
				for nn, vv := range s {
					vv := strings.ToLower(vv)
					//fmt.Printf("vv %v\r\n", vv)
					sn, vn := ocorp.Tag2str_attr_ru_int(vv)
					if len(vn) > 0 {
						// fmt.Printf("sn %v, vn %v\r\n", sn, vn)
						s[sn] = vn
					} else {
						// s[sn] = v
						if false {
							fmt.Printf("-> %v, %v\r\n", nn, vv)
						}
					}
				}
				wd.Feats = s
			case "start":
				wd.Start = v
			case "stop":
				wd.Stop = v
			case "text":
				wd.Text = v
			case "lemma":
				wd.Lemma = v
			case "id":
				wd.Id = v
				sl := strings.Split(v, "_")
				sn, err := strconv.ParseInt(sl[0], 10, 64)
				if err != nil {
					fmt.Printf("error %v\r\n", err)
					return nil, err
				}
				pn, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					fmt.Printf("error %v\r\n", err)
					return nil, err
				}
				wd.IdN = int(pn)
				wd.SidN = int(sn)
			case "head_id":
				wd.HeadID = v
				sl := strings.Split(v, "_")
				sn, err := strconv.ParseInt(sl[0], 10, 64)
				if err != nil {
					fmt.Printf("error %v\r\n", err)
					return nil, err
				}
				pn, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					fmt.Printf("error %v\r\n", err)
					return nil, err
				}
				wd.HeadIdN = int(pn)
				wd.SheadIdN = int(sn)
			}
			//fmt.Printf("%v", res[i])
		}
		wds = append(wds, wd)
	}

	return wds, nil
}
