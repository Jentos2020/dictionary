package handler

import (
	"encoding/json"
	"leetgo/internal/app/controller/converters"
	"leetgo/internal/entity"
	"leetgo/internal/errors"
	"leetgo/internal/gen"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) AddWord(ctx *fiber.Ctx) error {
	var reqWord gen.Word
	if err := json.Unmarshal(ctx.Body(), &reqWord); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	err := h.c.AddWord(ctx.UserContext(), converters.GenToEntityWord(&reqWord))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(reqWord)
}

func (h *Handler) RemoveWord(ctx *fiber.Ctx, word string) error {
	var body struct {
		Dictionary string `json:"dictionary"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
	}

	dict := body.Dictionary
	if dict == "" {
		return fiber.NewError(fiber.StatusBadRequest, "dictionary is required")
	}

	if err := h.c.RemoveWord(ctx.UserContext(), word, dict); err != nil {
		if err == errors.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Word not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpdateWord(ctx *fiber.Ctx, oldWord string) error {
	var req entity.Word
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}
	if req.Dictionary == "" {
		return fiber.NewError(fiber.StatusBadRequest, "dictionary is required")
	}
	if req.Data == "" {
		return fiber.NewError(fiber.StatusBadRequest, "word.data is required")
	}

	err := h.c.UpdateWord(ctx.UserContext(), oldWord, req.Data, req.Dictionary)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Word not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(req)
}
