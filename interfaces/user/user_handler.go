package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/application"
	"forum/domain/entity"
	"github.com/valyala/fasthttp"
	"net/http"
)

type UserInfo struct {
	userApp application.UserAppInterface
}

func NewUserInfo(userApp application.UserAppInterface) *UserInfo {
	return &UserInfo{
		userApp: userApp,
	}
}

func (userInfo *UserInfo) HandleCreateUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	userInput := new(entity.User)

	switch usernameInterface.(type) {
	case string:
		userInput.Nickname = usernameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(ctx.Request.Body(), userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err = userInfo.userApp.CreateUser(userInput)
	if err != nil {
		users, err := userInfo.userApp.GetUsersWithNicknameAndEmail(userInput.Nickname, userInput.Email)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(users)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal(userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
	return
}

func (userInfo *UserInfo) HandleGetUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	userInput := new(entity.User)

	switch usernameInterface.(type) {
	case string:
		userInput.Nickname = usernameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	profile, err := userInfo.userApp.GetUserByNickname(userInput.Nickname)
	if err != nil {
		msg := entity.Message{
			Text: fmt.Sprintf("Can't find user with id #%v\n", userInput.Nickname),
		}
		body, err := json.Marshal(msg)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal(profile)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}

func (userInfo *UserInfo) HandleUpdateUser(ctx *fasthttp.RequestCtx) {
	usernameInterface := ctx.UserValue("username")
	userInput := new(entity.User)

	switch usernameInterface.(type) {
	case string:
		userInput.Nickname = usernameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(ctx.Request.Body(), userInput)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	profileData, err := userInfo.userApp.UpdateUser(userInput)
	if err != nil {
		var msg entity.Message
		if errors.Is(err, entity.UserDoesntExistsError) {
			msg = entity.Message{
				Text: fmt.Sprintf("Can't find user with id #%v\n", userInput.Nickname),
			}
			body, err := json.Marshal(msg)
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(http.StatusNotFound)
			ctx.SetBody(body)
			return
		} else if errors.Is(err, entity.DataError) {
			emailOwnerNickname, err := userInfo.userApp.GetUserNicknameWithEmail(userInput.Email)
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}
			msg = entity.Message{
				Text: fmt.Sprintf("This email is already registered by user: %v", emailOwnerNickname),
			}
			body, err := json.Marshal(msg)
			if err != nil {
				ctx.SetStatusCode(http.StatusInternalServerError)
				return
			}

			ctx.SetContentType("application/json")
			ctx.SetStatusCode(http.StatusConflict)
			ctx.SetBody(body)
			return
		}
	}

	body, err := json.Marshal(profileData)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
