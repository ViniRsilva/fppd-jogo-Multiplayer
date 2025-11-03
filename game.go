package main

import (
	"fppd-jogo-Multiplayer/services"
)

type Game struct {
	coinService    *services.CoinService
	monsterService *services.MonsterService
	playerService  *services.PlayerService
	powerService   *services.PowerService
}

func NewGame() *Game {
	return &Game{
		coinService:    services.NewCoinService(),
		monsterService: services.NewMonsterService(),
		playerService:  services.NewPlayerService(),
		powerService:   services.NewPowerService(),
	}
}

func (g *Game) CoinService() *services.CoinService {
	return g.coinService
}

func (g *Game) MonsterService() *services.MonsterService {
	return g.monsterService
}

func (g *Game) PlayerService() *services.PlayerService {
	return g.playerService
}

func (g *Game) PowerService() *services.PowerService {
	return g.powerService
}
