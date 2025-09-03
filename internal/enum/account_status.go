package enum

type AccountStatus int

const (
	_ AccountStatus = iota
	ActiveAccount
	SuspendedAccount
)

var accountStatusState = map[AccountStatus]string{
	ActiveAccount:    "active",
	SuspendedAccount: "suspended",
}

func (s AccountStatus) String() string {
	return accountStatusState[s]
}
