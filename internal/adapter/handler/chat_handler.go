package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"tro-go/internal/domain"
	"tro-go/internal/port"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 // Tăng lên 1KB cho tin nhắn
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Cần cấu hình cụ thể khi chạy production
	},
}

type Client struct {
	hub    *ChatHub
	conn   *websocket.Conn
	send   chan domain.ChatPayload
	userID int64
}

type ChatHub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	usecase    port.ChatUseCase
}

func NewChatHub(usecase port.ChatUseCase) *ChatHub {
	return &ChatHub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		usecase:    usecase,
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var payload domain.ChatPayload
		if err := json.Unmarshal(message, &payload); err != nil {
			continue
		}

		if payload.Type == "message" && payload.ReceiverID != nil {
			msg, err := c.hub.usecase.SendMessage(context.Background(), c.userID, *payload.ReceiverID, payload.Content)
			if err != nil {
				log.Println("Error saving message:", err)
				continue
			}

			c.hub.mu.Lock()
			if receiver, ok := c.hub.clients[*payload.ReceiverID]; ok {
				receiver.send <- domain.ChatPayload{
					Type:     "message",
					SenderID: c.userID,
					Content:  payload.Content,
					Data:     msg,
				}
			}
			c.hub.mu.Unlock()
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case payload, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(payload); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type ChatHandler struct {
	hub     *ChatHub
	usecase port.ChatUseCase
}

func NewChatHandler(group *echo.Group, hub *ChatHub, usecase port.ChatUseCase) {
	h := &ChatHandler{
		hub:     hub,
		usecase: usecase,
	}
	group.GET("/ws/chat", h.HandleWebSocket)
}

func (h *ChatHandler) HandleWebSocket(c echo.Context) error {
	// 1. Lấy UserID (Thử từ Middleware trước, nếu không có thì thử từ Query)
	var userID int64
	userToken, ok := c.Get("user").(*jwt.Token)
	if ok {
		claims := userToken.Claims.(*jwt.MapClaims)
		userID = int64((*claims)["id"].(float64))
	} else {
		// Nếu Middleware không bắt được (do thiếu Header), check thử query ?token=
		tokenString := c.QueryParam("token")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		// Ở đây bạn có thể thêm logic verify token thủ công nếu muốn cực kỳ chặt chẽ
		// Nhưng hiện tại để đơn giản, ta yêu cầu token phải hợp lệ qua Middleware trước.
		// LƯU Ý: Cách tốt nhất là cập nhật main.go để Middleware tự xử lý cả 2 nguồn.
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token must be valid"})
	}

	// 2. Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &Client{
		hub:    h.hub,
		conn:   conn,
		send:   make(chan domain.ChatPayload, 256),
		userID: userID,
	}
	h.hub.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}
