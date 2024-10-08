package routes

import (
	"log/slog"
	"net/http"
	"time"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /mempool/api/orangeclock", handleMempoolApiOrangeclock)
	mux.HandleFunc("GET /datetime", handleDatetime)

	mux.HandleFunc("GET /", handleRoot)
}

func handleMempoolApiOrangeclock(w http.ResponseWriter, req *http.Request) {
	slog.Debug("got req from", slog.String("remoteAddr", req.RemoteAddr))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(LoadMempoolData()))
	if err != nil {
		slog.Error(err.Error())
	}
}

func handleDatetime(w http.ResponseWriter, req *http.Request) {
	slog.Debug("got req from", slog.String("remoteAddr", req.RemoteAddr))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(time.Now().Format(time.RFC3339)))
	if err != nil {
		slog.Error(err.Error())
	}
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	slog.Debug("got req from", slog.String("remoteAddr", req.RemoteAddr))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		slog.Error(err.Error())
	}
}
