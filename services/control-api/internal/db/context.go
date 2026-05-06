package db

import "context"

type contextKey string

const orgIDKey contextKey = "org_id"

func WithOrgID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, orgIDKey, orgID)
}

func OrgIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(orgIDKey).(string)
	return v
}
