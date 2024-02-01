package gins

import (
	gbi18n "ghostbb.io/i18n/gb_i18n"
)

// I18n returns an instance of gi18n.Manager.
// The parameter `name` is the name for the instance.
func I18n(name ...string) *gbi18n.Manager {
	return gbi18n.Instance(name...)
}
