package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) DeliveryItemList(ctx *gin.Context) {
	priceFrom := ctx.Query("price_from")
	priceTo := ctx.Query("price_to")
	if priceFrom == "" && priceTo == "" {
		cards, err := h.Repository.DeliveryItemList()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		fmt.Println("Found items:", len(*cards)) // Добавьте этот лог
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"NoCards":    "",
			"payload":    cards,
			"SearchFrom": priceFrom,
			"SearchUp":   priceTo,
		})
		return
	}
	cards, err := h.Repository.SearchDeliveryItem(priceFrom, priceTo)
	fmt.Println("Found items:", len(*cards)) // Добавьте этот лог
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"NoCards":    "нет подходящего",
		"payload":    cards,
		"SearchFrom": priceFrom,
		"SearchUp":   priceTo,
	})

}

func (h *Handler) DeliveryItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	card, err := h.Repository.GetDeliveryItemByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "cardDetails.html", gin.H{
		"title":   card.Title,
		"payload": card,
	})

}

func (h *Handler) DeleteDeliveryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	h.Repository.DeleteDeliveryItem(id)
	ctx.Redirect(http.StatusFound, "/")

}
