package adapter_ws

import (
	"context"
	"feeder/internal/entity"
	"feeder/pkg/tools"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

type openConnectionHandler func(ctx context.Context, ip string, token string) (id string, err error)
type closeConnectionHandler func(ctx context.Context, ip, id string)

type Handler struct {
	*http.ServeMux
	connections sync.Map // map[*websocket.Conn]context.CancelFunc

	openConnection  openConnectionHandler
	closeConnection closeConnectionHandler
}

func NewHandler(
	ctx context.Context,
	path string,
	openConnection openConnectionHandler,
	closeConnection closeConnectionHandler,
) *Handler {
	var result = &Handler{
		ServeMux:        http.NewServeMux(),
		connections:     sync.Map{},
		openConnection:  openConnection,
		closeConnection: closeConnection,
	}
	result.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		childCtx, cancel := context.WithCancel(ctx)
		result.handler(childCtx, cancel, w, r)
	})
	return result
}

func (h *Handler) PublishSingleTrade(ctx context.Context, trade *entity.SingleTrade) error {
	if trade == nil {
		return nil
	}
	h.connections.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)
		cancel := value.(context.CancelFunc)

		buff, err := protojson.MarshalOptions{
			UseProtoNames:     true,
			EmitDefaultValues: true,
		}.Marshal(trade)
		if err != nil {
			cancel()
			return true
		}

		if err := conn.Write(ctx, websocket.MessageText, buff); err != nil {
			cancel()
			return true
		}
		return true
	})
	return nil
}

func (h *Handler) addConnection(c *websocket.Conn, cancel context.CancelFunc) func() {
	h.connections.Store(c, cancel)
	return func() {
		h.connections.Delete(c)
	}
}

func (h *Handler) getIP(r *http.Request) string {
	var remoteIP = r.Header.Get("X-Real-IP")
	if remoteIP == "" {
		remoteIP = r.Header.Get("X-Forwarded-For")
	}
	if remoteIP == "" {
		remoteIP = r.Header.Get("X-Client-IP")
	}
	if remoteIP == "" {
		remoteIP = r.Header.Get("CF-Connecting-IP")
	}
	if remoteIP == "" {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}

func (h *Handler) handler(ctx context.Context, cancelCtx context.CancelFunc, w http.ResponseWriter, r *http.Request) {
	defer cancelCtx()

	authToken := r.URL.Query().Get("token")
	if authToken == "" {
		http.Error(w, "Missing auth token", http.StatusUnauthorized)
		return
	}

	ip := h.getIP(r)
	id, err := h.openConnection(ctx, ip, authToken)
	if err != nil {
		http.Error(w, "Failed to open connection", http.StatusUnauthorized)
		return
	}

	connection, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		},
	)
	if err != nil {
		return
	}
	defer func() {
		connection.Close(websocket.StatusNormalClosure, "closing")
		h.closeConnection(ctx, ip, id)
	}()

	connection.SetReadLimit(1 << 20)

	deleteConn := h.addConnection(connection, cancelCtx)
	defer deleteConn()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ctx = connection.CloseRead(ctx)
	for range tools.IteratorInt64WithContext(ctx) {
		select {
		case <-ticker.C:
			if err := connection.Ping(ctx); err != nil {
				return
			}
		}
	}
}
