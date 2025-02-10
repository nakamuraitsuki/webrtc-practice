#!/bin/sh

# プロセスのクリーンアップ用の関数
cleanup() {
  echo "Stopping server and frontend..."
  kill "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null
  wait "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null
  echo "Shutdown complete."
  exit 0
}

# SIGINT (Ctrl+C) を受け取ったら cleanup() を実行
trap cleanup INT

# 既に 8080 番ポートが使用されているか確認し、もしあれば停止
if lsof -i :8080 >/dev/null; then
  echo "Port 8080 is already in use. Stopping existing process..."
  kill -9 $(lsof -t -i :8080)
  sleep 2
fi

# バックエンド起動
echo "Starting server..."
cd ./backend
go run ./cmd/server/main.go &
BACKEND_PID=$!

# フロントエンド起動
cd ../frontend
echo "Starting frontend..."
npm run dev &
FRONTEND_PID=$!

# 子プロセスの終了を待つ
wait "$BACKEND_PID" "$FRONTEND_PID"
