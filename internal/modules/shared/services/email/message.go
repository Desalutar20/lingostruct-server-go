package email

import (
	"encoding/json"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/voutil"
)

type EmailMessage struct {
	Subject   valueobjects.NonEmptyString
	To        valueobjects.Email
	PlainText valueobjects.NonEmptyString
	HtmlText  *valueobjects.NonEmptyString
}

type emailMessageData struct {
	Subject   string  `json:"subject" db:"subject"`
	To        string  `json:"to" db:"to"`
	PlainText string  `json:"plainText" db:"plain_text"`
	HtmlText  *string `json:"htmlText,omitempty" db:"html_text"`
}

func (s EmailMessage) MarshalJSON() ([]byte, error) {
	data := emailMessageData{
		Subject:   s.Subject.Value(),
		To:        s.To.Value(),
		PlainText: s.PlainText.Value(),
		HtmlText:  ptrutil.ValuePtr(s.HtmlText),
	}

	return json.Marshal(data)
}

func (m *EmailMessage) UnmarshalJSON(data []byte) error {
	var raw emailMessageData

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	subject, _ := valueobjects.NewNonEmptyString(raw.Subject, "subject", "Subject", 0, 0, nil)
	to, _ := valueobjects.NewEmail(raw.To)
	plainText, _ := valueobjects.NewNonEmptyString(raw.PlainText, "plainText", "Plain text", 0, 0, nil)

	*m = EmailMessage{
		Subject:   subject,
		To:        to,
		PlainText: plainText,
		HtmlText: voutil.Restore(raw.HtmlText, func(val string) (valueobjects.NonEmptyString, error) {
			return valueobjects.NewNonEmptyString(val, "htmlText", "Html text", 0, 0, nil)
		}),
	}

	return nil
}
