package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "math/rand"
    "time"
    "encoding/json"
    "strconv"
)


func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

type Hub struct{
    clients map[*Client] int
    broadcast     chan [] byte
    addClient     chan *Client
    removeClient  chan *Client
    handle_req    chan []byte
    count         int
}
// initialize a new hub
var hub = Hub{
    broadcast:     make(chan []byte, 4096),
    addClient:     make(chan *Client),
    removeClient:  make(chan *Client),
    clients:       make(map[*Client]int),
    handle_req:    make(chan [] byte, 2048),
    count  :       0,
}

var MAX_WORKERS = 5

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}
// Runs forever as a goroutine
func  worker() {
    for {
        // one of these fires when a channel 
        // receives data
        select {
        case data  := <-hub.handle_req:
             var parsed map[string]interface{}
             err := json.Unmarshal(data, &parsed)
             client_id := parse["id"].(int)
             command := parse["command"]
             if err == nil {
                 for conn := range hub.clients {
                    if hub.clients[conn] == client_id {
                        if command == "MOVE_LEFT" {
                            conn.x -= 16
                            new_response := map[string]string{"status": "USER_MOVED_LEFT", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
                            msg,_ := json.Marshal(new_response)
                            hub.broadcast <- msg
                        }

                        if command == "MOVE_RIGHT" {
                            conn.x += 16
                            new_response := map[string]string{"status": "USER_MOVED_RIGHT", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
                            msg,_ := json.Marshal(new_response)
                            hub.broadcast <- msg
                        }
                        
                        if command == "MOVE_UP" {
                            conn.y -= 16
                            new_response := map[string]string{"status": "USER_MOVED_UP", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
                            msg,_ := json.Marshal(new_response)
                            hub.broadcast <- msg
                        }
                        
                        if command == "MOVE_DOWN" {
                            conn.y += 16
                            new_response := map[string]string{"status": "USER_MOVED_DOWN", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
                            msg,_ := json.Marshal(new_response)
                            hub.broadcast <- msg
                        }
                    }
                 }
             }
        case conn := <-hub.addClient:
            conn.x =  random(100, 700)
            conn.y =  random(100, 500)
            // add a new client
            hub.count += 1
            hub.clients[conn] = hub.count
            //let new client know it has been accepted. also, let it know the position and it's id
            response := map[string]string{"status": "OK", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
            msg,_ := json.Marshal(response)
            conn.send <- msg
            broadcast_response := map[string]string{"status": "NEW_USER", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(conn.x), "y":strconv.Itoa(conn.y)}
            msg2,_ := json.Marshal(broadcast_response)
            hub.broadcast <- msg2
        case conn := <-hub.removeClient:
             broadcast_response := map[string]string{"status": "USER_LEFT", "id": strconv.Itoa(hub.clients[conn]), "x":strconv.Itoa(0), "y":strconv.Itoa(0)}
            msg,_ := json.Marshal(broadcast_response)
            hub.broadcast <- msg
            // remove a client
            if _, ok := hub.clients[conn]; ok {
                delete(hub.clients, conn)
                close(conn.send)
            }
        case message := <-hub.broadcast:
            // broadcast a message to all clients
            for conn := range hub.clients {
                select {
                case conn.send <- message:
                default:
                    close(conn.send)
                    delete(hub.clients, conn)
                }
            }
        }
    }
}

type Client struct {
    ws *websocket.Conn
    // Hub passes broadcast messages to this channel
    send chan []byte
    x    int
    y    int
}

// Hub broadcasts a new message and this fires 
func (c *Client) write() {
    // make sure to close the connection incase the loop exits
    defer func() {
        c.ws.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.ws.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.ws.WriteMessage(websocket.TextMessage, message)
        }
    }
}

// New message received so pass it to the Hub 
func (c *Client) read() {
    defer func() {
        hub.removeClient <- c
        c.ws.Close()
    }()

    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            hub.removeClient <- c
            c.ws.Close()
            break
        }

        hub.handle_req <- message
    }
}


func ws_handler(res http.ResponseWriter, req *http.Request){
  conn, err := upgrader.Upgrade(res, req, nil)

  if(err != nil){
    http.NotFound(res, req)
    return
  }

  client := &Client{
        ws:   conn,
        send: make(chan []byte),
    }

    hub.addClient <- client

    go client.write()
    go client.read()
}
func main(){
   for w := 1; w <= MAX_WORKERS; w++ {
       go worker()
   }
  http.Handle("/", http.FileServer(http.Dir("./static")))
  http.HandleFunc("/game", ws_handler)
  http.ListenAndServe(":8080", nil)
}
