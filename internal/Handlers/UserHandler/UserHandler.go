package UserHandler

import (
	"SupChat/internal/Storage"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type request struct {
	Message  string `json:"message"`
	TicketID string `json:"tId"`
}

type DB interface {
	OpenTicket() (int, error)
	CreateMessage(message, from string, ticketsId int) error
	GetMessages(ticketId int) ([]Storage.Message, error)
	CheckOpenness(ticketId int) (bool, error)
}

func MainHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/MainUser.html"))
		if err := tmpl.ExecuteTemplate(w, "MainUser.html", ""); err != nil {
			log.Error("GetTodoHandler", "err", err)
			return
		}
	}
}

func SendMessageHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("SendMessageHandler", "err", err)
			return
		}
		var tID int
		var err error
		if req.TicketID == "new" {
			tID, err = db.OpenTicket()
			if err != nil {
				log.Error("SendMessageHandler", "err", err)
				return
			}
		} else {
			tID, err = strconv.Atoi(req.TicketID)
			if err != nil {
				log.Error("SendMessageHandler", "err", err)
				return
			}
		}
		if err = db.CreateMessage(req.Message, "user", tID); err != nil {
			log.Error("SendMessageHandler", "err", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"ticketId": tID,
		})
	}
}

func GetMessageHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketId, err := strconv.Atoi(chi.URLParam(r, "ticketID"))
		if err != nil {
			log.Error("GetMessageHandler", "err", err)
			return
		}
		messages, err := db.GetMessages(ticketId)
		if err != nil {
			log.Error("GetMessageHandler", "err", err)
			return
		}
		closed, err := db.CheckOpenness(ticketId)
		if err != nil {
			log.Error("GetMessageHandler", "err", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "success",
			"messages": messages,
			"closed":   closed,
		})
	}
}
