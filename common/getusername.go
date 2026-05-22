package common

import "context"

func GetUsername(ctx context.Context) string {
	username, ok := ctx.Value("username").(string)
	if !ok {
		return ""
	}
	return username
}
