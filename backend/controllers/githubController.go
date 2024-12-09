package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"area51/schemas"
	"area51/services"
	"area51/toolbox"
)

type GithubController interface {
	RedirectionToGithubService(ctx *gin.Context, path string) (string, error)
	ServiceGithubCallback(ctx *gin.Context, path string) (string, error)
	GetUserInfos(ctx *gin.Context) (userInfos schemas.GithubUserInfo, err error)
}

type githubController struct {
	service 		services.GithubService
	userService 	services.UserService
	serviceToken 	services.TokenService
	servicesService services.ServicesService
}

func NewGithubController(
	service services.GithubService,
	userService services.UserService,
	serviceToken services.TokenService,
	servicesService services.ServicesService,
) GithubController {
	return &githubController{
		service: service,
		userService: userService,
		serviceToken: serviceToken,
		servicesService: servicesService,
	}
}


func (controller *githubController) RedirectionToGithubService(ctx *gin.Context, path string) (string, error) {
	clientId := toolbox.GetInEnv("GITHUB_CLIENT_ID")
	appPort := toolbox.GetInEnv("APP_PORT")
	appAdressHost := toolbox.GetInEnv("APP_HOST_ADDRESS")

	state, err := toolbox.GenerateCSRFToken()
	if err != nil {
		return "", err
	}

	ctx.SetCookie("latestCSRFToken", state, 3600, "/", "localhost", false, true)
	redirectUri := appAdressHost + appPort + path
	authUrl := "https://github.com/login/oauth/authorize" +
		"?client_id=" + clientId +
		"&response_type=code" +
		"&scope=repo" +
		"&redirect_uri=" + redirectUri +
		"&state=" + state
	return authUrl, nil
}

func (controller *githubController) ServiceGithubCallback(ctx *gin.Context, path string) (string, error) {
	var isAlreadyRegistered bool = false
	code := ctx.Query("code")
	if code == "" {
		return "", nil
	}
	state := ctx.Query("state")
	// latestCSRFToken, err := ctx.Cookie("latestCSRFToken")
	if state == "" {
		return "", nil
	}
	// if state != latestCSRFToken {
	// 	return "", nil
	// }
	githubTokenResponse, err := controller.service.AuthGetServiceAccessToken(code, path)
	if err != nil {
		return "", err
	}

	githubService := controller.servicesService.FindByName(schemas.Github)

	userInfo, err := controller.service.GetUserInfo(githubTokenResponse.AccessToken)
	if err != nil {
		return "", fmt.Errorf("unable to get user info because %w", err)
	}
	var actualUser schemas.User
	if userInfo.Email == "" {
		actualUser = controller.userService.GetUserByUsername(userInfo.Login)
		if actualUser.Username != "" {
			isAlreadyRegistered = true
		}
	}
	if userInfo.Email != "" {
		actualUser = controller.userService.GetUserByEmail(userInfo.Email)
	}
	if actualUser.Email != "" {
		isAlreadyRegistered = true
	}
	var newGithubToken schemas.ServiceToken
	var newUser schemas.User
	if isAlreadyRegistered {
		newGithubToken = schemas.ServiceToken{
			Token:   githubTokenResponse.AccessToken,
			Service: githubService,
			UserId:  actualUser.Id,
		}
	} else {
		newUser = schemas.User{
			Username: userInfo.Login,
			Email:    userInfo.Email,
		}
		err := controller.userService.CreateUser(newUser)
		if err != nil {
			return "", fmt.Errorf("unable to create user because %w", err)
		}
		actualUser = controller.userService.GetUserByUsername(userInfo.Login)
		newGithubToken = schemas.ServiceToken{
			Token:   githubTokenResponse.AccessToken,
			RefreshToken: githubTokenResponse.RefreshToken,
			Service: githubService,
			UserId:  actualUser.Id,
		}
		isAlreadyRegistered = true
	}

	tokenId, _ := controller.serviceToken.SaveToken(newGithubToken)

	if newUser.Username == "" {
		newUser = schemas.User{
			Username: userInfo.Login,
			Email:    userInfo.Email,
			TokenId:  tokenId,
		}
	} else {
		tokens, _ := controller.serviceToken.GetTokenByUserId(actualUser.Id)
		for _, token := range tokens {
			if token.UserId == actualUser.Id {
				newUser = schemas.User{
					Username: userInfo.Login,
					Email:    userInfo.Email,
					TokenId: token.Id,
				}
				actualUser.TokenId = token.Id
				err := controller.userService.UpdateUserInfos(actualUser)
				if err != nil {
					return "", fmt.Errorf("unable to update user infos because %w", err)
				}
				break
			}
		}

	}
	// if userInfo.Email == "" {
	// 	newUser.Email = "no email"
	// }

	if isAlreadyRegistered {
		token, _ := controller.userService.Login(newUser)
		// if err != nil {
		// 	return "", fmt.Errorf("unable to login user because %w", err)
		// }
		return token, nil
	} else {
		token, err := controller.userService.Register(newUser)
		if err != nil {
			return "", fmt.Errorf("unable to register user because %w", err)
		}
		return token, nil
	}
}

func (controller *githubController) GetUserInfos(ctx *gin.Context) (userInfos schemas.GithubUserInfo, err error) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len("Bearer "):]

	user, err := controller.userService.GetUserInfos(tokenString)
	if err != nil {
		return schemas.GithubUserInfo{}, err
	}
	token, err := controller.serviceToken.GetTokenById(user.TokenId)
	if err != nil {
		return schemas.GithubUserInfo{}, err
	}

	githubUserInfos, err := controller.service.GetUserInfo(token.Token)
	if err != nil {
		return schemas.GithubUserInfo{}, err
	}
	return githubUserInfos, nil
}
