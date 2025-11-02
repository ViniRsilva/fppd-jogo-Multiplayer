package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	//Inicializar um objeto do tipo dos metodos exportaveis
	g := NewGame()

	rpc.RegisterName("CoinService", g.CoinService())
	rpc.RegisterName("MonsterService", g.MonsterService())
	rpc.RegisterName("PlayerService", g.PlayerService())
	rpc.RegisterName("PowerService", g.PowerService())
	//Permite que a biblioteca utilize http para comunicacao
	rpc.HandleHTTP()

	/*
		Inicializa um processo que escuta toda comunicacao em
		determinada porta, seguindo o protocolo tcp
	*/
	listener, err := net.Listen("tcp4", ":4040")

	if err != nil {
		log.Fatal("Listener error ", err)
	}
	log.Printf("Serving rpc on port: %d", 4040)

	/*
		Ativa o servidor na porta e com o protocolo definido
		pelo listener
	*/
	http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}
