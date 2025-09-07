// // package main

// // import (
// // 	"encoding/json"
// // 	"fmt"
// // 	"log"
// // 	"math"
// // 	"net/http"
// // 	"os"
// // 	"sort"
// // 	"time"
// // 	"bytes"

// // 	"drive_safe_server/db"
// // 	"github.com/go-resty/resty/v2"
// // "github.com/sfreiberg/gotwilio"  // Twilio SDK
// //     "firebase.google.com/go/messaging" // FCM
// // 	"github.com/joho/godotenv"
// // )

// // type VerbwireMintRequest struct {
// //     Blockchain   string `json:"blockchain"`   // e.g., "polygon"
// //     ContractAddress string `json:"contract_address"`
// //     Metadata struct {
// //         Name        string `json:"name"`
// //         Description string `json:"description"`
// //         Image       string `json:"image"` // URL to image
// //     } `json:"metadata"`
// //     RecipientAddress string `json:"recipient_address"`
// // }

// // func mintNFT(carID string, sosTimestamp string) {
// //     client := resty.New()
// //     apiKey := os.Getenv("VERBWIRE_API_KEY")

// //     req := VerbwireMintRequest{
// //         Blockchain:      "polygon",
// //         ContractAddress: os.Getenv("VERBWIRE_CONTRACT_ADDRESS"),
// //         RecipientAddress: os.Getenv("OWNER_WALLET_ADDRESS"),
// //     }

// //     req.Metadata.Name = fmt.Sprintf("SOS NFT - Car %s", carID)
// //     req.Metadata.Description = fmt.Sprintf("SOS triggered at %s for car %s", sosTimestamp, carID)
// //     req.Metadata.Image = "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0c/SOS.svg/512px-SOS.svg.png"
// // // placeholder image

// //     resp, err := client.R().
// //         SetHeader("Authorization", "Bearer "+apiKey).
// //         SetHeader("Content-Type", "application/json").
// //         SetBody(req).
// //         Post("https://api.verbwire.com/v1/nft/mint")

// //     if err != nil {
// //         log.Println("❌ Verbwire NFT minting failed:", err)
// //         return
// //     }

// //     log.Println("✅ Verbwire NFT minted! Response:", resp.String())
// // }

// // type ParkingSession struct {
// // 	SessionID string `json:"session_id"`
// // 	CarID     string `json:"car_id"`
// // 	StartedAt string `json:"started_at"`
// // 	StoppedAt string `json:"stopped_at,omitempty"`
// // 	Duration  int    `json:"duration_minutes,omitempty"`
// // }

// // type SOS struct {
// // 	CarID     string  `json:"car_id"`
// // 	Latitude  float64 `json:"latitude"`
// // 	Longitude float64 `json:"longitude"`
// // 	Timestamp string  `json:"timestamp"`
// // }

// // func distance(lat1, lon1, lat2, lon2 float64) float64 {
// // 	const R = 6371
// // 	dLat := (lat2 - lat1) * math.Pi / 180.0
// // 	dLon := (lon2 - lon1) * math.Pi / 180.0
// // 	lat1 = lat1 * math.Pi / 180.0
// // 	lat2 = lat2 * math.Pi / 180.0

// // 	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
// // 		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
// // 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
// // 	return R * c
// // }

// // func main() {
// // 	// Load environment variables
// // 	err := godotenv.Load()
// // 	if err != nil {
// // 		log.Println("⚠️ .env file not found, falling back to system env")
// // 	}

// // 	dbConn := db.ConnectAndInit()

// // 	// --- APIs ---

// // 	// Latest car location
// // 	http.HandleFunc("/api/location", func(w http.ResponseWriter, r *http.Request) {
// // 		var latitude, longitude float64
// // 		err := dbConn.QueryRow(`SELECT latitude, longitude FROM car_data ORDER BY updated_at DESC LIMIT 1`).Scan(&latitude, &longitude)
// // 		if err != nil {
// // 			http.Error(w, "No location data", http.StatusNotFound)
// // 			return
// // 		}
// // 		json.NewEncoder(w).Encode(map[string]float64{
// // 			"latitude":  latitude,
// // 			"longitude": longitude,
// // 		})
// // 	})

// // 	// Start parking
// // 	http.HandleFunc("/api/parking/start", func(w http.ResponseWriter, r *http.Request) {
// // 		carID := r.URL.Query().Get("car_id")
// // 		if carID == "" {
// // 			carID = "CAR123"
// // 		}
// // 		sessionID := fmt.Sprintf("PARK-%d", time.Now().Unix())

