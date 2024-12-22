package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"livestreamall/dao"
	"livestreamall/model"
	"net/http"
	"sync"
	"time"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域请求，生产环境需改为安全校验
	},
}

// BroadcastMessage 广播消息传递
type BroadcastMessage struct {
	LiveRoomID uint   `json:"liveRoomID"`
	Username   string `json:"username"`
	Content    string `json:"content"`
}

// 存储房间的 WebSocket 连接
var rooms = make(map[string]map[*websocket.Conn]bool)
var broadcast = make(chan BroadcastMessage)
var mutex = sync.Mutex{}

// HandleConnections WebSocket 连接处理
func HandleConnections(c *gin.Context) {
	// 获取直播间名称
	streamName := c.Param("stream_name")
	if streamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing stream_name"})
		return
	}

	// 查找直播间
	var existingRoom model.LiveRoom
	err := dao.DB.Where("stream_name = ?", streamName).First(&existingRoom).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "live room not found"})
		return
	}

	// 升级为 WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer ws.Close()

	// 将连接加入直播间
	mutex.Lock()
	if rooms[uintToString(existingRoom.ID)] == nil {
		rooms[uintToString(existingRoom.ID)] = make(map[*websocket.Conn]bool)
	}
	rooms[uintToString(existingRoom.ID)][ws] = true
	mutex.Unlock()

	for {
		// 读取客户端发来的消息
		var message struct {
			Token   string `json:"token"` // 消息中附带的 token
			Content string `json:"Content"`
		}

		if err := ws.ReadJSON(&message); err != nil {
			fmt.Println("读取消息失败:", err)
			// 连接断开，移除连接
			mutex.Lock()
			delete(rooms[uintToString(existingRoom.ID)], ws)
			if len(rooms[uintToString(existingRoom.ID)]) == 0 {
				delete(rooms, uintToString(existingRoom.ID))
			}
			mutex.Unlock()
			break
		}

		// 验证 token
		isValid, claims, err := validateToken(message.Token)
		if !isValid {
			// 如果 token 无效，发送错误消息
			errMsg := gin.H{"error": err.Error()}
			if err := ws.WriteJSON(errMsg); err != nil {
				fmt.Println("发送错误消息失败:", err)
				mutex.Lock()
				delete(rooms[uintToString(existingRoom.ID)], ws)
				if len(rooms[uintToString(existingRoom.ID)]) == 0 {
					delete(rooms, uintToString(existingRoom.ID))
				}
				mutex.Unlock()
				break
			}
			continue
		}

		msg := new(model.Message)

		msg.UserID = claims.UID
		msg.Content = message.Content
		msg.LiveRoomID = existingRoom.ID
		msg.CreatedAt = time.Now()

		// 存储到数据库
		if err := dao.DB.Create(&msg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
			continue
		}

		broadMessage := BroadcastMessage{
			LiveRoomID: existingRoom.ID,
			Username:   claims.Username,
			Content:    message.Content,
		}

		// 广播消息
		broadcast <- broadMessage
	}
}

// 广播消息到对应直播间
func HandleMessages() {
	for {
		msg := <-broadcast

		// 广播到对应直播间
		liveRoomID := msg.LiveRoomID
		mutex.Lock()
		for client := range rooms[uintToString(liveRoomID)] {

			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(rooms[uintToString(liveRoomID)], client)
			}
		}
		mutex.Unlock()
	}
}

func uintFromParam(param string) uint {
	// 将字符串参数转为 uint 类型
	var id uint
	_, _ = fmt.Sscanf(param, "%d", &id)
	return id
}

func uintToString(id uint) string {
	return fmt.Sprintf("%d", id)
}
