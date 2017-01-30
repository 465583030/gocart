package engine

type (
	// MailSender interface
	MailSender interface {
		Send(to []string, subject string, body []byte) error
	}

	// Validator interface
	Validator interface {
		CheckEmail(string) error
		CheckRequired(string, string) error
		CheckStringLen(s string, min int, max int, field string) error
	}

	// JWTSignParser interface
	JWTSignParser interface {
		Sign(claims map[string]interface{}, secret string) (string, error)
		Parse(tokenStr string, secret string) (map[string]interface{}, error)
	}

	// Factory engine factory interface
	Factory interface {
		NewCatalog() Catalog
		NewMail() Mailer
		NewUser() User
	}

	factory struct {
		StorageFactory
		ms  MailSender
		v   Validator
		jwt JWTSignParser
	}
)

// New instances new engine factory
func New(sf StorageFactory, ms MailSender, v Validator, jwt JWTSignParser) Factory {
	return &factory{
		StorageFactory: sf,
		ms:             ms,
		v:              v,
		jwt:            jwt,
	}
}
