package server

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/internal/handler"
	"example.com/webrtc-practice/internal/infrastructure/repository/sqlite3"
	"example.com/webrtc-practice/internal/infrastructure/service/hasher"
	"example.com/webrtc-practice/internal/infrastructure/service/jwt"
	"example.com/webrtc-practice/routes"
	
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"



	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}

	clients              = make(map[*websocket.Conn]string)
	clientsByID          = make(map[string]*websocket.Conn)
	broadcast            = make(chan []byte)
	offerId       string = ""
	functions            = make(map[string]interface{})
	sdpData              = make(map[string]string)
	candidateData        = make(map[string][]string)
)

func ServerStart(cfg *config.Config, db *sqlx.DB) {
	e := echo.New()

	functions["connect"] = connect
	functions["offer"] = offer
	functions["answer"] = answer
	functions["candidateAdd"] = candidateAdd

	http.HandleFunc("/ws", handleWebSocket)
	go handleMessages()
	
	// ミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},  // すべてのオリジンを許可
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE}, // 許可するHTTPメソッド
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization}, // 許可するHTTPヘッダー
	}))
	
	// ユーザーハンドラーの初期化
	userRepository := sqlite3.NewUserRepository(db)
	hasher := hasher.NewBcryptHasher()
	tokenService := jwt.NewJWTService(cfg.SecretKey, cfg.TokenExpiry)
	userHandler := handler.NewUserHandler(userRepository, hasher, tokenService)
	
	routes.SetupRoutes(e, cfg, userHandler)
	
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	
	for {
		_, message, err := conn.ReadMessage()
		if _, ok := clients[conn]; !ok {
			// 新規接続
			var jsonStr = string(message)
			var data map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &data)
			if err != nil {
				panic(err)
			}

			// idの登録
			id := data["id"].(string)
			clients[conn] = id
			clientsByID[id] = conn
		}

		if err != nil {
			log.Println(err)
			delete(clientsByID, clients[conn])
			delete(clients, conn)
			break
		}

		broadcast <- message
	}
	offerId = ""
	defer conn.Close()
}

func handleMessages() {
	for {
		message := <-broadcast
		// text -> json
		var jsonStr = string(message)
		fmt.Println(jsonStr)
		var data map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			panic(err)
		}

		// 処理分岐
		msgDataType := data["type"].(string)
		function := functions[msgDataType].(func(map[string]interface{}))
		function(data)
	}
}

func connect(data map[string]interface{}) {
	resultData := make(map[string]string)
	
	id := data["id"].(string)
	client := clientsByID[id]

	// offerを送ってもらう
	if len(offerId) == 0 {
		offerId = id
		resultData["type"] = "offer"
		bytes := jsonToBytes(resultData)
		sendMessage(client, bytes)
		return
	} else if id == offerId {
		// 重複
		return
	}

	// offerを送る
	resultData["type"] = "offer"
	resultData["sdp"] = sdpData[offerId]
	resultData["target_id"] = offerId
	bytes := jsonToBytes(resultData)
	sendMessage(client, bytes)
}

func offer(data map[string]interface{}) {
	fmt.Println("[Offer]")
	id := data["id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	sdpData[id] = string(sdp)
}

func answer(data map[string]interface{}) {
	// offerの送り主にanswerを返す
	sendAnswer(data)

	// answerの送り主にcandidateを送る
	sendCandidate(data)
}

func sendAnswer(data map[string]interface{}) {
	fmt.Println("[Answer]")
	resultData := make(map[string]string)
	resultData["type"] = "answer"
	target_id := data["target_id"].(string)
	sdp, _ := json.Marshal(data["sdp"])
	resultData["sdp"] = string(sdp)

	client := clientsByID[target_id]
	bytes := jsonToBytes(resultData)
	sendMessage(client, bytes)
}

func sendCandidate(data map[string]interface{}) {
	returnData := make(map[string]string)
	id := offerId
	if _, ok := candidateData[id]; !ok {
		return
	}

	answerId := data["id"].(string)
	client := clientsByID[answerId]
	fmt.Println("candidate受け取り")
	fmt.Println("[Candidate]")
	returnData["type"] = "candidate"
	returnData["candidate"] = strings.Join(candidateData[id], "|")
	bytes := jsonToBytes(returnData)
	sendMessage(client, bytes)

}

func candidateAdd(data map[string]interface{}) {
	fmt.Println("[Candidate Add]")
	resultData := make(map[string]string)

	// 相手が已經接続的話、candidateDataに入れずに直接送る
	id := data["id"].(string)
	candidateByte, _ := json.Marshal(data["candidate"])
	candidate := string(candidateByte)

	target_id := data["target_id"].(string)
	if target_id != "" {
		if client, ok2 := clientsByID[target_id]; ok2 {
			// 相手が有的話
			fmt.Println("[Candidate]")
			resultData["type"] = "candidate"
			resultData["candidate"] = candidate
			bytes := jsonToBytes(resultData)
			sendMessage(client, bytes)
			return
		}
	}

	// 相手が還沒來 -> 保存
	if _, ok := candidateData[id]; !ok {
		candidateData[id] = []string{candidate}
	} else {
		candidateData[id] = append(candidateData[id], candidate)
	}
}

// 訊息送信
func sendMessage(client *websocket.Conn, bytes []byte) {
	err := client.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println(err)
		client.Close()
		delete(clients, client)
	}
}

func jsonToBytes(result map[string]string) []byte {
	jsonText, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	bytes := []byte(jsonText)
	return bytes
}
