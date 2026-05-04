package common

import "context"

func GetEmail(ctx context.Context) string {
	email, ok := ctx.Value("email").(string)
	if !ok {
		return ""
	}
	return email
}
