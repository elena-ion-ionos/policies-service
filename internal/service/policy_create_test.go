package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ionos-cloud/policies-service/internal/api"
	"github.com/ionos-cloud/policies-service/internal/controller"
	mocks "github.com/ionos-cloud/policies-service/internal/mocks/port"
	"github.com/ionos-cloud/policies-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_PostPolicies(t *testing.T) {
	mc := gomock.NewController(t)
	writer := api.NewReaderWriter("http://localhost:8080")

	mockPolicyRepo := mocks.NewMockPolicyRepo(mc)
	createPolicyCtrl, err := controller.NewCreatePolicyCtrl(mockPolicyRepo)
	if err != nil {
		t.Fail()
	}
	policiesApi := &PoliciesApi{
		CreatePolicyController: createPolicyCtrl,
		Helper:                 writer,
	}

	tests := []struct {
		name           string
		policy         *model.Policy
		wantErr        bool
		responseStatus int
		expected       func()
	}{
		{
			name: "success",
			policy: &model.Policy{
				Name:      "test1",
				Prefix:    "",
				Action:    "",
				Time:      "",
				Metadata:  nil,
				CreatedAt: time.Time{},
			},
			responseStatus: http.StatusCreated,
			expected: func() {
				mockPolicyRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		tt.expected()
		recorder := httptest.NewRecorder()
		policyBody := map[string]interface{}{
			"properties": map[string]interface{}{
				"name":   tt.policy.Name,
				"prefix": tt.policy.Prefix,
				"action": tt.policy.Action,
				"time":   tt.policy.Time,
			},
		}
		policyBytes, _ := json.Marshal(policyBody)

		bodyReader := bytes.NewBuffer(policyBytes)
		req, err := http.NewRequestWithContext(context.Background(), "POST", "/", bodyReader)
		req.Header.Set("Content-Type", "application/json")

		require.Nil(t, err)

		policiesApi.PostPolicies(recorder, req)

		body, err := io.ReadAll(recorder.Body)
		require.Nil(t, err)
		assert.Equal(t, tt.responseStatus, recorder.Code)
		if tt.responseStatus != http.StatusCreated {
			return
		}
		var response struct {
			Properties struct {
				Name string `json:"name"`
			} `json:"properties"`
		}

		err = json.Unmarshal(body, &response)
		require.Nil(t, err)
		if !tt.wantErr {
			assert.Equal(t, tt.policy.Name, response.Properties.Name)
		}
		require.NotEmpty(t, response.Properties)
	}
}

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
