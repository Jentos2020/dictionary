package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) LoadDictionary(ctx *fiber.Ctx, name string) error {
	h.c.FillTrieWithWords(ctx.UserContext(), name)

	return ctx.JSON(fiber.Map{"message": "Dictionary loaded/updated"})
}
