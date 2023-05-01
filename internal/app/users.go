package app

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/oshankkumar/sockshop/api"
	"github.com/oshankkumar/sockshop/domain"
)

type UserService struct {
	userStore domain.UserStore
	domain    string
}

func (u *UserService) Login(ctx context.Context, username, password string) (*api.User, error) {
	user, err := u.userStore.GetUserByName(ctx, username)
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
		Links:     make(api.Links),
	}

	usr.Links.AddCustomer(u.domain, usr.ID.String())

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

	if err := u.userStore.CreateUser(ctx, userM); err != nil {
		return uuid.UUID{}, fmt.Errorf("UserService.Register(username=%s): %w", user.Username, err)
	}

	return userM.ID, nil
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func salt() string {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	return fmt.Sprintf("%x", h.Sum(nil))
}
