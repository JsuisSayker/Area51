package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	FindActionByName(name string) func(channel chan string, workflowId uint64, actionOption json.RawMessage)
	FindReactionByName(name string) func(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption json.RawMessage)
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

func (service *spotifyService) FindActionByName(name string) func(channel chan string, workflowId uint64, actionOption json.RawMessage) {
	switch name {
	case string(schemas.SpotifyAddTrackAction):
		return service.AddTrackAction
	default:
		return nil
	}
}

func (service *spotifyService) FindReactionByName(name string) func(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption json.RawMessage) {
	switch name {
	case string(schemas.SpotifyAddTrackReaction):
		return service.AddTrackReaction
	case string(schemas.SpotifyCreatePlaylist):
		return service.CreatePlaylist
	default:
		return nil
	}
}

func (service *spotifyService) AddTrackAction(channel chan string, workflowId uint64, actionOption json.RawMessage) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	workflow, err := service.workflowRepository.FindByIds(workflowId)
	if err != nil {
		fmt.Println(err)
		return
	}

	user := service.userService.GetUserById(workflow.UserId)
	accessToken := service.tokenRepository.FindByUserId(user)

	options := schemas.SpotifyActionOptions{}
	err = json.Unmarshal([]byte(actionOption), &options)
	if err != nil {
		fmt.Println(err)
		return
	}
	playlistId := ""
	parts := strings.Split(options.PlaylistURL, "?")
	_, err = fmt.Sscanf(parts[0], "https://open.spotify.com/playlist/%s", &playlistId)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		return
	}

	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistId, nil)
	if err != nil {
		fmt.Printf("unable to create request because: %s", err)
		return
	}
	client := &http.Client{}
	searchedService := service.serviceRepository.FindByName(schemas.Spotify)
	for _, token := range accessToken {
		if token.ServiceId == searchedService.Id {
			request.Header.Set("Authorization", "Bearer "+token.Token)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()
	result := schemas.SpotifyPlaylistInfos{}
	bodyBytes, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println(err)
		return
	}

	existingRecords := map[string]interface{}{}
	if string(workflow.Utils) != "" {
		err = json.Unmarshal([]byte(workflow.Utils), &existingRecords)
		if err != nil {
			fmt.Println("Error unmarshaliing existingRecords: ", err)
			return
		}
	}

	if existingRecords["Tracks"] == nil {
		existingRecords["Tracks"] = 0
		jsonData, err := json.Marshal(existingRecords)
		if err != nil {
			fmt.Println("Error marshalling existingRecords:", err)
			return
		}
		workflow.Utils = jsonData
		service.workflowRepository.Update(workflow)
	}

	var nbTracks uint64
	switch v := existingRecords["Tracks"].(type) {
	case float64:
		nbTracks = uint64(v)
	case uint64:
		nbTracks = v
	default:
		fmt.Println("Error asserting nbTracks to int or float64")
		return
	}
	existingRecords["Tracks"] = result.Tracks.Total
	jsonData, err := json.Marshal(existingRecords)
	if err != nil {
		fmt.Println(err)
		return
	}
	workflow.Utils = jsonData
	if result.Tracks.Total > nbTracks {
		workflow.ReactionTrigger = true
		service.workflowRepository.UpdateReactionTrigger(workflow)
	}
	service.workflowRepository.UpdateUtils(workflow)
	channel <- "Action workflow done"
}

func (service *spotifyService) AddTrackReaction(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption json.RawMessage) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	workflow, err := service.workflowRepository.FindByIds(workflowId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !workflow.ReactionTrigger {
		return
	}

	options := schemas.SpotifyReactionOptions{}
	err = json.Unmarshal([]byte(reactionOption), &options)
	if err != nil {
		fmt.Println(err)
		return
	}

	trackId := ""
	parts := strings.Split(options.TrackURL, "?")
	_, err = fmt.Sscanf(parts[0], "https://open.spotify.com/track/%s", &trackId)
	if err != nil {
		fmt.Printf("unable to create request 1 because: %s", err)
		return
	}
	playlistId := ""
	parts = strings.Split(options.PlaylistURL, "?")
	_, err = fmt.Sscanf(parts[0], "https://open.spotify.com/playlist/%s", &playlistId)
	if err != nil {
		fmt.Printf("unable to create request 2 because: %s", err)
		return
	}

	reqBody := fmt.Sprintf(`{"uris":["spotify:track:%s"],"position":0}`, trackId)
	request, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	searchedService := service.serviceRepository.FindByName(schemas.Spotify)

	for _, token := range accessToken {
		if token.ServiceId == searchedService.Id {
			request.Header.Set("Authorization", "Bearer "+token.Token)
		}
	}

	_, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	workflow.ReactionTrigger = false
	service.workflowRepository.UpdateReactionTrigger(workflow)
}

func (service *spotifyService) CreatePlaylist(channel chan string, workflowId uint64, accessToken []schemas.ServiceToken, reactionOption json.RawMessage) {
	service.mutex.Lock()
	defer service.mutex.Unlock()
	workflow, err := service.workflowRepository.FindByIds(workflowId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !workflow.ReactionTrigger {
		return
	}

	options := schemas.SpotifyPlaylistOptionsSchema{}
	err = json.Unmarshal([]byte(reactionOption), &options)
	if err != nil {
		fmt.Println(err)
		return
	}
	public, err := toolbox.StringToBoolean(options.Public)
	if err != nil {
		fmt.Println(err)
		return
	}
	collaborative, err := toolbox.StringToBoolean(options.Collaborative)
	if err != nil {
		fmt.Println(err)
		return
	}
	trueOptions := schemas.SpotifyPlaylistOptions{
		Name:          options.Name,
		Description:   options.Description,
		Public:        public,
		Collaborative: collaborative,
	}
	optionsJSON, err := json.Marshal(trueOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	request, err := http.NewRequest("POST", "https://api.spotify.com/v1/me/playlists", bytes.NewBuffer(optionsJSON))
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	searchedService := service.serviceRepository.FindByName(schemas.Spotify)

	for _, token := range accessToken {
		if token.ServiceId == searchedService.Id {
			request.Header.Set("Authorization", "Bearer "+token.Token)
		}
	}
	_, err = client.Do(request)
	if err != nil {
		fmt.Println(err)
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
