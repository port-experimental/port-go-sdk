package auth

import (
	"context"
	"fmt"
	"net/url"
)

// Doer matches client.Client for dependency injection.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Service exposes authentication-related routes.
type Service struct {
	doer Doer
}

// NewService constructs an auth Service.
func NewService(doer Doer) *Service {
	return &Service{doer: doer}
}

// AccessTokenRequest represents the payload sent to /v1/auth/access_token.
type AccessTokenRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// AccessTokenResponse is returned from /v1/auth/access_token.
type AccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

// RequestAccessToken exchanges client credentials for a Port bearer token.
func (s *Service) RequestAccessToken(ctx context.Context, req AccessTokenRequest) (AccessTokenResponse, error) {
	var resp AccessTokenResponse
	if err := s.doer.Do(ctx, "POST", "/v1/auth/access_token", req, &resp); err != nil {
		return AccessTokenResponse{}, err
	}
	return resp, nil
}

// RotateCredentials rotates credentials for the provided user email.
func (s *Service) RotateCredentials(ctx context.Context, userEmail string) error {
	if userEmail == "" {
		return fmt.Errorf("user email required for credential rotation")
	}
	path := fmt.Sprintf("/v1/rotate-credentials/%s", url.PathEscape(userEmail))
	return s.doer.Do(ctx, "POST", path, nil, nil)
}
