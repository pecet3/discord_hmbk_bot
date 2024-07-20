package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func handlePraise(praiseCh chan *discordgo.User, discord *discordgo.Session) {
	var users []*discordgo.User
	defer close(praiseCh)
	for {
		select {
		case dUser := <-praiseCh:
			isOnList := false
			for _, u := range users {
				if u.ID == dUser.ID {
					isOnList = true
				}
			}
			if !isOnList {
				users = append(users, dUser)
			}
		default:
			now := time.Now()
			if now.Second() == 0 && now.Minute() == 0 && now.Hour() == 9 || now.Hour() == 13 || now.Hour() == 20 {
				if len(users) <= 0 {
					continue
				}
				randomIndex := rand.Intn(len(users))

				winner := users[randomIndex]

				if winner == nil {
					continue
				}
				praise, err := generatePraise(winner.Username)
				if err != nil {
					log.Println(err)
					continue
				}
				lenName := len(winner.Username)
				discord.ChannelMessageSend(FONTANNA_ID, winner.Mention()+praise[lenName:])
				users = nil
			}
		}
		time.Sleep(70 * time.Millisecond)
	}
}

func generatePraise(name string) (string, error) {
	const apiEndpoint = "https://api.openai.com/v1/chat/completions"

	ctx := context.Background()
	apiKey := os.Getenv("GPT_KEY")
	if apiKey == "" {
		return "", errors.New("no api key")
	}
	system := "Jesteś botem, który chwali ludzi w sposób zabawny. Twoje odpowiedzi mają maksymalnie 3 zdania. Pochwałe zaczynasz zawsze od nicku."
	human := "Pochwal użytkownika o nicku: " + name
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reqBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []interface{}{
			map[string]interface{}{"role": "system", "content": system},
			map[string]interface{}{"role": "user", "content": human},
		},
		"max_tokens": 1200,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var data map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content, nil
}
