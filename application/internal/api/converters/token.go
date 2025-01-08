package converters

import (
	"github.com/DionisPalpatin/ppo-and-db/application/internal/api/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

func ToTransportToken(token *models.Token) transport_models.Token {
	return transport_models.Token{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}
}
