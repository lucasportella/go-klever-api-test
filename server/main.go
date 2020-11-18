package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/roneycharles/klever/model"
	postpb "github.com/roneycharles/klever/third_party/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServiceServer struct {
}

var db *mongo.Client
var postDB *mongo.Collection
var mongoCtx context.Context

func (s *PostServiceServer) CreatePost(ctx context.Context, request *postpb.CreatePostRequest) (*postpb.CreatePostResponse, error) {

	post := request.GetPost()

	data := model.Post{
		ID:      primitive.NewObjectID(),
		Title:   post.GetTitle(),
		Content: post.GetContent(),
	}

	result, err := postDB.InsertOne(mongoCtx, data)

	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid := result.InsertedID.(primitive.ObjectID)
	post.Id = oid.Hex()

	return &postpb.CreatePostResponse{Post: post}, nil
}

func (s *PostServiceServer) GetPost(ctx context.Context, request *postpb.GetPostRequest) (*postpb.GetPostResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := postDB.FindOne(ctx, bson.M{"_id": oid})

	data := model.Post{}

	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find post with Object Id %s: %v", request.GetId(), err))
	}

	response := &postpb.GetPostResponse{
		Post: &postpb.Post{
			Id:      oid.Hex(),
			Title:   data.Title,
			Content: data.Content,
		},
	}
	return response, nil
}

func (s *PostServiceServer) DeletePost(ctx context.Context, request *postpb.DeletePostRequest) (*postpb.DeletePostResponse, error) {
	oid, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	_, err = postDB.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete post with id %s: %v", request.GetId(), err))
	}
	return &postpb.DeletePostResponse{
		Success: true,
	}, nil
}

func (s *PostServiceServer) ListPosts(request *postpb.ListPostsRequest, stream postpb.PostService_ListPostsServer) error {

	data := &model.Post{}

	cursor, err := postDB.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {

		err := cursor.Decode(data)

		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}

		stream.Send(&postpb.ListPostsResponse{
			Post: &postpb.Post{
				Id:      data.ID.Hex(),
				Content: data.Content,
				Title:   data.Title,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

func (s *PostServiceServer) UpdatePost(ctx context.Context, request *postpb.UpdatePostRequest) (*postpb.UpdatePostResponse, error) {

	post := request.GetPost()

	oid, err := primitive.ObjectIDFromHex(post.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied post id to a MongoDB ObjectId: %v", err),
		)
	}

	// Convert the data to be updated into an unordered Bson document
	update := bson.M{
		"title":   post.GetTitle(),
		"content": post.GetContent(),
	}

	// Convert the oid into an unordered bson document to search by id
	filter := bson.M{"_id": oid}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := postDB.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := model.Post{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find post with supplied ID: %v", err),
		)
	}
	return &postpb.UpdatePostResponse{
		Post: &postpb.Post{
			Id:      decoded.ID.Hex(),
			Title:   decoded.Title,
			Content: decoded.Content,
		},
	}, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting server on port :50051...")

	listener, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	srv := &PostServiceServer{}

	postpb.RegisterPostServiceServer(s, srv)

	fmt.Println("Connecting to MongoDB...")
	mongoCtx = context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb+srv://roney:candy1989@cluster0.fpwgl.mongodb.net/kleverdb?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	postDB = db.Database("kleverdb").Collection("posts")

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)

	<-c

	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}
