package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"area51/repository"
	"area51/schemas"
)

type WorkflowService interface {
	FindAll() []schemas.Workflow
	CreateWorkflow(ctx *gin.Context) (string, error)
	InitWorkflow(workflowStartingPoint schemas.Workflow, githubServiceToken []schemas.ServiceToken)
	ExistWorkflow(workflowId uint64) bool
	GetWorkflowByName(name string) schemas.Workflow
	GetMostRecentReaction(ctx *gin.Context) ([]schemas.GithubListCommentsResponse, error)
}

type workflowService struct {
	repository                  repository.WorkflowRepository
	userService                 UserService
	actionService               ActionService
	reactionService             ReactionService
	servicesService             ServicesService
	serviceToken                TokenService
	reactionResponseDataService ReactionResponseDataService
}

func NewWorkflowService(
	repository repository.WorkflowRepository,
	userService UserService,
	actionService ActionService,
	reactionService ReactionService,
	servicesService ServicesService,
	serviceToken TokenService,
	reactionResponseDataService ReactionResponseDataService,
) WorkflowService {
	return &workflowService{
		repository:                  repository,
		userService:                 userService,
		actionService:               actionService,
		reactionService:             reactionService,
		servicesService:             servicesService,
		serviceToken:                serviceToken,
		reactionResponseDataService: reactionResponseDataService,
	}
}

func (service *workflowService) FindAll() []schemas.Workflow {
	return service.repository.FindAll()
}

func (service *workflowService) CreateWorkflow(ctx *gin.Context) (string, error) {
	var result schemas.WorkflowResult
	err := json.NewDecoder(ctx.Request.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len("Bearer "):]

	user, err := service.userService.GetUserInfos(tokenString)
	if err != nil {
		return "", err
	}

	workflowName := result.Name
	workflowValue := "1"
	if workflowName == "" {
		workflowName = "Workflow " + workflowValue
		fmt.Printf("Workflow name found: %+v", service.GetWorkflowByName(workflowName).Name)
		for service.GetWorkflowByName(workflowName).Name != "" {
			fmt.Printf("Workflow name found: %+v", service.GetWorkflowByName(workflowName).Name)
			workflowValueInt, _ := strconv.Atoi(workflowValue)
			workflowValueInt++
			workflowValue = strconv.Itoa(workflowValueInt)
			workflowName = "Workflow " + workflowValue
		}
		workflowValueInt, _ := strconv.Atoi(workflowValue)
		// workflowValueInt++
		// workflowValueInt++
		workflowValue = strconv.Itoa(workflowValueInt)

		workflowName = "Workflow " + workflowValue
	}
	// panic("Not implemented")

	githubServiceToken, _ := service.serviceToken.GetTokenByUserId(user.Id)
	newWorkflow := schemas.Workflow{
		UserId:     user.Id,
		User:       user,
		IsActive:   true,
		ActionId:   result.ActionId,
		ReactionId: result.ReactionId,
		Action:     service.actionService.FindById(result.ActionId),
		Reaction:   service.reactionService.FindById(result.ReactionId),
		Name:       workflowName,
	}
	actualWorkflow := service.repository.FindExistingWorkflow(newWorkflow)
	fmt.Printf("Workflow %+v", actualWorkflow)
	if actualWorkflow != (schemas.Workflow{}) {
		fmt.Print("\nMON TOTO\n")
		if actualWorkflow.IsActive {
			return "Workflow already exists and is active", nil
		} else {
			fmt.Print("OOOOOOUUUUUUUIIIIIIII\n")
			return "Workflow already exists and is not active", nil
		}
	}
	workflowId, err := service.repository.SaveWorkflow(newWorkflow)
	if err != nil {
		return "", err
	}
	newWorkflow.Id = workflowId
	service.InitWorkflow(newWorkflow, githubServiceToken)
	return "Workflow Created succesfully", nil

}

func (service *workflowService) InitWorkflow(workflowStartingPoint schemas.Workflow, githubServiceToken []schemas.ServiceToken) {
	workflowChannel := make(chan string)
	var workflowStateMutex sync.Mutex
	go service.WorkflowActionChannel(workflowStartingPoint, workflowChannel, workflowStateMutex)
	go service.WorkflowReactionChannel(workflowStartingPoint, workflowChannel, githubServiceToken, workflowStateMutex)
}

func (service *workflowService) WorkflowActionChannel(workflowStartingPoint schemas.Workflow, channel chan string, workflowStateMutex sync.Mutex) {
	go func(workflowStartingPoint schemas.Workflow, channel chan string) {
		fmt.Println("Start of WorkflowActionChannel")
		for service.ExistWorkflow(workflowStartingPoint.Id) {
			workflowStateMutex.Lock()
			defer workflowStateMutex.Unlock()
			workflow, err := service.repository.FindByIds(workflowStartingPoint.Id)
			if err != nil {
				fmt.Println("Error")
				return
			}
			action := service.servicesService.FindActionByName(workflow.Action.Name)
			if action == nil {
				fmt.Println("Action not found")
				return
			}
			if workflow.IsActive {
				action(channel, workflow.Action.Name, workflow.Id)
			} else {
				break
			}
			// workflowStateMutex.Unlock()
		}
		fmt.Println("Clear")
		channel <- "Workflow finished"
	}(workflowStartingPoint, channel)
}

func (service *workflowService) WorkflowReactionChannel(workflowStartingPoint schemas.Workflow, channel chan string, githubServiceToken []schemas.ServiceToken, workflowStateMutex sync.Mutex) {
	go func(workflowStartingPoint schemas.Workflow, channel chan string) {
		for service.ExistWorkflow(workflowStartingPoint.Id) {
			workflowStateMutex.Lock()
			defer workflowStateMutex.Unlock()
			workflow, err := service.repository.FindByIds(workflowStartingPoint.Id)
			if err != nil {
				fmt.Println("Error")
				return
			}
			reaction := service.servicesService.FindReactionByName(workflow.Reaction.Name)
			if reaction == nil {
				fmt.Println("Reaction not found")
				return
			}
			if workflow.IsActive {
				result := <-channel
				reaction(channel, workflow.Id, githubServiceToken)
				fmt.Printf("result value: %+v\n", result)
			} else {
				break
			}
		}

	}(workflowStartingPoint, channel)
}

func (service *workflowService) ExistWorkflow(workflowId uint64) bool {
	_, err := service.repository.FindByIds(workflowId)
	return err == nil
}

func (service *workflowService) GetWorkflowByName(name string) schemas.Workflow {
	return service.repository.FindByWorkflowName(name)
}

func (service *workflowService) GetMostRecentReaction(ctx *gin.Context) ([]schemas.GithubListCommentsResponse, error) {
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) <= len("Bearer ") {
		return []schemas.GithubListCommentsResponse{}, fmt.Errorf("no authorization header found")
	}
	tokenString := authHeader[len("Bearer "):]

	user, err := service.userService.GetUserInfos(tokenString)
	if err != nil {
		return nil, err
	}

	workflows := service.repository.FindByUserId(user.Id)
	var reactionResponse []schemas.GithubListCommentsResponse
	for _, workflow := range workflows {
		reactionResponseData := service.reactionResponseDataService.FindByWorkflowId(workflow.Id)
		for _, data := range reactionResponseData {
			err := json.Unmarshal(data.ApiResponse, &reactionResponse)
			if err != nil {
				return nil, err
			}
		}
	}
	return reactionResponse, nil
}
