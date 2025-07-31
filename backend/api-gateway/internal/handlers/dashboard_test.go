package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
    mock.Mock
}

func (m *MockDashboardService) GetDashboards(userID string) ([]Dashboard, error) {
    args := m.Called(userID)
    return args.Get(0).([]Dashboard), args.Error(1)
}

func TestGetDashboards(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockDashboardService)
    handler := NewDashboardHandler(mockService)
    
    expectedDashboards := []Dashboard{
        {
            ID:     "123",
            UserID: "user123",
            Name:   "My Dashboard",
        },
    }
    
    mockService.On("GetDashboards", "user123").Return(expectedDashboards, nil)
    
    router := gin.New()
    router.GET("/dashboards", func(c *gin.Context) {
        c.Set("user_id", "user123")
        handler.GetDashboards(c)
    })
    
    req, _ := http.NewRequest("GET", "/dashboards", nil)
    resp := httptest.NewRecorder()
    
    router.ServeHTTP(resp, req)
    
    assert.Equal(t, http.StatusOK, resp.Code)
    
    var response []Dashboard
    err := json.Unmarshal(resp.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedDashboards, response)
    
    mockService.AssertExpectations(t)
}

func TestCreateDashboard(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockDashboardService)
    handler := NewDashboardHandler(mockService)
    
    newDashboard := Dashboard{
        Name: "New Dashboard",
        Layout: map[string]interface{}{
            "columns": 2,
            "rows":    3,
        },
    }
    
    mockService.On("CreateDashboard", "user123", mock.Anything).Return(&newDashboard, nil)
    
    router := gin.New()
    router.POST("/dashboards", func(c *gin.Context) {
        c.Set("user_id", "user123")
        handler.CreateDashboard(c)
    })
    
    body, _ := json.Marshal(newDashboard)
    req, _ := http.NewRequest("POST", "/dashboards", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()
    
    router.ServeHTTP(resp, req)
    
    assert.Equal(t, http.StatusCreated, resp.Code)
    mockService.AssertExpectations(t)
}