package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var CoinNotFound = errors.New("moeda não encontrada")

type CreateCoinArgs struct {
	PlayerID  int
	RequestID int
	PosX      int
	PosY      int
}

type DeleteCoinByPositionArgs struct {
	PlayerID  int
	RequestID int
	PosX      int
	PosY      int
}

type DeleteCoinByIDArgs struct {
	PlayerID  int
	RequestID int
	CoinID    int
}

type CoinService struct {
	mu                   sync.RWMutex
	allCoins             map[int]*models.Coin
	nextCoinID           int
	lastProcessedRequest map[int]int
}

func NewCoinService() *CoinService {
	return &CoinService{
		allCoins:             make(map[int]*models.Coin),
		nextCoinID:           1,
		lastProcessedRequest: make(map[int]int),
	}
}

// CreateCoin cria uma nova moeda no jogo
func (service *CoinService) CreateCoin(args *CreateCoinArgs, reply *models.Coin) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.PlayerID]

	if args.RequestID > lastID {
		newCoin := &models.Coin{
			ID: service.nextCoinID,
			X:  args.PosX,
			Y:  args.PosY,
		}
		service.nextCoinID++
		service.allCoins[newCoin.ID] = newCoin
		service.lastProcessedRequest[args.PlayerID] = args.RequestID
		log.Printf("Created Coin: (%d,%d)", args.PosX, args.PosY)

		*reply = *newCoin
	}

	return nil
}

// DeleteCoinByPosition encontra e deleta a primeira moeda em uma dada posição.
func (service *CoinService) DeleteCoinByPosition(args *DeleteCoinByPositionArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.PlayerID]

	if args.RequestID > lastID {
		var foundID = -1
		for id, coin := range service.allCoins {
			if coin.X == args.PosX && coin.Y == args.PosY {
				foundID = id
				break
			}
		}

		if foundID == -1 {
			return CoinNotFound
		}

		delete(service.allCoins, foundID)
		log.Printf("Deleted Coin: (%d,%d)", args.PosX, args.PosY)
		service.lastProcessedRequest[args.PlayerID] = args.RequestID
	}

	*res = true
	return nil
}

// DeleteCoinByID deleta uma moeda pelo seu ID, com proteção contra reexecução (não utilizado).
func (service *CoinService) DeleteCoinByID(args *DeleteCoinByIDArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.PlayerID]
	if args.RequestID > lastID {
		_, found := service.allCoins[args.CoinID]
		if !found {
			return CoinNotFound
		}
		delete(service.allCoins, args.CoinID)
		service.lastProcessedRequest[args.PlayerID] = args.RequestID
	}

	*res = true
	return nil
}

// FindCoinByID devolve a moeda identificada pelo ID (não utilizado).
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

// ListAllCoins lista todas as moedas ativas no jogo.
func (service *CoinService) ListAllCoins(_ *struct{}, reply *[]models.Coin) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	coins := make([]models.Coin, 0, len(service.allCoins))
	for _, coin := range service.allCoins {
		coins = append(coins, *coin)
	}

	*reply = coins
	return nil
}
