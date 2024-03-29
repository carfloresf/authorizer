package storage

import (
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"authorizer/internal/app/model"
)

// InMemory is my way to simulate a Database,
// this version of the database has a table account and a table transaction
// The PK of Account is Id even though it is not needed for this example (we are using always ID 1)
// Id is also the FK in Transaction to relate the transactions to the Account

type InMemory struct {
	History map[int][]Transaction
	Account map[int]Account
}

// Account in this package represents the table of Accounts in the simulated DB
type Account struct {
	Id             int
	ActiveCard     bool
	AvailableLimit int
}

// Transaction in this package represents the table of Transactions in the simulated DB
type Transaction struct {
	Id       uuid.UUID
	Merchant string
	Amount   int
	Time     time.Time
}

// GenerateAccountID is the function to get the sequential ID for the accounts,
// for this example we always set this value to 1
func (im *InMemory) GenerateAccountID() int {
	return 1
}

// CreateAccount is the function needed to create an account,
// it creates the "initial" transaction on the Transaction Map and adds the new record to Account map
func (im *InMemory) CreateAccount(a model.Account) error {
	log.Debugf("creation account: %+v", a)

	t := Transaction{
		Id:       uuid.New(),
		Merchant: "initial",
		Amount:   a.AvailableLimit,
		Time:     time.Now(),
	}

	transactions := []Transaction{
		t,
	}

	account := Account{
		Id:             a.Id,
		ActiveCard:     a.ActiveCard,
		AvailableLimit: a.AvailableLimit,
	}

	im.History = make(map[int][]Transaction)

	im.History[a.Id] = transactions

	im.Account = make(map[int]Account)

	im.Account[a.Id] = account

	return nil
}

// ExecuteTransaction is the operation in storage that updates the availableLimit
// and registers a new transaction in the transactionHistory
func (im *InMemory) ExecuteTransaction(a model.Account, t model.Transaction) (model.Account, error) {
	transaction := Transaction{
		Id:       uuid.New(),
		Merchant: t.Merchant,
		Amount:   t.Amount,
		Time:     t.Time,
	}

	a.AvailableLimit -= t.Amount

	account := Account{
		Id:             a.Id,
		ActiveCard:     a.ActiveCard,
		AvailableLimit: a.AvailableLimit,
	}

	im.Account[a.Id] = account

	im.History[a.Id] = append(im.History[a.Id], transaction)

	return a, nil
}

// GetAccount gets the info of the account using the account ID
func (im *InMemory) GetAccount(accountID int) model.Account {
	account := model.Account{
		Id:             accountID,
		ActiveCard:     im.Account[accountID].ActiveCard,
		AvailableLimit: im.Account[accountID].AvailableLimit,
	}

	return account
}

// GetTransactions gets all the transactions related to an account ID
func (im *InMemory) GetTransactions(accountID int) []model.Transaction {
	response := []model.Transaction{}

	for _, v := range im.History[accountID] {
		tx := model.Transaction{
			Merchant: v.Merchant,
			Amount:   v.Amount,
			Time:     v.Time,
		}

		response = append(response, tx)
	}

	return response
}

// Close closes connection to DB (not really needed for this abstraction of a DB)
func (im *InMemory) Close() error {
	return nil
}
