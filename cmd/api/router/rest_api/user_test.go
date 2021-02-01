package rest_api

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-friend-mgmt/cmd/internal/services/mocks"
	"go-friend-mgmt/cmd/internal/services/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T)  {
	var jsonStr = []byte(`{"email":"andy@example"}`)
	var jsonStr2 = []byte(`{"email":andy@example}`)
	testCase:=[]struct{
		name string
		bodyRequest *bytes.Buffer
		expectedCode int
		expectedBody string
		mockresponse models.Response
		mockError error
	}{
		{
			name:"success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.Response{
				Success: true,
			},
			expectedCode: 200,
			expectedBody: "{\"success\":true}\n",
		},
		{
			name: "create failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/createuser", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("CreateUser", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=CreateUser(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}
