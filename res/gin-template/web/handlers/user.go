package handlers

import (
    "net/http"
    "{{ .GoModulePath }}/internal/user"

    "github.com/IfanTsai/go-lib/gin/middlewares"
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
)

func CreateUser(c *gin.Context) {
    var req user.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        middlewares.SetErr(c, http.StatusBadRequest, errors.Wrap(err, "invalid request parameters"))

        return
    }

    resp, err := user.Create(&req)
    if err != nil {
        middlewares.SetErrWithTraceBack(c, http.StatusInternalServerError, err)
    }

    middlewares.SetResp(c, resp)
}

func LoginUser(c *gin.Context) {
    var req user.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        middlewares.SetErr(c, http.StatusBadRequest, errors.Wrap(err, "invalid request parameters"))

        return
    }

    resp, err := user.Login(&req, middlewares.GetTokenMaker(c))
    if err != nil {
        middlewares.SetErrWithTraceBack(c, http.StatusInternalServerError, err)
    }

    middlewares.SetResp(c, resp)
}

func GetUser(c *gin.Context) {
    username, err := middlewares.GetUsername(c)
    if err != nil {
        middlewares.SetErrWithTraceBack(c, http.StatusUnauthorized, err)

        return
    }

    resp, err := user.Get(username)
    if err != nil {
        middlewares.SetErrWithTraceBack(c, http.StatusInternalServerError, err)

        return
    }

    middlewares.SetResp(c, resp)
}
