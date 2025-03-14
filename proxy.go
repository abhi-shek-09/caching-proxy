package main

import (
	"fmt"
	"io"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request){
	targetURL := *origin + r.URL.Path
	fmt.Printf("Received request: %s %s\n", r.Method, targetURL)

	// check if its in cache
	if cachedResponse, found := cache.Get(targetURL); found{
		// set the header as hit, write to the header that the status is ok and write the cached response and return
		fmt.Println("HIT: ", targetURL)
		w.Header().Set("X-Cache", "HIT")
		
		w.WriteHeader(http.StatusOK)
		w.Write(cachedResponse)
		return
	}

	// if not in cache, forward it to the server
	// create a client, create a request with the same method, body and targeturl
	// copy the request header into the new one
	// client will do the request and get the response, CLOSE THE RESPONSE BODY and read the response into a body
	// VV IMP, set the cache and set the header as miss, writeheader the response status and write the body and return
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
        return
	}

	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
        return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
        return
	}

	cache.Set(targetURL, body)

	w.Header().Set("X-Cache", "MISS")
	for key, values := range resp.Header {
        // Skip the Content-Disposition header to prevent downloads
        if key == "Content-Disposition" {
            continue
        }
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}