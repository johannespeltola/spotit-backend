package controllers

import (
	"gin-boilerplate/dao/userdao"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user userdao.UserDAO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}
	err := user.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}
	result, err := userdao.Create(&user, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": result.ID, "username": result.Username})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var request LoginRequest
	var data userdao.UserDAO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}
	data.HashPassword()
	user, err := userdao.GetOne(data, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	expirationTime := time.Now().Add(72 * time.Hour)
	tokenString, err := middleware.GenerateJWT(user.ID, user.Username, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}
	c.SetCookie("session", tokenString, int(expirationTime.Unix()), "", "", false, false)
	c.JSON(http.StatusOK, true)
}
