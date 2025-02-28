package suites

import (
	"fmt"
	"os"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
)

// BaseDomain the base domain.
var BaseDomain = "example.com:8080"

// PathPrefix the prefix/url_base of the login portal.
var PathPrefix = os.Getenv("PathPrefix")

// LoginBaseURL the base URL of the login portal.
var LoginBaseURL = fmt.Sprintf("https://login.%s", BaseDomain)

// SingleFactorBaseURL the base URL of the singlefactor domain.
var SingleFactorBaseURL = fmt.Sprintf("https://singlefactor.%s", BaseDomain)

// AdminBaseURL the base URL of the admin domain.
var AdminBaseURL = fmt.Sprintf("https://admin.%s", BaseDomain)

// MailBaseURL the base URL of the mail domain.
var MailBaseURL = fmt.Sprintf("https://mail.%s", BaseDomain)

// HomeBaseURL the base URL of the home domain.
var HomeBaseURL = fmt.Sprintf("https://home.%s", BaseDomain)

// PublicBaseURL the base URL of the public domain.
var PublicBaseURL = fmt.Sprintf("https://public.%s", BaseDomain)

// SecureBaseURL the base URL of the secure domain.
var SecureBaseURL = fmt.Sprintf("https://secure.%s", BaseDomain)

// DevBaseURL the base URL of the dev domain.
var DevBaseURL = fmt.Sprintf("https://dev.%s", BaseDomain)

// MX1MailBaseURL the base URL of the mx1.mail domain.
var MX1MailBaseURL = fmt.Sprintf("https://mx1.mail.%s", BaseDomain)

// MX2MailBaseURL the base URL of the mx2.mail domain.
var MX2MailBaseURL = fmt.Sprintf("https://mx2.mail.%s", BaseDomain)

// OIDCBaseURL the base URL of the oidc domain.
var OIDCBaseURL = fmt.Sprintf("https://oidc.%s", BaseDomain)

// DuoBaseURL the base URL of the Duo configuration API.
var DuoBaseURL = "https://duo.example.com"

// AutheliaBaseURL the base URL of Authelia service.
var AutheliaBaseURL = "https://authelia.example.com:9091"

const (
	t            = "true"
	testUsername = "john"
	testPassword = "password"
)

const (
	namespaceAuthelia  = "authelia"
	namespaceDashboard = "kubernetes-dashboard"
	namespaceKube      = "kube-system"
)

var (
	storageLocalTmpConfig = schema.Configuration{
		TOTP: schema.TOTPConfiguration{
			Issuer: "Authelia",
			Period: 6,
		},
		Storage: schema.StorageConfiguration{
			EncryptionKey: "a_not_so_secure_encryption_key",
			Local: &schema.LocalStorageConfiguration{
				Path: "/tmp/db.sqlite3",
			},
		},
	}
)
