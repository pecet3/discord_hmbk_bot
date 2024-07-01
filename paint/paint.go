package paint

import (
	"log"
	"net/http"
)

func Run(ps *PaintSessions) {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("POST /painting", func(w http.ResponseWriter, r *http.Request) {
		handlePainting(w, r, ps)
	})

	http.ListenAndServe(":8080", mux)
}

func handlePainting(w http.ResponseWriter, r *http.Request, ps *PaintSessions) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metoda niedozwolona", http.StatusMethodNotAllowed)
		return
	}

	session := r.FormValue("session")
	log.Println(session)
	s, isSession := ps.GetSession(session)

	s.FinishCh <- true
	if isSession {
		log.Println(s.Id)
	}

}
