package persist

import (
	"gorm.io/gorm"
)

type Post struct {
	Title   string
	Content string
	Author  int32
	gorm.Model
}
