package valueobjects

import (
	"net/url"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
)

type Url struct {
	NonEmptyString
}

func NewUrl(value string) (Url, error) {
	return NewNonEmptyStringBase[Url](value, "url", "Url", 0, 0,
		func(value string) error {
			u, err := url.Parse(value)
			if err != nil {
				return domain.ValidationError("url", []string{
					"Invalid URL",
				})
			}

			if u.Scheme != "http" || u.Scheme != "https" {
				return domain.ValidationError("url", []string{
					"Invalid URL",
				})
			}

			return nil
		},
		func(value NonEmptyString) Url {
			return Url{NonEmptyString: value}
		})
}
