package handler

import (
	"encoding/json"
	"leetgo/internal/consts"
	"leetgo/internal/errors"
	"leetgo/internal/gen"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PostDefinitions(c *fiber.Ctx) error {
	var req gen.Definition
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	if req.Word == "" {
		return fiber.NewError(fiber.StatusBadRequest, "word is required")
	}
	if req.Definition == "" {
		return fiber.NewError(fiber.StatusBadRequest, "definition is required")
	}

	dict := consts.DefaultDict
	if req.Dictionary != nil && *req.Dictionary != "" {
		dict = *req.Dictionary
	}

	if err := h.c.AddDefinition(c.UserContext(), req.Word, req.Definition, dict); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(req)
}

func (h *Handler) GetDefinitionsWord(c *fiber.Ctx, word string, params gen.GetDefinitionsWordParams) error {
	dict := consts.DefaultDict
	if params.Dictionary != nil && *params.Dictionary != "" {
		dict = *params.Dictionary
	}

	def, found, err := h.c.GetDefinition(c.UserContext(), word, dict)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !found {
		return fiber.NewError(fiber.StatusNotFound, "Definition not found")
	}

	return c.JSON(map[string]interface{}{
		"word":       word,
		"definition": def,
		"dictionary": dict,
	})
}

func (h *Handler) PutDefinitionsWord(c *fiber.Ctx, word string, params gen.PutDefinitionsWordParams) error {
	dict := consts.DefaultDict
	if params.Dictionary != nil && *params.Dictionary != "" {
		dict = *params.Dictionary
	}

	var body struct {
		Definition string `json:"definition"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}
	if body.Definition == "" {
		return fiber.NewError(fiber.StatusBadRequest, "definition is required")
	}

	if err := h.c.UpdateDefinition(c.UserContext(), word, body.Definition, dict); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Definition not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(map[string]interface{}{
		"word":       word,
		"definition": body.Definition,
		"dictionary": dict,
	})
}

func (h *Handler) DeleteDefinitionsWord(c *fiber.Ctx, word string, params gen.DeleteDefinitionsWordParams) error {
	dict := consts.DefaultDict
	if params.Dictionary != nil && *params.Dictionary != "" {
		dict = *params.Dictionary
	}

	if err := h.c.RemoveDefinition(c.UserContext(), word, dict); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Definition not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
