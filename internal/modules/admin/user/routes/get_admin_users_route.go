package adminuserroutes

import (
	"log/slog"
	"net/http"

	adminuserconst "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/const"
	adminusertypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/types"
	adminuserusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/ggicci/httpin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type GetAdminUsersRequest struct {
	Search     *string `in:"query=search" json:"search"`
	IsBanned   *bool   `in:"query=isBanned" json:"isBanned"`
	IsVerified *bool   `in:"query=isVerified" json:"isVerified"`
	routes.KeysetPaginationRequest
}

func (r GetAdminUsersRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Search, validation.By(routes.Trimmed), validation.Length(1, adminuserconst.GetAdminUsersSearchMaxLimit)),
		validation.Field(&r.Limit, validation.Max(adminuserconst.GetAdminUsersMaxLimit)),
		validation.Field(&r.PrevCursor, validation.By(routes.Trimmed), validation.Length(1, pagination.CursorMaxLength)),
		validation.Field(&r.NextCursor, validation.By(routes.Trimmed), validation.Length(1, pagination.CursorMaxLength)),
		validation.Field(&r.KeysetPaginationRequest, validation.By(routes.ValidateKeysetPagination)),
	)
}

func (r GetAdminUsersRequest) ToQuery() (*adminuserusecases.GetAdminUsersQuery, []error) {
	err := r.Validate()
	if err != nil {
		flattenErrors := validation.Flatten(err)
		errStrings := make([]error, 0, len(flattenErrors))

		for _, fe := range flattenErrors {
			errStrings = append(errStrings, domain.ValidationError(fe.Field(), []string{fe.Message}))
		}

		return nil, errStrings
	}

	var search *valueobjects.NonEmptyString
	var prevCursor *pagination.KeysetCursor[uuid.UUID]
	var nextCursor *pagination.KeysetCursor[uuid.UUID]
	limit := adminuserconst.GetAdminUsersDefaultLimit

	if r.Search != nil {
		s, err := valueobjects.NewNonEmptyString(*r.Search, "search", "Search", 0, adminuserconst.GetAdminUsersSearchMaxLimit, nil)
		if err != nil {
			return nil, []error{err}
		}

		search = &s
	}

	if r.PrevCursor != nil {
		c, err := pagination.NewKeysetCursor(*r.PrevCursor, func(s string) (uuid.UUID, error) {
			return uuid.Parse(s)
		})
		if err != nil {
			return nil, []error{err}
		}

		prevCursor = &c
	}

	if r.NextCursor != nil {
		c, err := pagination.NewKeysetCursor(*r.NextCursor, func(s string) (uuid.UUID, error) {
			return uuid.Parse(s)
		})
		if err != nil {
			return nil, []error{err}
		}

		nextCursor = &c
	}

	if r.Limit != nil {
		limit = *r.Limit
	}

	return &adminuserusecases.GetAdminUsersQuery{
		Pagination: &pagination.KeysetPagination[uuid.UUID]{
			PrevCursor: prevCursor,
			NextCursor: nextCursor,
			Limit:      limit,
		},
		Filters: &adminusertypes.AdminUsersFilters{
			Search:     search,
			IsVerified: r.IsVerified,
			IsBanned:   r.IsBanned,
		},
	}, nil
}

func (router *Router) GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value(httpin.Input).(*GetAdminUsersRequest)

	query, errors := request.ToQuery()
	if routes.WriteErrors(w, errors, router.logger) {
		return
	}

	result, err := router.getAdminUsersQueryHandler.Handle(r.Context(), query)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	if err := routes.WriteJSON(w, http.StatusOK, routes.NewKeysetPaginatedResponseDto(result)); err != nil {
		router.logger.Error(
			"Failed to write get-admin-users response",
			slog.Any("error", err),
		)
	}
}
