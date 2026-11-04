package adminuser_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	adminuserconst "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/const"
	adminuserdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/dto"
	adminuserroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/admin/user/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/Desalutar20/lingostruct-server-go/tests/integration/testapp"
)

func TestGetAdminUsers(t *testing.T) {
	t.Run("Should Return Ok When Request Is Valid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			cookies := app.CreateAndSignIn(t, &request, ptrutil.To(userdomain.UserRoleAdmin()))

			response := app.GetAdminUsers(t, nil, cookies)
			defer response.Body.Close()

			var body routes.KeysetPaginatedResponseDto[*adminuserdto.AdminUserDto]
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			app.AssertValidStatus(t, response, http.StatusOK)
		})
	})

	t.Run("Should Return BadRequest When Data Is Invalid", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			cookies := app.CreateAndSignIn(t, &request, ptrutil.To(userdomain.UserRoleAdmin()))

			var getUsersRequest = adminuserroutes.GetAdminUsersRequest{}

			data := []struct {
				request func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest
				field   string
				name    string
			}{
				{
					field: "search",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.Search = ptrutil.To("")
						return r
					},
					name: "empty search",
				},
				{
					field: "search",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.Search = ptrutil.To("   ")
						return r
					},
					name: "whitespace search",
				},

				{
					field: "search",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.Search = ptrutil.To(strings.Repeat("t", adminuserconst.GetAdminUsersSearchMaxLimit+1))
						return r
					},
					name: "search exceeds maximum length",
				},
				{
					field: "prevCursor",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.PrevCursor = ptrutil.To("")
						return r
					},
					name: "empty prev cursor",
				},
				{
					field: "prevCursor",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.PrevCursor = ptrutil.To("   ")
						return r
					},
					name: "whitespace prev cursor",
				},
				{
					field: "nextCursor",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.NextCursor = ptrutil.To("")
						return r
					},
					name: "empty next cursor",
				},
				{
					field: "nextCursor",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.NextCursor = ptrutil.To("   ")
						return r
					},
					name: "whitespace next cursor",
				},
				{
					field: "limit",
					request: func(r adminuserroutes.GetAdminUsersRequest) adminuserroutes.GetAdminUsersRequest {
						r.Limit = ptrutil.To(adminuserconst.GetAdminUsersMaxLimit + 1)
						return r
					},
					name: "",
				},
			}

			for _, data := range data {
				t.Run(data.name, func(t *testing.T) {
					t.Parallel()

					r := data.request(getUsersRequest)
					response := app.GetAdminUsers(t, &r, cookies)
					defer response.Body.Close()

					app.AssertValidStatus(t, response, http.StatusBadRequest)
					app.AssertValidationError(t, response, data.field)

				})
			}

			response := app.GetAdminUsers(t, nil, cookies)
			defer response.Body.Close()

			var body routes.KeysetPaginatedResponseDto[*adminuserdto.AdminUserDto]
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			app.AssertValidStatus(t, response, http.StatusOK)
		})
	})

	t.Run("Should Return Unauthorized When User Is Not Signed In", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			response := app.GetAdminUsers(t, nil, nil)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusUnauthorized)
		})
	})

	t.Run("Should Return AccessForbidden When User Is Not Admin", func(t *testing.T) {
		t.Parallel()

		testapp.Run(t, func(app *testapp.TestApp) {
			request := testapp.SignUpRequest()
			cookies := app.CreateAndSignIn(t, &request, ptrutil.To(userdomain.UserRoleRegular()))

			response := app.GetAdminUsers(t, nil, cookies)
			defer response.Body.Close()

			app.AssertValidStatus(t, response, http.StatusForbidden)

		})
	})
}
