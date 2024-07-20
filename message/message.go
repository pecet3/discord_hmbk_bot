package message

import (
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/pecet3/discord_hmbk_bot/discord"
)

func Run(mux *http.ServeMux, discord *discordgo.Session) {
	mux.HandleFunc("POST /message", func(w http.ResponseWriter, r *http.Request) {
		handleMessage(w, r, discord)
	})
}

func handleMessage(w http.ResponseWriter, r *http.Request, dc *discordgo.Session) {
	message := r.FormValue("message")
	log.Println("> Message from Web: ", message)
	dc.ChannelMessageSend(discord.FONTANNA_ID, message)

	w.WriteHeader(http.StatusOK)
}
