package transport

import (
	"BMSTU_IU5_53B_rip/internal/models"
	"BMSTU_IU5_53B_rip/internal/render"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func showIndexPage(c *gin.Context) {
	priceFrom := c.Query("price_from")
	priceUp := c.Query("price_up")

	if priceFrom != "" && priceUp != "" {
		// валидация. Надо ее раньше сделать или нет?
		cards, err := models.FindCallCardsByPrice(priceFrom, priceUp)
		if err != nil {
			render.Render(c, "index.html", gin.H{
				"NoCards":    "Некорректный запрос",
				"SearchFrom": priceFrom,
				"SearchUp":   priceUp,
			})
		} else {
			render.Render(c, "index.html", gin.H{
				"payload":    cards,
				"SearchFrom": priceFrom,
				"SearchUp":   priceUp,
			})
		}
	} else {
		cards := models.GetAllCards()
		render.Render(c, "index.html", gin.H{
			"payload":    cards,
			"SearchFrom": "",
			"SearchUp":   "",
		})
	}
}

func getCallCard(c *gin.Context) {
	if cardID, err := strconv.Atoi(c.Param("card_id")); err == nil {
		if card, err := models.GetCallCardByID(cardID); err == nil {
			render.Render(c, "cardDetails.html", gin.H{
				"title":   card.Title,
				"payload": card,
			})
			return
		}
	}
	// Если дошли сюда, значит либо ошибка парсинга ID, либо карточка не найдена
	allCards := models.GetAllCards()
	if len(allCards) > 0 {
		firstCard := allCards[0]
		render.Render(c, "cardDetails.html", gin.H{
			"title":   firstCard.Title,
			"payload": allCards[0],
		})
	} else {
		// Если карточек вообще нет, вернем ошибку
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("sorry! no cards available"))
	}
}

func getMyCallCards(c *gin.Context) {
	if callRequestId, err := strconv.Atoi(c.Param("callrequest_id")); err == nil {
		callRequest := models.GetMyCallCards(callRequestId)
		render.Render(c, "mycards.html", gin.H{
			"payload":      callRequest.Cards,
			"Data":         callRequest.Data,
			"Address":      callRequest.Address,
			"DeliveryName": callRequest.DeliveryName,
		})
	}
}
