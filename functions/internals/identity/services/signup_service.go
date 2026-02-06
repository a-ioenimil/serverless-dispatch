package services

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type SignupService struct {
	allowedDomains []string
}

func NewSignupService() *SignupService {
	// Initialize configuration once
	env := os.Getenv("ALLOWED_EMAIL_DOMAINS")
	var domains []string
	if env != "" {
		raw := strings.Split(env, ",")
		for _, d := range raw {
			domains = append(domains, strings.TrimSpace(strings.ToLower(d)))
		}
	}
	return &SignupService{
		allowedDomains: domains,
	}
}

// ValidateSignup enforces domain policies for new users
func (s *SignupService) ValidateSignup(ctx context.Context, email string) error {
	if len(s.allowedDomains) == 0 {
		return fmt.Errorf("signup configuration error: no allowed domains")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	domain := strings.ToLower(parts[1])

	for _, d := range s.allowedDomains {
		if domain == d {
			return nil
		}
	}

	return fmt.Errorf("sign-up restricted to authorized domains only")
}
