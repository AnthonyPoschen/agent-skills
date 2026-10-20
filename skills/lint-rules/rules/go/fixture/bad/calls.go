package bad

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func calls() {
	_ = os.Getenv("APP_ENV")
	_, _ = http.Get("https://example.com")
	fmt.Println("debug")
	log.Fatal("stop")
}
