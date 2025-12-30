package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sora/bot"
	"sora/config"

	_ "sora/commands"
	"sora/commands/gemini"
)

func main() {

	cfg := config.LoadConf()
	log.Println("Configuration loaded.")
	gemini.InitGemini()
	myBot := bot.NewBot(cfg)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Menyalakan bot...")
		err := myBot.Start()
		if err != nil {
			log.Fatalf("Gagal memulai bot : %v", err)
		}
	}()
	log.Println("Bot berhasil dinyalakan. Menunggu sinyal shutdown (Ctrl+C)...")

	<-c

	log.Println("Sinyal shutdown diterima, mematikan bot...")
	myBot.Disconnect()
	log.Println("Shutting down... Bye Bye!")
}
