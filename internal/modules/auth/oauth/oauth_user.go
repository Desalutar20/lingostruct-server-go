package oauth

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
)

type OAuthUser struct {
	ProviderId userdomain.ProviderId
	Email      valueobjects.Email
}
