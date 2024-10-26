package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ayobami6/todo_cli/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID          int
	Description string
	CreatedAt   time.Time
	IsComplete  bool
}

type User struct {
	Passcode string
}

func NewUser(passcode string) *User {
	return &User{
		Passcode: passcode,
	}
}

type UserRepo struct {
	db *mongo.Client
}

func NewUserRepo() (*UserRepo, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return &UserRepo{db: client}, nil
}

func (u *UserRepo) saveUser(passcode string) error {
	collection := u.db.Database("todo").Collection("users")
	_, err := collection.InsertOne(context.TODO(), NewUser(passcode))
	if err != nil {
		return err
	}
	return nil
}

func SaveUser(passcode string) {
	userRepo, err := NewUserRepo()
	if err != nil {
		log.Fatalf("error saving user")
	}
	err = userRepo.saveUser(passcode)
	if err != nil {
		panic("error saving user")
	}

}

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// GetClient returns a singleton instance of the MongoDB client
func GetClient() (*mongo.Client, error) {
	var err error
	clientOnce.Do(func() {
		mongoUrl := utils.GetEnv("MONGO_URL", "mongodb://localhost:27017")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(mongoUrl)
		clientInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
	})

	return clientInstance, err
}

func NewTask(description string) *Task {
	return &Task{
		ID:          1,
		Description: description,
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}
}

type TaskRepo struct {
	db *mongo.Client
}

func (t *TaskRepo) addTask(description string) error {
	collection := t.db.Database("todo").Collection("tasks")
	_, err := collection.InsertOne(context.Background(), NewTask(description))
	return err
}

//

func NewTaskRepo() (*TaskRepo, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return &TaskRepo{db: client}, nil
}

func AddTask(description string) error {
	repo, err := NewTaskRepo()
	if err != nil {
		return err
	}
	err = repo.addTask(description)
	if err != nil {
		return err
	}
	fmt.Println("Tasked added")
	return nil
}
