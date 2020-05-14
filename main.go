package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/youtous/odfe-alerts-handler/handlers"
	"gopkg.in/alecthomas/kingpin.v2"
)

const shutdownTimeout time.Duration = 60 * time.Second

func init() {
	log.SetHeader("${time_rfc3339} ${level}\t${short_file}:${line}\t")
}

func main() {
	hostname, _ := os.Hostname()

	var (
		listenAddress      = kingpin.Flag("web.listen-address", "Address to listen on for incoming HTTP requests.").Default(":8080").Envar("LISTEN_ADDRESS").String()
		smtpHost           = kingpin.Flag("smtp.host", "SMTP server hostname.").Default("localhost").Envar("SMTP_HOSTNAME").String()
		smtpPort           = kingpin.Flag("smtp.port", "SMTP server port.").Default("25").Envar("SMTP_PORT").Int()
		smtpUsername       = kingpin.Flag("smtp.username", "SMTP server login username.").Default("").Envar("SMTP_USERNAME").String()
		smtpPassword       = kingpin.Flag("smtp.password", "SMTP server login password.").Default("").Envar("SMTP_PASSWORD").String()
		smtpFrom           = kingpin.Flag("smtp.from", "SMTP from address.").Default(fmt.Sprintf("opendistro@%s", hostname)).Envar("SMTP_FROM").String()
		smtpDefaultSubject = kingpin.Flag("smtp.default-subject", "SMTP default subject.").Default("Opendistro Alert fired").Envar("SMTP_DEFAULT_SUBJECT").String()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	emailHandler := handlers.Email{
		Host:           *smtpHost,
		Port:           *smtpPort,
		Username:       *smtpUsername,
		Password:       *smtpPassword,
		From:           *smtpFrom,
		DefaultSubject: *smtpDefaultSubject,
	}

	e := echo.New()
	e.HideBanner = true

	e.POST("/email", emailHandler.EchoHandler)

	s := &http.Server{
		Addr: *listenAddress,

		// Good practice to set timeouts :)
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
	}

	go func() {
		if err := e.StartServer(s); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until the signals

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info("shutting down with graceful timeout of", shutdownTimeout)
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
