package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(title string, text string, author int32) (Post, error)
	Update(id int, title string, text string, author int32) (Post, error)
	Delete(id int) (Post, error)
	Get(id int) (Post, error)
	Index(page int, per_page int) ([]Post, error)
}

type Repository struct {
	Connection *gorm.DB
}

func (rep *Repository) Create(title string, text string, author int32) (Post, error) {
	post := &Post{Title: title, Content: text, Author: author}
	save := rep.Connection.Create(&post)
	if save.Error != nil {
		return *post, save.Error
	}
	return *post, nil
}

func (rep *Repository) Update(id int, title string, text string, author int32) (Post, error) {
	post := &Post{}
	rep.Connection.First(post, id)
	if post.ID == 0 {
		return *post, errors.New(fmt.Sprintf("The post with %d id can not be found", id))
	}

	post.Title = title
	post.Content = text
	post.Author = author

	rep.Connection.Save(post)

	return *post, nil
}

func (rep *Repository) Delete(id int) (Post, error) {
	post := &Post{}
	rep.Connection.First(post, id)
	if post.ID == 0 {
		return *post, errors.New(fmt.Sprintf("Post with the id %d can not be found\n", id))
	}

	exec := rep.Connection.Delete(post)
	if exec.Error != nil {
		return *post, exec.Error
	}

	return *post, nil
}

func (rep *Repository) Get(id int) (Post, error) {
	post := &Post{}
	rep.Connection.First(post, id)

	return *post, nil
}

func (rep *Repository) Index(page int, per_page int) ([]Post, error) {
	var posts []Post

	var offset int = (page - 1) * per_page

	res := rep.Connection.Offset(offset).Limit(per_page).Find(&posts, []string{})
	if res.Error != nil {
		return []Post{}, res.Error
	}

	return posts, nil
}
