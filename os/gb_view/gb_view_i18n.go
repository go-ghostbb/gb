package gbview

import (
	"context"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	gbconv "ghostbb.io/gb/util/gb_conv"
)

const (
	i18nLanguageVariableName = "I18nLanguage"
)

// i18nTranslate translate the content with i18n feature.
func (view *View) i18nTranslate(ctx context.Context, content string, variables Params) string {
	if view.config.I18nManager != nil {
		// Compatible with old version.
		if language, ok := variables[i18nLanguageVariableName]; ok {
			ctx = gbi18n.WithLanguage(ctx, gbconv.String(language))
		}
		return view.config.I18nManager.T(ctx, content)
	}
	return content
}

// setI18nLanguageFromCtx retrieves language name from context and sets it to template variables map.
func (view *View) setI18nLanguageFromCtx(ctx context.Context, variables map[string]interface{}) {
	if _, ok := variables[i18nLanguageVariableName]; !ok {
		if language := gbi18n.LanguageFromCtx(ctx); language != "" {
			variables[i18nLanguageVariableName] = language
		}
	}
}
