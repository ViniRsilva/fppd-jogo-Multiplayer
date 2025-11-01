package fppd_jogo_Multiplayer

import (
	"fppd-jogo-Multiplayer/services"
)

type Game struct {
	coinService *services.CoinService
	monsterService *services.MonsterService
	playerService *services.PlayerService
	powerService *services.PowerService
}

func NewGame() *Game {
	return &Game{
		coinService:    services.NewCoinService(),
		playerService:  services.NewPlayerService(),
		monsterService: services.NewMonsterService(),
		powerService:   services.NewPowerService(),
	}
}