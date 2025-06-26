// package main

// import (
// 	"context"
// 	"log"
// 	"time"

// )

// func main() {

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err := redisClinet.Redis.SAdd(ctx, "rbac:admin", "/url/shorten", "/url/redirect/:shortCode", "/url/stats/:shortCode").Err()
// 	if err != nil {
// 		log.Fatal("❌ Failed to seed admin routes:", err)
// 	}
// 	err = redisClinet.Redis.SAdd(ctx, "rbac:user", "/url/shorten", "/url/redirect/:shortCode").Err()
// 	if err != nil {
// 		log.Fatal("❌ Failed to seed user routes:", err)
// 	}
// }

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	redisClinet "github.com/prajwalzaware/go-urlShortner/pkg/redis"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Now REDIS_URL will be available to os.Getenv()
	ctx := context.Background()
	redisClinet.ConnectRedis()

	// Example: Allow "admin" role to access /admin/stats and /shorten
	errr := redisClinet.Redis.SAdd(ctx, "rbac:admin", "/url/shorten", "/url/redirect/:shortCode", "/url/stats/:shortCode").Err()
	if errr != nil {
		log.Fatal("❌ Failed to seed admin routes:", err)
	}
	err = redisClinet.Redis.SAdd(ctx, "rbac:user", "/url/shorten", "/url/redirect/:shortCode").Err()
	if err != nil {
		log.Fatal("❌ Failed to seed user routes:", err)
	}

	fmt.Println("✅ Roles seeded to Redis")
}
