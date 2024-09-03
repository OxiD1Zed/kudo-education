package main

import (
	"kode-education/internal/app"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	port, err := strconv.Atoi(os.Getenv("GOSERVERPORT"))
	if err != nil || port < 0 {
		log.Warn("serverport is incorrect", slog.String("error", err.Error()))
		panic(err)
	}

	app := app.New(log, uint(port))
	app.RestApp.Run()
}
