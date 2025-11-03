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

// CreateCoin creates a new coin
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

// DeleteCoinByPosition deletes a coin based on its position
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

// ListAllCoins lists all coins (used to render the complete map)
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
