package rest_api

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-friend-mgmt/cmd/internal/services/mocks"
	"go-friend-mgmt/cmd/internal/services/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*type MockService struct{
	MockId int64
	MockResponse *models.Response
}
*/
/*func (m MockService) CreateConnectionFriend(friend models.Friend) (int64, *models.Response){
	return m.MockId,m.MockResponse
}*/

func TestCreateConnectionFriend(t *testing.T)  {
	var jsonStr = []byte(`{"friends":["andy@example","john@example.com"]}`)
	var jsonStr2 = []byte(`{"friends":[andy@example]}`)
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
			name: "create failed by CreateConnectionFriend error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.Response{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/makefriend", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("CreateConnectionFriend", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=CreateConnectionFriend(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}

func TestReceiveFriendListByEmail(t *testing.T)  {
	var jsonStr = []byte(`{"email":"andy@example"}`)
	var jsonStr2 = []byte(`{"email":andy@example}`)
	testCase:=[]struct{
		name string
		bodyRequest *bytes.Buffer
		expectedCode int
		expectedBody string
		mockresponse models.ResponseFriend
		mockError error
	}{
		{
			name:"success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.ResponseFriend{
				Success: true,
				Friends: []string{"andy@gmail.com","david@gmail.com"},
				Count: 2,
			},
			expectedCode: 200,
			expectedBody: "{\"success\":true,\"friends\":[\"andy@gmail.com\",\"david@gmail.com\"],\"count\":2}\n",
		},
		{
			name: "create failed by ReceiveFriendListByEmail error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.ResponseFriend{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/friends", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("ReceiveFriendListByEmail", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=ReceiveFriendListByEmail(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}

func TestReceiveCommonFriendList(t *testing.T)  {
	var jsonStr = []byte(`{"friends":["andy@example","john@example.com"]}`)
	var jsonStr2 = []byte(`{"friends":[andy@example,"john@example.com"]}`)
	testCase:=[]struct{
		name string
		bodyRequest *bytes.Buffer
		expectedCode int
		expectedBody string
		mockresponse models.ResponseFriend
		mockError error
	}{
		{
			name:"success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.ResponseFriend{
				Success: true,
				Friends: []string{"andy@gmail.com","david@gmail.com"},
				Count: 2,
			},
			expectedCode: 200,
			expectedBody: "{\"success\":true,\"friends\":[\"andy@gmail.com\",\"david@gmail.com\"],\"count\":2}\n",
		},
		{
			name: "create failed by ReceiveCommonFriendList error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.ResponseFriend{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/commonfriends", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("ReceiveCommonFriendList", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=ReceiveCommonFriendList(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}

func TestSubscribeUpdateFromEmail(t *testing.T)  {
	var jsonStr = []byte(`{"requestor":"andy@example","target":"john@example"}`)
	var jsonStr2 = []byte(`{"requestor":andy@example,"target":"john@example"}`)
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
			name: "create failed by SubscribeUpdateFromEmail error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.Response{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/subscribe", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("SubscribeUpdateFromEmail", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=SubscribeUpdateFromEmail(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}

func TestBlockUpdateFromEmail(t *testing.T)  {
	var jsonStr = []byte(`{"requestor":"andy@example","target":"john@example"}`)
	var jsonStr2 = []byte(`{"requestor":andy@example,"target":"john@example"}`)
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
			name: "create failed by SubscribeUpdateFromEmail error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.Response{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/block", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("BlockUpdateFromEmail", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=BlockUpdateFromEmail(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}

func TestGetAllSubscribeUpdateByEmail(t *testing.T)  {
	var jsonStr = []byte(`{"sender":"andy@example","text":"Hello World! kate@example.com"}`)
	var jsonStr2 = []byte(`{"sender":andy@example,"text":"Hello World! kate@example.com"}`)
	testCase:=[]struct{
		name string
		bodyRequest *bytes.Buffer
		expectedCode int
		expectedBody string
		mockresponse models.SubscribeResponse
		mockError error
	}{
		{
			name:"success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.SubscribeResponse{
				Success: true,
				Recipients: []string{"andy@example","kate@example.com"},
			},
			expectedCode: 200,
			expectedBody: "{\"success\":true,\"recipients\":[\"andy@example\",\"kate@example.com\"]}\n",
		},
		{
			name: "create failed by SubscribeUpdateFromEmail error",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: models.SubscribeResponse{
				Success: false,
			},
			mockError: errors.New("db error"),
			expectedCode: 500,
			expectedBody: "{\"statusCode\":500,\"message\":\"db error\"}\n",
		},
		{
			name: "retrieve failed by incorrect input",
			bodyRequest: bytes.NewBuffer(jsonStr2),
			expectedCode: 400,
			expectedBody: "{\"statusCode\":400,\"message\":\"Bad request\"}\n",
		},
	}
	for _,tt:=range testCase{
		t.Run(tt.name, func(t *testing.T){
			req, err:=http.NewRequest(http.MethodPost,"/user/emailssubscribe", tt.bodyRequest)
			req.Header.Set("X-Custom-Header","myvalue")
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t,err)
			router:=chi.NewRouter()
			req=req.WithContext(context.WithValue(req.Context(),chi.RouteCtxKey,router))
			rr:=httptest.NewRecorder()
			serviceMock:=new(mocks.ServiceMock)
			serviceMock.On("GetAllSubscribeUpdateByEmail", mock.Anything,mock.Anything).Return(tt.mockresponse,tt.mockError)
			handler:=GetAllSubscribeUpdateByEmail(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedCode, rr.Code)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})

	}
}
