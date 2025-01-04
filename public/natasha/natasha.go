package natasha

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	natashaclient "github.com/wanderer69/MorphologicalSentenceParser/public/natasha_client"
	universaldependencies "github.com/wanderer69/MorphologicalSentenceParser/public/universal_dependencies"
	ocorp "github.com/wanderer69/OpCorpora/public/opcorpora"
)

type Natasha struct {
	isHaveInit           bool
	natashaClient        *natashaclient.NatashaClient
	universalDependecies *universaldependencies.UniversalDependencies
}

func NewNatasha() *Natasha {
	ud := universaldependencies.NewUniversalDependencies()
	natashaClient := natashaclient.NewNatashaClient()
	n := Natasha{
		natashaClient:        natashaClient,
		universalDependecies: ud,
	}

	ocorp.Global_Dict_Init()
	return &n
}

func (n *Natasha) Init() error {
	err := n.natashaClient.Init()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(2) * time.Second)
	n.isHaveInit = true
	return nil
}

func (n *Natasha) ExecParseSentence(sentence string) (string, error) {
	if !n.isHaveInit {
		return "", fmt.Errorf("not initialized")
	}
	ctx := context.Background()
	return n.natashaClient.ParsePhrase(ctx, sentence)
}

func (n *Natasha) Close() {
	n.natashaClient.Close()
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
	res, err := n.ExecParseSentence(str)
	/*
		res = res[3 : len(res)-2]
		res_, err := base64.StdEncoding.DecodeString(res)
		if err != nil {
			fmt.Printf("error %v\r\n", err)
			return nil, err
		}
		res = string(res_)
		//fmt.Printf("res %v\r\n", res)
	*/
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
				vn := n.universalDependecies.Tag2StrByUniversalDependecies(v)
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
