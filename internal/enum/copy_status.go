package enum

type CopyStatus int

const (
	_ CopyStatus = iota
	AvailableCopy
	LoanedCopy
	ReservedCopy
	DamagedCopy
	LostCopy
)

var copyStatusState = map[CopyStatus]string{
	AvailableCopy: "available",
	LoanedCopy:    "loaned",
	ReservedCopy:  "reserved",
	DamagedCopy:   "damaged",
	LostCopy:      "lost",
}

func (s CopyStatus) String() string {
	return copyStatusState[s]
}
