package entity

import (
	"arch/proto/pb"
	"github.com/samber/lo"
)

const errorMessage = "error: message is empty or type is unsupported"

type Message struct {
	ID          string       `json:"id"`
	MediaObject *MediaObject `json:"mediaObject,omitempty"`
	TextObject  *TextObject  `json:"textObject,omitempty"`
}

type MediaObject struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type TextObject struct {
	Body string `json:"body"`
}

func (m *Message) String() string {
	switch {
	case m.MediaObject != nil:
		return m.MediaObject.Filename
	case m.TextObject != nil:
		return m.TextObject.Body
	default:
		return errorMessage
	}
}

func MessageFromProto(p *pb.Message) Message {
	m := Message{
		ID:          p.Id,
		MediaObject: nil,
		TextObject:  nil,
	}
	if p.MediaObject != nil {
		m.MediaObject = lo.ToPtr(MediaObjectFromProto(p.MediaObject))
	}
	if p.TextObject != nil {
		m.TextObject = lo.ToPtr(TextObjectFromProto(p.TextObject))
	}
	return m
}

func MediaObjectFromProto(p *pb.MediaObject) MediaObject {
	return MediaObject{
		Filename: p.Filename,
		Path:     p.Path,
	}
}

func TextObjectFromProto(p *pb.TextObject) TextObject {
	return TextObject{
		Body: p.Body,
	}
}
