package service

import (
	"adorable-star/internal/pkg/crawler"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

var GPT = &GptService{}

type GptService struct{}

type GptRes struct {
	Message struct {
		Id      string `json:"id"`
		Content struct {
			Parts []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
	ConversationId string `json:"conversation_id"`
}

func (s *GptService) Conversation(msg string) (convId string, response string) {
	res := crawler.GptPage.MustEval(`(token) => {
		let xhr = new XMLHttpRequest();
		xhr.open("POST", "/backend-api/conversation", false);
		xhr.setRequestHeader('X-Authorization', 'Bearer ' + token);
		xhr.setRequestHeader('Content-type', 'application/json');

		xhr.send('{"model":"text-davinci-002-render-sha","action":"next","parent_message_id":"`+uuid.NewString()+`","messages":[{"id":"`+uuid.NewString()+`","author":{"role": "user"},"content":{"content_type":"text","parts":["`+msg+`"]}}]}');

		let res = {"responseText":xhr.responseText,"responseHeaders":xhr.getAllResponseHeaders(),"status":xhr.status};
		return JSON.stringify(res);
	}`, crawler.GptAccessToken).Str()

	resArr := strings.Split(res, "\\n\\n")

	var resStruct *GptRes
	json.Unmarshal([]byte(strings.ReplaceAll(strings.ReplaceAll(resArr[len(resArr)-5], "\\", ""), "data: ", "")), &resStruct)

	convId = resStruct.ConversationId
	response = resStruct.Message.Content.Parts[0]
	return
}
