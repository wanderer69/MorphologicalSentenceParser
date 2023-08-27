package grpc

import (
	"context"
	"fmt"

	grpcGoogle "google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/wanderer69/MorphologicalSentenceParser/internal/natasha"
	"github.com/wanderer69/MorphologicalSentenceParser/internal/relations"
	pb "github.com/wanderer69/MorphologicalSentenceParser/pkg/server/grpc/morphological_parser"
)

func GrpcInit(address string, port int) (*grpcGoogle.ClientConn, error) {
	opts := []grpcGoogle.DialOption{
		// grpcGoogle.WithInsecure(),
	}
	ss := fmt.Sprintf("%v:%v", address, port) // 5300
	conn, err := grpcGoogle.Dial(ss, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, nil
}

func PutSentenceToParsing(conn *grpcGoogle.ClientConn, query string) (string, error) {
	client := pb.NewMorphologicalSentenceParserClient(conn)
	request := &pb.PutSentenceToParsingRequest{
		Sentence: query,
	}
	response, err := client.PutSentenceToParsing(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return "", err
	}
	return response.TaskId, nil
}

func GrpcParsePhrase(conn *grpcGoogle.ClientConn, taskID string) ([]*relations.TranslateSentensesResultItem, error) {
	client := pb.NewMorphologicalSentenceParserClient(conn)
	request := &pb.GetResultSentenceParsingRequest{
		TaskId: taskID,
	}
	//	fmt.Printf("request %#v\r\n", request)
	response, err := client.GetResultSentenceParsing(context.Background(), request)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	results := []*relations.TranslateSentensesResultItem{}
	for i := range response.Result {
		r := response.Result[i]
		result := &relations.TranslateSentensesResultItem{
			Sentence: r.Sentence,
		}
		for j := range r.Relations {
			rr := r.Relations[j]
			res := &relations.Relation{
				Type:     rr.Type,
				ValuePtr: rr.ValuePtr,
				Value:    rr.Type,
				WordNum:  int(rr.WordNum),
			}
			result.Relations = append(result.Relations, res)
		}
		for j := range r.Relations {
			rr := r.Relations[j]
			if rr.Relation > -1 {
				result.Relations[j] = result.Relations[rr.Relation]
			}
		}
		for j := range r.WordDatas {
			wd := r.WordDatas[j]
			WordsData := natasha.WordData{
				Rel:      wd.Rel,
				Pos:      wd.Pos,
				Feats:    wd.Feats,
				Start:    wd.Start,
				Stop:     wd.Stop,
				Text:     wd.Text,
				Lemma:    wd.Lemma,
				Id:       wd.Id,
				HeadID:   wd.HeadId,
				IdN:      int(wd.IdN),
				SidN:     int(wd.SidN),
				HeadIdN:  int(wd.HeadIdN),
				SheadIdN: int(wd.SheadIdN),
			}
			result.WordsData = append(result.WordsData, WordsData)
		}
		results = append(results, result)
	}
	return results, nil
}
