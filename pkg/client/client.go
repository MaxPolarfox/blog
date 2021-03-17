package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/MaxPolarfox/blog/blogpb"
	"github.com/MaxPolarfox/blog/pkg/types"
	goToolsClient "github.com/MaxPolarfox/goTools/client"
)

type Client interface {
	CreateBlog(ctx context.Context, data types.Blog) (*types.Blog, error)
}

type BlogClient struct {
	client blogpb.BlogServiceClient
	Conn   *grpc.ClientConn
}

func NewBlogClient(options goToolsClient.Options) BlogClient {
	var conn *grpc.ClientConn

	serverAddress := options.URL

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		panic(fmt.Sprintf("conn failed: %v", err))
	}

	client := blogpb.NewBlogServiceClient(conn)

	return BlogClient{
		client,
		conn,
	}
}

func (i *BlogClient) CreateBlog(ctx context.Context, data types.Blog) (*string, error) {
	req := &blogpb.CreateBlogReq{
		Blog: &blogpb.Blog{
			AuthorId: data.AuthorId,
			Title:    data.Title,
			Content:  data.Content,
		},
	}

	res, err := i.client.CreateBlog(ctx, req)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	createdId := res.GetBlog().GetId()

	return &createdId, nil
}

func (i *BlogClient) ReadBlog(ctx context.Context, blogId string) (*types.Blog, error) {
	req := &blogpb.ReadBlogReq{
		BlogId: blogId,
	}

	res, err := i.client.ReadBlog(ctx, req)
	if err != nil {
		return nil, err
	}

	blog := res.GetBlog()

	blogRes := types.Blog{
		ID:       blog.GetId(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	return &blogRes, nil
}

func (i *BlogClient) UpdateBlog(ctx context.Context, blog types.Blog) error {
	req := &blogpb.UpdateBlogReq{
		Blog: &blogpb.Blog{
			Id:       blog.ID,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}

	_, err := i.client.UpdateBlog(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (i *BlogClient) DeleteBlog(ctx context.Context, blogId string) (*types.Blog, error) {
	req := &blogpb.DeleteBlogReq{
		BlogId: blogId,
	}

	res, err := i.client.DeleteBlog(ctx, req)
	if err != nil {
		return nil, err
	}

	blog := res.GetBlog()

	blogRes := types.Blog{
		ID:       blog.GetId(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	return &blogRes, nil
}

func (i *BlogClient) ListBlog(ctx context.Context) ([]types.Blog, error) {

	blogs := []types.Blog{}

	stream, err := i.client.ListBlog(ctx, &blogpb.ListBlogReq{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalf("error while calling receiving blog: %v", err)
			}
		}

		blog := res.GetBlog()
		blogRes := types.Blog{
			ID:       blog.GetId(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		}
		blogs = append(blogs, blogRes)
	}

	return blogs, nil
}
