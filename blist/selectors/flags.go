package selectors

type EventFlags uint64

func (t EventFlags) XOR(x EventFlags) EventFlags {
	return t ^ x
}

func (e EventFlags) Val() uint64 {
	return uint64(e)
}

const (
	EventFlagsNoFlags                   EventFlags = 0
	EventFlagsInvalid                   EventFlags = 1 << 0
	EventFlagsRegular                   EventFlags = 1 << 1
	EventFlagsGuestList                 EventFlags = 1 << 2
	EventFlagsFeatured                  EventFlags = 1 << 3
	EventFlagsBoxOffice                 EventFlags = 1 << 4
	EventFlagsQRCode                    EventFlags = 1 << 5
	EventFlagsBarCode                   EventFlags = 1 << 6
	EventFlagsOnSale                    EventFlags = 1 << 7
	EventFlagsOffSale                   EventFlags = 1 << 8
	EventFlagsNoWaitingList             EventFlags = 1 << 9
	EventFlagsCodeLocked                EventFlags = 1 << 10
	EventFlagsNoDayOfTheEvent           EventFlags = 1 << 11
	EventFlagsBranded                   EventFlags = 1 << 12
	EventFlagsGenerateNewCodeOnTransfer EventFlags = 1 << 13
	EventFlagsStripeConnect             EventFlags = 1 << 14
)
