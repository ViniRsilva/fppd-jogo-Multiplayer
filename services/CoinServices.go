package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"sync"
)

var CoinNotFound = errors.New("moeda não encontrada")

type CoinServiceArgs struct {
	PosX int
	PosY int
}

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

func (service *CoinService) DeleteCoinByPosition(args *CoinServiceArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	// procura a moeda pela posição
	for id, coin := range service.allCoins {
		if coin.X == args.PosX && coin.Y == args.PosY {
			delete(service.allCoins, id)
			*res = true
			return nil
		}
	}

	// se não encontrou nenhuma moeda na posição
	return CoinNotFound
}

/*
Cria uma nova moeda no jogo
*/
func (service *CoinService) CreateCoin(args *CoinServiceArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	newCoin := &models.Coin{
		ID: service.nextCoinID,
		X:  args.PosX,
		Y:  args.PosY,
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
func (service *CoinService) DeleteCoin(id int, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	_, found := service.allCoins[id]
	if !found {
		return CoinNotFound
	}

	delete(service.allCoins, id)
	*res = true
	return nil
}

/*
Lista todas as moedas ativas
*/
func (service *CoinService) ListAllCoins(args *struct{}, reply *[]models.Coin) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	coins := make([]models.Coin, 0, len(service.allCoins))
	for _, coin := range service.allCoins {
		coins = append(coins, *coin)
	}

	*reply = coins
	return nil
}
