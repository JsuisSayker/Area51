package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"area51/repository"
	"area51/schemas"
	"area51/toolbox"
)

type SpotifyService interface {
	AuthGetServiceAccessToken(code string, path string) (schemas.SpotifyResponseToken, error)
	// GetUserInfo(accessToken string) (schemas.SpotifyUserInfo, error)
	FindActionByName(name string) func(channel chan string, option string, workflowId uint64, actionOption string)
	FindReactionByName(name string) func(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption string)
	GetUserInfosByToken(accessToken string, serviceName schemas.ServiceName) func(*schemas.ServicesUserInfos)
}

type spotifyService struct {
	userService        UserService
	spotifyRepository  repository.SpotifyRepository
	workflowRepository repository.WorkflowRepository
	actionRepository   repository.ActionRepository
	reactionRepository repository.ReactionRepository
	tokenRepository    repository.TokenRepository
	serviceRepository  repository.ServiceRepository
	mutex              sync.Mutex
}

func NewSpotifyService(
	userService UserService,
	spotifyRepository repository.SpotifyRepository,
	workflowRepository repository.WorkflowRepository,
	actionRepository repository.ActionRepository,
	reactionRepository repository.ReactionRepository,
	tokenRepository repository.TokenRepository,
	serviceRepository repository.ServiceRepository,
) SpotifyService {
	return &spotifyService{
		userService:        userService,
		spotifyRepository:  spotifyRepository,
		workflowRepository: workflowRepository,
		actionRepository:   actionRepository,
		reactionRepository: reactionRepository,
		tokenRepository:    tokenRepository,
		serviceRepository:  serviceRepository,
	}
}

func (service *spotifyService) AuthGetServiceAccessToken(code string, path string) (schemas.SpotifyResponseToken, error) {
	clientId := toolbox.GetInEnv("SPOTIFY_CLIENT_ID")
	clientSecret := toolbox.GetInEnv("SPOTIFY_SECRET")
	appPort := toolbox.GetInEnv("FRONTEND_PORT")
	appAddressHost := toolbox.GetInEnv("APP_HOST_ADDRESS")

	redirectUri := fmt.Sprintf("%s%s/callback", appAddressHost, appPort)
	apiUrl := "https://accounts.spotify.com/api/token"

	options := url.Values{}
	options.Set("code", code)
	options.Set("redirect_uri", redirectUri)
	options.Set("grant_type", "authorization_code")

	header := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret))

	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(options.Encode()))
	if err != nil {
		return schemas.SpotifyResponseToken{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", header)
	client := &http.Client{
		Timeout: time.Second * 45,
	}

	res, err := client.Do(req)
	if err != nil {
		return schemas.SpotifyResponseToken{}, err
	}

	resultToken := schemas.SpotifyResponseToken{}

	err = json.NewDecoder(res.Body).Decode(&resultToken)
	if err != nil {
		return schemas.SpotifyResponseToken{}, err
	}
	res.Body.Close()
	return resultToken, nil
}

// func (service *spotifyService) GetUserInfo(accessToken string) (schemas.SpotifyUserInfo, error) {
// 	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
// 	if err != nil {
// 		return schemas.SpotifyUserInfo{}, err
// 	}

// 	request.Header.Set("Authorization", "Bearer "+accessToken)
// 	client := &http.Client{}

// 	response, err := client.Do(request)
// 	if err != nil {
// 		return schemas.SpotifyUserInfo{}, err
// 	}

// 	result := schemas.SpotifyUserInfo{}
// 	err = json.NewDecoder(response.Body).Decode(&result)
// 	if err != nil {
// 		return schemas.SpotifyUserInfo{}, err
// 	}
// 	response.Body.Close()
// 	return result, nil
// }

func (service *spotifyService) FindActionByName(name string) func(channel chan string, option string, workflowId uint64, actionOption string) {
	switch name {
	case string(schemas.SpotifyAddTrackAction):
		return service.AddTrackAction
	default:
		return nil
	}
}

func (service *spotifyService) FindReactionByName(name string) func(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption string) {
	switch name {
	case string(schemas.SpotifyAddTrackReaction):
		return service.AddTrackReaction
	default:
		return nil
	}
}

func (service *spotifyService) AddTrackAction(channel chan string, option string, workflowId uint64, actionOption string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	workflow, err := service.workflowRepository.FindByIds(workflowId)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}

	user := service.userService.GetUserById(workflow.UserId)
	accessToken := service.tokenRepository.FindByUserId(user)

	options := schemas.SpotifyActionOptions{}
	err = json.NewDecoder(strings.NewReader(option)).Decode(&options)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}
	playlistId := ""
	_, err = fmt.Sscanf(options.PlaylistURL, "https://open.spotify.com/playlist/%s", &playlistId)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		time.Sleep(30 * time.Second)
		return
	}

	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistId, nil)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		time.Sleep(30 * time.Second)
		return
	}
	client := &http.Client{}
	searchedService := service.serviceRepository.FindByName(schemas.Spotify)
	for _, token := range accessToken {
		if token.ServiceId == searchedService.Id {
			request.Header.Set("Authorization", "Bearer "+token.Token)
		}
	}
	// request.Header.Set("Authorization", "Bearer "+accessToken[len(accessToken)-1].Token)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}

	result := schemas.SpotifyPlaylistInfos{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}
	defer response.Body.Close()
	if options.IsOld {
		if options.NbSongs < result.Tracks.Total {
			options.NbSongs = result.Tracks.Total
			workflow.ReactionTrigger = true
			workflow.ActionOptions = toolbox.MustMarshal(options)
			service.workflowRepository.Update(workflow)
		}
	} else {
		fmt.Println("total: ", result.Tracks.Total)
		options.NbSongs = result.Tracks.Total
		options.IsOld = true
		workflow.ActionOptions = toolbox.MustMarshal(options)
		service.workflowRepository.Update(workflow)
	}
	channel <- "Action workflow done"
	time.Sleep(30 * time.Second)
}

