package repo

import (
	"arch/internal/viewer/entity"
	"context"
	"github.com/patrickmn/go-cache"
	"sort"
)

type Message struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Message {
	return &Message{
		cache: cache,
	}
}

func (m *Message) Set(_ context.Context, message entity.Message) {
	m.cache.Set("message:"+message.ID, message, cache.NoExpiration)
}

func (m *Message) Get(_ context.Context, id string) (entity.Message, error) {
	storedObject, found := m.cache.Get("message:" + id)
	if !found {
		return entity.Message{}, ErrNotFound
	}
	message, ok := storedObject.(entity.Message)
	if !ok {
		return entity.Message{}, ErrBrokenData
	}
	return message, nil
}

func (m *Message) GetAll(_ context.Context) ([]entity.Message, error) {
	items := m.cache.Items()

	messages := make([]entity.Message, 0, len(items))
	for _, item := range items {
		message, ok := item.Object.(entity.Message)
		if !ok {
			return []entity.Message{}, ErrBrokenData
		}
		messages = append(messages, message)
	}

	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	return messages, nil
}
