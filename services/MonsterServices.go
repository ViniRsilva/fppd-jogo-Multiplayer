package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var MonsterNotFound = errors.New("monstro nao encontrado")

type ListAllMonstersArgs struct{}
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

type DeleteAllMonstersArgs struct {
	PlayerID  int
	RequestID int
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

// ListAllMonsters returns all active monsters (used to render the complete map)
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
