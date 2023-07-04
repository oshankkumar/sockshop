package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/oshankkumar/sockshop/internal/db"
	"github.com/oshankkumar/sockshop/internal/domain"
)

const MySQLErrCodeDupe = 1062

func NewUserStore(db db.DBTx) *UserStore {
	return &UserStore{db: db}
}

type UserStore struct {
	db db.DBTx
}

func (u *UserStore) GetUserByName(ctx context.Context, uname string) (domain.User, error) {
	query := "SELECT customer.id, customer.first_name, customer.last_name, customer.email, customer.username, customer.password, customer.salt " +
		"FROM customer WHERE username=?;"

	var user domain.User
	if err := GetContext(ctx, u.db, &user, query, uname); err != nil {
		return user, fmt.Errorf("UserStore.GetUserByName(%s): %w", uname, err)
	}

	return user, u.addAttributes(ctx, &user)
}

func (u *UserStore) GetUser(ctx context.Context, id string) (domain.User, error) {
	query := "SELECT customer.id, customer.first_name, customer.last_name, customer.email, customer.username, customer.password, customer.salt " +
		"FROM customer WHERE customer.id=?;"

	var user domain.User
	if err := GetContext(ctx, u.db, &user, query, id); err != nil {
		return user, fmt.Errorf("UserStore.GetUser(%s): %w", id, err)
	}

	return user, u.addAttributes(ctx, &user)
}

func (u *UserStore) addAttributes(ctx context.Context, user *domain.User) error {
	query := "SELECT address_id FROM customer_address WHERE customer_id=?;"

	var addrs []string
	if err := SelectContext(ctx, u.db, &addrs, query, user.ID); err != nil {
		return fmt.Errorf("UserStore.addAttributes(%s): %w", user.ID, err)
	}

	query = "SELECT card_id FROM customer_card WHERE customer_id=?;"

	var cards []string
	if err := SelectContext(ctx, u.db, &cards, query, user.ID); err != nil {
		return fmt.Errorf("UserStore.addAttributes(%s): %w", user.ID, err)
	}

	user.AddressIDs = addrs
	user.CardIDs = cards

	return nil
}

func (u *UserStore) GetUsers(ctx context.Context) ([]domain.User, error) {
	query := "SELECT customer.id, customer.first_name, customer.last_name, customer.email, customer.username, customer.password, customer.salt " +
		"FROM customer;"

	var users []domain.User
	if err := SelectContext(ctx, u.db, &users, query); err != nil {
		return users, fmt.Errorf("UserStore.GetUser(): %w", err)
	}

	for _, user := range users {
		if err := u.addAttributes(ctx, &user); err != nil {
			return users, fmt.Errorf("UserStore.GetUser(): %w", err)
		}
	}

	return users, nil
}

func (u *UserStore) GetAddress(ctx context.Context, id string) (domain.Address, error) {
	query := "SELECT id, street, number, country, city, postcode FROM address WHERE id=?;"

	var addr domain.Address
	if err := GetContext(ctx, u.db, &addr, query, id); err != nil {
		return addr, fmt.Errorf("UserStore.GetAddress(%s): %w", id, err)
	}

	return addr, nil
}

func (u *UserStore) GetUserAddresses(ctx context.Context, userID string) ([]domain.Address, error) {
	query := "SELECT a.id, a.street, a.number, a.country, a.city, a.postcode " +
		"FROM customer_address ca JOIN address a ON ca.address_id=a.id " +
		"WHERE ca.customer_id=?;"

	var addrs []domain.Address
	if err := SelectContext(ctx, u.db, &addrs, query, userID); err != nil {
		return addrs, fmt.Errorf("UserStore.GetAddresses(): %w", err)
	}

	return addrs, nil
}

func (u *UserStore) GetCard(ctx context.Context, id string) (domain.Card, error) {
	query := "SELECT id, long_num, expires, ccv FROM card WHERE id=?;"

	var card domain.Card
	if err := GetContext(ctx, u.db, &card, query, id); err != nil {
		return card, fmt.Errorf("UserStore.GetCard(%s): %w", id, err)
	}

	return card, nil
}

func (u *UserStore) GetUserCards(ctx context.Context, userID string) ([]domain.Card, error) {
	query := "SELECT c.id, c.long_num, c.expires, c.ccv " +
		"FROM card c JOIN customer_card cc ON c.id=cc.card_id " +
		"WHERE cc.customer_id=?;"

	var cards []domain.Card
	if err := SelectContext(ctx, u.db, &cards, query, userID); err != nil {
		return cards, fmt.Errorf("UserStore.GetCards(): %w", err)
	}

	return cards, nil
}

func (u *UserStore) CreateUser(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO customer(id, first_name, last_name, email, username, password, salt) VALUES (?, ?, ?, ?, ?, ?, ?)"

	user.ID = uuid.New()

	_, err := u.db.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Username,
		user.Password,
		user.Salt,
	)

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == MySQLErrCodeDupe {
		return domain.DuplicateEntryError{Entity: "user", Err: err}
	}

	if err != nil {
		return fmt.Errorf("UserStore.CreateUser(%s): %w", user.Username, err)
	}

	return nil
}

func (u *UserStore) CreateAddress(ctx context.Context, addr *domain.Address, userID string) error {
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	addr.ID = uuid.New()
	query := "INSERT INTO address(id, street, number, country, city, postcode) VALUES (?, ?, ?, ?, ?, ?)"

	_, err = tx.ExecContext(ctx, query,
		addr.ID,
		addr.Street,
		addr.Number,
		addr.Country,
		addr.City,
		addr.PostCode,
	)
	if err != nil {
		return err
	}

	query = "INSERT INTO customer_address(customer_id, address_id) VALUES (?, ?)"

	if _, err = tx.ExecContext(ctx, query, userID, addr.ID); err != nil {
		return fmt.Errorf("UserStore.CreateAddress(userID=%s): %w", userID, err)
	}

	return tx.Commit()
}

func (u *UserStore) CreateCard(ctx context.Context, card *domain.Card, userID string) error {
	tx, err := u.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	card.ID = uuid.New()

	query := "INSERT INTO card(id, long_num, expires, ccv) VALUES (?, ?, ?, ?)"

	_, err = tx.ExecContext(ctx, query,
		card.ID,
		card.LongNum,
		card.Expires,
		card.CCV,
	)

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == MySQLErrCodeDupe {
		return domain.DuplicateEntryError{Entity: "card", Err: err}
	}

	if err != nil {
		return fmt.Errorf("UserStore.CreateCard(userID=%s): txn.ExecContext(card): %w", userID, err)
	}

	query = "INSERT INTO customer_card(customer_id, card_id) VALUES (?, ?)"

	if _, err = tx.ExecContext(ctx, query, userID, card.ID); err != nil {
		return fmt.Errorf("UserStore.CreateCard(userID=%s): txn.ExecContext(customer_card): %w", userID, err)
	}

	return tx.Commit()
}

func (u *UserStore) Delete(ctx context.Context, entity string, id string) error {
	panic("not implemented") // TODO: Implement
}
