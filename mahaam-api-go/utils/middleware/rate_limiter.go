package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]
	if !exists {
		// 5 requests/minute with burst of 5
		// once limit reached, client can try after 12 seconds
		limiter := rate.NewLimiter(rate.Every(12*time.Second), 5)
		clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func RateLimiterMW() gin.HandlerFunc {
	go cleanupClients()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getClient(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
