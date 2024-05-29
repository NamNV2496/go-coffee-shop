package ratelimit

import (
	"strconv"
	"strings"
	"time"

	"github.com/ulule/limiter/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Rate struct {
	limiter   *limiter.Limiter
	Formatted string
	Period    time.Duration
	Limit     int64
}

func NewRateFromFormatted(formatted string) (Rate, error) {
	rate := Rate{}

	values := strings.Split(formatted, "-")
	if len(values) != 2 {
		return rate, status.Error(codes.Internal, "failed to convert orders to newData template")
	}

	periods := map[string]time.Duration{
		"S": time.Second,    // Second
		"M": time.Minute,    // Minute
		"H": time.Hour,      // Hour
		"D": time.Hour * 24, // Day
	}

	limit, period := values[0], strings.ToUpper(values[1])

	p, ok := periods[period]
	if !ok {
		return rate, status.Error(codes.Internal, "incorrect period")
	}

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return rate, status.Error(codes.Internal, "incorrect limit")
	}

	rate = Rate{
		Formatted: formatted,
		Period:    p,
		Limit:     l,
	}

	return rate, nil
}
