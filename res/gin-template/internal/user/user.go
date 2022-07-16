package user

import (
    "{{ .GoModulePath }}/internal/db"

    "github.com/IfanTsai/go-lib/config"
    "github.com/IfanTsai/go-lib/user/token"
    "github.com/IfanTsai/go-lib/utils/passwordutils"
    "github.com/pkg/errors"
)

func Create(req *CreateRequest) (*CreateResponse, error) {
    hashedPassword, err := passwordutils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    user, err := db.CreateUser(req.Username, hashedPassword, req.FullName, req.Email)
    if err != nil {
        return nil, err
    }

    return &CreateResponse{
        Response: newUserResponse(user),
    }, nil
}

func Login(req *LoginRequest, tokenMaker token.Maker) (*LoginResponse, error) {
    user, err := db.GetUser(req.Username)
    if err != nil {
        return nil, err
    }

    if err := passwordutils.CheckPassword(req.Password, user.HashedPassword); err != nil {
        return nil, errors.Wrap(err, "invalid password")
    }

    token, tokenPayload, err := tokenMaker.CreateToken(0, user.Username, config.GetTokenAccessDuration())
    if err != nil {
        return nil, errors.WithMessage(err, "cannot create access token")
    }

    return &LoginResponse{
        AccessToken:          token,
        AccessTokenExpiresAt: tokenPayload.ExpiredAt,
    }, nil
}

func Get(username string) (*Response, error) {
    user, err := db.GetUser(username)
    if err != nil {
        return nil, err
    }

    return newUserResponse(user), nil
}

func newUserResponse(user *db.User) *Response {
    return &Response{
        Username:  user.Username,
        FullName:  user.FullName,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}
