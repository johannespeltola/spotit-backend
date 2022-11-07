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
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	err := user.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	result, err := userdao.Create(&user, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": result.ID, "username": result.Username})
}

func Login(c *gin.Context) {
	var request userdao.UserDAO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	user, err := userdao.GetOne(userdao.UserDAO{Username: request.Username}, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		c.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password.String)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}

	expirationTime := time.Now().Add(72 * time.Hour)
	tokenString, err := middleware.GenerateJWT(user.ID, user.Username, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, tokenString)
}

func AutoLogin(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}
