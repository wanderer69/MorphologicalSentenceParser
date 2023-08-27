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
	return response, nil
}

func (s *Server) GetResultSentenceParsing(ctx context.Context, request *pb.GetResultSentenceParsingRequest) (*pb.GetResultSentenceParsingResponse, error) {
	response := &pb.GetResultSentenceParsingResponse{}
	return response, nil
}
