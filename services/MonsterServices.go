package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"sync"
)

var MonsterNotFound = errors.New("Monstro não encontrado")

type MonsterService struct {
	mu            sync.RWMutex
	allMonsters   map[int]*models.Monster
	nextMonsterID int
	lastRequestID map[int]bool
}

func NewMonsterService() *MonsterService {
	return &MonsterService{
		allMonsters:   make(map[int]*models.Monster),
		nextMonsterID: 1,
	}
}

// Create a new monster
func (service *MonsterService) CreateMonster(args *MonsterServiceArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.lastRequestID[args.RequestID] {
		*reply = true
		return nil
	}

	newMonster := &models.Monster{
		ID: service.nextMonsterID,
		X:  args.PosX,
		Y:  args.PosY,
	}
	service.nextMonsterID++
	service.allMonsters[newMonster.ID] = newMonster
	service.lastRequestID[args.RequestID] = true

	*reply = true
	return nil
}

// Delete monster by position
func (service *MonsterService) DeleteMonsterByPosition(args *MonsterServiceArgs, res *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.lastRequestID[args.RequestID] {
		*res = true
		return nil
	}

	// procura o mosntro pela posição
	for id, monster := range service.allMonsters {
		if monster.X == args.PosX && monster.Y == args.PosY {
			delete(service.allMonsters, id)
			service.lastRequestID[args.RequestID] = true
			*res = true
			return nil
		}
	}

	// se não encontrou nenhuma monstro na posição
	return MonsterNotFound
}

// Delete all monsters (When player get a power)
func (service *MonsterService) DeleteAllMonsters(args *DeleteAllMonsterArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	for id := range service.allMonsters {
		delete(service.allMonsters, id)
	}

	*reply = true
	return nil

}

// get monster by ID
func (service *MonsterService) FindMonsterByID(args *FindMonsterByIdArgs, reply *models.Monster) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	monster, found := service.allMonsters[args.ID]
	if !found {
		return MonsterNotFound
	}

	*reply = *monster
	return nil
}

// get all monsters
func (service *MonsterService) ListAllMonsters(args *ListMonstersArgs, reply *[]models.Monster) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	monsters := make([]models.Monster, 0, len(service.allMonsters))
	for _, monster := range service.allMonsters {
		monsters = append(monsters, *monster)
	}

	*reply = monsters
	return nil
}

// arguments structs
type FindMonsterByIdArgs struct {
	ID int
}

type ListMonstersArgs struct {
}


type DeleteAllMonsterArgs struct {
}

type MonsterServiceArgs struct {
	PosX int
	PosY int
	RequestID int
}
