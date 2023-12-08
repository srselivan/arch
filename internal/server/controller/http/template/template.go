package template

import (
	"arch/internal/server/entity"
)

const errorMessage = "error: message is empty or type is unsupported"

type TemplateContent struct {
	Messages []messageObject
}

type messageObject struct {
	IsImage   bool
	ImagePath string
	Text      string
}

func TemplateContentFromServiceModels(messages []entity.Message) TemplateContent {
	messageObjects := make([]messageObject, 0, len(messages))
	for _, message := range messages {
		msg := messageObject{}
		switch {
		case message.MediaObject != nil:
			msg.IsImage = true
			msg.ImagePath = message.MediaObject.Path
		case message.TextObject != nil:
			msg.Text = message.TextObject.Body
		default:
			msg.Text = errorMessage
		}
		messageObjects = append(messageObjects, msg)
	}
	return TemplateContent{messageObjects}
}
