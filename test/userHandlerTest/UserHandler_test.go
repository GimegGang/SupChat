package userHandlerTest

import (
	"SupChat/internal/Handlers/UserHandler"
	"SupChat/internal/Storage"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct{}

func (m *mockDB) OpenTicket() (int, error) {
	return 1, nil
}

func (m *mockDB) CreateMessage(message, from string, ticketsId int) error {
	return nil // Здесь можно добавить логику для имитации ошибки
}

func (m *mockDB) GetMessages(ticketId int) ([]Storage.Message, error) {
	return nil, nil
}

func TestSendMessageHandler(t *testing.T) {
	db := &mockDB{}
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	reqBody := []byte(`{"message": "Test message", "tId": "new"}`)
	req, err := http.NewRequest("POST", "/sendUserMessage", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := UserHandler.SendMessageHandler(log, db)
	handler.ServeHTTP(rr, req)

	// Проверка статуса ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверка тела ответа
	expected := `{"status":"success","ticketId":1}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
