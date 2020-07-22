package controllers

import (
	"encoding/json"
	"github.com/MoonSHRD/logger"
	"net/http"

	"github.com/MoonSHRD/shortify/app"
	httpModels "github.com/MoonSHRD/shortify/models/http"
	"github.com/MoonSHRD/shortify/services"
)

type LinksController struct {
	app          *app.App
	linksService *services.LinksService
}

func NewLinksController(a *app.App, ls *services.LinksService) *LinksController {
	return &LinksController{
		app:          a,
		linksService: ls,
	}
}

func (lc *LinksController) CreateLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var createLinkRequest httpModels.CreateLinkRequest
	err := decoder.Decode(&createLinkRequest)
	if err != nil {
		logger.Error(err)
		ReturnHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := lc.linksService.Put(&createLinkRequest)
	if err != nil {
		if err == services.ErrEmptyLinkValue {
			ReturnHTTPError(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Error(err)
		ReturnHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, _ := json.Marshal(res)
	w.Write(json)
}

func (lc *LinksController) GetLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	linkIDParam, ok := r.URL.Query()["linkID"]
	if !ok || len(linkIDParam[0]) < 1 {
		logger.Error("Url Param 'linkID' is missing")
		return
	}
	linkID := linkIDParam[0]

	res, err := lc.linksService.GetByLinkID(linkID)
	if err != nil {
		if err == services.ErrNoSuchLink {
			ReturnHTTPError(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Error(err)
		ReturnHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, _ := json.Marshal(res)
	w.Write(json)
}
