package paint

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func Run(mux *http.ServeMux, ps *PaintSessions) {

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("POST /painting", func(w http.ResponseWriter, r *http.Request) {
		handlePainting(w, r, ps)
	})

	http.ListenAndServe(":8080", mux)
}

type PaintData struct {
	Image     string `json:"image"`
	SessionId string `json:"session_id"`
}

func handlePainting(w http.ResponseWriter, r *http.Request, ps *PaintSessions) {
	var data PaintData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	s, isExisting := ps.GetSession(data.SessionId)

	if !isExisting {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if s.ExpiresAt.Before(time.Now()) {
		ps.RemoveSession(data.SessionId)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	imageData := strings.TrimPrefix(data.Image, "data:image/png;base64,")

	imgBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		http.Error(w, "Unable to decode image", http.StatusInternalServerError)
		return
	}

	s.ImgBytesCh <- imgBytes

}
