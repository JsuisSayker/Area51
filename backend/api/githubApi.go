package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"area51/controllers"
	"area51/schemas"
)

type GithubApi struct {
	controller controllers.GithubController
}

func NewGithubApi(controller controllers.GithubController) *GithubApi {
	return &GithubApi{
		controller: controller,
	}
}

func (api *GithubApi) RedirectToGithub(ctx *gin.Context, path string) {
	if authURL, err := api.controller.RedirectionToGithubService(ctx, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, schemas.OAuthConnectionResponse{
			ServiceAuthenticationUrl: authURL,
		})
	}
}

func (api *GithubApi) HandleGithubTokenCallback(ctx *gin.Context, path string) {
	if github_token, err := api.controller.ServiceGithubCallback(ctx, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"access_token": github_token})
	}
}
