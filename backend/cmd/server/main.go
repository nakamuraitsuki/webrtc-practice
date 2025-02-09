package main

import (
	"example.com/webrtc-practice/config"
	"example.com/webrtc-practice/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// データベースの初期化
	db, err := sqlx.Open("sqlite3", "practice.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 設定の読み込み
	cfg := config.LoadConfig()

	// サーバーの起動
	server.ServerStart(cfg, db)
}
