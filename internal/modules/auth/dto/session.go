package authdto

import (
	"encoding/json"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/voutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	userdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/dto"
	"github.com/google/uuid"
)

type Session struct {
	UserId    userdomain.UserId
	Email     valueobjects.Email
	FirstName *userdomain.FirstName
	LastName  *userdomain.LastName
	Role      userdomain.UserRole
	AvatarUrl *valueobjects.Url
}

type sessionData struct {
	UserId    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Role      string    `json:"role"`
	AvatarUrl *string   `json:"avatarUrl"`
}

func (s Session) MarshalJSON() ([]byte, error) {
	data := sessionData{
		UserId:    uuid.UUID(s.UserId),
		Email:     s.Email.Value(),
		FirstName: ptrutil.ValuePtr(s.FirstName),
		LastName:  ptrutil.ValuePtr(s.LastName),
		Role:      s.Role.Value(),
		AvatarUrl: ptrutil.ValuePtr(s.AvatarUrl),
	}

	return json.Marshal(data)
}

func (s *Session) UnmarshalJSON(data []byte) error {
	var raw sessionData

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	userEmail, _ := valueobjects.NewEmail(raw.Email)
	userRole, _ := userdomain.NewUserRole(raw.Role)

	firstName := voutil.Restore(raw.FirstName, userdomain.NewFirstName)
	lastName := voutil.Restore(raw.LastName, userdomain.NewLastName)
	avatarUrl := voutil.Restore(raw.AvatarUrl, valueobjects.NewUrl)

	*s = Session{
		UserId:    userdomain.UserId(raw.UserId),
		Email:     userEmail,
		FirstName: firstName,
		LastName:  lastName,
		Role:      userRole,
		AvatarUrl: avatarUrl,
	}

	return nil
}

func (s *Session) ToUserDto() userdto.UserDto {
	return userdto.UserDto{
		Id:        uuid.UUID(s.UserId),
		Email:     s.Email.Value(),
		FirstName: ptrutil.ValuePtr(s.FirstName),
		LastName:  ptrutil.ValuePtr(s.LastName),
		Role:      s.Role.Value(),
		AvatarUrl: ptrutil.ValuePtr(s.AvatarUrl),
	}
}