// // 		_, err := dbConn.Exec(`INSERT INTO parking_sessions (session_id, car_id, started_at) VALUES ($1, $2, $3)`,
// // 			sessionID, carID, time.Now().UTC())
// // 		if err != nil {
// // 			http.Error(w, "Failed to start parking session", http.StatusInternalServerError)
// // 			return
// // 		}

// // 		json.NewEncoder(w).Encode(ParkingSession{
// // 			SessionID: sessionID,
// // 			CarID:     carID,
// // 			StartedAt: time.Now().Format(time.RFC3339),
// // 		})
// // 	})

// // 	// Stop parking
// // 	http.HandleFunc("/api/parking/stop", func(w http.ResponseWriter, r *http.Request) {
// // 		sessionID := r.URL.Query().Get("session_id")
// // 		if sessionID == "" {
// // 			http.Error(w, "session_id required", http.StatusBadRequest)
// // 			return
// // 		}

// // 		var startedAt time.Time
// // 		err := dbConn.QueryRow(`SELECT started_at FROM parking_sessions WHERE session_id=$1`, sessionID).Scan(&startedAt)
// // 		if err != nil {
// // 			http.Error(w, "Session not found", http.StatusNotFound)
// // 			return
// // 		}

// // 		duration := int(time.Since(startedAt.UTC()).Minutes())

// // 		_, err = dbConn.Exec(`UPDATE parking_sessions SET stopped_at=$1, duration_minutes=$2 WHERE session_id=$3`,
// // 			time.Now(), duration, sessionID)
// // 		if err != nil {
// // 			http.Error(w, "Failed to stop session", http.StatusInternalServerError)
// // 			return
// // 		}

// // 		json.NewEncoder(w).Encode(ParkingSession{
// // 			SessionID: sessionID,
// // 			Duration:  duration,
// // 		})
// // 	})

// // 	// Nearest parking
// // 	http.HandleFunc("/api/parking/nearest", func(w http.ResponseWriter, r *http.Request) {
// // 		latStr := r.URL.Query().Get("lat")
// // 		lngStr := r.URL.Query().Get("lng")
// // 		if latStr == "" || lngStr == "" {
// // 			http.Error(w, "lat and lng query parameters required", http.StatusBadRequest)
// // 			return
// // 		}

// // 		var carLat, carLng float64
// // 		fmt.Sscanf(latStr, "%f", &carLat)
// // 		fmt.Sscanf(lngStr, "%f", &carLng)

// // 		parkingSpots := []map[string]interface{}{
// // 			{"name": "Central Parking", "lat": 28.614, "lng": 77.210, "available_slots": 5},
// // 			{"name": "East Side Parking", "lat": 28.615, "lng": 77.211, "available_slots": 2},
// // 			{"name": "West End Parking", "lat": 28.612, "lng": 77.208, "available_slots": 3},
// // 		}

// // 		for _, spot := range parkingSpots {
// // 			sLat := spot["lat"].(float64)
// // 			sLng := spot["lng"].(float64)
// // 			spot["distance_km"] = distance(carLat, carLng, sLat, sLng)
// // 		}

// // 		sort.Slice(parkingSpots, func(i, j int) bool {
// // 			return parkingSpots[i]["distance_km"].(float64) < parkingSpots[j]["distance_km"].(float64)
// // 		})

// // 		json.NewEncoder(w).Encode(parkingSpots)
// // 	})

// // 	// SOS
// // 	// http.HandleFunc("/api/sos", func(w http.ResponseWriter, r *http.Request) {
// // 	// 	carID := r.URL.Query().Get("car_id")
// // 	// 	lat := r.URL.Query().Get("lat")
// // 	// 	lng := r.URL.Query().Get("lng")

// // 	// 	if carID == "" || lat == "" || lng == "" {
// // 	// 		http.Error(w, "car_id, lat, lng required", http.StatusBadRequest)
// // 	// 		return
// // 	// 	}

// // 	// 	var latitude, longitude float64
// // 	// 	fmt.Sscanf(lat, "%f", &latitude)
// // 	// 	fmt.Sscanf(lng, "%f", &longitude)

