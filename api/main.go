package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
)

var ctx = context.Background()
var rdb *redis.Client

// Initialize Redis client
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis container address
		Password: "",           // No password
		DB:       0,            // Default DB
	})
}

// Circuit Breaker Configuration
func newCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,                // Allow up to 5 requests in the half-open state
		Interval:    10 * time.Second, // Reset the stats after 10 seconds
		Timeout:     5 * time.Second,  // Stay open for 5 seconds before transitioning to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.6 // Open if 60% of requests fail
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit Breaker '%s' changed from %v to %v\n", name, from, to)
		},
	})
}

// Check user from external API with Circuit Breaker
func checkUser(userID string, cb *gobreaker.CircuitBreaker) (bool, error) {
	var isValid bool

	_, execErr := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(fmt.Sprintf("https://external-api.com/check-user/%s", userID))
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return false, err
		}

		isValid, ok := result["valid"].(bool)
		if !ok {
			return false, fmt.Errorf("invalid response format")
		}

		return isValid, nil
	})

	if execErr != nil {
		return false, execErr
	}

	return isValid, nil
}

// Check balance from external API with Circuit Breaker
func checkBalance(userID string, cb *gobreaker.CircuitBreaker) (float64, error) {
	var balance float64

	_, execErr := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(fmt.Sprintf("https://external-api.com/check-balance/%s", userID))
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return 0, err
		}

		balance, ok := result["balance"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid response format")
		}

		return balance, nil
	})

	if execErr != nil {
		return 0, execErr
	}

	return balance, nil
}

func main() {
	initRedis()

	// Initialize Circuit Breakers
	userCircuitBreaker := newCircuitBreaker("user-check")
	balanceCircuitBreaker := newCircuitBreaker("balance-check")

	router := gin.Default()

	router.GET("/check/:userID", func(c *gin.Context) {
			userID := c.Param("userID")

			// Check cache for user validation
			cacheKey := fmt.Sprintf("user:%s", userID)
			isValid, err := rdb.Get(ctx, cacheKey).Bool()
			if err == redis.Nil {
					// Cache miss, call external API with Circuit Breaker
					valid, err := checkUser(userID, userCircuitBreaker)
					if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check user"})
							return
					}
					isValid = valid
					rdb.Set(ctx, cacheKey, isValid, 5*time.Minute) // Cache for 5 minutes
			} else if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "cache error"})
					return
			}

			if !isValid {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
					return
			}

			// Check cache for balance
			balanceCacheKey := fmt.Sprintf("balance:%s", userID)
			balance, err := rdb.Get(ctx, balanceCacheKey).Float64()
			if err == redis.Nil {
					// Cache miss, call external API with Circuit Breaker
					bal, err := checkBalance(userID, balanceCircuitBreaker)
					if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check balance"})
							return
					}
					balance = bal
					rdb.Set(ctx, balanceCacheKey, balance, 5*time.Minute) // Cache for 5 minutes
			} else if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "cache error"})
					return
			}

			c.JSON(http.StatusOK, gin.H{
					"userID": userID,
					"valid":  isValid,
					"balance": balance,
			})
	})

	router.Run(":8080")
}
