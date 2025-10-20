package enum

type ReservStatus int

const (
	_ ReservStatus = iota
	PendingReserv
	FulfilledReserv
	CancelledReserv
)

var reservString = map[ReservStatus]string{
	PendingReserv:   "pending",
	FulfilledReserv: "fulfilled",
	CancelledReserv: "cancelled",
}

func (s ReservStatus) String() string {
	return reservString[s]
}
