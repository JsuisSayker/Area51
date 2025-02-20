package services

import (
	"fmt"

	"area51/repository"
	"area51/schemas"
)

type TokenService interface {
	SaveToken(token schemas.ServiceToken) (tokenId uint64, err error)
	Update(token schemas.ServiceToken) error
	Delete(token schemas.ServiceToken) error
	FindAll() ([]schemas.ServiceToken, error)
	GetTokenById(tokenId uint64) (schemas.ServiceToken, error)
	GetTokenByUserId(userId uint64) ([]schemas.ServiceToken, error)
	GetTokenByUserIdAndServiceId(userId uint64, serviceId uint64) (schemas.ServiceToken, error)
}

type tokenService struct {
	repository  repository.TokenRepository
	userService UserService
}

func NewTokenService(
	repository repository.TokenRepository,
	userService UserService,
) TokenService {
	newService := tokenService{
		repository:  repository,
		userService: userService,
	}
	return &newService
}

func (service *tokenService) SaveToken(
	token schemas.ServiceToken,
) (tokenID uint64, err error) {
	tokens := service.repository.FindByToken(token.Token)
	for _, t := range tokens {
		if t.Token == token.Token {
			return t.Id, fmt.Errorf("token already exists")
		}
	}

	service.repository.Save(token)
	tokens = service.repository.FindByToken(token.Token)

	for _, t := range tokens {
		if t.Token == token.Token {
			return t.Id, nil
		}
	}
	return 0, fmt.Errorf("unable to save token")
}

func (service *tokenService) Update(token schemas.ServiceToken) error {
	service.repository.Update(token)
	return nil
}

func (service *tokenService) Delete(token schemas.ServiceToken) error {
	service.repository.Delete(token)
	return nil
}

func (service *tokenService) FindAll() ([]schemas.ServiceToken, error) {
	return service.repository.FindAll(), nil
}

func (service *tokenService) GetTokenById(tokenId uint64) (schemas.ServiceToken, error) {
	return service.repository.FindById(tokenId), nil
}

func (service *tokenService) GetTokenByUserId(userId uint64) ([]schemas.ServiceToken, error) {
	user := service.userService.GetUserById(userId)
	if user.Id == 0 {
		return nil, schemas.ErrUserNotFound
	}
	return service.repository.FindByUserId(user), nil
}

func (service *tokenService) GetTokenByUserIdAndServiceId(
	userId uint64,
	serviceId uint64,
) (schemas.ServiceToken, error) {
	return service.repository.FindByUserIdAndServiceId(userId, serviceId), nil
}
