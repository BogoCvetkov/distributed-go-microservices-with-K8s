package helpers

import (
	"broker-service/cmd/api/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JsonResponse struct {
	Err     bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func DecodeJSON(src io.Reader, v any) error {
	err := json.NewDecoder(src).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data *JsonResponse, status ...int) error {
	payload := data

	statusCode := http.StatusAccepted
	if len(status) > 0 {
		statusCode = status[0]
	}

	out, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)

	return nil
}

func ParseJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func ErrJson(w http.ResponseWriter, msg string, status ...int) {

	payload := JsonResponse{
		Err:     true,
		Message: msg,
	}
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	err := WriteJSON(w, &payload, statusCode)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte("Error while parsing JSON"))
	}
}

func CallModule(w http.ResponseWriter, info *types.MethodCallInfo) {
	body, err := json.Marshal(info.Body)
	if err != nil {
		ErrJson(w, "Failed to marshal body")
		return
	}

	request, err := http.NewRequest(info.Method, info.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		ErrJson(w, "Failed to call the auth module")
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrJson(w, fmt.Sprintf("Failed to call -->  %s", info.Endpoint))
		return
	}

	defer response.Body.Close()

	var result types.JsonResponse

	err = DecodeJSON(response.Body, &result)
	if err != nil {
		fmt.Println(err)
		ErrJson(w, "Failed to parse service response", http.StatusBadRequest)
		return
	}
	fmt.Println(response.StatusCode, result)

	if response.StatusCode >= 400 {
		ErrJson(w, string(body), response.StatusCode)
		return
	}

	WriteJSON(w, (*JsonResponse)(&result))

}
