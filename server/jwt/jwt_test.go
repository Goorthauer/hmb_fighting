package jwt

import (
	"testing"

	"hmb_fighting/server/entities"
	"hmb_fighting/server/types"
)

func TestRefreshTokenPreservesClientID(t *testing.T) {
	pair, err := GenerateTokenPair(entities.User{ID: "client-1", Email: "user@example.com"}, types.UserRoleSpectator)
	if err != nil {
		t.Fatal(err)
	}

	refreshed, err := RefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := ValidateToken(refreshed.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if claims.ClientID != "client-1" {
		t.Fatalf("got ClientID %q, want %q", claims.ClientID, "client-1")
	}
}

func TestRefreshTokenRejectsAccessToken(t *testing.T) {
	pair, err := GenerateTokenPair(entities.User{ID: "client-1", Email: "user@example.com"}, types.UserRoleSpectator)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := RefreshToken(pair.AccessToken); err == nil {
		t.Fatal("expected access token to be rejected as refresh token")
	}
}
