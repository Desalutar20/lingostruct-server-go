package email

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"

type EmailMessage struct {
	Subject   valueobjects.NonEmptyString  `json:"subject"`
	To        valueobjects.Email           `json:"to"`
	PlainText valueobjects.NonEmptyString  `json:"plainText"`
	HtmlText  *valueobjects.NonEmptyString `json:"htmlText,omitempty"`
}
