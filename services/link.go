package services

import (
	"fmt"
	"github.com/MoonSHRD/shortify/app"
	"github.com/MoonSHRD/shortify/models"
	httpModels "github.com/MoonSHRD/shortify/models/http"
	"github.com/MoonSHRD/shortify/repositories"
	"github.com/MoonSHRD/shortify/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrNoSuchLink     = fmt.Errorf("no such link")
	ErrEmptyLinkValue = fmt.Errorf("empty link value")
)

type LinksService struct {
	app             *app.App
	linksRepository *repositories.LinksRepository
}

func NewLinksService(a *app.App, ur *repositories.LinksRepository) *LinksService {
	return &LinksService{
		app:             a,
		linksRepository: ur,
	}
}

func (ls *LinksService) Put(createLinkRequest *httpModels.CreateLinkRequest) (*models.Link, error) {
	link := &models.Link{}
	if createLinkRequest.TTL == -1 {
		link.LinkID = utils.GenerateAlphanumericString(8)
		link.ExpiresAt = nil
	} else {
		link.LinkID = utils.GenerateAlphanumericString(6)
		now := time.Now().UTC()
		expire := now.Add(time.Duration(createLinkRequest.TTL) * time.Second)
		link.ExpiresAt = &expire
	}
	if createLinkRequest.LinkValue == "" {
		return nil, ErrEmptyLinkValue
	}
	link.LinkValue = createLinkRequest.LinkValue
	err := ls.linksRepository.Put(link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (ls *LinksService) GetByLinkID(linkID string) (*models.Link, error) {
	link, err := ls.linksRepository.GetByLinkID(linkID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchLink
		}
		return nil, err
	}
	return link, nil
}
