package auth

import (
	"net/http"
	"strings"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
	"github.com/swaggest/openapi-go/openapi3"
	"google.golang.org/api/idtoken"
)

// AuthorizationHeader const for header used for authorization token
const AuthorizationHeader = "Authorization"

// GCloudServiceAccount is a middleware that will check AuthorizationHeader for incoming request.
// The expected header value is an OpenID Token generated by Google Cloud for given ServiceAccount.
type GCloudServiceAccount struct {
	ServiceAccount string
}

// Wrap implements the handler.Middleware interface
func (m GCloudServiceAccount) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		token := r.Header.Get(AuthorizationHeader)
		if token == "" {
			return nil, errors.Forbidden("invalid_token", "missing token")
		}

		splitAuthHeader := strings.Split(token, " ")
		if len(splitAuthHeader) == 0 {
			return nil, errors.Forbidden("invalid_token", "missing token")
		}

		if err := m.VerifyServiceAccount(r, splitAuthHeader[1]); err != nil {
			return nil, err
		}

		return h.Serve(w, r)
	})
}

// VerifyServiceAccount check if the token was sent by a gcloud service account
func (m GCloudServiceAccount) VerifyServiceAccount(r *http.Request, token string) error {
	payload, err := idtoken.Validate(r.Context(), token, "")
	if err != nil {
		// invalid token
		return errors.Forbidden("invalid_token", "failed to validate token")
	}

	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		return errors.Forbidden("invalid_token", "invalid issuer")
	}

	if payload.Claims == nil {
		return errors.Forbidden("invalid_token", "missing claims")
	}

	if emailVerified := payload.Claims["email_verified"].(bool); !emailVerified && payload.Claims["email_verified"] != "true" {
		return errors.Forbidden("invalid_token", "invalid token email")
	}

	if payload.Claims["email"] != m.ServiceAccount {
		return errors.Forbidden("invalid_token", "invalid token email")
	}
	return nil
}

// IsEligible checks if the request is eligible for the middleware
func (m GCloudServiceAccount) IsEligible(r *http.Request) bool {
	token := r.Header.Get(AuthorizationHeader)
	if token == "" {
		return false
	}

	splitAuthHeader := strings.Split(token, " ")
	if len(splitAuthHeader) == 0 {
		return false
	}

	if len(splitAuthHeader) > 1 {
		payload, err := idtoken.Validate(r.Context(), splitAuthHeader[1], "")
		if err != nil {
			// invalid token
			return false
		}

		if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
			return false
		}

		return true
	}

	return false
}

const securitySchemeKey = "gcloud_service_account"

// Doc implements the openapi.Documented interface
func (m GCloudServiceAccount) Doc(builder *openapi.DocBuilder) error {
	_, ok := builder.Reflector().
		SpecEns().
		ComponentsEns().
		SecuritySchemesEns().
		MapOfSecuritySchemeOrRefValues[securitySchemeKey]

	if !ok {
		openIDScheme := new(openapi3.SecurityScheme).
			WithOpenIDConnectSecurityScheme(
				*new(openapi3.OpenIDConnectSecurityScheme).
					WithDescription(`
This auth middleware is useful for endpoint called by a GCloud service for example:
- Cloud Scheduler: https://cloud.google.com/scheduler/docs/http-target-auth
- Pub/Sub: https://cloud.google.com/pubsub/docs/push#authentication
`),
			)

		authSchema := *new(openapi3.SecuritySchemeOrRef).WithSecurityScheme(*openIDScheme)
		builder.Reflector().
			SpecEns().
			ComponentsEns().
			SecuritySchemesEns().
			WithMapOfSecuritySchemeOrRefValuesItem(
				securitySchemeKey,
				authSchema,
			)
	}

	builder.Operation().
		WithSecurity(map[string][]string{
			securitySchemeKey: {
				m.ServiceAccount,
			},
		})

	return nil
}
