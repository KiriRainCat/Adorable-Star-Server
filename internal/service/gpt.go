package service

import (
	"adorable-star/internal/pkg/crawler"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

var GPT = &GptService{}

type GptService struct{}

func (s *GptService) Conversation(stream chan string, convId string, msgs []map[string]string) {
	// Define GPT response struct
	type GptResponse struct {
		Message struct {
			Id      string `json:"id"`
			Content struct {
				Parts []string `json:"parts"`
			} `json:"content"`
		} `json:"message"`
		ConversationId string `json:"conversation_id"`
	}

	// Define base request body
	reqBody := map[string]any{
		"action":            "next",
		"model":             "text-davinci-002-render-sha",
		"parent_message_id": uuid.NewString(),
		"messages":          []map[string]any{},
	}

	// Add messages to request body
	var prev string
	for _, msg := range msgs {
		prev += msg["text"]
		reqBody["messages"] = append(reqBody["messages"].([]map[string]any), map[string]any{
			"id": uuid.NewString(),
			"author": map[string]any{
				"role": msg["role"],
			},
			"content": map[string]any{
				"content_type": "text",
				"parts":        []string{msg["text"]},
			},
		})
	}

	// If the request aim to continue a chat with a conversation ID
	if len(convId) != 0 {
		reqBody["conversation_id"] = convId
	}

	// Marshal request body to JSON
	reqData, _ := json.Marshal(reqBody)

	// Send request to GPT
	crawler.GptPage.MustEval(`async (token, data) => {
		const res = await fetch('/backend-api/conversation', {
			method: 'POST',
			headers: {
				'Authorization': token,
				'Content-Type': 'application/json',
			},
			body: data,
		});
		window.reader = res.body.getReader();
	}`, crawler.GptAccessToken, string(reqData))

	// Read response stream from GPT
	for {
		// Get data from stream reader
		res := crawler.GptPage.MustEval(`async () => {
        const { value, done } = await window.reader.read();
        if (done) {
          return "DONE";
        }
        return new TextDecoder().decode(value);
      }`).Str()

		// When stream is over, stop
		if res == "DONE" {
			break
		}

		// Loop through the data got from stream and evaluate them
		for _, data := range strings.Split(strings.ReplaceAll(res, "data: ", ""), "\n\n") {
			var gptRes GptResponse
			err := json.Unmarshal([]byte(data), &gptRes)
			if err != nil {
				continue
			}

			// Store the conversation ID when needed
			if len(convId) == 0 && gptRes.ConversationId != "" {
				convId = gptRes.ConversationId
			}

			// Only read response when it's not empty
			if len(gptRes.Message.Content.Parts) != 0 {
				// Skip responses that are the same as the previous one
				text := gptRes.Message.Content.Parts[0]
				if strings.Contains(prev, text) || len(text) == 0 {
					continue
				}

				// Send the response to the stream
				prev = text
				stream <- text
			}
		}
	}

	// Send the conversation ID to the stream as the end
	stream <- convId
}
