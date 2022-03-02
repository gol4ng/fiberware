package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

const (
	XAuthorizationHeader = "X-Authorization"
)

func FromHeader(request *fiber.Request) CredentialProvider {
	return func() Credential {
		return ExtractFromHeader(request)
	}
}

func ExtractFromHeader(request *fiber.Request) Credential {
	if request == nil {
		return ""
	}

	tokenHeader := request.Header.Peek(fiber.HeaderAuthorization)
	if len(tokenHeader) == 0 {
		tokenHeader = request.Header.Peek(XAuthorizationHeader)
	}

	return utils.CopyString(utils.UnsafeString(tokenHeader))
}

func AddHeader(request *fiber.Request) CredentialSetter {
	return func(credential Credential) {
		if request == nil {
			return
		}
		if creds, ok := credential.(string); ok {
			request.Header.Set(fiber.HeaderAuthorization, creds)
			request.Header.Set(XAuthorizationHeader, creds)
		}
	}
}
