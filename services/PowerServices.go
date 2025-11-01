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
func (service *PowerService) CreatePower(args *CreatePowerArgs, reply *models.Power) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	newPower := &models.Power{
		ID: service.nextPowerID,
		X:  args.X,
		Y:  args.Y,
	}
	service.nextPowerID++
	service.allPowers[newPower.ID] = newPower

	*reply = *newPower
	return nil
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
Deleta o poder
*/
func (service *PowerService) DeletePower(args *DeletePowerArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	_, found := service.allPowers[args.ID]
	if !found {
		return PowerNotFound
	}

	delete(service.allPowers, args.ID)
	*reply = true
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
}

type FindPowerByIdArgs struct {
	ID int
}

type DeletePowerArgs struct {
	ID int
}

type ListPowersArgs struct {
}
