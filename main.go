package main

import (
    "flag"
    "fmt"
    "net/http"
    "os"
)

// Ensure origin is globally accessible
var origin *string

func main() {
    portNumber := flag.Int("port", 3000, "port number")
    origin = flag.String("origin", "", "origin server URL to forward requests")
    clearCache := flag.Bool("clear-cache", false, "Clear the cache and exit")

    flag.Parse()

    // If --clear-cache flag is set, clear cache and exit
    if *clearCache {
        cache.Clear() // ✅ Now cache is recognized
        fmt.Println("Cache cleared successfully.")
        os.Exit(0)
    }

    if *origin == "" {
        fmt.Println("Error: --origin is required. Example usage:")
        fmt.Println("  caching-proxy --port 3000 --origin http://dummyjson.com")
        os.Exit(1) // Exit with an error
    }

    fmt.Printf("Starting caching proxy on port %d, forwarding to %s\n", *portNumber, *origin)

    addr := fmt.Sprintf(":%d", *portNumber)
    http.HandleFunc("/", proxyHandler) // ✅ Now proxyHandler is recognized

    err := http.ListenAndServe(addr, nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
        os.Exit(1)
    }
}