// // 	// 	_, err := dbConn.Exec(`INSERT INTO sos_events (car_id, latitude, longitude) VALUES ($1, $2, $3)`,
// // 	// 		carID, latitude, longitude)
// // 	// 	if err != nil {
// // 	// 		http.Error(w, "Failed to send SOS", http.StatusInternalServerError)
// // 	// 		return
// // 	// 	}

// // 	// 	json.NewEncoder(w).Encode(SOS{
// // 	// 		CarID:     carID,
// // 	// 		Latitude:  latitude,
// // 	// 		Longitude: longitude,
// // 	// 		Timestamp: time.Now().Format(time.RFC3339),
// // 	// 	})
// // 	// })
// // 	///////////////////////////////////////////////////////////////////////////////////////////
// // 	http.HandleFunc("/api/sos", func(w http.ResponseWriter, r *http.Request) {
// // 		carID := r.URL.Query().Get("car_id")
// // 		lat := r.URL.Query().Get("lat")
// // 		lng := r.URL.Query().Get("lng")
// // 		sosTimestamp := time.Now().Format(time.RFC3339)
// //     mintNFT(carID, sosTimestamp)

// // 		if carID == "" || lat == "" || lng == "" {
// // 			http.Error(w, "car_id, lat, lng required", http.StatusBadRequest)
// // 			return
// // 		}

// // 		var latitude, longitude float64
// // 		fmt.Sscanf(lat, "%f", &latitude)
// // 		fmt.Sscanf(lng, "%f", &longitude)

// // 		// Save SOS in database
// // 		_, err := dbConn.Exec(`INSERT INTO sos_events (car_id, latitude, longitude) VALUES ($1, $2, $3)`,
// // 			carID, latitude, longitude)
// // 		if err != nil {
// // 			http.Error(w, "Failed to save SOS", http.StatusInternalServerError)
// // 			return
// // 		}

// // 		// Send Twilio SMS
// // 		twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
// // 		twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
// // 		twilioFrom := os.Getenv("TWILIO_PHONE_NUMBER")
// // 		emergencyNumber := os.Getenv("EMERGENCY_CONTACT")

// // 		client := twilio.NewClient(twilioSID, twilioToken, nil)
// // 		sendRealSMS := false // set true only for final demo....as we cant but twilio now lol

// // 		message := fmt.Sprintf(
// // 			"SOS! Car %s needs help at https://maps.google.com/?q=%f,%f",
// // 			carID, latitude, longitude,
// // 		)

// // 		if sendRealSMS {
// // 			_, err = client.Messages.SendMessage(twilioFrom, emergencyNumber, message, nil)
// // 			if err != nil {
// // 				log.Println("❌ Twilio SMS failed:", err)
// // 			}
// // 		} else {
// // 			log.Println("📱 Placeholder SMS:", message)
// // 		}

// // 		// _, err = client.Messages.SendMessage(twilioFrom, emergencyNumber,
// // 		//     fmt.Sprintf("SOS! Car %s needs help at https://maps.google.com/?q=%f,%f", carID, latitude, longitude), nil)
// // 		// if err != nil {
// // 		//     log.Println("❌ Twilio SMS failed:", err)
// // 		// }

// // 		// Send FCM push notification ....need to get the api after front end is done
// // 		fcmKeyFile := os.Getenv("FCM_SERVICE_ACCOUNT")
// // 		fcmClient, err := fcm.NewClient(fcmKeyFile)
// // 		if err != nil {
// // 			log.Println("❌ FCM client error:", err)
// // 		} else {
// // 			token := "DEVICE_REGISTRATION_TOKEN_FROM_APP"
// // 			msg := &fcm.Message{
// // 				To: token,
// // 				Notification: &fcm.Notification{
// // 					Title: "SOS Alert",
// // 					Body:  fmt.Sprintf("Car %s needs help at https://maps.google.com/?q=%f,%f", carID, latitude, longitude),
// // 				},
// // 			}
// // 			_, err := fcmClient.Send(msg)
// // 			if err != nil {
// // 				log.Println("❌ FCM push failed:", err)
// // 			}
// // 		}

// // 		// Return JSON response
// // 		json.NewEncoder(w).Encode(map[string]interface{}{
// // 			"car_id":    carID,
// // 			"latitude":  latitude,
// // 			"longitude": longitude,
// // 			"timestamp": time.Now().Format(time.RFC3339),
// // 		})
// // 	})

// // 	//////////////////////////////////////////////////////////////////////////////////////////////////////
// // 	// Get port from env or fallback to 8080
// // 	port := os.Getenv("PORT")
// // 	if port == "" {
// // 		port = "8080"
// // 	}

