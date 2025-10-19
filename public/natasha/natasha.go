package natasha

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wanderer69/MorphologicalSentenceParser/public/entity"
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

	ocorp.GlobalDictInit()
	return &n
}

func (n *Natasha) Init() error {
	err := n.natashaClient.Init()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(5) * time.Second)
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

func (n *Natasha) ParseSentence(str string) ([]entity.WordData, error) {
	res, err := n.ExecParseSentence(str)
	if err != nil {
		fmt.Printf("ParseSentence: error %v\r\n", err)
		return nil, err
	}
	var s []map[string]string
	wds := []entity.WordData{}
	err = json.Unmarshal([]byte(res), &s)
	if err != nil {
		fmt.Printf("ParseSentence: error %v\r\n", err)
		return nil, err
	}
	for i := range s {
		wd := entity.WordData{}
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
				_, vn := ocorp.Tag2StrAttrInt(v)
				if len(vn) > 0 {
					wd.Pos = vn
				} else {
					wd.Pos = v
				}
			case "feats":
				var s map[string]string
				err := json.Unmarshal([]byte(v), &s)
				if err != nil {
					fmt.Printf("ParseSentence: error %v\r\n", err)
					return nil, err
				}
				for nn, vv := range s {
					vv := strings.ToLower(vv)
					sn, vn := ocorp.Tag2StrAttrRuInt(vv)
					if len(vn) > 0 {
						s[sn] = vn
					} else {
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
					fmt.Printf("ParseSentence: error %v\r\n", err)
					return nil, err
				}
				pn, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					fmt.Printf("ParseSentence: error %v\r\n", err)
					return nil, err
				}
				wd.IdN = int(pn)
				wd.SidN = int(sn)
			case "head_id":
				wd.HeadID = v
				sl := strings.Split(v, "_")
				sn, err := strconv.ParseInt(sl[0], 10, 64)
				if err != nil {
					fmt.Printf("ParseSentence: error %v\r\n", err)
					return nil, err
				}
				pn, err := strconv.ParseInt(sl[1], 10, 64)
				if err != nil {
					fmt.Printf("ParseSentence: error %v\r\n", err)
					return nil, err
				}
				wd.HeadIdN = int(pn)
				wd.SheadIdN = int(sn)
			}
		}
		wds = append(wds, wd)
	}

	return wds, nil
}
