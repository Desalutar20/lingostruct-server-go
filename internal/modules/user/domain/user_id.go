package userdomain

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type UserId valueobjects.Id

func (id UserId) String() string {
	return valueobjects.Id(id).String()
}
