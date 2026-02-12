package cognito

import (
	"context"
	"fmt"
	"strings"

	identityprovider "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type UserDirectory struct {
	client     *identityprovider.Client
	userPoolID string
}

func NewUserDirectory(client *identityprovider.Client, userPoolID string) *UserDirectory {
	return &UserDirectory{client: client, userPoolID: userPoolID}
}

func (d *UserDirectory) ResolveUserEmailByUsername(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", nil
	}

	out, err := d.client.AdminGetUser(ctx, &identityprovider.AdminGetUserInput{
		UserPoolId: &d.userPoolID,
		Username:   &username,
	})
	if err != nil {
		return "", fmt.Errorf("admin get user failed: %w", err)
	}

	return attributeValue(out.UserAttributes, "email"), nil
}

func (d *UserDirectory) ListAdminEmails(ctx context.Context) ([]string, error) {
	var (
		emails []string
		token  *string
	)

	for {
		out, err := d.client.ListUsersInGroup(ctx, &identityprovider.ListUsersInGroupInput{
			UserPoolId: &d.userPoolID,
			GroupName:  strPtr("Admins"),
			NextToken:  token,
		})
		if err != nil {
			return nil, fmt.Errorf("list admins failed: %w", err)
		}

		for _, user := range out.Users {
			email := attributeValue(user.Attributes, "email")
			if email != "" {
				emails = append(emails, email)
			}
		}

		if out.NextToken == nil || strings.TrimSpace(*out.NextToken) == "" {
			break
		}
		token = out.NextToken
	}

	return unique(emails), nil
}

func attributeValue(attrs []types.AttributeType, name string) string {
	for _, attr := range attrs {
		if attr.Name != nil && attr.Value != nil && *attr.Name == name {
			return strings.TrimSpace(*attr.Value)
		}
	}
	return ""
}

func unique(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func strPtr(v string) *string {
	return &v
}
