package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const methodNotAllowedMsg = "Method not allowed"
const emptyBodyMsg = "Body Required"
const contentNotJsonMsg = "Content-Type must be \"" + ContentTypeJSON + "\""
const unexpectedEOFMsg = "Bad Input(JSON) format, unexpected EOF"

const contentNotUrlEncMsg = "Content-Type must be \"" + ContentTypeURLENCODED + "\""
const multipleKeysMsg = "Multiple values detected, the only valid input is 1 \"ID=X\""
const keyNotIdMsg = "Invalid key detected, the only valid input is 1 \"ID=X\""
const multipleIdMsg = "Only one ID is allowed, and it must be an int"
const idNotIntMsg = "ID must be a positive int"

type testReqArgs struct {
	w       *httptest.ResponseRecorder
	payload string
	method  string
	ctype   string
}

func prepareTestRequest(args testReqArgs, path string) *http.Request {
	var req *http.Request
	if args.payload != "nil" {
		payload := strings.NewReader(args.payload)
		req = httptest.NewRequest(args.method, path, payload)
	} else {
		req = httptest.NewRequest(args.method, path, nil)
	}
	if args.ctype != "" {
		req.Header.Add("Content-Type", args.ctype)
	}
	return req
}

func Test_statusHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		args   args
		status int
	}{
		// TODO: Add test cases.
		{"MethodGET", args{httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/status", nil)}, http.StatusOK},
		{"MethodNotGET", args{httptest.NewRecorder(), httptest.NewRequest(http.MethodPut, "/status", nil)}, http.StatusMethodNotAllowed},
	}
	handler := http.HandlerFunc(statusHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.ServeHTTP(tt.args.w, tt.args.r)
			if tt.args.w.Code != tt.status {
				t.Fatalf("(Expected) %d != %d (Returned)", tt.status, tt.args.w.Code)
			}
		})
	}
}

