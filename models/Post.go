package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title     string    `json:"title"`
	Alt       string    `json:"alt_title"`
	Chapter   []Chapter `gorm:"foreignKey:PostRefer"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Released  string    `json:"released"`
	Desc      string    `json:"description"`
	Genre     []Genre   `gorm:"many2many:user_genre;"`
	UserRefer uint
}

type Genre struct {
	gorm.Model
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Chapter struct {
	gorm.Model
	Url       []Url `gorm:"foreignKey:ChapterRefer"`
	PostRefer uint
}

type Url struct {
	Link         string
	ChapterRefer uint
}
