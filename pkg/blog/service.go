package blog

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

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
	log.Println("Create Blog request")
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

func (s *Service) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	log.Println("Read Blog request")
	blogID := req.GetBlogId()

	res := types.Blog{}

	filter := bson.M{"id": blogID}

	err := s.db.blog.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       res.ID,
			AuthorId: res.AuthorId,
			Title:    res.Title,
			Content:  res.Content,
		},
	}, nil
}

func (s *Service) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	log.Println("Update Blog request")
	blog := req.GetBlog()

	filter := bson.M{"id": blog.Id}

	update := bson.M{
		"$set": bson.M{
			"authorid": blog.GetAuthorId(),
			"title":    blog.GetTitle(),
			"content":  blog.GetContent(),
		},
	}

	updateRes, err := s.db.blog.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected error: %v", err),
		)
	}

	if updateRes.MatchedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Not found: %v", err),
		)
	}

	return &blogpb.UpdateBlogRes{}, nil
}

func (s *Service) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	log.Println("Delete Blog request")
	blogID := req.GetBlogId()

	blog := types.Blog{}

	filter := bson.M{"id": blogID}

	err := s.db.blog.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	_, err = s.db.blog.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected error: %v", err),
		)
	}

	res := blogpb.DeleteBlogRes{
		Blog: &blogpb.Blog{
			Id:       blog.ID,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}

	return &res, nil
}

func (s *Service) ListBlog(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogServer) error {
	log.Println("List Blog request")
	ctx := context.Background()

	cursor, err := s.db.blog.Find(ctx, bson.M{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected find error: %v", err),
		)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		blog := &types.Blog{}
		err := cursor.Decode(blog)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Unexpected cursor error: %v", err))
		}

		stream.Send(&blogpb.ListBlogRes{Blog: &blogpb.Blog{
			Id:       blog.ID,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		}})
	}
	if err = cursor.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unexpected error: %v", err))
	}
	return nil
}
