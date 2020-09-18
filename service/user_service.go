package service

import (
	"github.com/nguyenphucthienan/book-store-oauth-service/utils/crypto_utils"
	"github.com/nguyenphucthienan/book-store-user-service/domain/user"
	"github.com/nguyenphucthienan/book-store-user-service/util/crypto_util"
	"github.com/nguyenphucthienan/book-store-user-service/util/date_util"
	"github.com/nguyenphucthienan/book-store-utils-go/errors"
)

var (
	UserService userServiceInterface = &userService{}
)

type userService struct {
}

type userServiceInterface interface {
	GetUser(int64) (*user.User, errors.RestError)
	CreateUser(user.User) (*user.User, errors.RestError)
	UpdateUser(bool, user.User) (*user.User, errors.RestError)
	DeleteUser(int64) errors.RestError
	SearchUserByStatus(string) (user.Users, errors.RestError)
	LoginUser(user.LoginRequest) (*user.User, errors.RestError)
}

func (s *userService) GetUser(id int64) (*user.User, errors.RestError) {
	existingUser := user.User{Id: id}
	if err := existingUser.Get(); err != nil {
		return nil, err
	}
	return &existingUser, nil
}

func (s *userService) CreateUser(createUser user.User) (*user.User, errors.RestError) {
	if err := createUser.Validate(); err != nil {
		return nil, err
	}

	createUser.Password = crypto_util.GetMd5(createUser.Password)
	createUser.Status = user.StatusActive
	createUser.DateCreated = date_util.GetNowDBFormat()
	if err := createUser.Save(); err != nil {
		return nil, err
	}
	return &createUser, nil
}

func (s *userService) UpdateUser(isPartial bool, updateUser user.User) (*user.User, errors.RestError) {
	existingUser, err := UserService.GetUser(updateUser.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if updateUser.FirstName != "" {
			existingUser.FirstName = updateUser.FirstName
		}
		if updateUser.LastName != "" {
			existingUser.LastName = updateUser.LastName
		}
		if updateUser.Email != "" {
			existingUser.Email = updateUser.Email
		}
	} else {
		existingUser.FirstName = updateUser.FirstName
		existingUser.LastName = updateUser.LastName
		existingUser.Email = updateUser.Email
	}

	if err := existingUser.Validate(); err != nil {
		return nil, err
	}
	if err := existingUser.Update(); err != nil {
		return nil, err
	}
	return existingUser, nil
}

func (s *userService) DeleteUser(id int64) errors.RestError {
	existingUser := user.User{Id: id}
	return existingUser.Delete()
}

func (s *userService) SearchUserByStatus(status string) (user.Users, errors.RestError) {
	dao := &user.User{}
	return dao.SearchByStatus(status)
}

func (s *userService) LoginUser(request user.LoginRequest) (*user.User, errors.RestError) {
	dao := &user.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindUserByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
