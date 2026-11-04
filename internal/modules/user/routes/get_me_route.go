package userroutes

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/dto"
)

func (router *UserRouter) GetMe(w http.ResponseWriter, r *http.Request) {
	session, _, ok := routes.GetSession(r.Context())
	if !ok {
		routes.WriteError(w, domain.UnauthorizedError(), router.logger)
		return
	}

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[userdto.UserDto]{
		Data: session.ToUserDto(),
	}); err != nil {
		router.logger.Error(
			"Failed to write get-me response",
			slog.Any("error", err),
		)
	}

}
