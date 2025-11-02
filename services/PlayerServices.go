package services

import (
	"errors"
	"fppd-jogo-Multiplayer/models"
	"log"
	"sync"
)

var PlayerNotFound = errors.New("Jogador não encontrado!")

type PlayerService struct {
	mu           sync.RWMutex
	allPlayers   map[int]*models.Player
	nextPlayerID int
}

func NewPlayerService() *PlayerService {
	return &PlayerService{
		allPlayers:   make(map[int]*models.Player),
		nextPlayerID: 1,
	}
}

/*
Cria um novo jogador com posição inicial X, Y
*/
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

/*
Busca jogador pelo ID
*/
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

/*
Remove jogador do mapa
*/
func (service *PlayerService) DeletePlayer(args *DeleteArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.lastRequestID[args.RequestID] {
		*reply = true
		return nil
	}

	_, found := service.allPlayers[args.ID]
	if !found {
		return PlayerNotFound
	}

	delete(service.allPlayers, args.ID)
	service.lastRequestID[args.RequestID] = true
	*reply = true
	return nil
}

/*
Move jogador alterando sua posição
*/
func (service *PlayerService) MovePlayer(args *MoveArgs, reply *bool) error {
	service.mu.Lock()
	defer service.mu.Unlock()
	
	if service.lastRequestID[args.RequestID] {
		*reply = true
		return nil
	}

	log.Printf("MovePlayer: ID=%d -> (%d,%d) | totalPlayers=%d", args.ID, args.X, args.Y, len(service.allPlayers))

	player, found := service.allPlayers[args.ID]
	if !found {
		return PlayerNotFound
	}

	player.X = args.X
	player.Y = args.Y
	service.lastRequestID[args.RequestID] = true
	*reply = true
	return nil
}

/*
Atualiza a pontuação do jogador
*/
func (service *PlayerService) UpdateScore(args *ScoreArgs, reply *models.Player) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	player, found := service.allPlayers[args.ID]
	if !found {
		return PlayerNotFound
	}

	player.Score += args.Delta
	*reply = *player
	return nil
}

/*
Ativa ou desativa o Boost.
*/
func (service *PlayerService) SetBoost(args *BoostArgs, reply *models.Player) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	player, found := service.allPlayers[args.ID]
	if !found {
		return PlayerNotFound
	}

	player.Boost = args.Active
	*reply = *player
	return nil
}

/*
Lista todos os jogadores ativos
*/
func (service *PlayerService) ListAllPlayers(args *struct{}, reply *[]models.Player) error {
	service.mu.RLock()
	defer service.mu.RUnlock()

	players := make([]models.Player, 0, len(service.allPlayers))
	for _, player := range service.allPlayers {
		players = append(players, *player)
	}

	*reply = players
	return nil
}

/*
Tipos auxiliares para chamadas RPC
*/

type MoveArgs struct {
	ID int
	X  int
	Y  int
	RequestID int
}

type ScoreArgs struct {
	ID    int
	Delta int
}

type BoostArgs struct {
	ID     int
	Active bool
}

type CreatePlayerArgs struct {
	X int
	Y int
}

type FindPlayerByIdArgs struct {
	ID int
}

type DeleteArgs struct {
	ID int
	RequestID int
}
