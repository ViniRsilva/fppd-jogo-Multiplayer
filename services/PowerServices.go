package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var PowerNotFound = errors.New("poder nao encontrado")

type CreatePowerArgs struct {
	PlayerID  int
	RequestID int
	X         int
	Y         int
}

type DeleteAllPowersArgs struct {
	PlayerID  int
	RequestID int
}

type FindPowerByIdArgs struct {
	ID int
}

type ListPowersArgs struct {
}

type PowerService struct {
	mu                   sync.RWMutex
	allPowers            map[int]*models.Power
	nextPowerID          int
	lastProcessedRequest map[int]int
}

// NewPowerService é o construtor do serviço.
func NewPowerService() *PowerService {
	return &PowerService{
		allPowers:            make(map[int]*models.Power),
		nextPowerID:          1,
		lastProcessedRequest: make(map[int]int),
	}
}

// CreatePower cria um novo poder no jogo.
func (service *PowerService) CreatePower(args *CreatePowerArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.PlayerID]

	if args.RequestID > lastID {
		newPower := &models.Power{
			ID: service.nextPowerID,
			X:  args.X,
			Y:  args.Y,
		}
		service.nextPowerID++
		service.allPowers[newPower.ID] = newPower
		service.lastProcessedRequest[args.PlayerID] = args.RequestID
		log.Printf("Create Power: (%d,%d)", args.X, args.Y)
	}

	*reply = true
	return nil
}

// DeleteAllPowers deleta TODOS os poderes do mapa.
func (service *PowerService) DeleteAllPowers(args *DeleteAllPowersArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.PlayerID]

	if args.RequestID > lastID {
		service.allPowers = make(map[int]*models.Power)

		log.Printf("Player %d deletou todos os poderes.", args.PlayerID)

		service.lastProcessedRequest[args.PlayerID] = args.RequestID
	}

	*res = true
	return nil
}

// FindPowerByID devolve o poder identificado pelo ID.
func (service *PowerService) FindPowerByID(args *FindPowerByIdArgs, reply *models.Power) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	power, found := service.allPowers[args.ID]
	if !found {
		return PowerNotFound
	}

	*reply = *power
	return nil
}

// ListAllPowers lista todos os poderes ativos.
func (service *PowerService) ListAllPowers(_ *ListPowersArgs, reply *[]models.Power) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	powers := make([]models.Power, 0, len(service.allPowers))
	for _, power := range service.allPowers {
		powers = append(powers, *power)
	}

	*reply = powers
	return nil
}
