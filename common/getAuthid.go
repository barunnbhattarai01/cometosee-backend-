package common

import "context"

func GetAuthid(ctx context.Context) float64 {
	authid, ok := ctx.Value("authId").(float64)
	if !ok {
		return 0
	}
	return authid
}
