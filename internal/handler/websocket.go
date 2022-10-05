package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

//channel for sending everyone message realtime
var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response sent back to websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// WsPayload wsPayload send data to websocket
type WsPayload struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Message  string `json:"message"`
	//current connection
	Conn WebSocketConnection `json:"-"`
}

// WsEndpint WsEndpoint upgrades connrction to websocket
func WsEndpint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</smal></em>`

	//add new person connection to clients map
	conn := WebSocketConnection{Conn : ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WebSocketConnection)  {
	////if go routine stop for any reason comes back
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Println("Error",fmt.Sprintf("%v",r))
	//	}
	//}()

	var payload WsPayload
	
	//infinet lop for all request
	for  {
		err := conn.ReadJSON(&payload)
		if err  != nil{
			//
		}else{
			payload.Conn = *conn
			wsChan <- payload
		}

	}
}
func ListenToWsChannel()  {
	var response WsJsonResponse

	for  {
		//each event
		e := <-wsChan

		switch e.Action {
		case "username":
			//handle user go  to chat with enter username
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)
			//handle when user left chat
		case "left":
			response.Action = "list_users"
			delete(clients,e.Conn)
			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)
			//handle send message
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong> : %s" , e.Username , e.Message)
			broadcastToAll(response)
			

		}
		
		//response.Action = "Got here"
		//response.Message = fmt.Sprintf("Some message and action was %s",e.Action)
		//broadcastToAll(response)
	}
}

func getUserList() []string {
	var userList []string
	for _ , user := range clients {
		if user != "" {
		userList = append(userList , user)
		}
	}
	sort.Strings(userList)
	return userList
}

//broadcast to all
func broadcastToAll( response WsJsonResponse)  {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket err")
			//close connection for user
			_ = client.Close()
			delete(clients,client)
		}
	}
}