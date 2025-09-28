package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// LogMiddleware registra informações sobre cada requisição
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()

		// Criar wrapper para capturar status code
		wrapper := &ResponseWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log da requisição
		fmt.Printf("📥 %s %s - %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		// Executar handler
		next(wrapper, r)

		// Log da resposta
		duracao := time.Since(inicio)
		fmt.Printf("📤 %s %s - Status: %d - Duração: %v\n",
			r.Method, r.URL.Path, wrapper.statusCode, duracao)
	}
}

// ResponseWrapper captura o status code da resposta
type ResponseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWrapper) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
