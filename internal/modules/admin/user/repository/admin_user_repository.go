package adminuserrepo

import (
	"context"
	"fmt"
	"time"

	adminuserdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/dto"
	adminusertypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/types"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	pg "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUserRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *AdminUserRepository {
	return &AdminUserRepository{pool}
}

func (r *AdminUserRepository) GetAll(ctx context.Context, pagination *pg.KeysetPagination[uuid.UUID], filters *adminusertypes.AdminUsersFilters) (*pg.KeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID], error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			email,
			first_name,
			last_name,
			role,
			is_banned,
			is_verified,
			avatar_url,
			google_id,
			github_id
	  FROM users
		WHERE
			  ($1::text IS NULL OR email ILIKE $1 || '%')
				AND ($2::boolean IS NULL OR is_verified = $2)
				AND ($3::boolean IS NULL OR is_banned = $3)
				AND (
            $4::timestamptz IS NULL
            OR (
                CASE
                    WHEN $6::boolean THEN
                        (created_at, id) > ($4, $5::uuid)
                    ELSE
                        (created_at, id) < ($4, $5::uuid)
                END
            )
        )
		 ORDER BY
        CASE WHEN $6 THEN created_at END ASC,
        CASE WHEN $6 THEN id END ASC,
        CASE WHEN NOT $6 THEN created_at END DESC,
        CASE WHEN NOT $6 THEN id END DESC
		LIMIT $7
	`

	cursor := pagination.Cursor()

	var cursorCreatedAt *time.Time
	var cursorId *uuid.UUID

	if cursor != nil {
		cursorCreatedAt = &cursor.CreatedAt
		cursorId = &cursor.Id
	}

	isBackward := pagination.Direction() == pg.Backward

	limit := pagination.Limit + 1

	rows, err := r.pool.Query(ctx, query, filters.Search, filters.IsVerified, filters.IsBanned, cursorCreatedAt,
		cursorId,
		isBackward,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adminUserDtos := make([]*adminuserdto.AdminUserDto, 0, 10)

	for rows.Next() {
		adminUserDto := &adminuserdto.AdminUserDto{}

		if err := rows.Scan(
			&adminUserDto.Id,
			&adminUserDto.CreatedAt,
			&adminUserDto.UpdatedAt,
			&adminUserDto.Email,
			&adminUserDto.FirstName,
			&adminUserDto.LastName,
			&adminUserDto.Role,
			&adminUserDto.IsBanned,
			&adminUserDto.IsVerified,
			&adminUserDto.AvatarUrl,
			&adminUserDto.GoogleId,
			&adminUserDto.GithubId,
		); err != nil {
			return nil, domain.InternalError(
				fmt.Sprintf("Failed to get users: %v", err),
			)
		}

		adminUserDtos = append(adminUserDtos, adminUserDto)
	}

	if err := rows.Err(); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("Failed to scan user: %v", err),
		)
	}

	return pg.NewKeysetPaginated[*adminuserdto.AdminUserDto, uuid.UUID](
		adminUserDtos,
		pagination,
		func(dto *adminuserdto.AdminUserDto) uuid.UUID { return dto.Id },
		func(dto *adminuserdto.AdminUserDto) time.Time { return dto.CreatedAt }), nil
}