func (service *spotifyService) AddTrackReaction(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption string) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	workflow, err := service.workflowRepository.FindByIds(workflowId)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}
	if !workflow.ReactionTrigger {
		time.Sleep(30 * time.Second)
		return
	}

	options := schemas.SpotifyReactionOptions{}
	err = json.NewDecoder(strings.NewReader(workflow.ReactionOptions)).Decode(&options)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}

	trackId := ""
	_, err = fmt.Sscanf(options.TrackURL, "https://open.spotify.com/track/%s", &trackId)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		time.Sleep(30 * time.Second)
		return
	}
	playlistId := ""
	_, err = fmt.Sscanf(options.PlaylistURL, "https://open.spotify.com/playlist/%s", &playlistId)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		time.Sleep(30 * time.Second)
		return
	}

	reqBody := fmt.Sprintf(`{"uris":["spotify:track:%s"],"position":0}`, trackId)
	request, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}

	client := &http.Client{}
	searchedService := service.serviceRepository.FindByName(schemas.Spotify)

	for _, token := range accessToken {
		if token.ServiceId == searchedService.Id {
			request.Header.Set("Authorization", "Bearer "+token.Token)
		}
	}
	// request.Header.Set("Authorization", "Bearer "+accessToken[len(accessToken)-1].Token)

	_, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
		time.Sleep(30 * time.Second)
		return
	}
	workflow.ReactionTrigger = false
	service.workflowRepository.UpdateReactionTrigger(workflow)
}

func (service *spotifyService) GetUserInfosByToken(accessToken string, serviceName schemas.ServiceName) func(*schemas.ServicesUserInfos) {
	return func(userInfos *schemas.ServicesUserInfos) {
		request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
		if err != nil {
			return
		}

		request.Header.Set("Authorization", "Bearer "+accessToken)
		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			return
		}

		err = json.NewDecoder(response.Body).Decode(&userInfos.SpotifyUserInfos)
		if err != nil {
			return
		}
		response.Body.Close()
	}
}
