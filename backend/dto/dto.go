package dto

import (
	"time"
)

// CreatePost to create blog post
type CreatePost struct {
	Title      string  `json:"title" binding:"required"`
	Content    string  `json:"content" binding:"required"`
	AuthorID   *int64  `json:"author_id,omitempty"`
	AuthorName *string `json:"author_name,omitempty"`
	Slug       string  `json:"slug,omitempty"`
	Published  *bool   `json:"published,omitempty"`
}

// UpdatePost is used for updating an existing post
type UpdatePost struct {
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"content,omitempty"`
	AuthorID   *int64  `json:"author_id,omitempty"`
	AuthorName *string `json:"author_name,omitempty"`
	Slug       *string `json:"slug,omitempty"`
	Published  *bool   `json:"published,omitempty"`
}

// PostResponse represents a post as returned to API clients and templates.
type PostResponse struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Slug       string     `json:"slug,omitempty"`
	AuthorID   *int64     `json:"author_id,omitempty"`
	AuthorName *string    `json:"author_name,omitempty"`
	Published  bool       `json:"published"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// PostsListResponse wraps multiple posts for list endpoints.
type PostsListResponse struct {
	Posts []PostResponse `json:"posts"`
}

// NewPostResponse constructs a PostResponse from raw values.
func NewPostResponse(id int64, title, content, slug string, authorID *int64, authorName *string, published bool, createdAt time.Time, updatedAt *time.Time) PostResponse {
	return PostResponse{
		ID:         id,
		Title:      title,
		Content:    content,
		Slug:       slug,
		AuthorID:   authorID,
		AuthorName: authorName,
		Published:  published,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// ToPostsList converts a slice of PostResponse into PostsListResponse.
func ToPostsList(posts []PostResponse) PostsListResponse {
	return PostsListResponse{
		Posts: posts,
	}
}
