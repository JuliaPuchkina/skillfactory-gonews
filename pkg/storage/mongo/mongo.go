package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "data"
	collectionName = "posts"
)

// хранилище данных
type Store struct {
	db *mongo.Client
}

// конструктор объекта хранилища
func New(constr string) (*Store, error) {
	// подключение к СУБД MongoDB
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	defer client.Disconnect(context.Background())
	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	s := Store{
		db: client,
	}

	return &s, err
}

// Posts выводит все существующие публикации
func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	var data []storage.Post
	for cur.Next(context.Background()) {
		var l storage.Post
		err := cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		data = append(data, l)
	}
	return data, cur.Err()
}

// AddPost создает новую публикацию
func (s *Store) AddPost(doc storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePost обновляет публикацию
func (s *Store) UpdatePost(doc storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	_, err := collection.UpdateOne(context.Background(), filter, doc)
	if err != nil {
		return err
	}
	return nil
}

// DeletePost удаляет публикацию
func (s *Store) DeletePost(doc storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteOne(context.Background(), bson.M{"ID": doc.ID})
	if err != nil {
		return err
	}
	return nil
}
