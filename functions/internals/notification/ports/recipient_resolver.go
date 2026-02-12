package ports

import "context"

type RecipientResolver interface {
	ResolveUserEmailByUsername(ctx context.Context, username string) (string, error)
	ListAdminEmails(ctx context.Context) ([]string, error)
}
