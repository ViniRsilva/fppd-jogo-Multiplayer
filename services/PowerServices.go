package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"sync"
)

var PowerNotFound = errors.New("poder não encontrado")

type PowerService struct {
	mu          sync.RWMutex
	allPowers   map[int]*models.Power
	nextPowerID int
	lastRequestID map[int]bool
}

func NewPowerService() *PowerService {
	return &PowerService{
		allPowers:   make(map[int]*models.Power),
		nextPowerID: 1,
	}
}

/*
Cria um novo poder no jogo
*/
func (service *PowerService) CreatePower(args *PowerServiceArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.lastRequestID[args.RequestID] {
		*reply = true
		return nil
	}

	newPower := &models.Power{
		ID: service.nextPowerID,
		X:  args.PosX,
		Y:  args.PosY,
	}
	service.nextPowerID++
	service.allPowers[newPower.ID] = newPower
	service.lastRequestID[args.RequestID] = true

	*reply = true
	return nil
}

// Delete power by position
func (service *PowerService) DeletePowerByPosition(args *PowerServiceArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.lastRequestID[args.RequestID] {
		*res = true
		return nil
	}

	// procura o mosntro pela posição
	for id, power := range service.allPowers {
		if power.X == args.PosX && power.Y == args.PosY {
			delete(service.allPowers, id)
			service.lastRequestID[args.RequestID] = true
			*res = true
			return nil
		}
	}

	// se não encontrou nenhuma poder na posição
	return PowerNotFound
}

/*
Devolve o poder identificado pelo ID
*/
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

/*
Lista todos os poderes ativos
*/
func (service *PowerService) ListAllPowers(args *ListPowersArgs, reply *[]models.Power) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	powers := make([]models.Power, 0, len(service.allPowers))
	for _, power := range service.allPowers {
		powers = append(powers, *power)
	}

	*reply = powers
	return nil
}

/*
Tipos auxiliares para chamadas RPC
*/
type CreatePowerArgs struct {
	X int
	Y int
	RequestID int
}

type FindPowerByIdArgs struct {
	ID int
}

type DeletePowerArgs struct {
	ID int
	RequestID int
}

type ListPowersArgs struct {
}


type PowerServiceArgs struct {
	PosX int
	PosY int
	RequestID int
}
