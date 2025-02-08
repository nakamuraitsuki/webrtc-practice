package signaling

import "github.com/gorilla/websocket"

type SignalingHandler struct {
	manager *SignalingManager
}

func NewSignalingHandler(manager *SignalingManager) *SignalingHandler {
	return &SignalingHandler{
		manager: manager,
	}
}

func (h *SignalingHandler) HandleMessages(conn *websocket.Conn) {
	defer func() {
		h.manager.RemoveClient(conn)//クライアント削除
		conn.Close()				//接続を閉じる
	}()
	
	//メッセージ待機ループとメッセージ処理
	for {
		//メッセージの読み込み
	}
}