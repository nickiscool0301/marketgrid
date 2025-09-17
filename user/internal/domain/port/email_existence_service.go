package port

import "context"

type EmailExistenceService interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	AddEmail(ctx context.Context, email string) error
	InitializeFromEmails(ctx context.Context, emails []string) error
}