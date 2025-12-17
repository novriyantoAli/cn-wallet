package oauth

import (
	"github.com/novriyantoAli/cn-wallet/internal/application/oauth/handler"
	"github.com/novriyantoAli/cn-wallet/internal/application/oauth/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides OAuth-related dependencies
var Module = fx.Options(
	fx.Provide(
		NewOAuthConfig,
		service.NewOAuthService,
		handler.NewOAuthHandler,
	),
)

// NewOAuthConfig creates a new OAuth configuration
// This would be replaced with actual configuration loading in production
func NewOAuthConfig(logger *zap.Logger) service.OAuthConfig {
	logger.Info("Initializing OAuth configuration")

	return service.OAuthConfig{
		Google: service.GoogleOAuthConfig{
			ClientID:     "1098309707454-6nep66gts4i0e55enhg1aobl5vp1kjqr.apps.googleusercontent.com",
			ClientSecret: "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",
			RedirectURI:  "http://localhost:8080/api/v1/oauth/callback/google",
			Scopes:       []string{"openid", "profile", "email"},
		},
		Github: service.GithubOAuthConfig{
			ClientID:     "your-github-client-id",
			ClientSecret: "your-github-client-secret",
			RedirectURI:  "http://localhost:8080/api/v1/oauth/callback/github",
			Scopes:       []string{"user:email"},
		},
		Gitlab: service.GitlabOAuthConfig{
			ClientID:     "your-gitlab-client-id",
			ClientSecret: "your-gitlab-client-secret",
			RedirectURI:  "http://localhost:8080/api/v1/oauth/callback/gitlab",
			Scopes:       []string{"openid", "profile", "email"},
		},
		Microsoft: service.MicrosoftOAuthConfig{
			ClientID:     "your-microsoft-client-id",
			ClientSecret: "your-microsoft-client-secret",
			RedirectURI:  "http://localhost:8080/api/v1/oauth/callback/microsoft",
			Scopes:       []string{"openid", "profile", "email"},
			TenantID:     "common",
		},
	}
}
