package enum

type LoanStatus int

const (
	_ LoanStatus = iota
	ActiveLoan
	OverdueLoan
	ReturnedLoan
)

var loanStatusState = map[LoanStatus]string{
	ActiveLoan:   "active",
	OverdueLoan:  "overdue",
	ReturnedLoan: "returned",
}

func (s LoanStatus) String() string {
	return loanStatusState[s]
}