// //		fmt.Println("🚀 Backend running on http://localhost:" + port)
// //		log.Fatal(http.ListenAndServe(":"+port, nil))
// //	}
// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math"
// 	"net/http"
// 	"os"
// 	"sort"
// 	"time"

// 	"drive_safe_server/db"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/sfreiberg/gotwilio"

// 	firebase "firebase.google.com/go"
// 	"firebase.google.com/go/messaging"
// 	"google.golang.org/api/option"
// )

// // ---------------- Verbwire NFT ----------------

// type VerbwireMintRequest struct {
// 	Blockchain      string `json:"blockchain"`
// 	ContractAddress string `json:"contract_address"`
// 	Metadata        struct {
// 		Name        string `json:"name"`
// 		Description string `json:"description"`
// 		Image       string `json:"image"`
// 	} `json:"metadata"`
// 	RecipientAddress string `json:"recipient_address"`
// }

// func mintNFT(carID string, sosTimestamp string) {
// 	client := resty.New()
// 	apiKey := os.Getenv("VERBWIRE_API_KEY")

// 	req := VerbwireMintRequest{
// 		Blockchain:       "polygon",
// 		ContractAddress:  os.Getenv("VERBWIRE_CONTRACT_ADDRESS"),
// 		RecipientAddress: os.Getenv("OWNER_WALLET_ADDRESS"),
// 	}

// 	req.Metadata.Name = fmt.Sprintf("SOS NFT - Car %s", carID)
// 	req.Metadata.Description = fmt.Sprintf("SOS triggered at %s for car %s", sosTimestamp, carID)
// 	req.Metadata.Image = "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0c/SOS.svg/512px-SOS.svg.png"

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Bearer "+apiKey).
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(req).
// 		Post("https://api.verbwire.com/v1/nft/mint")

// 	if err != nil {
// 		log.Println("❌ Verbwire NFT minting failed:", err)
// 		return
// 	}

// 	if resp.StatusCode() != 200 {
// 		log.Println("❌ NFT mint failed:", resp.String())
// 		return
// 	}

// 	log.Println("✅ Verbwire NFT minted! Response:", resp.String())
// }

// // ---------------- Data Models ----------------

// type ParkingSession struct {
// 	SessionID string `json:"session_id"`
// 	CarID     string `json:"car_id"`
// 	StartedAt string `json:"started_at"`
// 	StoppedAt string `json:"stopped_at,omitempty"`
// 	Duration  int    `json:"duration_minutes,omitempty"`
// }

// type SOS struct {
// 	CarID     string  `json:"car_id"`
// 	Latitude  float64 `json:"latitude"`
// 	Longitude float64 `json:"longitude"`
// 	Timestamp string  `json:"timestamp"`
// }

// // ---------------- Helpers ----------------

// func distance(lat1, lon1, lat2, lon2 float64) float64 {
// 	const R = 6371
// 	dLat := (lat2 - lat1) * math.Pi / 180.0
// 	dLon := (lon2 - lon1) * math.Pi / 180.0
// 	lat1 = lat1 * math.Pi / 180.0
// 	lat2 = lat2 * math.Pi / 180.0

// 	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
// 		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
// 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
// 	return R * c
// }

// // ---------------- Main ----------------

// func main() {
// 	// Load environment variables
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("⚠️ .env file not found, falling back to system env")
// 	}

// 	dbConn := db.ConnectAndInit()

// 	// --- APIs ---

// 	// Latest car location
// 	http.HandleFunc("/api/location", func(w http.ResponseWriter, r *http.Request) {
// 		var latitude, longitude float64
// 		err := dbConn.QueryRow(`SELECT latitude, longitude FROM car_data ORDER BY updated_at DESC LIMIT 1`).Scan(&latitude, &longitude)
// 		if err != nil {
// 			http.Error(w, "No location data", http.StatusNotFound)
// 			return
// 		}
// 		json.NewEncoder(w).Encode(map[string]float64{
// 			"latitude":  latitude,
// 			"longitude": longitude,
// 		})
// 	})

// 	// Start parking
// 	http.HandleFunc("/api/parking/start", func(w http.ResponseWriter, r *http.Request) {
// 		carID := r.URL.Query().Get("car_id")
// 		if carID == "" {
// 			carID = "CAR123"
// 		}
// 		sessionID := fmt.Sprintf("PARK-%d", time.Now().Unix())