func Test_carsHandler(t *testing.T) {
	const missingIdOrSeatsMsg = "Bad Input(JSON) format, id and seats are required for all cars"
	const idIs0Msg = "Bad Input(JSON) format, id must be different from 0"
	const idNotMinSeatsMsg = "Bad Input(JSON) format, seats must be > 3"
	const idNotMaxSeatsMsg = "Bad Input(JSON) format, seats must be < 7"
	const idRepeatedMsg = "Bad Input(JSON) format, cars Ids must be unique"
	tests := []struct {
		name   string
		args   testReqArgs
		status int
		tstMsg string
	}{
		// TODO: Add test cases.
		{"MethodNotPUT", testReqArgs{httptest.NewRecorder(), "nil", http.MethodGet, ""}, http.StatusMethodNotAllowed, methodNotAllowedMsg},
		{"MethodPUTEmptyBody", testReqArgs{httptest.NewRecorder(), "nil", http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, emptyBodyMsg},
		{"MethodPUTNotJSON", testReqArgs{httptest.NewRecorder(), `[ { "id": 2, "seats": 4 } ]`, http.MethodPut, ContentTypeURLENCODED}, http.StatusBadRequest, contentNotJsonMsg},
		{"MethodPUTMissId", testReqArgs{httptest.NewRecorder(), `[ { "seats": 4 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, missingIdOrSeatsMsg},
		{"MethodPUTMissSeats", testReqArgs{httptest.NewRecorder(), `[ { "id": 4 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, missingIdOrSeatsMsg},
		{"MethodPUTId0", testReqArgs{httptest.NewRecorder(), `[ { "id": 0, "seats": 4 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, idIs0Msg},
		{"MethodPUTSeatBelow1", testReqArgs{httptest.NewRecorder(), `[ { "id": 4, "seats": 0 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, idNotMinSeatsMsg},
		{"MethodPUTSeatOver6", testReqArgs{httptest.NewRecorder(), `[ { "id": 4, "seats": 7 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, idNotMaxSeatsMsg},
		{"MethodPUTIdRepeated", testReqArgs{httptest.NewRecorder(), `[ { "id": 4, "seats": 4 },{ "id": 4, "seats": 5 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, idRepeatedMsg},
		{"MethodPUTInvalidJson", testReqArgs{httptest.NewRecorder(), `[`, http.MethodPut, ContentTypeJSON}, http.StatusBadRequest, unexpectedEOFMsg},
		{"MethodPUT", testReqArgs{httptest.NewRecorder(), `[ { "id": 2, "seats": 4 } ]`, http.MethodPut, ContentTypeJSON}, http.StatusOK, ""},
	}
	startStorage()
	handler := http.HandlerFunc(carsHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := prepareTestRequest(tt.args, "/cars")
			handler.ServeHTTP(tt.args.w, req)
			if tt.args.w.Code != tt.status {
				t.Fatalf("(Expected) %d != %d (Returned)", tt.status, tt.args.w.Code)
			}
			if tt.args.w.Body.String() != tt.tstMsg {
				t.Fatalf("(Expected) %s != %s (Returned)", tt.tstMsg, tt.args.w.Body.String())
			}
		})
	}
}

func Test_journeyHandler(t *testing.T) {
	const MinPeopleMsg = "Bad Input(JSON) format, number of people should be between 1 or 6"
	const MaxPeopleMsg = "Bad Input(JSON) format, number of people should be 6 at most"
	const IdOrPeopleMissMsg = "Bad Input(JSON) format, id and people are required for all groups"
	const repeatedIdMsg = "Error, group Id already exists"
	tests := []struct {
		name   string
		args   testReqArgs
		status int
		tstMsg string
	}{
		// TODO: Add test cases.
		{"MethodNotPost", testReqArgs{httptest.NewRecorder(), "nil", http.MethodGet, ""}, http.StatusMethodNotAllowed, methodNotAllowedMsg},
		{"MethodPostEmptyBody", testReqArgs{httptest.NewRecorder(), "nil", http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, emptyBodyMsg},
		{"MethodPostNotJSON", testReqArgs{httptest.NewRecorder(), `{ "id": 2, "people": 4 }`, http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, contentNotJsonMsg},
		{"MethodPostMissId", testReqArgs{httptest.NewRecorder(), `{ "people": 4 }`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, IdOrPeopleMissMsg},
		{"MethodPostMissPeople", testReqArgs{httptest.NewRecorder(), `{ "id": 2 }`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, IdOrPeopleMissMsg},
		{"MethodPostNoPeople", testReqArgs{httptest.NewRecorder(), `{ "id": 2, "people": 0 }`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, MinPeopleMsg},
		{"MethodPostManyPeople", testReqArgs{httptest.NewRecorder(), `{ "id": 2, "people": 7 }`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, MaxPeopleMsg},
		{"MethodPostRepeatedId", testReqArgs{httptest.NewRecorder(), `{ "id": 2, "people": 5 }`, http.MethodPost, ContentTypeJSON}, http.StatusInternalServerError, repeatedIdMsg},
		{"MethodPostInvalidJson", testReqArgs{httptest.NewRecorder(), `{`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, unexpectedEOFMsg},
		{"MethodPost", testReqArgs{httptest.NewRecorder(), `{ "id": 2, "people": 4 }`, http.MethodPost, ContentTypeJSON}, http.StatusAccepted, ""},
	}
	startStorage()
	handler := http.HandlerFunc(journeyHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "MethodPostRepeatedId" {
				simulateTestCall(t, reqArgs{`{ "id": 2, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			}
			req := prepareTestRequest(tt.args, "/journey")
			handler.ServeHTTP(tt.args.w, req)
			if tt.args.w.Code != tt.status {
				t.Fatalf("(Expected) %d != %d (Returned)", tt.status, tt.args.w.Code)
			}
			if tt.args.w.Body.String() != tt.tstMsg {
				t.Fatalf("(Expected) %s != %s (Returned)", tt.tstMsg, tt.args.w.Body.String())
			}
			if tt.name == "MethodPostRepeatedId" {
				simulateTestCall(t, reqArgs{"ID=2", "POST", "/dropoff", dropoffHandler, ContentTypeURLENCODED})
			}
		})
	}
}

func Test_dropoffHandler(t *testing.T) {
	tests := []struct {
		name   string
		args   testReqArgs
		status int
		tstMsg string
	}{
		// TODO: Add test cases.

		{"NotPost", testReqArgs{httptest.NewRecorder(), "nil", http.MethodGet, ""}, http.StatusMethodNotAllowed, methodNotAllowedMsg},
		{"PostEmptyBody", testReqArgs{httptest.NewRecorder(), "nil", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, emptyBodyMsg},
		{"PostNotUrlEnc", testReqArgs{httptest.NewRecorder(), `{}`, http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, contentNotUrlEncMsg},
		{"PostMultipleKeys", testReqArgs{httptest.NewRecorder(), "ID=7&IDX=8", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, multipleKeysMsg},
		{"PostKeyIsNotId", testReqArgs{httptest.NewRecorder(), "IDX=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, keyNotIdMsg},
		{"PostMultipleId", testReqArgs{httptest.NewRecorder(), "ID=7&ID=8", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, multipleIdMsg},
		{"PostIdNotInt", testReqArgs{httptest.NewRecorder(), "ID=X", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, idNotIntMsg},
		{"PostNonexistentGroup", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusNotFound, ""},
		{"PostGroupWithoutCar", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusNoContent, ""},
		{"PostGroupWithCar", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusOK, ""},
		{"PostGroupWithCarAssign", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusOK, ""},
		{"PostGroupWithCarAssignRemoveDropped", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusOK, ""},
	}
	startStorage()
	handler := http.HandlerFunc(dropoffHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "PostGroupWithoutCar" {
				simulateTestCall(t, reqArgs{`[]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 5 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			} else if tt.name == "PostGroupWithCar" {
				simulateTestCall(t, reqArgs{`[ { "id": 3, "seats": 5 } ]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 5 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			} else if tt.name == "PostGroupWithCarAssign" {
				simulateTestCall(t, reqArgs{`[ { "id": 3, "seats": 5 } ]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 4, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			} else if tt.name == "PostGroupWithCarAssignRemoveDropped" {
				simulateTestCall(t, reqArgs{`[ { "id": 3, "seats": 5 } ]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
				// add 2 waiting groups
				simulateTestCall(t, reqArgs{`{ "id": 4, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 5, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
				// drop a waiting group, waiting list is not updated until we try to assign a new waiting group
				simulateTestCall(t, reqArgs{"ID=4", "POST", "/dropoff", dropoffHandler, ContentTypeURLENCODED})
			}
			req := prepareTestRequest(tt.args, "/dropoff")
			handler.ServeHTTP(tt.args.w, req)
			if tt.args.w.Code != tt.status {
				t.Fatalf("(Expected) %d != %d (Returned)", tt.status, tt.args.w.Code)
			}
			if tt.args.w.Body.String() != tt.tstMsg {
				t.Fatalf("(Expected) %s != %s (Returned)", tt.tstMsg, tt.args.w.Body.String())
			}
			if tt.name == "PostGroupWithCarAssign" {
				w := simulateTestCall(t, reqArgs{"ID=4", "POST", "/locate", locateHandler, ContentTypeURLENCODED})
				expected := `{ "id": 3, "seats": 5 }`
				if w.Body.String() != expected {
					t.Fatalf("(Expected) %s != %s (Returned)", expected, w.Body.String())
				}
			} else if tt.name == "PostGroupWithCarAssignRemoveDropped" {
				w := simulateTestCall(t, reqArgs{"ID=5", "POST", "/locate", locateHandler, ContentTypeURLENCODED})
				expected := `{ "id": 3, "seats": 5 }`
				if w.Body.String() != expected {
					t.Logf("car %v", journeysMap[5])
					t.Fatalf("(Expected) %s != %s (Returned)", expected, w.Body.String())
				}
				if len(waitingGroups) != 0 {
					t.Fatalf("Waiting group was not updated")
				}
			}
		})
	}

}

func Test_locateHandler(t *testing.T) {
	tests := []struct {
		name   string
		args   testReqArgs
		status int
		tstMsg string
	}{
		// TODO: Add test cases.
		{"MethodNotPost", testReqArgs{httptest.NewRecorder(), "nil", http.MethodGet, ""}, http.StatusMethodNotAllowed, methodNotAllowedMsg},
		{"MethodPostEmptyBody", testReqArgs{httptest.NewRecorder(), "nil", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, emptyBodyMsg},
		{"MethodPostNotUrlEnc", testReqArgs{httptest.NewRecorder(), "{}", http.MethodPost, ContentTypeJSON}, http.StatusBadRequest, contentNotUrlEncMsg},
		{"MethodPostMultipleKeys", testReqArgs{httptest.NewRecorder(), "ID=7&IDX=8", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, multipleKeysMsg},
		{"MethodPostKeyIsNotId", testReqArgs{httptest.NewRecorder(), "IDX=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, keyNotIdMsg},
		{"MethodPostMultipleId", testReqArgs{httptest.NewRecorder(), "ID=7&ID=8", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, multipleIdMsg},
		{"MethodPostIdNotInt", testReqArgs{httptest.NewRecorder(), "ID=X", http.MethodPost, ContentTypeURLENCODED}, http.StatusBadRequest, idNotIntMsg},
		{"MethodPostNonexistentGroup", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusNotFound, ""},
		{"MethodPostGroupWithoutCar", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusNoContent, ""},
		{"MethodPostGroupToSameSizeCar", testReqArgs{httptest.NewRecorder(), "ID=7", http.MethodPost, ContentTypeURLENCODED}, http.StatusOK, "{ \"id\": 3, \"seats\": 5 }"},
		{"MethodPostGroupToDiffSizeCar", testReqArgs{httptest.NewRecorder(), "ID=8", http.MethodPost, ContentTypeURLENCODED}, http.StatusOK, "{ \"id\": 5, \"seats\": 5 }"},
	}
	startStorage()
	handler := http.HandlerFunc(locateHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "MethodPostGroupWithoutCar" {
				simulateTestCall(t, reqArgs{`[]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 3 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			} else if tt.name == "MethodPostGroupToSameSizeCar" {
				simulateTestCall(t, reqArgs{`[ { "id": 3, "seats": 5 } ]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 7, "people": 5 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			} else if tt.name == "MethodPostGroupToDiffSizeCar" {
				simulateTestCall(t, reqArgs{`[ { "id": 5, "seats": 5 } ]`, "PUT", "/cars", carsHandler, ContentTypeJSON})
				simulateTestCall(t, reqArgs{`{ "id": 8, "people": 4 }`, "POST", "/journey", journeyHandler, ContentTypeJSON})
			}
			req := prepareTestRequest(tt.args, "/locate")
			handler.ServeHTTP(tt.args.w, req)
			if tt.args.w.Code != tt.status {
				t.Fatalf("(Expected) %d != %d (Returned)", tt.status, tt.args.w.Code)
			}
			if tt.args.w.Body.String() != tt.tstMsg {
				t.Fatalf("(Expected) %s != %s (Returned)", tt.tstMsg, tt.args.w.Body.String())
			}
			if tt.args.w.Header().Get("Content-Type") != ContentTypeJSON {
				t.Fatalf("(Expected) %s != %s (Returned)", tt.args.w.Header().Get("Content-Type"), ContentTypeJSON)
			}
		})
	}

	// This line is basicly added for code coverage
	//cleanJourneysAndCars()
}

type reqArgs struct {
	body    string
	method  string
	path    string
	handler func(w http.ResponseWriter, r *http.Request)
	ctype   string
}

func simulateTestCall(t *testing.T, args reqArgs) *httptest.ResponseRecorder {
	W := httptest.NewRecorder()
	payload := strings.NewReader(args.body)
	req := httptest.NewRequest(args.method, args.path, payload)
	req.Header.Add("Content-Type", args.ctype)
	http.HandlerFunc(args.handler).ServeHTTP(W, req)
	if W.Code != http.StatusOK && W.Code != http.StatusAccepted && W.Code != http.StatusNoContent {
		t.Log("Failed in populate journeys")
		t.Fatalf("(Expected) %d or %d != %d (Returned)", http.StatusOK, http.StatusAccepted, W.Code)
	}
	return W
}
