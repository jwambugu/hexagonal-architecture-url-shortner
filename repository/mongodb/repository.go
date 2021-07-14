package mongodb

import (
	"context"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	errs "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// Find fetches and returns shortener.Redirect using the code provided
func (m *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}

	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"code": code}

	if err := collection.FindOne(ctx, filter).Decode(&redirect); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}

		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

// Store creates a new shortener.Redirect in the DB
func (m *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("redirects")

	_, err := collection.InsertOne(ctx, bson.M{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	})

	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}

// newMongoClient returns a client to connect to the db instance
func newMongoClient(url string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// NewMongoRepository returns an interface to interact with the DB
func NewMongoRepository(url, db string, timeout int) (shortener.RedirectRepository, error) {
	repo := &mongoRepository{
		database: db,
		timeout:  time.Duration(timeout) * time.Second,
	}

	client, err := newMongoClient(url, timeout)

	if err != nil {
		return nil, errs.Wrap(err, "repository.NewMongoRepository")
	}

	repo.client = client
	return repo, nil
}
