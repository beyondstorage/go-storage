package onedrive

// Onedrive only use OAuth2.0 authentication.

// ref:
// https://docs.microsoft.com/en-us/graph/auth/auth-concepts
// https://docs.microsoft.com/en-us/advertising/guides/authentication-oauth-get-tokens?view=bingads-13

import (
	"context"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

// getClient generate new onedrive client with oauth2 token.
func getClient(ctx context.Context, token string) *onedrive.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := onedrive.NewClient(tc)

	return client
}
