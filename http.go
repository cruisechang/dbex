package dbex

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	goHTTP "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type httpParameter struct {
	Address        string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
	Handler        goHTTP.Handler //nil for gorilla mux
}
type httpServer struct {
	server *goHTTP.Server
	router *mux.Router
	done   chan int
}

//NewService make a new SocketManager
func NewHTTPServer(params *httpParameter) (s *httpServer, err error) {

	defer func() {
		if r := recover(); r != nil {

			err = errors.New("httpAPI NewService panic")
		}
	}()

	if params.Port == "" {
		return nil, errors.New("port error ")
	}

	r := mux.NewRouter()

	sv := &httpServer{
		router: r,
		done:   make(chan int),
		server: &goHTTP.Server{
			Handler:        r,
			Addr:           params.Address + ":" + params.Port,
			ReadTimeout:    params.ReadTimeout,
			WriteTimeout:   params.WriteTimeout,
			IdleTimeout:    params.IdleTimeout,
			MaxHeaderBytes: params.MaxHeaderBytes,
		},
	}

	return sv, err

}

//Start starts service.
func (s *httpServer) Start() error {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("dbex start panic=%v", r)
			}
		}()

		if err := s.server.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("server listenAndServer error =%s", err.Error()))
		}
	}()

	return nil
}

//Stop is used to remove lobbyClient
func (s *httpServer) Stop() {
	s.done <- 0
}
func (s *httpServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
	}
	os.Exit(0)
}
func (s *httpServer) handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.done <- 0
	}()
}

func (s *httpServer) GetRouter() *mux.Router {
	return s.router
}

//Hold 如果主程式沒有hold住，就呼叫這個，主程式有就不需要。
func (s *httpServer) Hold() {
	s.handleCtrlC()

	defer func() {
		if r := recover(); r != nil {
			//nx.logger.Log(nxlog.LevelPanic, fmt.Sprintf("mainFlow panic:%v", r))
		}
	}()

	for {
		select {

		//ctrl-c to shtdown
		case <-s.done:
			s.shutdown()
		}
	}

}

func (s *httpServer) GetRealAddr(r *goHTTP.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}
