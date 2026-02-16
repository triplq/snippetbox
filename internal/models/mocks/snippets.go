package mocks

import (
	"time"

	"github.com/triplq/snippetbox/internal/models"
)

var MockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now().Add(1 * time.Hour),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return MockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{MockSnippet}, nil
}
