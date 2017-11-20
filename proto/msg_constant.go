package proto

import "sibo/share"

var (
	CREATE_PLAYER, _ = share.EncodeMessageID(1, 1000)
	LOGIN, _         = share.EncodeMessageID(1, 1001)
	CHANGENAME, _    = share.EncodeMessageID(1, 1002)
)