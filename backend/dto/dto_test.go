package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rohitsh16/go-blog/backend/dto"
)

func TestNewPostResponse_AllFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	updated := now.Add(2 * time.Hour)
	authorID := int64(42)
	authorName := "Alice Engineer"

	resp := dto.NewPostResponse(
		101,
		"Test Article",
		"Content of the article",
		"test-article",
		&authorID,
		&authorName,
		true,
		now,
		&updated,
	)

	if resp.ID != 101 {
		t.Errorf("expected ID 101, got %d", resp.ID)
	}
	if resp.Title != "Test Article" {
		t.Errorf("expected title 'Test Article', got %s", resp.Title)
	}
	if resp.Slug != "test-article" {
		t.Errorf("expected slug 'test-article', got %s", resp.Slug)
	}
	if resp.AuthorID == nil || *resp.AuthorID != 42 {
		t.Errorf("expected AuthorID 42, got %v", resp.AuthorID)
	}
	if resp.AuthorName == nil || *resp.AuthorName != "Alice Engineer" {
		t.Errorf("expected AuthorName 'Alice Engineer', got %v", resp.AuthorName)
	}
	if !resp.Published {
		t.Errorf("expected Published true")
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, resp.CreatedAt)
	}
	if resp.UpdatedAt == nil || !resp.UpdatedAt.Equal(updated) {
		t.Errorf("expected UpdatedAt %v, got %v", updated, resp.UpdatedAt)
	}
}

func TestNewPostResponse_NilFields(t *testing.T) {
	now := time.Now()

	resp := dto.NewPostResponse(
		202,
		"Minimal Post",
		"Some body",
		"",
		nil,
		nil,
		false,
		now,
		nil,
	)

	if resp.AuthorID != nil {
		t.Errorf("expected nil AuthorID, got %v", resp.AuthorID)
	}
	if resp.AuthorName != nil {
		t.Errorf("expected nil AuthorName, got %v", resp.AuthorName)
	}
	if resp.UpdatedAt != nil {
		t.Errorf("expected nil UpdatedAt, got %v", resp.UpdatedAt)
	}
	if resp.Published != false {
		t.Errorf("expected Published false")
	}
}

func TestToPostsList(t *testing.T) {
	now := time.Now()
	p1 := dto.NewPostResponse(1, "Post 1", "Body 1", "post-1", nil, nil, true, now, nil)
	p2 := dto.NewPostResponse(2, "Post 2", "Body 2", "post-2", nil, nil, true, now, nil)

	list := dto.ToPostsList([]dto.PostResponse{p1, p2})

	if len(list.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(list.Posts))
	}
	if list.Posts[0].Title != "Post 1" || list.Posts[1].Title != "Post 2" {
		t.Errorf("unexpected post titles in list: %+v", list.Posts)
	}
}

func TestCreatePostJSON(t *testing.T) {
	jsonData := `{
		"title": "New Blog Post",
		"content": "Deep dive into Go channels",
		"slug": "new-blog-post",
		"author_name": "Bob"
	}`

	var cp dto.CreatePost
	if err := json.Unmarshal([]byte(jsonData), &cp); err != nil {
		t.Fatalf("failed to unmarshal CreatePost: %v", err)
	}

	if cp.Title != "New Blog Post" {
		t.Errorf("expected title 'New Blog Post', got %s", cp.Title)
	}
	if cp.Content != "Deep dive into Go channels" {
		t.Errorf("expected content 'Deep dive into Go channels', got %s", cp.Content)
	}
	if cp.AuthorName == nil || *cp.AuthorName != "Bob" {
		t.Errorf("expected author_name Bob, got %v", cp.AuthorName)
	}
}

func TestUpdatePostJSON(t *testing.T) {
	newTitle := "Updated Title"
	published := false

	up := dto.UpdatePost{
		Title:     &newTitle,
		Published: &published,
	}

	data, err := json.Marshal(up)
	if err != nil {
		t.Fatalf("failed to marshal UpdatePost: %v", err)
	}

	var decoded dto.UpdatePost
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal UpdatePost: %v", err)
	}

	if decoded.Title == nil || *decoded.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %v", decoded.Title)
	}
	if decoded.Published == nil || *decoded.Published != false {
		t.Errorf("expected published false, got %v", decoded.Published)
	}
	if decoded.Content != nil {
		t.Errorf("expected nil content, got %v", decoded.Content)
	}
}
