package responses

import (
	"encoding/json"
	"net/http"
)

// JSONResponse provides setting JSON-encoded Responses.
type JSONResponse struct {
	w http.ResponseWriter
}

// NewJSON instantiates a JSONResponse.
func NewJSON(w http.ResponseWriter) JSONResponse {
	w.Header().Add("Content-Type", "application/json")
	return JSONResponse{
		w: w,
	}
}

// Success sets the JSONResponse to include a http.Status and a marshalled body.
func (j JSONResponse) Success(status int, body interface{}) JSONResponse {
	j.w.WriteHeader(status)

	resBytes, err := json.Marshal(body)
	if err != nil {
		return j.Fail(http.StatusBadRequest, err)
	}

	_, err = j.w.Write(resBytes)
	if err != nil {
		return j.Fail(http.StatusBadRequest, err)
	}

	return j
}

// Fail sets the JSONResponse to include a http.Status and the error that was raised.
func (j JSONResponse) Fail(status int, err error) JSONResponse {
	j.w.WriteHeader(status)

	encodedErr := Error{
		Status: status,
		Detail: err.Error(),
	}

	resBytes, er := json.Marshal(encodedErr)
	if er != nil {
		return j.Fail(status, er)
	}

	_, er = j.w.Write(resBytes)
	if er != nil {
		return j.Fail(status, er)
	}

	return j
}
