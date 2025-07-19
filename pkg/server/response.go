package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type RFC7808 struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func ListResponse(w http.ResponseWriter, items []interface{}) {
	SuccessResponse(w, struct {
		Items  []interface{} `json:"items"`
		Count  int           `json:"count"`
		Limit  int           `json:"limit"`
		Offset int           `json:"offset"`
		Total  int           `json:"total"`
	}{
		Items:  items,
		Offset: 0, // pagination is not yet implemented
		Count:  len(items),
		Limit:  0, // limiting results is not yet implemented
		Total:  len(items),
	})
}

func CollectMap[M ~map[K]V, K comparable, V any](m M) []interface{} {
	output := make([]interface{}, 0, len(m))
	for _, v := range m {
		output = append(output, v)
	}

	return output
}

func CodeResponse(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := fmt.Fprintf(w, "%s", fmt.Sprintf(`{"status": "%s"}`, http.StatusText(code))); err != nil {
		slog.Error("failed to write response", "err", err.Error())

		return
	}
}

func SuccessResponse(w http.ResponseWriter, item interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	content, err := json.Marshal(item)
	if err != nil {
		slog.Error("failed to write response", "err", err.Error())

		return
	}

	if _, err := fmt.Fprint(w, string(content)); err != nil {
		slog.Error("failed to write response", "err", err.Error())

		return
	}
}

func NotFoundResponse(w http.ResponseWriter, detail string) {
	ErrorResponse(w, http.StatusNotFound, "not found", detail)
}

func ErrorResponse(w http.ResponseWriter, code int, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(code)

	msg, err := json.Marshal(RFC7808{Type: "about:blank", Title: title, Detail: detail})
	if err != nil {
		slog.Error("failed to write response", "err", err.Error())

		return
	}

	if _, err := fmt.Fprint(w, string(msg)); err != nil {
		slog.Error("failed to write response", "err", err.Error())

		return
	}
}
