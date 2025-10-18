package server

import (
	"fmt"
	"net/http"
)

// SendCustomErrorMessage отправляет клиенту сообщение вместе с соответствующим HTTP-статусом.
// Параметр statusCode определяет какой именно статус отправить, а также какое сообщение вернуть.
// Используется для передачи сообщений об ошибках и особых ситуациях.
func SendCustomErrorMessage(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	switch statusCode {

	case http.StatusNoContent: // 204
		w.Write([]byte("Order is not registered"))

	case http.StatusTooManyRequests: // 429
		w.Write([]byte("No more than N requests per minute allowed"))

	case http.StatusInternalServerError: // 500
		w.Write([]byte("Server error"))

	default:
		w.Write([]byte(fmt.Sprintf("Unknown error (%d)", statusCode)))
	}
}
