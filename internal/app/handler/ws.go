package handler

import (
	"encoding/json"
	"leetgo/internal/entity"

	"github.com/gofiber/websocket/v2"
)

func WSHandler(h *Handler) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		clientAddr := c.RemoteAddr().String()
		h.c.Logger.Debug("WebSocket connection established", "client", clientAddr)

		defer func() {
			h.c.Logger.Debug("WebSocket connection closed", "client", clientAddr)
			c.Close()
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				h.c.Logger.Error("WebSocket read error", "client", clientAddr, "error", err)
				break
			}
			h.c.Logger.Debug("WebSocket received", "client", clientAddr, "msg", string(msg))

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
				h.c.Logger.Error("WebSocket write error", "client", clientAddr, "error", err)
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
