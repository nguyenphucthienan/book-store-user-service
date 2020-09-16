package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenphucthienan/book-store-user-service/domain/user"
	"github.com/nguyenphucthienan/book-store-user-service/service"
	"github.com/nguyenphucthienan/book-store-user-service/util/errors"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	returnedUser, getErr := service.UserService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, returnedUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func Create(c *gin.Context) {
	var newUser user.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		restErr := errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	createdUser, savedError := service.UserService.CreateUser(newUser)
	if savedError != nil {
		c.JSON(savedError.Status, savedError)
		return
	}

	c.JSON(http.StatusOK, createdUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var updateUser user.User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		restErr := errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	isPartial := c.Request.Method == http.MethodPatch
	updateUser.Id = userId

	updatedUser, savedError := service.UserService.UpdateUser(isPartial, updateUser)
	if savedError != nil {
		c.JSON(savedError.Status, savedError)
		return
	}

	c.JSON(http.StatusOK, updatedUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	if deleteErr := service.UserService.DeleteUser(userId); deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "Deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := service.UserService.SearchUserByStatus(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request user.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("Invalid JSON body")
		c.JSON(restErr.Status, restErr)
		return
	}

	returnedUser, err := service.UserService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, returnedUser.Marshall(c.GetHeader("X-Public") == "true"))
}

func getUserId(userIdParam string) (int64, *errors.RestError) {
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		return 0, errors.NewBadRequestError("User ID should be a number")
	}
	return userId, nil
}
