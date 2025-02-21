package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

// Check user from external API
func checkUser(userID string) (bool, error) {
	// resp, err := http.Get(fmt.Sprintf("https://external-api.com/check-user/%s", userID))
	// if err != nil {
	// 	return false, err
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return false, err
	// }

	// var result map[string]interface{}
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	return false, err
	// }

	// isValid, ok := result["valid"].(bool)
	// if !ok {
	// 	return false, fmt.Errorf("invalid response format")
	// }

	return true, nil
}

// Check balance from external API
func checkBalance(userID string) (float64, error) {
	// 	resp, err := http.Get(fmt.Sprintf("https://external-api.com/check-balance/%s", userID))
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	var result map[string]interface{}
	// 	err = json.Unmarshal(body, &result)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	balance, ok := result["balance"].(float64)
	// 	if !ok {
	// 		return 0, fmt.Errorf("invalid response format")
	// 	}

	return 100, nil
}

func main() {
	initRedis()

	router := gin.Default()

	router.GET("/check/:userID", func(c *gin.Context) {
		userID := c.Param("userID")

		// Check cache for user validation
		cacheKey := fmt.Sprintf("user:%s", userID)
		isValid, err := rdb.Get(ctx, cacheKey).Bool()
		if err == redis.Nil {
			// Cache miss, call external API
			valid, err := checkUser(userID)
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
			// Cache miss, call external API
			bal, err := checkBalance(userID)
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
			"userID":  userID,
			"valid":   isValid,
			"balance": balance,
		})
	})

	router.Run(":8080")
}
