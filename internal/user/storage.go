package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindOne(ctx context.Context, username string, id string) (User, error)
	Delete(ctx context.Context, user User) error
	Update(ctx context.Context, id string) error
}
