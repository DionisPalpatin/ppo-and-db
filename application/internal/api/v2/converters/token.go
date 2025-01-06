package converters

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

func ToTransportToken(token *models.Token) transport_models.Token {
	return transport_models.Token{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresAt:   token.ExpiresAt,
	}
}
