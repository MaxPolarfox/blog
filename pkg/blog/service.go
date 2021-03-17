package blog

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/MaxPolarfox/blog/blogpb"
	"github.com/MaxPolarfox/blog/pkg/types"
	"github.com/MaxPolarfox/goTools/mongoDB"
)

type Service struct {
	options types.Options
	db      DB
}

type DB struct {
	blog mongoDB.Mongo
}

func NewService(options types.Options, blogCollection mongoDB.Mongo) *Service {
	return &Service{
		options: options,
		db: DB{
			blog: blogCollection,
		},
	}
}

func (s *Service) Start() {

	// listen to the appropriate signals, and notify a channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.options.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(server, s)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stopChan // wait for a signal to exit

	log.Println("shutting down the server")
	server.Stop()
	log.Println("Stopping listener")
	lis.Close()

	log.Println("End of program")
}

func (s *Service) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()

	data := types.Blog{
		ID:       uuid.NewV4().String(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	_, err := s.db.blog.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	return &blogpb.CreateBlogRes{
		Blog: &blogpb.Blog{
			Id:       data.ID,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}
