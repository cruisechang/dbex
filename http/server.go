package http

import (
	"context"
	"errors"
	goHTTP "net/http"
	"time"
	"github.com/gorilla/mux"
	"os"
	"strings"
	"net"
	"os/signal"
	"syscall"
	"fmt"
)

//Service is a httpApi service for http request operation.
//type Server interface {
//	Start()error
//	Stop()
//	GetRealAddr(r *goHTTP.Request) string
//	RegisterHandler(pattern string, handler goHTTP.Handler)
//
//}
type ServerParameters struct {
	Address        string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
	Handler        goHTTP.Handler //nil for gorilla mux
}
type Server struct {
	server *goHTTP.Server
	router *mux.Router
	done   chan int
}

//NewService make a new SocketManager
func NewServer(params *ServerParameters) (s *Server, err error) {

	defer func() {
		if r := recover(); r != nil {

			err = errors.New("httpAPI NewService panic")
		}
	}()

	if params.Port == "" {
		return nil, errors.New("port error ")
	}

		r := mux.NewRouter()

	sv := &Server{
		router: r,
		done:   make(chan int),
		server: &goHTTP.Server{
			Handler:        r,
			Addr:           params.Address + ":" + params.Port,
			ReadTimeout:    params.ReadTimeout,
			WriteTimeout:   params.WriteTimeout,
			IdleTimeout:    params.IdleTimeout,
			MaxHeaderBytes: params.MaxHeaderBytes,
			//MaxHeaderBytes: 1 << 20,
		},
	}

	return sv, err

}

//Start starts service.
func (s *Server) Start() error{
	go func(){
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		if err:=s.server.ListenAndServe();err!=nil{
			panic(fmt.Sprintf("server listenAndServer error =%s",err.Error()))
		}
	}()
	return nil
}

//Stop is used to remove lobbyClient
func (s *Server) Stop() {
	s.done <- 0
}
func (s *Server) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
	}
	os.Exit(0)
}
func (s *Server) handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.done <- 0
	}()
}

func (s *Server)GetRouter()*mux.Router{
	return s.router
}

//func (s *Server) HandleFunc(path string, f func(goHTTP.ResponseWriter, *goHTTP.Request)) *mux.Route {
//
//	return s.router.HandleFunc(path, f)
//
//}
//
//func (s *Server) Handle(path string, handler goHTTP.Handler) *mux.Route {
//
//	return s.router.Handle(path, handler)
//
//}
//func (s *Server)Methods(method string)*mux.Route{
//	return s.router.Methods(method)
//}


//Hold 如果主程式沒有hold住，就呼叫這個，主程式有就不需要。
func (s *Server) Hold() {
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

/*
func requestHandler(w goHTTP.ResponseWriter, r *goHTTP.Request) {

	if r := recover(); r != nil {
	}

	var params = ""
	var path = ""

	if r.Method == goHTTP.MethodGet {
		params = r.URL.Query().Get("params")
	} else {
		params = strings.TrimSpace(r.FormValue("params"))
	}

	fmt.Printf("Method : %s %s \n", r.Method, params)

	if len(params) == 0 {
		goHTTP.Error(w, "httpAPI params len==0", goHTTP.StatusBadRequest)
		return
	}

	//b, err := encode.Base64Decode(params)
	b, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		goHTTP.Error(w, "httpAPI params base64 decode error", goHTTP.StatusBadRequest)
		return
	}

	// fmt.Printf("params : %v \n", string(b))

	path = r.URL.Path
	path = strings.TrimPrefix(r.URL.Path, "/")

	re := Response{}
	if err := json.Unmarshal(b, &re);err!=nil{

	}

	//encode.JSONDecodeToStruct(b, &re)

	fmt.Printf("path :%s Response : %+v \n", path, re)

	//pass path(grabRedEnv) and Response
	//ex: response=map["grabRedEnv"]
	writeParamsChan(map[string]Response{path: re})

	// w.Write(b)

}
*/

func (s *Server) GetRealAddr(r *goHTTP.Request) string {

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
