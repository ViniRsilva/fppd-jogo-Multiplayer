package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var PlayerNotFound = errors.New("jogador nao encontrado")

type CreatePlayerArgs struct {
	X int
	Y int
}

type FindPlayerByIdArgs struct {
	ID int
}

type DeletePlayerArgs struct {
	ID        int
	RequestID int
}

type MovePlayerArgs struct {
	ID        int
	RequestID int
	X         int
	Y         int
}

type UpdateScoreArgs struct {
	ID        int
	RequestID int
	Delta     int
}

type SetBoostArgs struct {
	ID        int
	RequestID int
	Active    bool
}

type PlayerService struct {
	mu                   sync.RWMutex
	allPlayers           map[int]*models.Player
	nextPlayerID         int
	lastProcessedRequest map[int]int
}

func NewPlayerService() *PlayerService {
	return &PlayerService{
		allPlayers:           make(map[int]*models.Player),
		nextPlayerID:         1,
		lastProcessedRequest: make(map[int]int),
	}
}

// CreatePlayer registra um novo jogador no serviço e retorna seus dados, incluindo o ID atribuído.
func (service *PlayerService) CreatePlayer(args *CreatePlayerArgs, reply *models.Player) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	newPlayer := &models.Player{
		ID:    service.nextPlayerID,
		X:     args.X,
		Y:     args.Y,
		Score: 0,
		Boost: false,
	}

	service.allPlayers[newPlayer.ID] = newPlayer
	service.nextPlayerID++

	*reply = *newPlayer
	return nil
}

// DeletePlayer remove um jogador do mapa.
func (service *PlayerService) DeletePlayer(args *DeletePlayerArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.ID]
	if args.RequestID > lastID {
		_, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		delete(service.allPlayers, args.ID)
		service.lastProcessedRequest[args.ID] = args.RequestID
	}

	*reply = true
	return nil
}

// MovePlayer altera a posição de um jogador.
func (service *PlayerService) MovePlayer(args *MovePlayerArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.ID]
	if args.RequestID > lastID {
		log.Printf("MovePlayer: ID=%d -> (%d,%d)", args.ID, args.X, args.Y)

		player, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		player.X = args.X
		player.Y = args.Y
		service.lastProcessedRequest[args.ID] = args.RequestID
	}

	*reply = true
	return nil
}

// UpdateScore atualiza a pontuação do jogador.
func (service *PlayerService) UpdateScore(args *UpdateScoreArgs, reply *models.Player) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.ID]
	if args.RequestID > lastID {
		player, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		player.Score += args.Delta
		service.lastProcessedRequest[args.ID] = args.RequestID
		*reply = *player
	} else {
		player, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		*reply = *player
	}

	return nil
}

// SetBoost ativa ou desativa o Boost de um jogador.
func (service *PlayerService) SetBoost(args *SetBoostArgs, reply *models.Player) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	lastID := service.lastProcessedRequest[args.ID]
	if args.RequestID > lastID {
		player, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		player.Boost = args.Active
		service.lastProcessedRequest[args.ID] = args.RequestID
		*reply = *player
	} else {
		player, found := service.allPlayers[args.ID]
		if !found {
			return PlayerNotFound
		}
		*reply = *player
	}

	return nil
}

// FindPlayerByID busca um jogador pelo seu ID.
func (service *PlayerService) FindPlayerByID(args *FindPlayerByIdArgs, reply *models.Player) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	player, found := service.allPlayers[args.ID]
	if !found {
		return PlayerNotFound
	}

	*reply = *player
	return nil
}

// ListAllPlayers lista todos os jogadores ativos na sessão.
func (service *PlayerService) ListAllPlayers(_ *struct{}, reply *[]models.Player) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	players := make([]models.Player, 0, len(service.allPlayers))
	for _, player := range service.allPlayers {
		players = append(players, *player)
	}

	*reply = players
	return nil
}
