package types

// PasswordView is the scene shown in the password minigame.
type PasswordView int

const (
	PasswordOffice PasswordView = iota
	PasswordDesk
	PasswordNone
)

// String returns the BEM modifier for the view.
func (v PasswordView) String() string {
	if v == PasswordDesk {
		return "desk"
	}
	return "office"
}

// PasswordViewFromInt maps a raw route id to a PasswordView.
func PasswordViewFromInt(n int) (PasswordView, bool) {
	if n == int(PasswordOffice) || n == int(PasswordDesk) {
		return PasswordView(n), true
	}
	return 0, false
}
