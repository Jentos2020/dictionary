package handler

import (
	"encoding/json"
	"leetgo/internal/entity"

	"github.com/gofiber/websocket/v2"
)

type WSRequest struct {
	Type       string `json:"type"`
	Prefix     string `json:"prefix,omitempty"`
	Word       string `json:"word,omitempty"`
	Dictionary string `json:"dictionary,omitempty"`
}

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

			var req WSRequest
			if err := json.Unmarshal(msg, &req); err != nil {
				h.sendErrorWithType(c, "error", "Invalid JSON")
				continue
			}

			if req.Type == "" {
				h.sendErrorWithType(c, "error", "Missing 'type' in request")
				continue
			}

			if req.Type == "search" {
				if req.Prefix == "" {
					h.sendErrorWithType(c, "search", "Missing prefix")
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
			} else if req.Type == "get_definition" {
				if req.Word == "" {
					h.sendErrorWithType(c, "get_definition", "Missing word")
					continue
				}

				def, found, err := h.c.GetDefinition(h.c.Ctx, req.Word, req.Dictionary)
				resp := map[string]interface{}{
					"type": "definition",
					"word": req.Word,
				}
				if err != nil {
					resp["error"] = err.Error()
				} else if found {
					resp["definition"] = def
				} else {
					resp["error"] = "Definition not found"
				}
				respJSON, _ := json.Marshal(resp)
				if err := c.WriteMessage(websocket.TextMessage, respJSON); err != nil {
					h.c.Logger.Error("WebSocket write error", "client", clientAddr, "error", err)
					break
				}
			} else {
				h.sendErrorWithType(c, req.Type, "Unknown type")
			}
		}
	}
}

func (h *Handler) sendError(c *websocket.Conn, msg string) {
	h.sendErrorWithType(c, "error", msg)
}

func (h *Handler) sendErrorWithType(c *websocket.Conn, typ, msg string) {
	resp := map[string]interface{}{"type": typ, "error": msg}
	respJSON, _ := json.Marshal(resp)
	c.WriteMessage(websocket.TextMessage, respJSON)
}
