package userdomain

import (
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/option"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/voutil"
	"github.com/google/uuid"
)

type User struct {
	domain.Entity[UserId]
	firstName      *FirstName
	lastName       *LastName
	email          valueobjects.Email
	hashedPassword *HashedPassword
	role           UserRole
	isBanned       bool
	isVerified     bool
	avatarUrl      *valueobjects.Url
	googleId       *ProviderId
	githubId       *ProviderId
}

func (u *User) FirstName() *FirstName {
	return u.firstName
}

func (u *User) LastName() *LastName {
	return u.lastName
}

func (u *User) Email() valueobjects.Email {
	return u.email
}

func (u *User) HashedPassword() *HashedPassword {
	return u.hashedPassword
}

func (u *User) Role() UserRole {
	return u.role
}

func (u *User) IsBanned() bool {
	return u.isBanned
}

func (u *User) IsVerified() bool {
	return u.isVerified
}

func (u *User) AvatarUrl() *valueobjects.Url {
	return u.avatarUrl
}

func (u *User) GoogleId() *ProviderId {
	return u.googleId
}

func (u *User) GithubId() *ProviderId {
	return u.githubId
}

func (u *User) IsValid() bool {
	return !u.isBanned && u.isVerified
}

func (u *User) Verify() error {
	if u.isBanned {
		return domain.FailureError("User is banned")
	}

	u.isVerified = true
	return nil
}

func (u *User) UpdatePassword(newPassword HashedPassword) {
	u.hashedPassword = &newPassword
}

func (u *User) LinkOAuthProvider(provider OAuthProvider, providerId ProviderId) error {
	if provider == OAuthProviderGoogle() {
		if u.googleId != nil {
			return nil
		}

		u.googleId = &providerId
		return nil
	}

	if provider == OAuthProviderGithub() {
		if u.githubId != nil {
			return nil
		}

		u.githubId = &providerId
		return nil
	}

	return domain.FailureError("Unsupported OAuth provider: " + provider.Value())
}

func (u *User) Update(firstName option.Option[*FirstName], lastName option.Option[*LastName], avatarUrl option.Option[*valueobjects.Url], isBanned option.Option[bool]) bool {
	isUpdated := false

	if firstName.IsSome() && firstName.Unwrap() != u.firstName {
		u.firstName = firstName.Unwrap()
		isUpdated = true
	}

	if lastName.IsSome() && lastName.Unwrap() != u.lastName {
		u.lastName = lastName.Unwrap()
		isUpdated = true
	}

	if avatarUrl.IsSome() && avatarUrl.Unwrap() != u.avatarUrl {
		u.avatarUrl = avatarUrl.Unwrap()
		isUpdated = true
	}

	if isBanned.IsSome() && isBanned.Unwrap() != u.isBanned {
		u.isBanned = isBanned.Unwrap()
		isUpdated = true
	}

	return isUpdated
}

func NewUser(email valueobjects.Email, hashedPassword *HashedPassword, firstName *FirstName, lastName *LastName) User {
	return User{
		domain.NewEntity(UserId(uuid.New())),
		firstName,
		lastName,
		email,
		hashedPassword,
		UserRoleRegular(),
		false,
		false,
		nil,
		nil,
		nil,
	}
}

func RestoreUser(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	email string,
	hashedPassword *string,
	firstName *string,
	lastName *string,
	role string,
	isBanned bool,
	isVerified bool,
	avatarURL *string,
	googleID *string,
	githubID *string,
) User {
	userID := UserId(id)
	userEmail, _ := valueobjects.NewEmail(email)
	userRole, _ := NewUserRole(role)

	return User{
		Entity:         domain.RestoreEntity(userID, createdAt, updatedAt),
		firstName:      voutil.Restore(firstName, NewFirstName),
		lastName:       voutil.Restore(lastName, NewLastName),
		email:          userEmail,
		hashedPassword: voutil.Restore(hashedPassword, NewHashedPassword),
		role:           userRole,
		isBanned:       isBanned,
		isVerified:     isVerified,
		avatarUrl:      voutil.Restore(avatarURL, valueobjects.NewUrl),
		googleId:       voutil.Restore(googleID, NewProviderId),
		githubId:       voutil.Restore(githubID, NewProviderId),
	}
}
