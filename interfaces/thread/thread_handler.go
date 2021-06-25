package thread

import (
	"fmt"
	"forum/application"
	"forum/domain/entity"
	json "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ThreadInfo struct {
	ThreadApp application.ThreadAppInterface
	userApp application.UserAppInterface
}

func NewThreadInfo(
	ThreadApp application.ThreadAppInterface,
	userApp application.UserAppInterface,
	) *ThreadInfo {
	return &ThreadInfo{
		ThreadApp: ThreadApp,
		userApp: userApp,
	}
}

func (threadInfo *ThreadInfo) HandleCreateThread(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("threadnameOrID")
	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	thread, err := threadInfo.ThreadApp.GetThreadForumAndID(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find post thread by id: %v", slug),
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

	posts := entity.Posts(make([]entity.Post, 0))
	err = json.Unmarshal(ctx.Request.Body(), &posts)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if len(posts) == 0 {
		body, err := json.Marshal(posts)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusCreated)
		ctx.SetBody(body)
		return
	}

	for _, post := range posts {
		_, err = threadInfo.userApp.CheckIfUserExists(post.Author)
		if err != nil {
			msg := entity.Message {
				Text: fmt.Sprintf("Can't find post author by nickname: %v", post.Author),
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
	}

	err = threadInfo.ThreadApp.CreatePosts(thread, posts)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Parent post was created in another thread"),
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

	body, err := json.Marshal(posts)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
	return
}

func (threadInfo *ThreadInfo) HandleGetThreadDetails(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("threadnameOrID")
	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	threads, err := threadInfo.ThreadApp.GetThread(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find thread by slug: %v", slug),
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

	body, err := json.Marshal(threads)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}

func (threadInfo *ThreadInfo) HandleUpdateThread(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("threadnameOrID")
	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err := threadInfo.ThreadApp.CheckThread(slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find thread by slug: %v", slug),
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

	thread := &entity.Thread{}
	err = json.Unmarshal(ctx.Request.Body(), thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if thread.Title == "" && thread.Message == "" {
		thread, err = threadInfo.ThreadApp.GetThread(slug)
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	} else {
		err = threadInfo.ThreadApp.UpdateThread(slug, thread)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
	}

	body, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}

func (threadInfo *ThreadInfo) HandleGetThreadPosts(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("threadnameOrID")
	threadInput := new(entity.Thread)

	switch forumnameInterface.(type) {
	case string:
		slug := forumnameInterface.(string)
		threadInput.Slug = &slug
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}


	err := threadInfo.ThreadApp.CheckThread(*threadInput.Slug)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find thread by slug: %v", *threadInput.Slug),
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

	sortParam := string(queryParams.Peek(string(entity.SortKey)))
	sort := sortParam

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

	posts, err := threadInfo.ThreadApp.GetThreadPosts(*threadInput.Slug, int32(limit), since, sort, desc)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(entity.Posts(posts))
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}

func (threadInfo *ThreadInfo) HandleVoteForThread(ctx *fasthttp.RequestCtx) {
	forumnameInterface := ctx.UserValue("threadnameOrID")
	var slug string
	switch forumnameInterface.(type) {
	case string:
		slug = forumnameInterface.(string)
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	vote := &entity.Vote{}
	err := json.Unmarshal(ctx.Request.Body(), vote)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	vote.Slug = slug
	id, err := strconv.Atoi(slug)
	if err != nil {
		id = 0
	}
	vote.ID = id

	thread, err := threadInfo.ThreadApp.VoteForThread(vote)
	if err != nil {
		msg := entity.Message {
			Text: fmt.Sprintf("Can't find thread by slug: %v", slug),
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

	body, err := json.Marshal(thread)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}
