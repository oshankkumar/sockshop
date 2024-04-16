package app

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/internal/db"
	"github.com/oshankkumar/sockshop/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	UserStore    domain.UserStore
	CardStore    domain.CardStore
	AddressStore domain.AddressStore
	TxBeginner   db.TxBeginner
	Domain       string
}

func (u *UserService) Login(ctx context.Context, username, password string) (*api.User, error) {
	user, err := u.UserStore.GetUserByName(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("UserService.Login(username=%s): %w", username, err)
	}

	passHash := calculatePassHash(password, user.Salt)

	if user.Password != passHash {
		return nil, fmt.Errorf("UserService.Login(password=%s): %w", password, api.ErrUnauthorized)
	}

	usr := &api.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  username,
		Email:     user.Email,
		ID:        user.ID,
		Links:     api.NewCustomerLinks(u.Domain, user.ID.String()),
	}

	return usr, nil
}

func (u *UserService) Register(ctx context.Context, user api.User) (uuid.UUID, error) {
	userM := &domain.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
		Salt:      salt(),
	}

	userM.Password = calculatePassHash(user.Password, userM.Salt)

	if err := u.UserStore.CreateUser(ctx, userM); err != nil {
		return uuid.UUID{}, fmt.Errorf("UserService.Register(username=%s): %w", user.Username, err)
	}

	return userM.ID, nil
}

func (u *UserService) GetUser(ctx context.Context, id string) (*api.User, error) {
	user, err := u.UserStore.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUser(id=%s): %w", id, err)
	}

	usr := &api.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		ID:        user.ID,
		Links:     api.NewCustomerLinks(u.Domain, user.ID.String()),
	}

	return usr, nil
}

func (u *UserService) GetUsers(ctx context.Context) ([]api.User, error) {
	users, err := u.UserStore.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUser: %w", err)
	}

	var usrs []api.User

	for _, user := range users {
		usrs = append(usrs, api.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			ID:        user.ID,
			Links:     api.NewCustomerLinks(u.Domain, user.ID.String()),
		})
	}

	return usrs, nil
}

func (u *UserService) CreateAddress(ctx context.Context, addr api.Address, userID string) (uuid.UUID, error) {
	addrM := &domain.Address{
		Street:   addr.Street,
		Number:   addr.Number,
		Country:  addr.Country,
		City:     addr.City,
		PostCode: addr.PostCode,
	}

	err := db.RunInTransaction(ctx, u.TxBeginner, func(ctx context.Context, tx *sqlx.Tx) error {
		addrStore, userStore := u.AddressStore.WithTx(tx), u.UserStore.WithTx(tx)

		if err := addrStore.CreateAddress(ctx, addrM); err != nil {
			return err
		}

		if err := userStore.CreateAddress(ctx, addrM.ID.String(), userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: UserService.CreateAddress(userID=%s)", err, userID)
	}

	return addrM.ID, nil
}

func (u *UserService) CreateCard(ctx context.Context, card api.Card, userID string) (uuid.UUID, error) {
	cardM := &domain.Card{
		LongNum: card.LongNum,
		Expires: card.Expires,
		CCV:     card.CCV,
	}

	err := db.RunInTransaction(ctx, u.TxBeginner, func(ctx context.Context, tx *sqlx.Tx) error {
		cardStore, userStore := u.CardStore.WithTx(tx), u.UserStore.WithTx(tx)

		if err := cardStore.CreateCard(ctx, cardM); err != nil {
			return err
		}

		if err := userStore.CreateCard(ctx, cardM.ID.String(), userID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("UserService.CreateCard(userID=%s): %w", userID, err)
	}

	return cardM.ID, nil
}

func (u *UserService) GetAddresses(ctx context.Context, id string) (*api.Address, error) {
	addrM, err := u.UserStore.GetAddress(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetAddresses(id=%s): %w", id, err)
	}

	addr := &api.Address{
		ID:       addrM.ID,
		Street:   addrM.Street,
		Number:   addrM.Number,
		Country:  addrM.Country,
		City:     addrM.City,
		PostCode: addrM.PostCode,
		Links:    api.NewAddressLinks(u.Domain, addrM.ID.String()),
	}

	return addr, nil
}

func (u *UserService) GetCard(ctx context.Context, id string) (*api.Card, error) {
	cardM, err := u.UserStore.GetCard(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetCard(id=%s): %w", id, err)
	}

	card := &api.Card{
		ID:      cardM.ID,
		LongNum: cardM.LongNum,
		Expires: cardM.Expires,
		CCV:     cardM.CCV,
		Links:   api.NewCardLinks(u.Domain, cardM.ID.String()),
	}
	card.MaskCC()

	return card, nil
}

func (u *UserService) GetUserCards(ctx context.Context, userID string) ([]api.Card, error) {
	cardsM, err := u.UserStore.GetUserCards(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserCards(userID=%s): %w", userID, err)
	}

	var cards []api.Card
	for _, c := range cardsM {
		card := api.Card{
			ID:      c.ID,
			LongNum: c.LongNum,
			Expires: c.Expires,
			CCV:     c.CCV,
			Links:   api.NewCardLinks(u.Domain, c.ID.String()),
		}
		card.MaskCC()
		cards = append(cards, card)
	}

	return cards, nil
}

func (u *UserService) GetUserAddresses(ctx context.Context, userID string) ([]api.Address, error) {
	addrsM, err := u.UserStore.GetUserAddresses(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserAddresses(userID=%s): %w", userID, err)
	}

	var addresses []api.Address
	for _, adr := range addrsM {
		addresses = append(addresses, api.Address{
			ID:       adr.ID,
			Street:   adr.Street,
			Number:   adr.Number,
			Country:  adr.Country,
			City:     adr.City,
			PostCode: adr.PostCode,
			Links:    api.NewAddressLinks(u.Domain, adr.ID.String()),
		})
	}

	return addresses, nil
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, salt)
	_, _ = io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func salt() string {
	h := sha1.New()
	_, _ = io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	return fmt.Sprintf("%x", h.Sum(nil))
}
