package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github-real-lb/tutor-management-web/db/mocks"
	db "github.com/github-real-lb/tutor-management-web/db/sqlc"
	"github.com/github-real-lb/tutor-management-web/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLessonSubjectAPIs(t *testing.T) {
	tests := tests{
		"Test_createLessonSubjectAPI": createLessonSubjectTestCasesBuilder(),
		"Test_getLessonSubject":       getLessonSubjectTestCasesBuilder(),
		"Test_listLessonSubjects":     listLessonSubjectsTestCasesBuilder(),
		"Test_updateLessonSubjects":   updateLessonSubjectTestCasesBuilder(),
	}

	for key, tcs := range tests {
		t.Run(key, func(t *testing.T) {
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					// start mock db and build the stub
					mockStore := mocks.NewMockStore(t)
					tc.buildStub(mockStore)

					// send test request to server
					recorder := tc.sendRequestToServer(t, mockStore)

					// check response
					tc.checkResponse(t, mockStore, recorder)
				})
			}

		})
	}
}

// randomLessonSubject creates a new random LessonSubject struct.
func randomLessonSubject() db.LessonSubject {
	return db.LessonSubject{
		SubjectID: util.RandomInt64(1, 1000),
		Name:      util.RandomName(),
	}
}

// createLessonSubjectTestCasesBuilder creates a slice of test cases for the createLessonSubject API
func createLessonSubjectTestCasesBuilder() testCases {
	var testCases testCases

	lessonSubject := randomLessonSubject()

	arg := struct {
		Name string `json:"name"`
	}{
		Name: lessonSubject.Name,
	}

	methodName := "CreateLessonSubject"
	url := "/lesson_subjects"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, lessonSubject.Name).
				Return(lessonSubject, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonSubject)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodPost,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(db.LessonSubject{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid Body Data response by passing no arguments
	testCases = append(testCases, testCase{
		name:       "Invalid Body Data",
		httpMethod: http.MethodPost,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// getLessonSubjectTestCasesBuilder creates a slice of test cases for the getLessonSubject API
func getLessonSubjectTestCasesBuilder() testCases {
	var testCases testCases

	lessonSubject := randomLessonSubject()
	id := lessonSubject.SubjectID
	methodName := "GetLessonSubject"
	url := fmt.Sprintf("/lesson_subjects/%d", id)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(lessonSubject, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonSubject)

		},
	})

	// create a test case for Not Found response
	testCases = append(testCases, testCase{
		name:       "Not Found",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, id).
				Return(db.LessonSubject{}, sql.ErrNoRows).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNotFound, recorder.Code)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(db.LessonSubject{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid ID response by passing url with id=0
	testCases = append(testCases, testCase{
		name:       "Invalid ID",
		httpMethod: http.MethodGet,
		url:        "/lesson_subjects/0",
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// listLessonSubjectsTestCasesBuilder creates a slice of test cases for the listLessonSubjects API
func listLessonSubjectsTestCasesBuilder() testCases {
	var testCases testCases

	n := 5
	lessonSubjects := make([]db.LessonSubject, n)
	for i := 0; i < n; i++ {
		lessonSubjects[i] = randomLessonSubject()
	}

	arg := db.ListLessonSubjectsParams{
		Limit:  int32(n),
		Offset: 0,
	}

	methodName := "ListLessonSubjects"
	url := fmt.Sprintf("/lesson_subjects?page_id=%d&page_size=%d", 1, n)

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodGet,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(lessonSubjects, nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchStruct(t, recorder.Body, lessonSubjects)

		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodGet,
		url:        url,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return([]db.LessonSubject{}, sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		},
	})

	// create a test case for Invalid PageID response by passing url with page_id=-1
	testCases = append(testCases, testCase{
		name:       "Invalid Page_ID Parameter",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/lesson_subjects?page_id=%d&page_size=%d", -1, n),
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	// create a test case for Invalid PageSize response by passing url with page_size=10000
	testCases = append(testCases, testCase{
		name:       "Invalid Page_Size Parameter",
		httpMethod: http.MethodGet,
		url:        fmt.Sprintf("/lesson_subjects?page_id=%d&page_size=%d", 1, 10000),
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()

		},
	})

	return testCases
}

// updateLessonSubjectTestCasesBuilder creates a slice of test cases for the updateLessonSubject API
func updateLessonSubjectTestCasesBuilder() testCases {
	var testCases testCases

	arg := db.UpdateLessonSubjectParams{
		SubjectID: util.RandomInt64(1, 1000),
		Name:      util.RandomName(),
	}

	methodName := "UpdateLessonSubject"
	url := "/lesson_subjects"

	// create a test case for StatusOK response
	testCases = append(testCases, testCase{
		name:       "OK",
		httpMethod: http.MethodPut,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, arg).
				Return(nil).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
		},
	})

	// create a test case for Internal Server Error response
	testCases = append(testCases, testCase{
		name:       "Internal Error",
		httpMethod: http.MethodPut,
		url:        url,
		body:       arg,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).
				Return(sql.ErrConnDone).
				Once()
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		},
	})

	// create a test case for Invalid Body Data response by passing no arguments
	testCases = append(testCases, testCase{
		name:       "Invalid Body Data",
		httpMethod: http.MethodPut,
		url:        url,
		body:       nil,
		buildStub: func(mockStore *mocks.MockStore) {
			mockStore.On(methodName, mock.Anything, mock.Anything).Times(0)
		},
		checkResponse: func(t *testing.T, mockStore *mocks.MockStore, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			mockStore.On(methodName, mock.Anything, mock.Anything).Unset()
		},
	})

	return testCases
}
