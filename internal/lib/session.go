package lib

import "time"

var SESSION_EXPIRED_TIME = ParseExpiredTime(ParseTimeDuration(REFRESH_TOKEN_TTL, 7, time.Hour*24))