// 		_, err := dbConn.Exec(`INSERT INTO parking_sessions (session_id, car_id, started_at) VALUES ($1, $2, $3)`,
// 			sessionID, carID, time.Now().UTC())
// 		if err != nil {
// 			http.Error(w, "Failed to start parking session", http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(ParkingSession{
// 			SessionID: sessionID,
// 			CarID:     carID,
// 			StartedAt: time.Now().Format(time.RFC3339),
// 		})
// 	})

// 	// Stop parking
// 	http.HandleFunc("/api/parking/stop", func(w http.ResponseWriter, r *http.Request) {
// 		sessionID := r.URL.Query().Get("session_id")
// 		if sessionID == "" {
// 			http.Error(w, "session_id required", http.StatusBadRequest)
// 			return
// 		}

// 		var startedAt time.Time
// 		err := dbConn.QueryRow(`SELECT started_at FROM parking_sessions WHERE session_id=$1`, sessionID).Scan(&startedAt)
// 		if err != nil {
// 			http.Error(w, "Session not found", http.StatusNotFound)
// 			return
// 		}

// 		duration := int(time.Since(startedAt.UTC()).Minutes())

// 		_, err = dbConn.Exec(`UPDATE parking_sessions SET stopped_at=$1, duration_minutes=$2 WHERE session_id=$3`,
// 			time.Now(), duration, sessionID)
// 		if err != nil {
// 			http.Error(w, "Failed to stop session", http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(ParkingSession{
// 			SessionID: sessionID,
// 			Duration:  duration,
// 		})
// 	})

// 	// Nearest parking
// 	http.HandleFunc("/api/parking/nearest", func(w http.ResponseWriter, r *http.Request) {
// 		latStr := r.URL.Query().Get("lat")
// 		lngStr := r.URL.Query().Get("lng")
// 		if latStr == "" || lngStr == "" {
// 			http.Error(w, "lat and lng query parameters required", http.StatusBadRequest)
// 			return
// 		}

// 		var carLat, carLng float64
// 		fmt.Sscanf(latStr, "%f", &carLat)
// 		fmt.Sscanf(lngStr, "%f", &carLng)

// 		parkingSpots := []map[string]interface{}{
// 			{"name": "Central Parking", "lat": 28.614, "lng": 77.210, "available_slots": 5},
// 			{"name": "East Side Parking", "lat": 28.615, "lng": 77.211, "available_slots": 2},
// 			{"name": "West End Parking", "lat": 28.612, "lng": 77.208, "available_slots": 3},
// 		}

// 		for _, spot := range parkingSpots {
// 			sLat := spot["lat"].(float64)
// 			sLng := spot["lng"].(float64)
// 			spot["distance_km"] = distance(carLat, carLng, sLat, sLng)
// 		}

// 		sort.Slice(parkingSpots, func(i, j int) bool {
// 			return parkingSpots[i]["distance_km"].(float64) < parkingSpots[j]["distance_km"].(float64)
// 		})

// 		json.NewEncoder(w).Encode(parkingSpots)
// 	})

// 	// SOS
// 	http.HandleFunc("/api/sos", func(w http.ResponseWriter, r *http.Request) {
// 		carID := r.URL.Query().Get("car_id")
// 		lat := r.URL.Query().Get("lat")
// 		lng := r.URL.Query().Get("lng")

// 		if carID == "" || lat == "" || lng == "" {
// 			http.Error(w, "car_id, lat, lng required", http.StatusBadRequest)
// 			return
// 		}

// 		var latitude, longitude float64
// 		fmt.Sscanf(lat, "%f", &latitude)
// 		fmt.Sscanf(lng, "%f", &longitude)

// 		sosTimestamp := time.Now().Format(time.RFC3339)

// 		// Save SOS in database
// 		_, err := dbConn.Exec(`INSERT INTO sos_events (car_id, latitude, longitude) VALUES ($1, $2, $3)`,
// 			carID, latitude, longitude)
// 		if err != nil {
// 			http.Error(w, "Failed to save SOS", http.StatusInternalServerError)
// 			return
// 		}

// 		// Mint NFT
// 		mintNFT(carID, sosTimestamp)

// 		// Twilio SMS (placeholder unless you enable)
// 		twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
// 		twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
// 		twilioFrom := os.Getenv("TWILIO_PHONE_NUMBER")
// 		emergencyNumber := os.Getenv("EMERGENCY_CONTACT")

