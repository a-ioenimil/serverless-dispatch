package services

import (
	"errors"
	"strings"
)

var (
	ErrUnauthorizedDomain = errors.New("sign-up restricted to authorized domains only")
)

type AuthPolicyService struct {
	allowedDomains []string
}

func NewAuthPolicyService(allowedDomainsStr string) *AuthPolicyService {
	// Parse config once during initialization
	var domains []string
	if allowedDomainsStr != "" {
		parts := strings.Split(allowedDomainsStr, ",")
		for _, d := range parts {
			domains = append(domains, strings.TrimSpace(d))
		}
	}

	return &AuthPolicyService{
		allowedDomains: domains,
	}
}

// ValidateSignup enforces organization-specific access policies
func (s *AuthPolicyService) ValidateSignup(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	domain := strings.ToLower(parts[1])

	if len(s.allowedDomains) == 0 {
		return errors.New("no domains allowed by configuration")
	}

	for _, d := range s.allowedDomains {
		if domain == d {
			return nil
		}
	}

	return ErrUnauthorizedDomain
}
