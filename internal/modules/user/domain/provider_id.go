package userdomain

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type ProviderId struct {
	valueobjects.NonEmptyString
}

const ProviderIdMaxLength = 100

func NewProviderId(value string) (ProviderId, error) {
	return valueobjects.NewNonEmptyStringBase[ProviderId](value, "providerId", "Provider ID", 0, ProviderIdMaxLength,
		nil,
		func(value valueobjects.NonEmptyString) ProviderId {
			return ProviderId{NonEmptyString: value}
		},
	)
}
