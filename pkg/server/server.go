package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/cluttrdev/showdown/pkg/content"
)

type ServerConfig struct {
	Title string
	Style string
}

type Server struct {
	Config   ServerConfig
	Renderer content.Renderer

	sockets map[*websocket.Conn]context.CancelFunc
	cancel  context.CancelFunc
}

func (s *Server) Serve(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.sockets = make(map[*websocket.Conn]context.CancelFunc)
	s.cancel = cancel

	mux, err := s.routes()
	if err != nil {
		return fmt.Errorf("error creating routes: %w", err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	srv.RegisterOnShutdown(s.closeWebSockets)

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()

	// gracefully shutdown server
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("HTTP shutdown error: %w", err)
	}

	return http.ErrServerClosed
}

func (s *Server) routes() (*http.ServeMux, error) {
	sfs, err := staticFileServer()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("/", rootHandler(s.Config.Style))
	mux.Handle("/static/", http.StripPrefix("/static", sfs))
	mux.Handle("/ws", websocket.Handler(s.handleWebSocket))
	mux.HandleFunc("/shutdown", s.handleShutdown)

	return mux, nil
}

func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("error writing http response: %v\n", err)
	}

	// cancel server context on request
	s.cancel()
}

func (s *Server) Update() {
	bytes, err := s.Renderer.Render()
	if err != nil {
		log.Printf("error rendering content: %v", err)
		return
	}

	content := string(bytes)
	for ws := range s.sockets {
		if err := sendMessage(ws, MessageTypeContent, content); err != nil {
			log.Printf("error sending content: %v\n", err)
		}
	}
}
