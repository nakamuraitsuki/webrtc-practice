package usecase_test

import (
	"fmt"
	"sync"
	"testing"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_repository "example.com/webrtc-practice/mocks/repository"
	mock_service "example.com/webrtc-practice/mocks/service"
)

// 新しいUsecaseインスタンスを作成するテスト
func TestNewWebsocketUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	// Test（インスタンスが作成できているか確認）
	assert.NotNil(t, usecase)
}

// RegisterClientメソッドのテスト
func TestRegisterClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	mockWm.EXPECT().RegisterConnection(mockConn).Return(nil)

	// Test（RegisterClientメソッドの呼び出し）
	err := usecase.RegisterClient(mockConn)
	assert.NoError(t, err)
}

// ゴルーチンで呼ばれるListenForMessagesメソッドのテスト
func TestListenForMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	testMessage := entity.Message{
		ID:        "testID",
		Type:      "connect",
		SDP:       "testSDP",
		Candidate: []string{"testCandidate"},
		TargetID:  "targetID",
	}

	t.Run("正常系", func(t *testing.T) {
		mockConn.EXPECT().
			ReadMessage().
			Return(1, testMessage, nil).
			Times(1)

		mockWm.EXPECT().
			ExistsByID(testMessage.ID).
			Return(false).
			Times(1)

		mockWm.EXPECT().
			RegisterID(mockConn, testMessage.ID).
			Times(1)
		mockBr.EXPECT().
			Send(testMessage).
			Times(1)

		mockRepo.EXPECT().
			CreateClient(testMessage.ID).
			Times(1)

		mockConn.EXPECT().
			ReadMessage().
			Return(0, entity.Message{}, assert.AnError).
			Times(1)

		mockWm.EXPECT().
			DeleteConnection(mockConn).
			Times(1)

		mockRepo.EXPECT().
			DeleteClient(testMessage.ID).
			Times(1)

		mockO.EXPECT().
			ClearOffer().
			Times(1)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			usecase.ListenForMessages(mockConn)
		}()
		wg.Wait()
	})

	t.Run("ID既登録時の接続拒否", func(t *testing.T) {
		mockConn.EXPECT().
			ReadMessage().
			Return(0, testMessage, nil).
			Times(1)

		mockWm.EXPECT().
			ExistsByID(testMessage.ID).
			Return(true).
			Times(1)

		mockWm.EXPECT().
			DeleteConnection(mockConn).
			Times(1)

		mockO.EXPECT().
			ClearOffer().
			Times(1)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			usecase.ListenForMessages(mockConn)
		}()
		wg.Wait()
	})

	t.Run("ReadMessageがError", func(t *testing.T) {
		mockConn.EXPECT().
			ReadMessage().
			Return(0, entity.Message{}, assert.AnError).
			Times(1)

		mockWm.EXPECT().
			DeleteConnection(mockConn).
			Times(1)

		mockRepo.EXPECT().
			DeleteClient("").
			Times(1)

		mockO.EXPECT().
			ClearOffer().
			Times(1)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			usecase.ListenForMessages(mockConn)
		}()
		wg.Wait()
	})

	t.Run("ID登録後のループ", func(t *testing.T) {
		gomock.InOrder(
			// 初回メッセージ（ID登録）
			mockConn.EXPECT().
				ReadMessage().
				Return(1, testMessage, nil),

			mockWm.EXPECT().
				ExistsByID(testMessage.ID).
				Return(false),

			mockWm.EXPECT().
				RegisterID(mockConn, testMessage.ID),

			mockRepo.EXPECT().
				CreateClient(testMessage.ID),

			mockBr.EXPECT().
				Send(testMessage),

			// 2回目のメッセージ（すでにID登録済み）
			mockConn.EXPECT().
				ReadMessage().
				Return(1, testMessage, nil),
			mockBr.EXPECT().
				Send(testMessage),

			// 終了条件
			mockConn.EXPECT().
				ReadMessage().
				Return(0, entity.Message{}, assert.AnError),

			mockWm.EXPECT().
				DeleteConnection(mockConn),

			mockRepo.EXPECT().
				DeleteClient(testMessage.ID),

			mockO.EXPECT().
				ClearOffer(),
		)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			usecase.ListenForMessages(mockConn)
		}()
		wg.Wait()
	})
}

func TestConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)
	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	testConnectMessage := entity.Message{
		ID:        "testID",
		Type:      "connect",
		SDP:       "",
		Candidate: nil,
		TargetID:  "",
	}

	t.Run("誰もOfferしていない場合、接続者にoffer要求を送る", func(t *testing.T) {
		msgToSend := entity.Message{
			ID:   testConnectMessage.ID,
			Type: "offer",
		}

		// 期待される動作
		mockWm.EXPECT().
			GetConnectionByID(testConnectMessage.ID).
			Return(mockConn, nil)

		mockO.EXPECT().
			IsOffer().
			Return(false)

		mockO.EXPECT().
			SetOffer(testConnectMessage.ID)

		mockConn.EXPECT().
			WriteMessage(msgToSend).
			Times(1)

		// テスト実行
		usecase.Connect(testConnectMessage)
	})

	t.Run("offerが自分自身だった場合、何もしない", func(t *testing.T) {
		mockWm.EXPECT().
			GetConnectionByID(testConnectMessage.ID).
			Return(mockConn, nil)

		mockO.EXPECT().
			IsOffer().
			Return(true)

		mockO.EXPECT().
			IsOfferID(testConnectMessage.ID).
			Return(true)

		// WriteMessageは呼ばれない

		usecase.Connect(testConnectMessage)
	})

	t.Run("offerが他の人だった場合、offerを送信", func(t *testing.T) {
		msgToSend := entity.Message{
			ID:       testConnectMessage.ID,
			Type:     "offer",
			SDP:      "otherSDP",
			TargetID: "otherID",
		}
		mockWm.EXPECT().
			GetConnectionByID(testConnectMessage.ID).
			Return(mockConn, nil).
			Times(1)

		mockO.EXPECT().
			IsOffer().
			Return(true).
			Times(1)

		mockO.EXPECT().
			IsOfferID(testConnectMessage.ID).
			Return(false).
			Times(1)

		mockO.EXPECT().
			GetOffer().
			Return("otherID").
			Times(1)

		mockRepo.EXPECT().
			GetSDPByID("otherID").
			Return("otherSDP", nil).
			Times(1)

		mockO.EXPECT().
			GetOffer().
			Return("otherID").
			Times(1)

		mockConn.EXPECT().
			WriteMessage(msgToSend).
			Times(1)

		// テスト実行
		usecase.Connect(testConnectMessage)
	})
}

func TestOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	testOfferMessage := entity.Message{
		ID:        "testID",
		Type:      "offer",
		SDP:       "testSDP",
		Candidate: nil,
		TargetID:  "",
	}

	t.Run("SDPが正常に保存されること", func(t *testing.T) {
		mockRepo.EXPECT().
			SaveSDP(testOfferMessage.ID, testOfferMessage.SDP).
			Times(1)

		usecase.Offer(testOfferMessage)
	})
}

func TestSendAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	mockConn := mock_service.NewMockWebSocketConnection(ctrl)
	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)
	testMessage := entity.Message{
		ID:        "senderID",
		Type:      "answer",
		SDP:       "testSDP",
		Candidate: nil,
		TargetID:  "receiverID",
	}

	t.Run("正常系", func(t *testing.T) {
		msgToSend := entity.Message{
			ID:       testMessage.ID,
			Type:     "answer",
			SDP:      testMessage.SDP,
			TargetID: testMessage.TargetID,
		}

		mockWm.EXPECT().
			GetConnectionByID(testMessage.TargetID).
			Return(mockConn, nil).
			Times(1)

		mockConn.EXPECT().
			WriteMessage(msgToSend).
			Times(1)

		usecase.Answer(testMessage)
	})

	t.Run("クライアントが見つからない場合", func(t *testing.T) {
		mockWm.EXPECT().
			GetConnectionByID(testMessage.TargetID).
			Return(nil, assert.AnError).
			Times(1)

		usecase.Answer(testMessage)
	})
}

func TestCandidateAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	t.Run("直接通信＋Candidateを返答", func(t *testing.T) {
		testMessage := entity.Message{
			ID:       "senderID",
			Type:     "candidate",
			SDP:      "",
			Candidate: []string{"testCandidate1", "testCandidate2"},
			TargetID: "receiverID",
		}

		msgToSend := entity.Message{
			ID:       testMessage.ID,
			Type:     "candidate",
			Candidate: testMessage.Candidate,
		}

		mockConn := mock_service.NewMockWebSocketConnection(ctrl)

		// GetConnectionByIDが正常に動作し、クライアントが取得できる
		mockWm.EXPECT().
			GetConnectionByID(testMessage.TargetID).
			Return(mockConn, nil).
			Times(1)

		// WriteMessageが呼ばれる
		mockConn.EXPECT().
			WriteMessage(msgToSend).
			Times(1)

		mockRepo.EXPECT().
			ExistsCandidateByID(testMessage.ID).
			Return(true).
			Times(1)

		mockRepo.EXPECT().
			AddCandidate(testMessage.ID, testMessage.Candidate).
			Return(nil).
			Times(1)

		mockO.EXPECT().
			IsOfferID(testMessage.TargetID).
			Return(true).
			Times(1)

		// CandidateAddを呼び出し
		result := usecase.CandidateAdd(testMessage)

		// 結果としてtrueを期待
		if !result {
			t.Errorf("Expected result to be true, but got false")
		}
	})

	t.Run("候補者が保存されていない場合（保存処理）", func(t *testing.T) {
		testMessage := entity.Message{
			ID:       "senderID",
			Type:     "candidate",
			SDP:      "",
			Candidate: []string{"testCandidate1", "testCandidate2"},
			TargetID: "",
		}

		// SaveCandidateが呼ばれる
		mockRepo.EXPECT().
			ExistsCandidateByID(testMessage.ID).
			Return(false).
			Times(1)

		mockRepo.EXPECT().
			SaveCandidate(testMessage.ID, testMessage.Candidate).
			Return(nil).
			Times(1)

		mockO.EXPECT().
			IsOfferID(testMessage.TargetID).
			Return(false).
			Times(1)

		// CandidateAddを呼び出し
		result := usecase.CandidateAdd(testMessage)

		// 結果としてfalseを期待（targetIDが空なので送信はしない）
		if result {
			t.Errorf("Expected result to be false, but got true")
		}
	})

	t.Run("候補者の保存でエラーが発生した場合", func(t *testing.T) {
		testMessage := entity.Message{
			ID:       "senderID",
			Type:     "candidate",
			SDP:      "",
			Candidate: []string{"testCandidate1", "testCandidate2"},
			TargetID: "",
		}

		// SaveCandidateでエラーが発生するシナリオ
		mockRepo.EXPECT().
			ExistsCandidateByID(testMessage.ID).
			Return(false).
			Times(1)

		mockRepo.EXPECT().
			SaveCandidate(testMessage.ID, testMessage.Candidate).
			Return(fmt.Errorf("save error")).
			Times(1)

		// CandidateAddを呼び出し、falseが返る
		result := usecase.CandidateAdd(testMessage)

		if result {
			t.Errorf("Expected result to be false, but got true")
		}
	})
}

func TestSendCandidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockIWebsocketRepository(ctrl)
	mockWm := mock_service.NewMockWebsocketManager(ctrl)
	mockBr := mock_service.NewMockWebSocketBroadcastService(ctrl)
	mockO := mock_service.NewMockOfferService(ctrl)

	usecase := usecase.NewWebsocketUsecase(mockRepo, mockWm, mockBr, mockO)

	testMessage := entity.Message{
		ID: 	 "testID",
		Type:   "candidate",
		SDP:    "",
		Candidate: []string{"testCandidate1"},
		TargetID: "offerID",
	}

	offerID := "offerID"
	mockConn := mock_service.NewMockWebSocketConnection(ctrl)

	t.Run("正常系", func(t *testing.T) {
		msgToSend := entity.Message{
			Type:	 "candidate",
			Candidate: []string{"testCandidate2"},
		}

		mockO.EXPECT().
			GetOffer().
			Return(offerID).
			Times(1)
		
		mockRepo.EXPECT().
			ExistsCandidateByID(offerID).
			Return(true).
			Times(1)

		mockWm.EXPECT().
			GetConnectionByID(testMessage.ID).
			Return(mockConn, nil).
			Times(1)
		
		mockRepo.EXPECT().
			GetCandidatesByID(offerID).
			Return(msgToSend.Candidate, nil).
			Times(1)

		mockConn.EXPECT().
			WriteMessage(msgToSend).
			Times(1)

		usecase.SendCandidate(testMessage)
	})

	t.Run("Candidateが存在しない場合", func(t *testing.T) {
		mockO.EXPECT().
			GetOffer().
			Return(offerID).
			Times(1)

		mockRepo.EXPECT().
			ExistsCandidateByID(offerID).
			Return(false).
			Times(1)

		usecase.SendCandidate(testMessage)
	})

	t.Run("Connection取得失敗", func(t *testing.T) {
		mockO.EXPECT().
			GetOffer().
			Return(offerID).
			Times(1)

		mockRepo.EXPECT().
			ExistsCandidateByID(offerID).
			Return(true).
			Times(1)

		mockWm.EXPECT().
			GetConnectionByID(testMessage.ID).
			Return(nil, fmt.Errorf("connection not found")).
			Times(1)

		usecase.SendCandidate(testMessage)
	})

	t.Run("Candidate取得失敗", func(t *testing.T) {
		mockO.EXPECT().
			GetOffer().
			Return(offerID).
			Times(1)

		mockRepo.EXPECT().
			ExistsCandidateByID(offerID).
			Return(true).
			Times(1)

		mockWm.EXPECT().
			GetConnectionByID(testMessage.ID).
			Return(mockConn, nil).
			Times(1)

		mockRepo.EXPECT().
			GetCandidatesByID(offerID).
			Return(nil, fmt.Errorf("failed to get candidate")).
			Times(1)

		usecase.SendCandidate(testMessage)
	})
}

