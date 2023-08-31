package grpc

import (
	"context"

	"github.com/wanderer69/MorphologicalSentenceParser/internal/worker"
	pb "github.com/wanderer69/MorphologicalSentenceParser/pkg/server/grpc/morphological_parser"
)

type Server struct {
	pb.MorphologicalSentenceParserServer
	proc *worker.Processor
}

func NewServer() *Server {
	return &Server{
		proc: worker.NewProcessor(),
	}
}

func (s *Server) PutSentenceToParsing(ctx context.Context, request *pb.PutSentenceToParsingRequest) (*pb.PutSentenceToParsingResponse, error) {
	response := &pb.PutSentenceToParsingResponse{}

	taskId, err := s.proc.Send(request.Sentence)
	if err != nil {
		return nil, err
	}
	response.TaskId = taskId
	response.ClientId = request.ClientId
	return response, nil
}

func (s *Server) GetResultSentenceParsing(ctx context.Context, request *pb.GetResultSentenceParsingRequest) (*pb.GetResultSentenceParsingResponse, error) {
	response := &pb.GetResultSentenceParsingResponse{}

	tsris, err := s.proc.Check(request.TaskId)
	if err == nil {
		return nil, err
	}
	results := []*pb.TranslateSentensesResultItem{}
	for i := range tsris {
		r := response.Result[i]
		result := &pb.TranslateSentensesResultItem{
			Sentence: r.Sentence,
		}
		for j := range r.Relations {
			rr := r.Relations[j]
			res := &pb.Relation{
				Type:     rr.Type,
				ValuePtr: rr.ValuePtr,
				Value:    rr.Type,
				WordNum:  int32(rr.WordNum),
			}
			result.Relations = append(result.Relations, res)
		}
		for j := range r.Relations {
			rr := r.Relations[j]
			if rr.Relation > -1 {
				result.Relations[j] = result.Relations[rr.Relation]
			}
		}
		WordDatas := []*pb.WordData{}
		for j := range r.WordDatas {
			wd := r.WordDatas[j]
			WordsData := pb.WordData{
				Rel:      wd.Rel,
				Pos:      wd.Pos,
				Feats:    wd.Feats,
				Start:    wd.Start,
				Stop:     wd.Stop,
				Text:     wd.Text,
				Lemma:    wd.Lemma,
				Id:       wd.Id,
				HeadId:   wd.HeadId,
				IdN:      int32(wd.IdN),
				SidN:     int32(wd.SidN),
				HeadIdN:  int32(wd.HeadIdN),
				SheadIdN: int32(wd.SheadIdN),
			}
			WordDatas = append(WordDatas, &WordsData)
		}
		result.WordDatas = WordDatas
		results = append(results, result)
	}
	response.Result = results
	return response, nil
}
