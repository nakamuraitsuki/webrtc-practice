package websocketmanager_test

import (
	"encoding/json"
	"errors"
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/infrastructure/dto"
	websocketmanager "example.com/webrtc-practice/internal/infrastructure/service_impl/websocket_manager"
	mock_adapter "example.com/webrtc-practice/mocks/interface/adapter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReadMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_adapter.NewMockConnAdapter(ctrl)
	conn := websocketmanager.NewWebsocketConnection(mockConn)

	t.Run("正常にメッセージを読み込める", func(t *testing.T) {
		testMsg, _ := entity.NewMessage("123", "connection", "sdp", []string{"candidate", "candidate2"}, "456")
		testMsgDto := dto.WebsocketMessageDTO{}
		testMsgDto.FromEntity(testMsg)
		response, _ := json.Marshal(&testMsgDto)
		mockConn.EXPECT().ReadMessageFunc().Return(1, response, nil).Times(1)

		// テスト実行
		msgType, msg, err := conn.ReadMessage()

		assert.NoError(t, err)
		assert.Equal(t, 1, msgType)
		assert.Equal(t, testMsg.GetID(), msg.GetID())
		assert.Equal(t, testMsg.GetType(), msg.GetType())
		assert.Equal(t, testMsg.GetSDP(), msg.GetSDP())
		assert.Equal(t, testMsg.GetCandidate(), msg.GetCandidate())
		assert.Equal(t, testMsg.GetTargetID(), msg.GetTargetID())
	})

	t.Run("ReadMessageFuncがエラーを返す場合", func(t *testing.T) {
		mockConn.EXPECT().ReadMessageFunc().Return(0, nil, errors.New("read error")).Times(1)

		_, _, err := conn.ReadMessage()
		assert.Error(t, err)
	})

	t.Run("不正なJSONを返された場合", func(t *testing.T) {
		invalidJSON := []byte(`{"invalid": "json"`)
		mockConn.EXPECT().ReadMessageFunc().Return(1, invalidJSON, nil).Times(1)

		_, _, err := conn.ReadMessage()
		assert.Error(t, err)
	})
}

func TestWriteMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_adapter.NewMockConnAdapter(ctrl)
	conn := websocketmanager.NewWebsocketConnection(mockConn)

	t.Run("正常にメッセージを書き込める", func(t *testing.T) {
		testMsg, _ := entity.NewMessage("123", "connection", "sdp", []string{"candidate", "candidate2"}, "456")

		mockConn.EXPECT().WriteMessageFunc(gomock.Any(), gomock.Any()).Do(func(messageType int, data []byte) {
			msg := dto.WebsocketMessageDTO{}
			err := json.Unmarshal(data, &msg)

			assert.NoError(t, err)
			assert.Equal(t, testMsg.GetID(), msg.ID)
			assert.Equal(t, testMsg.GetType(), msg.Type)
			assert.Equal(t, testMsg.GetSDP(), msg.SDP)
			assert.Equal(t, testMsg.GetCandidate(), msg.Candidate)
			assert.Equal(t, testMsg.GetTargetID(), msg.TargetID)
		}).Return(nil).Times(1)

		err := conn.WriteMessage(*testMsg)

		assert.NoError(t, err)
	})

	t.Run("WriteMessageFuncがエラーを返す場合", func(t *testing.T) {
		testMsg, _ := entity.NewMessage("", "", "", nil, "")

		mockConn.EXPECT().WriteMessageFunc(gomock.Any(), gomock.Any()).Return(errors.New("write error")).Times(1)

		err := conn.WriteMessage(*testMsg)
		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mock_adapter.NewMockConnAdapter(ctrl)
	conn := websocketmanager.NewWebsocketConnection(mockConn)

	t.Run("正常にCloseできる", func(t *testing.T) {
		mockConn.EXPECT().CloseFunc().Return(nil).Times(1)
		err := conn.Close()
		assert.NoError(t, err)
	})

	t.Run("CloseFuncがエラーを返す場合", func(t *testing.T) {
		mockConn.EXPECT().CloseFunc().Return(errors.New("close error")).Times(1)
		err := conn.Close()
		assert.Error(t, err)
	})
}
