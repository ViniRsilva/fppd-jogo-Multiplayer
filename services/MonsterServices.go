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
}

func NewMonsterService() *MonsterService {
	return &MonsterService{
		allMonsters:   make(map[int]*models.Monster),
		nextMonsterID: 1,
	}
}

// Create a new monster
func (service *MonsterService) CreateMonster(args *CreateMonsterArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	newMonster := &models.Monster{
		ID: service.nextMonsterID,
		X:  args.X,
		Y:  args.Y,
	}
	service.nextMonsterID++
	service.allMonsters[newMonster.ID] = newMonster

	*reply = true
	return nil
}

// Delete monster by ID
func (service *MonsterService) DeleteMonster(args *DeleteMonsterArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	_, found := service.allMonsters[args.ID]
	if !found {
		return MonsterNotFound
	}

	delete(service.allMonsters, args.ID)
	*reply = true
	return nil
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
type CreateMonsterArgs struct {
	X int
	Y int
}

type FindMonsterByIdArgs struct {
	ID int
}

type ListMonstersArgs struct {
}

type DeleteMonsterArgs struct {
	ID int
}

type DeleteAllMonsterArgs struct {
}
