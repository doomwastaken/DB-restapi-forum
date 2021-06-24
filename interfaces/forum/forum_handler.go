package forum

import (
	"encoding/json"
	"fmt"
	"forum/application"
	"forum/domain/entity"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ForumInfo struct {
	ForumApp  application.ForumAppInterface
	UserApp   application.UserAppInterface
	ThreadApp application.ThreadAppInterface
}

func NewForumInfo(
	ForumApp application.ForumAppInterface,
	UserApp application.UserAppInterface,
	ThreadApp application.ThreadAppInterface,
	) *ForumInfo {
	return &ForumInfo{
		ForumApp:  ForumApp,
		UserApp:   UserApp,
		ThreadApp: ThreadApp,
	}
}

func (forumInfo *ForumInfo) HandleCreateForum(ctx *fasthttp.RequestCtx) {
	forum := &entity.Forum{}

	err := json.Unmarshal(ctx.Request.Body(), forum)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	nickname, err := forumInfo.UserApp.CheckIfUserExists(forum.User)
	if err != nil {
		msg := entity.Message{
			Text: fmt.Sprintf("Can't find user with id #%v\n", forum.User),
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

	forum.User = nickname

	err = forumInfo.ForumApp.CreateForum(forum)
	if err != nil {

		existingForum, err := forumInfo.ForumApp.GetForumDetails(forum.Slug)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}


		body, err := json.Marshal(existingForum)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal(forum)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)

}

func (forumInfo *ForumInfo) HandleGetForumDetails(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	forum, err := forumInfo.ForumApp.GetForumDetails(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find user with id #%v\n", slug),
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

	body, err := json.Marshal(forum)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (forumInfo *ForumInfo) HandleCreateForumThread(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	thread := &entity.Thread{}
	err := json.Unmarshal(ctx.Request.Body(), thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	thread.Forum = slug

	nickname, err := forumInfo.UserApp.CheckIfUserExists(thread.Author)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find user with id #%v\n", thread.Author),
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
	thread.Author = nickname

	err = forumInfo.ThreadApp.CreateThread(thread)
	if err != nil {
		if err == entity.ForumNotExistError {
			msg := entity.Message {
				Text: fmt.Sprintf("Can't find thread forum by slug: %v", thread.Forum),
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

		existedThread, err := forumInfo.ThreadApp.GetThread(*thread.Slug)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(existedThread)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	body, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (forumInfo *ForumInfo) HandleGetForumUsers(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	_, err := forumInfo.ForumApp.CheckForumCase(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find forum by slug: %v", slug),
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
	queryParams := ctx.QueryArgs()

	limitParam := string(queryParams.Peek(string(entity.LimitKey)))
	limit := 0
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	}

	descParam := string(queryParams.Peek(string(entity.DescKey)))
	desc := false
	if descParam == "" {
		desc = false
	} else {
		if descParam == "true" {
			desc = true
		}
	}

	sinceParam := string(queryParams.Peek(string(entity.SinceKey)))
	since := sinceParam

	users, err := forumInfo.ForumApp.GetForumUsers(slug, int32(limit), since, desc)
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
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (forumInfo *ForumInfo) HandleGetForumThreads(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("forumname")

	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	_, err := forumInfo.ForumApp.CheckForumCase(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find forum by slug: %v", slug),
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

	queryParams := ctx.QueryArgs()

	limitParam := string(queryParams.Peek(string(entity.LimitKey)))
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	descParam := string(queryParams.Peek(string(entity.DescKey)))
	desc := false
	if descParam == "" {
		desc = false
	} else {
		if descParam == "true" {
			desc = true
		}
	}

	sinceParam := string(queryParams.Peek(string(entity.SinceKey)))
	since := sinceParam

	threads, err := forumInfo.ThreadApp.GetThreadsByForumSlug(slug, int32(limit), since, desc)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(threads)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
