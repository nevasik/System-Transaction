package httpserver

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"golangProject/internal/entity"
	"golangProject/internal/storage/postgresql"
	"log"
	"net/http"
	"time"
)

const (
	QueueName = "Transactional"
)

type HTTPServer struct {
	db       *postgresql.Database
	NatsConn *nats.Conn
	router   *mux.Router
}

func NewHTTPServer(db *postgresql.Database, natsConn *nats.Conn) *HTTPServer {
	return &HTTPServer{
		db:       db,
		NatsConn: natsConn,
		router:   mux.NewRouter(),
	}
}
func (s *HTTPServer) Start(port string) {
	s.routes()
	log.Fatal(http.ListenAndServe(":"+port, s.router))
}

func (s *HTTPServer) routes() {
	s.router.HandleFunc("/invoice", s.handleInvoice).Methods("POST")
	s.router.HandleFunc("/withdraw", s.handleWithdraw).Methods("POST")
	s.router.HandleFunc("/balance", s.handleBalance).Methods("GET")
}

func (s *HTTPServer) handleInvoice(writer http.ResponseWriter, request *http.Request) {
	var transaction entity.Transaction

	if err := json.NewDecoder(request.Body).Decode(&transaction); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.Status = "Created"
	createdAt := time.Now()
	transaction.CreatedAt = &createdAt

	if err := s.depositFunds(&transaction); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(transaction)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *HTTPServer) handleWithdraw(writer http.ResponseWriter, request *http.Request) {
	var transaction entity.Transaction

	if err := json.NewDecoder(request.Body).Decode(&transaction); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.withdrawFunds(&transaction); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(transaction)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *HTTPServer) handleBalance(writer http.ResponseWriter, request *http.Request) {
	walletOrCardNumber := request.URL.Query().Get("wallet_or_card_number")
	currency := request.URL.Query().Get("currency")

	balance, err := s.db.GetBalance(walletOrCardNumber)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{"balance": balance}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)

	createdAt := time.Now()

	transaction := entity.Transaction{
		Currency:     currency,
		Amount:       balance,
		WalletOrCard: walletOrCardNumber,
		Status:       "Success",
		CreatedAt:    &createdAt,
	}

	if err := s.publishTransactionalMessage(&transaction); err != nil {
		log.Println("Ошибка при отправке сообщения в NATS:", err)
	}
}

func (s *HTTPServer) depositFunds(transaction *entity.Transaction) error {
	_, err := s.db.CreateTransaction(transaction.Currency, transaction.Amount, transaction.WalletOrCard)
	if err != nil {
		return err
	}

	if err := s.publishTransactionalMessage(transaction); err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) publishTransactionalMessage(transaction *entity.Transaction) error {
	message, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	return s.NatsConn.Publish(QueueName, message)
}

func (s *HTTPServer) withdrawFunds(transaction *entity.Transaction) error {
	if err := s.db.WithdrawFunds(transaction.WalletOrCard, transaction.Amount); err != nil {
		return err
	}

	return nil
}
