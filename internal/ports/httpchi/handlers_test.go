package httpchi_test

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
	"github.com/vlasashk/todo-manager/internal/adapters/pgrepo"
	"github.com/vlasashk/todo-manager/internal/models/mocks"
	"github.com/vlasashk/todo-manager/internal/models/todo"
	"github.com/vlasashk/todo-manager/internal/ports/httpchi"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
	service  httpchi.Service
	storage  todo.Repo
	taskReq  todo.TaskReq
	testTask todo.Task
}

func (suite *UnitTestSuite) SetupTest() {
	suite.storage = mocks.NewRepo(suite.T())
	stat := false
	suite.taskReq = todo.TaskReq{
		Title:       "test",
		Description: "test",
		DueDate:     "2024-10-26",
		Status:      &stat,
	}
	suite.testTask = todo.New(suite.taskReq)
	suite.testTask.ID = "test"
}

type TestCase struct {
	testName      string
	storageOutput func()
	expectedCode  int
	expectedResp  string
	urlParamID    string
	reqBody       string
	reqMethod     string
	reqTarget     string
}

func (suite *UnitTestSuite) TestCreateTask() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("CreateTask", suite.taskReq).Return(suite.testTask, nil).Once()
			},
			expectedCode: http.StatusCreated,
			expectedResp: `{"id":"test","title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqMethod:    "POST",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("CreateTask", suite.taskReq).Return(todo.Task{}, errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"param":"id","error":"action fail"}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqMethod:    "POST",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("CreateTask", suite.taskReq).Return(todo.Task{}, errors.New(pgrepo.DateErr)).Once()
			},
			expectedCode: http.StatusConflict,
			expectedResp: `{"param":"date","value":"2024-10-26","error":"bad date"}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqMethod:    "POST",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"param":"date","value":"12345","error":"bad date format"}`,
			reqBody:       `{"title":"test","description":"test","due_date":"12345","status":false}`,
			reqMethod:     "POST",
			reqTarget:     "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedResp:  `{"error":"invalid JSON"}`,
			reqBody:       `{"fail":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqMethod:     "POST",
			reqTarget:     "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"error":"bad JSON"}`,
			reqBody:       `{"fail"}`,
			reqMethod:     "POST",
			reqTarget:     "/task",
		},
	}
	for _, tc := range testCases {
		tc.storageOutput()
		suite.service = httpchi.NewService(suite.storage)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		suite.service.CreateTask(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func (suite *UnitTestSuite) TestGetSingleTask() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("GetTask", "test").Return(suite.testTask, nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"id":"test","title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			urlParamID:   "test",
			reqMethod:    "GET",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("GetTask", "test").Return(todo.Task{}, errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"param":"id","value":"test","error":"action fail"}`,
			urlParamID:   "test",
			reqMethod:    "GET",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("GetTask", "test").Return(todo.Task{}, errors.New(pgrepo.InvalidIdErr)).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"message":"invalid task id"}`,
			urlParamID:   "test",
			reqMethod:    "GET",
			reqTarget:    "/task",
		},
	}
	for _, tc := range testCases {
		tc.storageOutput()
		suite.service = httpchi.NewService(suite.storage)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.urlParamID)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()

		suite.service.GetSingleTask(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func (suite *UnitTestSuite) TestDeleteTask() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("DeleteTask", "test").Return(nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"message":"success"}`,
			urlParamID:   "test",
			reqMethod:    "DELETE",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("DeleteTask", "test").Return(errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"param":"id","value":"test","error":"action fail"}`,
			urlParamID:   "test",
			reqMethod:    "DELETE",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("DeleteTask", "test").Return(errors.New(pgrepo.InvalidIdErr)).Once()
			},
			expectedCode: http.StatusNotFound,
			expectedResp: `{"message":"invalid task id"}`,
			urlParamID:   "test",
			reqMethod:    "DELETE",
			reqTarget:    "/task",
		},
	}
	for _, tc := range testCases {
		tc.storageOutput()
		suite.service = httpchi.NewService(suite.storage)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.urlParamID)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()

		suite.service.DeleteTask(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func (suite *UnitTestSuite) TestUpdateTask() {
	testCases := []TestCase{
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("UpdateTask", suite.taskReq, "test").Return(suite.testTask, nil).Once()
			},
			expectedCode: http.StatusOK,
			expectedResp: `{"id":"test","title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			urlParamID:   "test",
			reqMethod:    "PUT",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("UpdateTask", suite.taskReq, "test").Return(todo.Task{}, errors.New("any err")).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedResp: `{"param":"id","value":"test","error":"action fail"}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			urlParamID:   "test",
			reqMethod:    "PUT",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {
				suite.storage.(*mocks.Repo).On("UpdateTask", suite.taskReq, "test").Return(todo.Task{}, errors.New(pgrepo.DateErr)).Once()
			},
			expectedCode: http.StatusConflict,
			expectedResp: `{"param":"date","value":"2024-10-26","error":"bad date"}`,
			reqBody:      `{"title":"test","description":"test","due_date":"2024-10-26","status":false}`,
			urlParamID:   "test",
			reqMethod:    "PUT",
			reqTarget:    "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"param":"date","value":"12345","error":"bad date format"}`,
			reqBody:       `{"title":"test","description":"test","due_date":"12345","status":false}`,
			urlParamID:    "test",
			reqMethod:     "PUT",
			reqTarget:     "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusUnprocessableEntity,
			expectedResp:  `{"error":"invalid JSON"}`,
			reqBody:       `{"fail":"test","description":"test","due_date":"2024-10-26","status":false}`,
			urlParamID:    "test",
			reqMethod:     "PUT",
			reqTarget:     "/task",
		},
		{
			storageOutput: func() {},
			expectedCode:  http.StatusBadRequest,
			expectedResp:  `{"error":"bad JSON"}`,
			reqBody:       `{"fail"}`,
			urlParamID:    "test",
			reqMethod:     "PUT",
			reqTarget:     "/task",
		},
	}
	for _, tc := range testCases {
		tc.storageOutput()
		suite.service = httpchi.NewService(suite.storage)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.urlParamID)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget, strings.NewReader(tc.reqBody))
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		w := httptest.NewRecorder()

		suite.service.UpdateTask(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func (suite *UnitTestSuite) TestListTasks() {
	type listTestCase struct {
		date   string
		status string
		page   string
		TestCase
	}
	tasks := make([]todo.Task, 0, 2)
	task2 := suite.testTask
	task2.ID += "2"
	task2.Description += "2"
	task2.Title += "2"
	stat := true
	task2.Status = &stat
	tasks = append(tasks, suite.testTask, task2)
	testCases := []listTestCase{
		{
			date:   "",
			status: "",
			page:   "",
			TestCase: TestCase{
				storageOutput: func() {

					suite.storage.(*mocks.Repo).On("ListTasks", uint(0), "", "").Return(tasks, nil).Once()
				},
				expectedCode: http.StatusOK,
				expectedResp: `[{"id":"test","title":"test","description":"test","due_date":"2024-10-26","status":false},{"id":"test2","title":"test2","description":"test2","due_date":"2024-10-26","status":true}]`,
				reqMethod:    "GET",
				reqTarget:    "/tasks?",
			},
		},
	}
	for _, tc := range testCases {
		tc.storageOutput()
		suite.service = httpchi.NewService(suite.storage)
		params := url.Values{}
		params.Add("page", tc.page)
		params.Add("date", tc.date)
		params.Add("status", tc.status)
		req := httptest.NewRequest(tc.reqMethod, tc.reqTarget+params.Encode(), strings.NewReader(tc.reqBody))
		w := httptest.NewRecorder()

		suite.service.ListTasks(w, req)

		body, err := io.ReadAll(w.Body)
		bodyStr := strings.TrimSpace(string(body))
		suite.NoError(err)
		suite.Equal(tc.expectedCode, w.Code)
		suite.Equal(tc.expectedResp, bodyStr)
	}
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
