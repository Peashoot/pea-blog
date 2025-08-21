package model

import (
	"time"
)

type User struct {
	ID          int       `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password_hash"`
	Avatar      *string   `json:"avatar" db:"avatar"`
	Role        string    `json:"role" db:"role"`
	Fingerprint *string   `json:"fingerprint,omitempty" db:"fingerprint"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Article struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Content      string     `json:"content" db:"content"`
	Summary      string     `json:"summary" db:"summary"`
	Tags         []string   `json:"tags" db:"tags"`
	AuthorID     int        `json:"-" db:"author_id"`
	Author       *User      `json:"author,omitempty"`
	Status       string     `json:"status" db:"status"`
	ViewCount    int        `json:"view_count" db:"view_count"`
	LikeCount    int        `json:"like_count" db:"like_count"`
	CommentCount int        `json:"comment_count" db:"comment_count"`
	CoverImage   *string    `json:"cover_image" db:"cover_image"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	PublishedAt  *time.Time `json:"published_at" db:"published_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Comment struct {
	ID          int        `json:"id" db:"id"`
	Content     string     `json:"content" db:"content"`
	AuthorID    int        `json:"-" db:"author_id"`
	Author      *User      `json:"author,omitempty"`
	ArticleID   int        `json:"article_id" db:"article_id"`
	ParentID    *int       `json:"parent_id" db:"parent_id"`
	Replies     []Comment  `json:"replies,omitempty"`
	ReplyCount  int        `json:"reply_count" db:"reply_count"`
	LatestReply *Comment   `json:"latest_reply,omitempty"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Like struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ArticleID int       `json:"article_id" db:"article_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateArticleRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Content     string     `json:"content" binding:"required,min=1,max=104857600"`
	Summary     string     `json:"summary" binding:"required,min=1,max=500"`
	Tags        []string   `json:"tags" binding:"max=10,dive,max=50"`
	Status      string     `json:"status" binding:"required,oneof=draft published scheduled"`
	CoverImage  *string    `json:"cover_image" binding:"omitempty,url"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateArticleRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=200"`
	Content     *string    `json:"content" binding:"omitempty,min=1,max=104857600"`
	Summary     *string    `json:"summary" binding:"omitempty,min=1,max=500"`
	Tags        []string   `json:"tags" binding:"omitempty,max=10,dive,max=50"`
	Status      *string    `json:"status" binding:"omitempty,oneof=draft published scheduled"`
	CoverImage  *string    `json:"cover_image" binding:"omitempty,url"`
	PublishedAt *time.Time `json:"published_at"`
}

type CreateCommentRequest struct {
	Content     string  `json:"content" binding:"required,min=1,max=1000"`
	ArticleID   int     `json:"article_id" binding:"required,min=1"`
	ParentID    *int    `json:"parent_id" binding:"omitempty,min=1"`
	Fingerprint *string `json:"fingerprint" binding:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type SearchParams struct {
	Keyword       string `form:"keyword"`
	Tags          string `form:"tags"`
	Page          int    `form:"page,default=1"`
	PageSize      int    `form:"page_size,default=10"`
	SortBy        string `form:"sort_by,default=created_at"`
	SortOrder     string `form:"sort_order,default=desc"`
	IncludeDrafts bool   `form:"include_drafts,default=false"`
}

type ArticleListResponse struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}

type CommentListResponse struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
