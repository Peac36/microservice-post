package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"post/persist"

	definition "github.com/peac36/microservice-definition"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	Boot()
}

func Boot() {
	fmt.Printf("The blog service is booting\n")

	var dbUsername string = os.Getenv("DB_USERNAME")
	var dbPassword string = os.Getenv("DB_PASSWORD")
	var dbAddress string = os.Getenv("DB_ADDRESS")
	var dbPort string = os.Getenv("DB_PORT")
	var dbName string = os.Getenv("DB_NAME")

	var network string = os.Getenv("GRPC_PROTOCOL")
	var address string = os.Getenv("GRPC_PROTOCOL")

	fmt.Printf("Connect to Database\n")
	connection := bootDatabase(dbUsername, dbPassword, dbAddress, dbPort, dbName)

	fmt.Printf("Run Migrations \n")
	connection.AutoMigrate(&persist.Post{})

	service := Service{repo: &persist.Repository{Connection: connection}}
	bootService(service, network, address)
}

func bootDatabase(username string, password string, address string, port string, databaseName string) *gorm.DB {
	dsnTemplate := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(dsnTemplate, username, password, address, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can not connect to database %s \n", err)
	}

	return db
}

func bootService(service Service, network string, address string) {
	server := grpc.NewServer()
	definition.RegisterPostServiceServer(server, service)

	listen, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln("Can bind the address")
	}
	fmt.Printf("Booting at: %s %s \n", network, address)
	server.Serve(listen)
}

type Service struct {
	repo persist.RepositoryInterface
	*definition.UnimplementedPostServiceServer
}

func (s Service) Create(ctx context.Context, req *definition.CreatePostRequest) (*definition.CreatePostResponse, error) {
	title := req.GetPost().GetTitle()
	content := req.GetPost().GetText()
	author := req.GetPost().GetAuthor()

	post, err := s.repo.Create(title, content, author)
	if err != nil {
		return &definition.CreatePostResponse{}, err
	}

	return &definition.CreatePostResponse{
		Post: &definition.Post{
			Title:  post.Title,
			Text:   post.Content,
			Id:     int32(post.ID),
			Author: post.Author,
		},
	}, nil
}

func (s Service) Update(ctx context.Context, req *definition.UpdatePostRequest) (*definition.UpdatePostResponse, error) {
	id := req.GetPost().GetId()
	title := req.GetPost().GetTitle()
	content := req.GetPost().GetText()
	author := req.GetPost().GetAuthor()

	post, err := s.repo.Update(int(id), title, content, author)
	if err != nil {
		return &definition.UpdatePostResponse{Post: &definition.Post{}}, err
	}

	return &definition.UpdatePostResponse{Post: &definition.Post{
		Id:     int32(post.ID),
		Title:  post.Title,
		Text:   post.Content,
		Author: post.Author,
	}}, nil
}

func (s Service) Delete(ctx context.Context, req *definition.DeletePostRequest) (*definition.DeletePostResponse, error) {
	var id int32 = req.GetId()
	_, err := s.repo.Delete(int(id))
	if err != nil {
		return nil, err
	}
	return &definition.DeletePostResponse{Id: id}, nil
}

func (s Service) Get(ctx context.Context, req *definition.GetPostRequest) (*definition.GetPostResponse, error) {
	var id int32 = req.GetId()

	post, err := s.repo.Get(int(id))
	if err != nil {
		return nil, err
	}

	return &definition.GetPostResponse{Post: &definition.Post{
		Id:     int32(post.ID),
		Title:  post.Title,
		Text:   post.Content,
		Author: post.Author,
	}}, nil
}

func (s Service) Index(ctx context.Context, req *definition.IndexPostRequest) (*definition.IndexPostResponse, error) {
	var posts []*definition.Post
	result, err := s.repo.Index(int(req.GetPage()), int(req.GetPerPage()))
	if err != nil {
		return nil, err
	}

	for _, value := range result {
		posts = append(posts, &definition.Post{
			Id:     int32(value.ID),
			Title:  value.Title,
			Text:   value.Content,
			Author: value.Author,
		})
	}
	return &definition.IndexPostResponse{Posts: posts}, nil
}
