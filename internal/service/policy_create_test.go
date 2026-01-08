package service

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/ionos-cloud/policies-service/internal/api"
	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_LoadRequestBody(t *testing.T) {
	writer := api.NewReaderWriter("http://localhost:8080")
	policyJSON := `{"properties": {"action": "test","name": "aa","prefix": "e","time": "12"}}`
	invalidJSON := `{"properties": "action": "test","name": "aa","prefix": "e","time": "12"}}`

	tests := []struct {
		name    string
		want    *model.Policy
		request *http.Request
	}{
		{
			name: "no body provided",
			want: nil,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "", nil)
				return req
			}(),
		},
		{
			name: "body provided",
			want: &model.Policy{
				ID:        uuid.UUID{},
				Name:      "aa",
				Prefix:    "",
				Action:    "",
				Time:      "",
				CreatedAt: time.Time{},
			},
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(policyJSON)))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
		},
		{
			name: "invalid body provided",
			want: nil,
			request: func() *http.Request {
				req, _ := http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(invalidJSON)))
				req.Header.Set("Content-Type", "application/json")
				return req
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PoliciesApi{Helper: writer}
			body, _ := p.loadRequestBody(context.Background(), tt.request)
			if tt.want != nil {
				assert.Equal(t, tt.want.Name, body.Name)
			} else {
				assert.Nil(t, body)
			}
		})
	}
}
