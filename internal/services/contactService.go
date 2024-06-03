package services

/*
 - uma coisa é a separação de rotas na hora de fazer
o mux.HandleFunc("verbo rota", função de um determinado domínio).
Essa parte vc pode fazer um método "Setup<resource>Routes" no app

- type UserRouter struct {
	getAllUsers	(func) (w http.ResponseWritter, req *http.Request)
}

func (rt *UserRouter) getAllUsers
 - outra coisa é a separação dos nomes das queries.
 Ce não tá usando um orm, as queries não são construídas,
 só são chamadas
*/

import (
	"context"
	"encoding/json"
	"github.com/Grsaiago/gonebook_server/internal/database"
	"log/slog"
	"net/http"
)

// dessa forma e usando o Repo da ContactService eu dou um 'narrow down' nos métodos quer
// essa model tem acesso dentro do universo de queries do database.Queries
type ContactRepository interface {
	GetAllContacts(ctx context.Context) ([]database.Contact, error)
	CreateContact(ctx context.Context, arg database.CreateContactParams) (database.Contact, error)
}

type ContactService struct {
	Repo       ContactRepository
	AppContext *context.Context
}

func New(methodScope ContactRepository, appCtx *context.Context) ContactService {
	return ContactService{
		Repo:       methodScope,
		AppContext: appCtx,
	}
}

func (ctts *ContactService) GetAllContacts(w http.ResponseWriter, req *http.Request) {
	slog.Info("Hit: Get")
	users, err := ctts.Repo.GetAllContacts(*ctts.AppContext)
	if err != nil {
		slog.Error("db: query failed", "error", err.Error())
		http.Error(w, "db: query failed", http.StatusBadRequest)
		return
	}
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		slog.Error("http: Json encode error", "error", err.Error())
		http.Error(w, "http: Json encode error", http.StatusInternalServerError)
		return
	}
}

func (ctts *ContactService) CreateContact(w http.ResponseWriter, req *http.Request) {
	slog.Info("Hit: Post")
	newContact :=
		database.CreateContactParams{}

	err := json.NewDecoder(req.Body).Decode(&newContact)
	if err != nil {
		slog.Error("json: invalid fields", "error", err.Error())
		http.Error(w, "json: invalid fields", http.StatusBadRequest)
		return
	}
	addedUser, err := ctts.Repo.CreateContact(*ctts.AppContext, newContact)
	if err != nil {
		slog.Error("db: query failed", "error", err.Error())
		http.Error(w, "db: query failed", http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(addedUser)
	if err != nil {
		slog.Error("json: encoding Failed but user added", "error", err.Error())
		http.Error(w, "json: encoding Failed but user added", http.StatusInternalServerError)
		return
	}
}
