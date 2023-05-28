package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"robinhood/model"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Interviews interface {
	AddInterview(ctx context.Context, Interview *model.Interview) error
	ListInterviews(ctx context.Context) ([]model.Interview, error)
	GetInterview(ctx context.Context, id string) (model.Interview, error)
	UpdateInterview(ctx context.Context, id string, interview model.Interview) (int, error)
	DeleteInterview(ctx context.Context, id string) (int, error)
	AddComment(ctx context.Context, id string, comment model.Comment) (int, error)
	UpdateComment(ctx context.Context, id string, cid int, comment model.Comment) (int, error)
	DeleteComment(ctx context.Context, id string, cid int) (int, error)
}

type InterviewsClient struct {
	client        *mongo.Client
	cfg           *viper.Viper
	interviewsCol *mongo.Collection
}

func (c *InterviewsClient) InitInterviews(ctx context.Context) {
	if err := loadInterviewsStaticData(ctx, c.interviewsCol); err != nil {
		log.Fatal(fmt.Errorf("could not insert static data: %w\n", err))
	}
}

func NewInterviewsClient(client *mongo.Client, cfg *viper.Viper) *InterviewsClient {
	return &InterviewsClient{
		client:        client,
		cfg:           cfg,
		interviewsCol: getCollection(cfg, client, "mongodb.dbcollections.interviews"),
	}
}

func (c *InterviewsClient) ListInterviews(ctx context.Context) ([]model.Interview, error) {
	interviews := make([]model.Interview, 0)
	cur, err := c.interviewsCol.Find(ctx, bson.M{})
	if err != nil {
		log.Print(fmt.Errorf("could not get all interviews: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &interviews); err != nil {
		log.Print(fmt.Errorf("could marshall the interviews results: %w", err))
		return nil, err
	}

	return interviews, nil
}

func (c *InterviewsClient) GetInterview(ctx context.Context, id string) (model.Interview, error) {
	var interview model.Interview
	objID, _ := primitive.ObjectIDFromHex(id)
	res := c.interviewsCol.FindOne(ctx, bson.M{"_id": objID})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return interview, nil
		}
		log.Print(fmt.Errorf("error when finding the interview [%s]: %q", id, res.Err()))
		return interview, res.Err()
	}

	if err := res.Decode(&interview); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
		return interview, err
	}
	return interview, nil
}
func (c *InterviewsClient) AddInterview(ctx context.Context, interview *model.Interview) error {
	interview.ID = primitive.NewObjectID()
	interview.Status = "Pending"
	interview.CreateDate = time.Now()
	interview.UpdatedBy = interview.CreateBy
	interview.UpdatedDate = time.Now()
	interview.Comments = make([]model.Comment, 0)
	_, err := c.interviewsCol.InsertOne(ctx, interview)
	if err != nil {
		log.Print(fmt.Errorf("could not add new interview: %w", err))
		return err
	}
	return nil
}
func (c *InterviewsClient) UpdateInterview(ctx context.Context, id string, interview model.Interview) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	interview.UpdatedDate = time.Now()
	res, err := c.interviewsCol.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "subject", Value: interview.Subject},
			{Key: "detail", Value: interview.Detail},
			{Key: "status", Value: interview.Status},
			{Key: "updated_by", Value: interview.UpdatedBy},
			{Key: "updated_date", Value: interview.UpdatedDate},
		},
	}})
	if err != nil {
		log.Print(fmt.Errorf("could not update interview with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}
func (c *InterviewsClient) DeleteInterview(ctx context.Context, id string) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.interviewsCol.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		log.Print(fmt.Errorf("error marshalling the interviews results: %w", err))
		return 0, err
	}

	return int(res.DeletedCount), nil
}
func (c *InterviewsClient) AddComment(ctx context.Context, id string, comment model.Comment) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	comment.CreateDate = time.Now()
	comment.UpdatedDate = time.Now()
	res, err := c.interviewsCol.UpdateOne(ctx, bson.M{"_id": objID}, bson.D{{
		Key: "$push", Value: bson.D{
			{Key: "comments", Value: bson.D{
				{Key: "id", Value: comment.ID},
				{Key: "comment", Value: comment.Comment},
				{Key: "create_by", Value: comment.CreateBy},
				{Key: "create_date", Value: comment.CreateDate},
				{Key: "updated_date", Value: comment.UpdatedDate},
			}},
		},
	}})
	if err != nil {
		log.Print(fmt.Errorf("could not add comment with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func GetInterview(ctx context.Context, id string) {
	panic("unimplemented")
}
func (c *InterviewsClient) UpdateComment(ctx context.Context, id string, cid int, comment model.Comment) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	comment.UpdatedDate = time.Now()
	res, err := c.interviewsCol.UpdateOne(ctx, bson.M{"_id": objID, "comments._id": cid}, bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "comments", Value: bson.D{
				{Key: "comment", Value: comment.Comment},
				{Key: "updated_date", Value: comment.UpdatedDate},
			}},
		},
	}})
	if err != nil {
		log.Print(fmt.Errorf("could not add comment with id [%s]: %w", id, err))
		return 0, err
	}
	return int(res.ModifiedCount), nil
}
func (c *InterviewsClient) DeleteComment(ctx context.Context, id string, cid int) (int, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := c.interviewsCol.UpdateOne(ctx, bson.M{"_id": objID},
		bson.M{"$pull": bson.M{"comments": bson.M{"id": cid}}},
	)
	if err != nil {
		log.Print(fmt.Errorf("could not add comment with id [%s]: %w", id, err))
		return 0, err
	}

	return int(res.ModifiedCount), nil
}

func getCollection(cfg *viper.Viper, client *mongo.Client, colKey string) *mongo.Collection {
	db := cfg.GetString("mongodb.dbname")
	col := cfg.GetString(colKey)

	return client.Database(db).Collection(col)
}

func loadInterviewsStaticData(ctx context.Context, collection *mongo.Collection) error {
	file, err := os.Open("default_data/interviews.json")
	if err != nil {
		return err
	}
	defer file.Close()

	var interviews []model.Interview
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&interviews)
	if err != nil {
		return err
	}
	var b []interface{}
	for _, interview := range interviews {
		b = append(b, interview)
	}
	result, err := collection.InsertMany(ctx, b)
	if err != nil {
		if mongoErr, ok := err.(mongo.BulkWriteException); ok {
			if len(mongoErr.WriteErrors) > 0 && mongoErr.WriteErrors[0].Code == 11000 {
				return nil
			}
		}
		return err
	}
	log.Printf("Inserted interviews: %d\n", len(result.InsertedIDs))

	return nil
}
