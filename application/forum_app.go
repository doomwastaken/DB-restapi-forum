package application

import (
	"forum/domain/entity"
	"forum/domain/repository"
)

type ForumApp struct {
	f repository.ForumRepository
}

func NewForumApp(f repository.ForumRepository) *ForumApp {
	return &ForumApp{f: f}
}

type ForumAppInterface interface {
	CreateForum(forumInput *entity.Forum) error
	GetForumDetails(slug string) (*entity.Forum, error)
	GetForumUsers(slug string, limit int32, since string, desc bool) ([]entity.User, error)
	CheckForum(slug string) error
}

func (f *ForumApp) CreateForum(forumInput *entity.Forum) error {
	return f.f.CreateForum(forumInput)
}

func (f *ForumApp) GetForumDetails(slug string) (*entity.Forum, error) {
	return f.f.GetForumDetails(slug)
}

func (f *ForumApp) GetForumUsers(slug string, limit int32, since string, desc bool) ([]entity.User, error) {
	order := "ASC"
	switch desc {
	case true:
		order = "DESC"
	}

	return f.f.GetForumUsers(slug, limit, since, order)
}

func (f *ForumApp) CheckForum(slug string) error {
	return f.f.CheckForum(slug)
}

