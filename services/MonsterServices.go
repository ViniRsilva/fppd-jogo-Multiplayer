package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var MonsterNotFound = errors.New("monstro nao encontrado")

type MonsterService struct {
	mu               sync.RWMutex
	allMonsters      map[int]*models.Monster
	nextMonsterID    int
	lastKnownRequest map[int]int
}

func NewMonsterService() *MonsterService {
	return &MonsterService{
		allMonsters:      make(map[int]*models.Monster),
		nextMonsterID:    1,
		lastKnownRequest: make(map[int]int),
	}
}

type CreateMonsterArgs struct {
	PlayerID  int
	RequestID int
	X         int
	Y         int
}

type DeleteMonsterByPositionArgs struct {
	PlayerID  int
	RequestID int
	X         int
	Y         int
}

type DeleteAllMonstersArgs struct {
	PlayerID  int
	RequestID int
}

type FindMonsterByIdArgs struct {
	ID int
}

// CreateMonster creates a new monster
func (service *MonsterService) CreateMonster(args *CreateMonsterArgs, reply *models.Monster) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastRequestId := service.lastKnownRequest[args.PlayerID]

	if args.RequestID > lastRequestId {
		newMonster := &models.Monster{
			ID: service.nextMonsterID,
			X:  args.X,
			Y:  args.Y,
		}
		service.nextMonsterID++
		service.allMonsters[newMonster.ID] = newMonster
		service.lastKnownRequest[args.PlayerID] = args.RequestID
		log.Printf("Created Monster: (%d,%d)", args.X, args.Y)

		*reply = *newMonster
	}

	return nil
}

// DeleteMonsterByPosition deletes a monster by its position
func (service *MonsterService) DeleteMonsterByPosition(args *DeleteMonsterByPositionArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastRequestId := service.lastKnownRequest[args.PlayerID]

	if args.RequestID > lastRequestId {
		var foundMonsterID = -1
		for id, monster := range service.allMonsters {
			if monster.X == args.X && monster.Y == args.Y {
				foundMonsterID = id
				break
			}
		}

		if foundMonsterID == -1 {
			return MonsterNotFound
		}

		delete(service.allMonsters, foundMonsterID)
		log.Printf("Deleted Monster: (%d,%d)", args.X, args.Y)

		service.lastKnownRequest[args.PlayerID] = args.RequestID
	}

	*reply = true
	return nil
}

// DeleteAllMonsters deletes all monsters (executed when player gets a power)
func (service *MonsterService) DeleteAllMonsters(args *DeleteAllMonstersArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastRequestId := service.lastKnownRequest[args.PlayerID]

	if args.RequestID > lastRequestId {
		service.allMonsters = make(map[int]*models.Monster)
		service.lastKnownRequest[args.PlayerID] = args.RequestID
	}

	*reply = true
	return nil
}

// FindMonsterByID finds a monster by ID (not being used)
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

type ListAllMonstersArgs struct{}

func (service *MonsterService) ListAllMonsters(args *ListAllMonstersArgs, reply *[]models.Monster) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	monsters := make([]models.Monster, 0, len(service.allMonsters))
	for _, monster := range service.allMonsters {
		monsters = append(monsters, *monster)
	}

	*reply = monsters
	return nil
}
