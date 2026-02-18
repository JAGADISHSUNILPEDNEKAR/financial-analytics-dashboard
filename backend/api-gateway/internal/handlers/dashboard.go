package handlers

import (
	"time"
)

type Dashboard struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"title"` // mapped to title in request
	Description string                 `json:"description,omitempty"`
	Layout      map[string]interface{} `json:"layout"`
	Widgets     []Widget               `json:"widgets"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Widget struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Title  string                 `json:"title"`
	Config map[string]interface{} `json:"config"`
}

type DashboardHandler struct {
	// dependencies
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) CreateDashboard(req *CreateDashboardRequest) (*Dashboard, error) {
	return &Dashboard{
		ID:        "test-id",
		Name:      req.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type CreateDashboardRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Layout      map[string]interface{} `json:"layout"`
}