// 		twilio := gotwilio.NewTwilioClient(twilioSID, twilioToken)

// 		sendRealSMS := false
// 		message := fmt.Sprintf(
// 			"SOS! Car %s needs help at https://maps.google.com/?q=%f,%f",
// 			carID, latitude, longitude,
// 		)

// 		if sendRealSMS {
// 			_, _, err := twilio.SendSMS(twilioFrom, emergencyNumber, message, "", "")
// 			if err != nil {
// 				log.Println("❌ Twilio SMS failed:", err)
// 			}
// 		} else {
// 			log.Println("📱 Placeholder SMS:", message)
// 		}

// 		// Firebase Cloud Messaging
// 		fcmKeyFile := os.Getenv("FCM_SERVICE_ACCOUNT")
// 		ctx := context.Background()
// 		app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(fcmKeyFile))
// 		if err != nil {
// 			log.Println("❌ FCM init failed:", err)
// 		} else {
// 			fcmClient, err := app.Messaging(ctx)
// 			if err != nil {
// 				log.Println("❌ FCM client error:", err)
// 			} else {
// 				token := "DEVICE_REGISTRATION_TOKEN_FROM_APP"
// 				_, err := fcmClient.Send(ctx, &messaging.Message{
// 					Token: token,
// 					Notification: &messaging.Notification{
// 						Title: "SOS Alert",
// 						Body:  message,
// 					},
// 				})
// 				if err != nil {
// 					log.Println("❌ FCM push failed:", err)
// 				} else {
// 					log.Println("✅ FCM push sent")
// 				}
// 			}
// 		}

// 		// Return JSON response
// 		json.NewEncoder(w).Encode(SOS{
// 			CarID:     carID,
// 			Latitude:  latitude,
// 			Longitude: longitude,
// 			Timestamp: sosTimestamp,
// 		})
// 	})

// 	// Start server
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	fmt.Println("🚀 Backend running on http://localhost:" + port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	"drive_safe_server/db"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type VerbwireMintRequest struct {
	Blockchain      string `json:"blockchain"`
	ContractAddress string `json:"contract_address"`
	Metadata        struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	} `json:"metadata"`
	RecipientAddress string `json:"recipient_address"`
}

func mintNFT(carID string, sosTimestamp string) {
	apiKey := os.Getenv("VERBWIRE_API_KEY")
	contract := os.Getenv("VERBWIRE_CONTRACT_ADDRESS")
	wallet := os.Getenv("OWNER_WALLET_ADDRESS")

	if apiKey == "" || contract == "" || wallet == "" {
		log.Println("⚠️ Verbwire NFT keys missing, skipping NFT minting")
		return
	}

	client := resty.New()
	req := VerbwireMintRequest{
		Blockchain:      "polygon",
		ContractAddress: contract,
		RecipientAddress: wallet,
	}
	req.Metadata.Name = fmt.Sprintf("SOS NFT - Car %s", carID)
	req.Metadata.Description = fmt.Sprintf("SOS triggered at %s for car %s", sosTimestamp, carID)
	req.Metadata.Image = "https://upload.wikimedia.org/wikipedia/commons/thumb/0/0c/SOS.svg/512px-SOS.svg.png"

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post("https://api.verbwire.com/v1/nft/mint")

	if err != nil {
		log.Println("❌ Verbwire NFT minting failed:", err)
		return
	}

	log.Println("✅ Verbwire NFT minted! Response:", resp.String())
}

type ParkingSession struct {
	SessionID string `json:"session_id"`
	CarID     string `json:"car_id"`
	StartedAt string `json:"started_at"`
	StoppedAt string `json:"stopped_at,omitempty"`
	Duration  int    `json:"duration_minutes,omitempty"`
}

