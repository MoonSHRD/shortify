package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MoonSHRD/logger"
	"github.com/MoonSHRD/shortify/app"
	"github.com/MoonSHRD/shortify/controllers"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ValidateTokenEndpoint = "/api/v1/validateAccessToken"
)

type AuthMiddleware struct {
	app *app.App
}

func NewAuthMiddleware(a *app.App) *AuthMiddleware {
	return &AuthMiddleware{
		app: a,
	}
}

func (am *AuthMiddleware) ProcessRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		authServerValidationEndpoint := am.app.Config.AuthServerURL + ValidateTokenEndpoint

		// parse access token from http header (Bearer Auth)
		accessToken := req.Header.Get("Authorization")
		if len(accessToken) == 0 {
			err := fmt.Errorf("unauthorized")
			logger.Warningf("Unauthorized access on route %s", req.URL.Path)
			controllers.ReturnHTTPError(res, err.Error(), http.StatusForbidden)
			return
		}
		splitToken := strings.Split(accessToken, " ")
		if len(splitToken) == 0 {
			err := fmt.Errorf("unauthorized")
			logger.Warningf("Unauthorized access on route %s", req.URL.Path)
			controllers.ReturnHTTPError(res, err.Error(), http.StatusForbidden)
			return
		}
		accessToken = splitToken[1]

		validateAccessTokenBodyRequest, err := json.Marshal(map[string]string{
			"accessToken": accessToken,
		})
		resp, err := http.Post(authServerValidationEndpoint, "application/json", bytes.NewBuffer(validateAccessTokenBodyRequest))
		if err != nil {
			logger.Error(err)
			controllers.ReturnHTTPError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Error(err)
				controllers.ReturnHTTPError(res, err.Error(), http.StatusInternalServerError)
				return
			}
			resp.Body.Close()

			var authServerResponse struct {
				Valid bool `json:"isValid"`
			}

			err = json.Unmarshal(bodyBytes, &authServerResponse)
			if err != nil {
				logger.Error(err)
				controllers.ReturnHTTPError(res, err.Error(), http.StatusInternalServerError)
				return
			}

			if authServerResponse.Valid {
				next(res, req)
				return
			} else {
				err := fmt.Errorf("unauthorized")
				logger.Warningf("Unauthorized access on route %s", req.URL.Path)
				controllers.ReturnHTTPError(res, err.Error(), http.StatusForbidden)
				return
			}
		} else {
			err = fmt.Errorf("Incorrect HTTP code from auth server - %d", resp.StatusCode)
			logger.Warningf(err.Error())
			controllers.ReturnHTTPError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
