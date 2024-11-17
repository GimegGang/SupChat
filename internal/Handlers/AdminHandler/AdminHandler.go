package AdminHandler

import (
	"SupChat/internal/Storage"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type DB interface {
	CloseTicket(id int) error
	CreateMessage(message, from string, ticketsId int) error
	GetTickets() ([]int, error)
	GetMessages(ticketId int) ([]Storage.Message, error)
}

type request struct {
	Message string `json:"message"`
}
type response struct {
	TicketId int
	Messages []Storage.Message
}

func MainHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tickets, err := db.GetTickets()
		if err != nil {
			log.Error("Error getting tickets", "error", err)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/MainAdmin.html"))
		if err := tmpl.ExecuteTemplate(w, "MainAdmin.html", tickets); err != nil {
			log.Error("GetTodoHandler", "err", err)
			return
		}
	}
}

func ChatHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := strconv.Atoi(chi.URLParam(r, "chatId"))
		if err != nil {
			log.Error("Error getting chatId", "error", err)
			return
		}
		messages, err := db.GetMessages(chatId)

		data := response{chatId, messages}

		if err != nil {
			log.Error("Error getting messages", "error", err)
		}
		tmpl := template.Must(template.ParseFiles("templates/ChatAdmin.html"))
		if err = tmpl.ExecuteTemplate(w, "ChatAdmin.html", data); err != nil {
			log.Error("GetTodoHandler", "err", err)
			return
		}
	}
}

func SendAdminMessageHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error decode request", "err", err)
			return
		}
		chatId, err := strconv.Atoi(chi.URLParam(r, "chatId"))
		if err != nil {
			log.Error("Error getting chatId", "error", err)
			return
		}
		if err = db.CreateMessage(req.Message, "admin", chatId); err != nil {
			log.Error("Error creating message", "err", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
		})
	}
}

func CloseTicketHandler(log *slog.Logger, db DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ticketId, err := strconv.Atoi(chi.URLParam(r, "chatId"))
		if err != nil {
			log.Error("Error getting ticketId", "error", err)
		}
		if err = db.CloseTicket(ticketId); err != nil {
			log.Error("Error closing ticket", "err", err)
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
