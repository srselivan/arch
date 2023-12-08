package entity

import (
	"arch/proto/pb"
)

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

func (m *Message) ToProto() *pb.Message {
	p := &pb.Message{
		Id:          m.ID,
		MediaObject: nil,
		TextObject:  nil,
	}
	if m.MediaObject != nil {
		p.MediaObject = m.MediaObject.ToProto()
	}
	if m.TextObject != nil {
		p.TextObject = m.TextObject.ToProto()
	}
	return p
}

func (o *MediaObject) ToProto() *pb.MediaObject {
	return &pb.MediaObject{
		Filename: o.Filename,
		Path:     o.Path,
	}
}

func (o *TextObject) ToProto() *pb.TextObject {
	return &pb.TextObject{
		Body: o.Body,
	}
}
