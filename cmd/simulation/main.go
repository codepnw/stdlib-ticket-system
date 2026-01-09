package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	URL        = "http://localhost:8080/bookings"
	NumUsers   = 20
	EventID    = 13
	TargetSeat = 63
)

func main() {
	fmt.Println("Start Race Condition Simulation...")
	fmt.Printf("%d users fighting for Seat ID: %d\n", NumUsers, EventID)
	fmt.Println("---------------------------------")

	var wg sync.WaitGroup
	startSignal := make(chan struct{})

	var successCount int64
	var failCount int64

	// Users Bookings
	for i := 1; i <= NumUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// 1. Wait Signal
			<-startSignal

			// 2. Create Payload
			reqBody := map[string]any{
				"event_id": EventID,
				"seat_ids": []int{TargetSeat},
				"user_id":  userID,
			}
			jsonValue, _ := json.Marshal(reqBody)

			// 3. Send Request
			resp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Printf("User %d Network Error\n", userID)
				atomic.AddInt64(&failCount, 1)
				return
			}
			defer resp.Body.Close()

			// bodyBytes, _ := io.ReadAll(resp.Body)
			// bodyString := string(bodyBytes)

			// 4. Result
			if resp.StatusCode == 200 || resp.StatusCode == 201 {
				fmt.Printf("-----> User %d Success!\n", userID)
				atomic.AddInt64(&successCount, 1)
			} else {
				fmt.Printf("User %d Failed (%d)\n", userID, resp.StatusCode)
				atomic.AddInt64(&failCount, 1)
			}
		}(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Done! (All requests sent successfully)")

	close(startSignal)

	wg.Wait()

	fmt.Println("------------------------------------------------")
	fmt.Println("📊 Summary Report")
	fmt.Printf("Total Requests: %d\n", NumUsers)
	fmt.Printf("✅ Success:      %d  (Should be 1 only)\n", successCount)
	fmt.Printf("❌ Failed:       %d\n", failCount)
	fmt.Println("------------------------------------------------")

	if successCount == 1 {
		fmt.Println("🏆 TEST PASSED: System handled race condition correctly!")
	} else if successCount > 1 {
		fmt.Println("😱 TEST FAILED: Double Booking Detected! (Check your Lock)")
	} else {
		fmt.Println("🤔 TEST WEIRD: No one got the ticket? (Check Seat ID / Logic)")
	}
}
