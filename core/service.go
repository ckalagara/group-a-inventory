package core

import (
	"context"
	"log"
	"time"

	pb "github.com/ckalagara/group-a-inventory/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	databaseName   = "group-a"
	collectionName = "inventory"
)

type Service struct {
	pb.UnimplementedServiceServer
	store *Store
}

type Store struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewService(ctx context.Context, mongoURI string) *Service {
	// Set client options for MongoDB connection
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create a new MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
		return nil
	}

	// Ping the MongoDB server to check if the connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil
	}

	// Initialize the collection to interact with 'items' in the 'inventorydb' database
	itemsCollection := client.Database(databaseName).Collection(collectionName)

	log.Println("Successfully connected to MongoDB")

	return &Service{store: &Store{client: client, collection: itemsCollection}}
}

func (s *Service) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	item := req.GetItem()

	// Insert the item into MongoDB
	_, err := s.store.collection.InsertOne(ctx, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add item: %v", err)
	}

	return &pb.AddItemResponse{Item: item}, nil
}

func (s *Service) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	filter := bson.D{{Key: "id", Value: req.GetId()}}
	var item pb.Item
	err := s.store.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Item not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to fetch item: %v", err)
	}

	return &pb.GetItemResponse{Item: &item}, nil
}

func (s *Service) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	cursor, err := s.store.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list items: %v", err)
	}
	defer cursor.Close(ctx)

	var items []*pb.Item
	for cursor.Next(ctx) {
		var item pb.Item
		if err := cursor.Decode(&item); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to decode item: %v", err)
		}
		items = append(items, &item)
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating over cursor: %v", err)
	}

	return &pb.ListItemsResponse{Items: items}, nil
}

func (s *Service) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	filter := bson.D{{Key: "id", Value: req.GetId()}}
	result, err := s.store.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete item: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "Item not found")
	}

	return &pb.DeleteItemResponse{Success: true}, nil
}

func (s *Service) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "Service is healthy"}, nil
}

// StreamItems is a server-streaming RPC
func (s *Service) StreamItems(req *pb.GetItemRequest, stream pb.Service_StreamItemsServer) error {
	filter := bson.D{{Key: "id", Value: req.GetId()}}

	cursor, err := s.store.collection.Find(context.Background(), filter)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to query items: %v", err)
	}
	defer cursor.Close(context.Background())

	// Stream each item to the client
	for cursor.Next(context.Background()) {
		var item pb.Item
		if err := cursor.Decode(&item); err != nil {
			return status.Errorf(codes.Internal, "Error decoding item: %v", err)
		}

		// Send each item to the client over the stream
		if err := stream.Send(&item); err != nil {
			return status.Errorf(codes.Internal, "Error sending item to client: %v", err)
		}

		// Simulate delay between responses (optional)
		time.Sleep(1 * time.Second)
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, "Error iterating over cursor: %v", err)
	}

	return nil
}
