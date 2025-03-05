package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/maphash"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
)

const version = "1.0"

var (
	bufferQueue [][]byte
	queueLock   sync.Mutex
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	// Kubernetes check if app is ok
	http.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "up")
	})

	// Kubernetes check if app can serve requests
	http.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "yes")
	})

	http.HandleFunc("/", handler)
	http.HandleFunc("/clear", clearHandler)
	fmt.Println("Listening now at port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Seed the random number generator with the current time
	myRandom := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))

	// Create a buffer of 1MB
	buffer := make([]byte, 1024*1024)

	// Fill the buffer with completely random data
	_, err := myRandom.Read(buffer)
	if err != nil {
		http.Error(w, "Failed to generate random data", http.StatusInternalServerError)
		return
	}

	// Run the SHA1 algorithm on the buffer
	hash := sha1.New()
	hash.Write(buffer)
	sha1Result := hex.EncodeToString(hash.Sum(nil))

	// Get memory and CPU usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memUsed := float64(memStats.Alloc) / (1024 * 1024) // Convert bytes to MB

	// Add the buffer to the global queue
	queueLock.Lock()
	bufferQueue = append(bufferQueue, buffer)
	queueLock.Unlock()

	// Respond to the request
	fmt.Fprintf(w, "Application version: %s\n", version)
	fmt.Fprintf(w, "Memory used by the golang runtime: %.2f MB\n", memUsed)
	fmt.Fprintf(w, "SHA1 result of the buffer: %s\n", sha1Result)
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the global queue
	queueLock.Lock()
	bufferQueue = nil
	queueLock.Unlock()

	// Respond to the request
	fmt.Fprintln(w, "Buffer queue cleared")
}
