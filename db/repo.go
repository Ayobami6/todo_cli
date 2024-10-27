package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Ayobami6/todo_cli/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID           uuid.UUID
	Description  string
	CreatedAt    time.Time
	IsComplete   bool
	UserPasscode string
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

func (u *UserRepo) findOne(passcode string) (*User, error) {
	collection := u.db.Database("todo").Collection("users")
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.D{{Key: "passcode", Value: passcode}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	// jsonify the data
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	return NewUser(result["passcode"].(string)), nil
}

func FetchUser(passcode string) (*User, error) {
	userRepo, err := NewUserRepo()
	if err != nil {
		return nil, err
	}
	user, err := userRepo.findOne(passcode)
	if err != nil {
		return nil, err
	}
	return user, nil

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

func NewTask(description string, userPasscode string) *Task {
	return &Task{
		ID:           uuid.New(),
		Description:  description,
		CreatedAt:    time.Now(),
		IsComplete:   false,
		UserPasscode: userPasscode,
	}
}

type TaskRepo struct {
	db *mongo.Client
}

func (t *TaskRepo) addTask(description string, userPasscode string) error {
	collection := t.db.Database("todo").Collection("tasks")
	_, err := collection.InsertOne(context.Background(), NewTask(description, userPasscode))
	return err
}

func (t *TaskRepo) findAll(userPasscode string) ([]Task, error) {
	collection := t.db.Database("todo").Collection("tasks")
	cursor, err := collection.Find(context.Background(), bson.D{{"userpasscode", userPasscode}})
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for cursor.Next(context.Background()) {
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func NewTaskRepo() (*TaskRepo, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return &TaskRepo{db: client}, nil
}

// FindAllUserTasks, finds all user tasks
func FindAllUserTasks(userPasscode string) ([]Task, error) {
	repo, err := NewTaskRepo()
	if err != nil {
		return nil, err
	}
	tasks, err := repo.findAll(userPasscode)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func AddTask(description string, userPasscode string) error {
	repo, err := NewTaskRepo()
	if err != nil {
		return err
	}
	err = repo.addTask(description, userPasscode)
	if err != nil {
		return err
	}
	fmt.Println("Tasked added")
	return nil
}
