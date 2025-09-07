# 🚗 TriVision-Drive-Safe

A Smart car safety and parking management system built for hackathons.  
It provides live car tracking, SOS alerts with SMS + push notifications, parking session management, and blockchain-based SOS proof using NFTs.
    
---

## ✨ Features
- 📍 Show current location on a map
- 🅿️ Start/Stop parking and track duration
- 🚨 SOS button → sends Twilio SMS + Firebase push notification
- 🗄️ Stores data in SQLite/Postgres
- 🎨 Bonus: Mints SOS proof NFT using Verbwire.

---

## ⚙️ Tech Stack
- **Backend:** Golang (with REST APIs)  
- **Frontend:** Flutter (cross-platform app)  
- **Database:** PostgreSQL (Supabase)  
- **Notifications:** Twilio (SMS), Firebase Cloud Messaging (Push)  
- **Blockchain:** Verbwire API (NFT minting for SOS proof)  
- **Frameworks/Libs:** gofr (optional), godotenv, lib/pq  

---

## 🚀 Features Implemented
- ✅ Real-time location API (`/api/location`)  
- ✅ Parking APIs (`/api/parking/start`, `/api/parking/stop`, `/api/parking/nearest`)  
- ✅ SOS API (`/api/sos`) → Sends Twilio SMS + FCM push + stores in DB  
- ✅ NFT Proof of SOS (`Verbwire integration`)  
- ✅ Secure DB with Supabase  

---

## 🚀 Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/<your-username>/TriVision-Drive-Safe.git
cd TriVision-Drive-Safein the zone 
