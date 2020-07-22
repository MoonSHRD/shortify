package repositories

import (
	"context"
	"github.com/MoonSHRD/shortify/app"
	"github.com/MoonSHRD/shortify/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	LinksCollectionName = "links"
)

type LinksRepository struct {
	app          *app.App
	dbCollection *mongo.Collection
	ctx          context.Context
}

func NewLinksRepository(a *app.App) (*LinksRepository, error) {
	ctx := context.Background()
	dbCol := a.DB.Collection(LinksCollectionName)
	ur := &LinksRepository{
		app:          a,
		dbCollection: dbCol,
		ctx:          ctx,
	}
	err := ur.initMongo()
	if err != nil {
		return nil, err
	}
	return ur, nil
}

func (lr *LinksRepository) initMongo() error {
	linkIDIndex := mongo.IndexModel{
		Keys: bson.M{
			"linkID": 1,
		}, Options: options.Index().SetUnique(true),
	}

	expiresAtIndex := mongo.IndexModel{
		Keys: bson.M{
			"expiresAt": 1,
		}, Options: options.Index().SetExpireAfterSeconds(0),
	}

	exists, err := isCollectionExists(lr.ctx, lr.app.DB, lr.dbCollection.Name())
	if err != nil {
		return err
	}
	if exists {
		_, err := lr.dbCollection.Indexes().DropAll(lr.ctx)
		if err != nil {
			return err
		}
	}

	_, err = lr.dbCollection.Indexes().CreateMany(lr.ctx, []mongo.IndexModel{linkIDIndex, expiresAtIndex})
	return err
}

func (lr *LinksRepository) Put(link *models.Link) error {
	link.ID = primitive.NewObjectID()
	_, err := lr.dbCollection.InsertOne(lr.ctx, link)
	return err
}

func (lr *LinksRepository) GetByID(id primitive.ObjectID) (*models.Link, error) {
	var link models.Link
	err := lr.dbCollection.FindOne(lr.ctx, bson.M{"_id": id}).Decode(&link)
	return &link, err
}

func (lr *LinksRepository) GetByLinkID(linkID string) (*models.Link, error) {
	var link models.Link
	err := lr.dbCollection.FindOne(lr.ctx, bson.M{"linkID": linkID}).Decode(&link)
	return &link, err
}
