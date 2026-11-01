package outbox

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type OutboxId valueobjects.Id

func (id OutboxId) String() string {
	return valueobjects.Id(id).String()
}