type SOS struct {
	CarID     string  `json:"car_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, falling back to system env")
	}

	dbConn := db.ConnectAndInit()

	// Latest car location
	http.HandleFunc("/api/location", func(w http.ResponseWriter, r *http.Request) {
		var latitude, longitude float64
		err := dbConn.QueryRow(`SELECT latitude, longitude FROM car_data ORDER BY updated_at DESC LIMIT 1`).Scan(&latitude, &longitude)
		if err != nil {
			http.Error(w, "No location data", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{
			"latitude":  latitude,
			"longitude": longitude,
		})
	})

	// Start parking
	http.HandleFunc("/api/parking/start", func(w http.ResponseWriter, r *http.Request) {
		carID := r.URL.Query().Get("car_id")
		if carID == "" {
			carID = "CAR123"
		}
		sessionID := fmt.Sprintf("PARK-%d", time.Now().Unix())

		_, err := dbConn.Exec(`INSERT INTO parking_sessions (session_id, car_id, started_at) VALUES ($1, $2, $3)`,
			sessionID, carID, time.Now().UTC())
		if err != nil {
			http.Error(w, "Failed to start parking session", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ParkingSession{
			SessionID: sessionID,
			CarID:     carID,
			StartedAt: time.Now().Format(time.RFC3339),
		})
	})

	// Stop parking
	http.HandleFunc("/api/parking/stop", func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "session_id required", http.StatusBadRequest)
			return
		}

		var startedAt time.Time
		err := dbConn.QueryRow(`SELECT started_at FROM parking_sessions WHERE session_id=$1`, sessionID).Scan(&startedAt)
		if err != nil {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}

		duration := int(time.Since(startedAt.UTC()).Minutes())
		_, err = dbConn.Exec(`UPDATE parking_sessions SET stopped_at=$1, duration_minutes=$2 WHERE session_id=$3`,
			time.Now(), duration, sessionID)
		if err != nil {
			http.Error(w, "Failed to stop session", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ParkingSession{
			SessionID: sessionID,
			Duration:  duration,
		})
	})

	// Nearest parking
	http.HandleFunc("/api/parking/nearest", func(w http.ResponseWriter, r *http.Request) {
		latStr := r.URL.Query().Get("lat")
		lngStr := r.URL.Query().Get("lng")
		if latStr == "" || lngStr == "" {
			http.Error(w, "lat and lng query parameters required", http.StatusBadRequest)
			return
		}

		var carLat, carLng float64
		fmt.Sscanf(latStr, "%f", &carLat)
		fmt.Sscanf(lngStr, "%f", &carLng)

		parkingSpots := []map[string]interface{}{
			{"name": "Central Parking", "lat": 28.614, "lng": 77.210, "available_slots": 5},
			{"name": "East Side Parking", "lat": 28.615, "lng": 77.211, "available_slots": 2},
			{"name": "West End Parking", "lat": 28.612, "lng": 77.208, "available_slots": 3},
		}

		for _, spot := range parkingSpots {
			sLat := spot["lat"].(float64)
			sLng := spot["lng"].(float64)
			spot["distance_km"] = distance(carLat, carLng, sLat, sLng)
		}

		sort.Slice(parkingSpots, func(i, j int) bool {
			return parkingSpots[i]["distance_km"].(float64) < parkingSpots[j]["distance_km"].(float64)
		})

		json.NewEncoder(w).Encode(parkingSpots)
	})

	// SOS API
	http.HandleFunc("/api/sos", func(w http.ResponseWriter, r *http.Request) {
		carID := r.URL.Query().Get("car_id")
		lat := r.URL.Query().Get("lat")
		lng := r.URL.Query().Get("lng")
		sosTimestamp := time.Now().Format(time.RFC3339)

		if carID == "" || lat == "" || lng == "" {
			http.Error(w, "car_id, lat, lng required", http.StatusBadRequest)
			return
		}

		var latitude, longitude float64
		fmt.Sscanf(lat, "%f", &latitude)
		fmt.Sscanf(lng, "%f", &longitude)

		// Mint NFT only if ENABLE_NFT_MINT=true
		if os.Getenv("ENABLE_NFT_MINT") == "true" {
			mintNFT(carID, sosTimestamp)
		} else {
			log.Println("⚠️ NFT minting skipped (ENABLE_NFT_MINT not true)")
		}

		// Save SOS in database
		_, err := dbConn.Exec(`INSERT INTO sos_events (car_id, latitude, longitude) VALUES ($1, $2, $3)`,
			carID, latitude, longitude)
		if err != nil {
			http.Error(w, "Failed to save SOS", http.StatusInternalServerError)
			return
		}

		// Placeholder SMS
		message := fmt.Sprintf("SOS! Car %s needs help at https://maps.google.com/?q=%f,%f",
			carID, latitude, longitude)
		log.Println("📱 Placeholder SMS:", message)

		// Return JSON response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"car_id":    carID,
			"latitude":  latitude,
			"longitude": longitude,
			"timestamp": sosTimestamp,
		})
	})

	// Get port from env or fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Backend running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
