package middleware

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/ulule/limiter/v3"
"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit(rps string) gin.HandlerFunc {
rate, _ := limiter.NewRateFromFormatted(rps)
store := memory.NewStore()
instance := limiter.New(store, rate)
return func(c *gin.Context) {
ctx, err := instance.Get(c, c.ClientIP())
if err != nil || ctx.Reached {
c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
return
}
c.Next()
}
}
