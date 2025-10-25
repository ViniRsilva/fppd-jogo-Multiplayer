package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"sync"
)

var CoinNotFound = errors.New("moeda não encontrada")

type CoinService struct {
	mu         sync.RWMutex
	allCoins   map[int]*models.Coin
	nextCoinID int
}

func NewCoinService() *CoinService {
	return &CoinService{
		allCoins:   make(map[int]*models.Coin),
		nextCoinID: 1,
	}
}

/*
Cria uma nova moeda no jogo
*/
func (service *CoinService) createCoin(posX, posY int, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	newCoin := &models.Coin{
		ID: service.nextCoinID,
		X:  posX,
		Y:  posY,
	}
	service.nextCoinID++
	service.allCoins[newCoin.ID] = newCoin

	*reply = true
	return nil
}

/*
Devolve a moeda identificada pelo ID
*/
func (service *CoinService) FindCoinByID(id int, reply *models.Coin) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	coin, found := service.allCoins[id]
	if !found {
		return CoinNotFound
	}

	*reply = *coin
	return nil
}

/*
Deleta a moeda
*/
func (service *CoinService) deleteCoin(ID int, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	_, found := service.allCoins[ID]
	if !found {
		return CoinNotFound
	}

	delete(service.allCoins, ID)
	*res = true
	return nil
}
