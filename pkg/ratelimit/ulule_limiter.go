package ratelimit

/*
// Create a rate with the given limit (number of requests) for the given
// period (a time.Duration of your choice).
import "github.com/ulule/limiter/v3"

rate := limiter.Rate{
    Period: 1 * time.Hour,
    Limit:  1000,
}


// You can also use the simplified format "<limit>-<period>"", with the given
// periods:
//
// * "S": second
// * "M": minute
// * "H": hour
// * "D": day
//
// Examples:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
//
rate, err := limiter.NewRateFromFormatted("1000-H")
if err != nil {
    panic(err)
}

// Then, create a store. Here, we use the bundled Redis store. Any store
// compliant to limiter.Store interface will do the job. The defaults are
// "limiter" as Redis key prefix and a maximum of 3 retries for the key under
// race condition.
import "github.com/ulule/limiter/v3/drivers/store/redis"

store, err := redis.NewStore(client)
if err != nil {
    panic(err)
}

// Alternatively, you can pass options to the store with the "WithOptions"
// function. For example, for Redis store:
import "github.com/ulule/limiter/v3/drivers/store/redis"

store, err := redis.NewStoreWithOptions(pool, limiter.StoreOptions{
    Prefix:   "your_own_prefix",
})
if err != nil {
    panic(err)
}

// Or use a in-memory store with a goroutine which clears expired keys.
import "github.com/ulule/limiter/v3/drivers/store/memory"

store := memory.NewStore()

// Then, create the limiter instance which takes the store and the rate as arguments.
// Now, you can give this instance to any supported middleware.
instance := limiter.New(store, rate)

// Alternatively, you can pass options to the limiter instance with several options.
instance := limiter.New(store, rate, limiter.WithClientIPHeader("True-Client-IP"), limiter.WithIPv6Mask(mask))

// Finally, give the limiter instance to your middleware initializer.
import "github.com/ulule/limiter/v3/drivers/middleware/stdlib"

middleware := stdlib.NewMiddleware(instance)
*/
