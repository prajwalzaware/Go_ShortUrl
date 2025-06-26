package background

import (
	"context"
	"log"
	"time"

	"github.com/prajwalzaware/go-urlShortner/utils"
)

func ClickFlusher() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	log.Println("🚀 Click flusher started...")

	for {
		select {
		case <-ticker.C:
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				log.Println("📝 Flushing click counts to DB...")
				utils.FlushClicksToDB(ctx)
			}()

		}
	}
}
