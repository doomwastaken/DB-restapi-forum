package post

import (
	"fmt"
	"forum/application"
	"forum/domain/entity"
	json "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

type PostInfo struct {
	PostApp   application.PostAppInterface
	UserApp   application.UserAppInterface
	ThreadApp application.ThreadAppInterface
	ForumApp  application.ForumAppInterface
}

func NewPostInfo(
	PostApp application.PostAppInterface,
	UserApp application.UserAppInterface,
	ThreadApp application.ThreadAppInterface,
	ForumApp application.ForumAppInterface,
) *PostInfo {
	return &PostInfo{
		PostApp:   PostApp,
		UserApp:   UserApp,
		ThreadApp: ThreadApp,
		ForumApp:  ForumApp,
	}
}

func (postInfo *PostInfo) HandleGetPostDetails(ctx *fasthttp.RequestCtx) {
	postIDInterface := ctx.UserValue("postID")
	postID := 0

	var err error
	switch postIDInterface.(type) {
	case string:
		postID, err = strconv.Atoi(postIDInterface.(string))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	post, err := postInfo.PostApp.GetPostDetails(postID)
	if err != nil {
		msg := entity.Message{
			Text: fmt.Sprintf("Can't find post with id: %v", postID),
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

	postInformation := entity.PostOutput{
		Post: post,
	}

	queryParams := ctx.QueryArgs()

	relatedParam := string(queryParams.Peek(string(entity.RelatedKey)))
	related := relatedParam

	if strings.Contains(related, "user") {
		author, err := postInfo.UserApp.GetUserByNickname(post.Author)
		if err != nil {
			msg := entity.Message{
				Text: fmt.Sprintf("Can't find user with id #%v\n", post.Author),
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
		postInformation.Author = author
	}

	if strings.Contains(related, "thread") {
		thread, err := postInfo.ThreadApp.GetThread(strconv.Itoa(post.Thread))
		if err != nil {
			msg := entity.Message{
				Text: fmt.Sprintf("Can't find thread forum by slug: %v", post.Thread),
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
		postInformation.Thread = thread
	}

	if strings.Contains(related, "forum") {
		forum, err := postInfo.ForumApp.GetForumDetails(post.Forum)
		if err != nil {
			msg := entity.Message{
				Text: fmt.Sprintf("Can't find forum by slug: %v", post.Forum),
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
		postInformation.Forum = forum
	}

	body, err := json.Marshal(postInformation)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}

func (postInfo *PostInfo) HandleChangePost(ctx *fasthttp.RequestCtx) {
	postIDInterface := ctx.UserValue("postID")
	postID := 0

	var err error
	switch postIDInterface.(type) {
	case string:
		postID, err = strconv.Atoi(postIDInterface.(string))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	default:
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	post := &entity.Post{}
	err = json.Unmarshal(ctx.Request.Body(), post)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if post.Message == "" {
		post, err = postInfo.PostApp.GetPostDetails(postID)
		if err != nil {
			msg := entity.Message{
				Text: fmt.Sprintf("Can't find post with id: %v", postID),
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

		body, err := json.Marshal(post)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(body)
		return
	}
	post.ID = postID

	post, err = postInfo.PostApp.ChangePostMessage(post)
	if err != nil {
		msg := entity.Message{
			Text: fmt.Sprintf("Can't find post with id: %v", postID),
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

	body, err := json.Marshal(post)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
	return
}
