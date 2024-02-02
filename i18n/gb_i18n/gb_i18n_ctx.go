// Package gbi18n implements internationalization and localization.
package gbi18n

import (
	"context"
	gbctx "ghostbb.io/gb/os/gb_ctx"
)

const (
	ctxLanguage gbctx.StrKey = "I18nLanguage"
)

// WithLanguage append language setting to the context and returns a new context.
func WithLanguage(ctx context.Context, language string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	return context.WithValue(ctx, ctxLanguage, language)
}

// LanguageFromCtx retrieves and returns language name from context.
// It returns an empty string if it is not set previously.
func LanguageFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(ctxLanguage)
	if v != nil {
		return v.(string)
	}
	return ""
}
