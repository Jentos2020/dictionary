package handler

import (
	"encoding/json"
	"fmt"
	"leetgo/internal/entity"

	"github.com/gofiber/websocket/v2"
)

func WSHandler(h *Handler) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				h.c.Logger.Error(fmt.Sprintf("WebSocket read error %s", err))
				break
			}

			var req entity.SearchRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				h.sendError(c, "Invalid JSON")
				continue
			}

			currentTrie := h.c.Trie.Load().(*entity.Trie)
			wordsStr := currentTrie.GetWordsByPrefix(req.Prefix)
			words := make(entity.Words, len(wordsStr))
			for i, s := range wordsStr {
				words[i] = entity.Word{Data: s}
			}
			resp := entity.SearchResponse{Words: words}

			respJSON, _ := json.Marshal(resp)
			if err := c.WriteMessage(websocket.TextMessage, respJSON); err != nil {
				h.c.Logger.Error(fmt.Sprintf("WS write error %s", err))
				break
			}
		}
	}
}

func (h *Handler) sendError(c *websocket.Conn, msg string) {
	resp := entity.SearchResponse{Error: msg}
	respJSON, _ := json.Marshal(resp)
	c.WriteMessage(websocket.TextMessage, respJSON)
}
