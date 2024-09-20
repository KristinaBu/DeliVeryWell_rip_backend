package models

import (
	"errors"
	"strconv"
	"strings"
)

type callcenterCard struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Price   int    `json:"price"`
	Content string `json:"content"`
	Count   int    `json:"count"`
}

type callcenterRequest struct {
	Cards        []callcenterCard `json:"cards"`
	Data         string           `json:"data"`
	Address      string           `json:"address"`
	DeliveryName string           `json:"delivery-name"`
}

var callcenterCardlist = []callcenterCard{
	{ID: 1, Title: "Посылки", Price: 1075, Content: "Доставка посылок весом менее 1 кг в Новокузнецке"},
	{ID: 2, Title: "Спортивные товары", Price: 430, Content: "Доставка спортивных товаров: спортивные принадлежности (резинки, гантели) и спортивное оборудование (шведские стенки) от 430 рублей"},
	{ID: 3, Title: "Пицца", Price: 149, Content: "Доставка пиццы из ресторанов ДоДо, Maestrello, FoodBand от 149 рублей. При покупке от 4х штук - доставка 59 рублей"},
	{ID: 4, Title: "Цветы", Price: 799, Content: "Доставка цветов, букетов, упаковочных материалов в Москве"},
	{ID: 5, Title: "Суши", Price: 200, Content: "Доставка суши в Москве 200 рублей. При покупке товара от 500 рублей - доставка 100 рублей (Суши Мастер)"},
}

var callcenterMyRequestCardlist = map[int]callcenterRequest{
	1: {
		Cards: []callcenterCard{
			{ID: 2, Title: "Спортивные товары", Price: 430, Content: "Доставка спортивных товаров: спортивные принадлежности (резинки, гантели) и спортивное оборудование (шведские стенки) от 430 рублей", Count: 3},
			{ID: 3, Title: "Пицца", Price: 149, Content: "Доставка пиццы из ресторанов ДоДо, Maestrello, FoodBand от 149 рублей. При покупке от 4х штук - доставка 59 рублей", Count: 1},
		},
		Data:         "12.12.2024",
		Address:      "Москва, ул Бауманская, д 4, кв 3",
		DeliveryName: "Буйдина К.А.",
	},
}

func GetAllCards() []callcenterCard {
	return callcenterCardlist
}

func GetCallCardByID(id int) (*callcenterCard, error) {
	for _, a := range callcenterCardlist {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Card not found")
}

func GetMyCallCards(callRequestId int) callcenterRequest {
	return callcenterMyRequestCardlist[callRequestId]
}

func FindCallCards(title string) ([]callcenterCard, error) {
	findCards := []callcenterCard{}
	for _, card := range callcenterCardlist {
		if strings.Contains(card.Title, title) {
			findCards = append(findCards, card)
		}

	}
	if len(findCards) > 0 {
		return findCards, nil
	}
	return nil, errors.New("Find no cards")
}

func FindCallCardsByPrice(priceFrom, priceUp string) ([]callcenterCard, error) {
	findCards := []callcenterCard{}
	priceFromInt, _ := strconv.Atoi(priceFrom)
	priceUpInt, _ := strconv.Atoi(priceUp)
	// они неверные
	if priceUpInt < priceFromInt {
		return nil, errors.New("неверные параметры поиска")
	}
	for _, card := range callcenterCardlist {
		if card.Price <= priceUpInt && card.Price >= priceFromInt {
			findCards = append(findCards, card)
		}
	}
	if len(findCards) > 0 {
		return findCards, nil
	}
	return nil, errors.New("find no cards")
}
