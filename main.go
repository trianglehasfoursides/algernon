package main

import (
	"fmt"
	"net"

	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	kp, err := keygen.New("awesome", keygen.WithPassphrase("halo"), keygen.WithKeyType(keygen.Ed25519))
	if err != nil {
		log.Error(err.Error())
		return
	}

	server, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("localhost", "7000")),
		wish.WithHostKeyPEM(kp.RawPrivateKey()),
		wish.WithMiddleware(
			bubbletea.Middleware(WishForm),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error(err.Error())
		return
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err.Error())
		return
	}
}
