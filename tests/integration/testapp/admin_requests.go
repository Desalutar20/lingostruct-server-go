package testapp

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	adminuserroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/routes"
)

func (a *TestApp) GetAdminUsers(
	t *testing.T,
	request *adminuserroutes.GetAdminUsersRequest,
	cookies []*http.Cookie,
) *http.Response {
	t.Helper()

	query := url.Values{}

	if request != nil {
		if request.Search != nil {
			query.Set("search", *request.Search)
		}
		if request.IsBanned != nil {
			query.Set("isBanned", strconv.FormatBool(*request.IsBanned))
		}
		if request.IsVerified != nil {
			query.Set("isVerified", strconv.FormatBool(*request.IsVerified))
		}
		if request.Limit != nil {
			query.Set("limit", strconv.FormatUint(uint64(*request.Limit), 10))
		}
		if request.PrevCursor != nil {
			query.Set("prevCursor", *request.PrevCursor)
		}
		if request.NextCursor != nil {
			query.Set("nextCursor", *request.NextCursor)
		}
	}

	return a.Get(t, "/api/v1/admin/users", query, cookies)
}
