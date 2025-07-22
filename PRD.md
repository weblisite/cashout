# Product Requirements Document (PRD)
## Project Name: Cashout
**Version:** 1.0  
**Date:** July 22, 2025  
**Author:** Based on Grok 3, xAI conversation  
**Target Audience:** Kenyan smartphone users (60% penetration, GSMA 2023), agents, and businesses  
**Purpose:** To develop a cross-platform payment app named Cashout (iOS, Android, web) offering P2P transfers, cash-in/cash-out via agents, and business payments, using a native solution.

---

## 1. Overview

### 1.1 Product Description
Cashout is a fintech application designed to provide a secure, scalable, and user-friendly payment solution for Kenya. It enables peer-to-peer (P2P) money transfers, commission-free cash-in and fee-based cash-out services through a network of agents, and QR-based business payment processing. The platform leverages a custom transaction engine, Supabase for database/storage/authentication, and Flutter for cross-platform deployment, targeting Kenya's 30M+ smartphone users (CA 2024).

### 1.2 Goals
- Provide competitive mobile money transfer rates with transparent fee structures on Cashout
- Offer a commission-free cash-in experience to boost user adoption on Cashout
- Provide a scalable platform supporting transactions up to 1,000,000 KES on Cashout
- Ensure compliance with CBK's KYC/AML regulations for Cashout users
- Implement advanced fraud detection using both on-device and cloud-based technologies

### 1.3 Key Features
- **User Registration & Authentication:** Phone-based signup with OTP via Supabase Auth on Cashout
- **P2P Transfers:** Instant transfers with a platform margin of 100% of the fee on Cashout
- **Cash-In/Cash-Out:** Agent-mediated with a 75%/25% fee split (platform/agent) on Cashout
- **Business Payments:** QR-based transactions with a 50/50 fee split between user and business on Cashout
- **Real-Time Notifications:** WebSocket-based updates on Cashout
- **Security:** TLS 1.3, AES-256, OAuth 2.0, and on-device fraud detection for Cashout

---

## 2. Technical Architecture

### 2.1 Technology Stack

#### 2.1.1 Frontend (Cross-Platform)
**Framework:** Flutter  
**Language:** Dart  
**Use Case:** iOS, Android, and web apps for P2P transfers, cash-in/cash-out, and business payments  
**Why:** Single codebase, high performance with Skia engine, and native feature support (e.g., QR scanning, biometrics) via packages like `qr_code_scanner` and `local_auth`

**Implementation:**
- **Mobile:** Native iOS (Swift-based) and Android (Kotlin-based) apps
- **Web:** Flutter for Web, hosted via Supabase's storage or a custom server
- **Tools:** VS Code or Android Studio with Flutter SDK

#### 2.1.2 Backend
**Languages/Frameworks:**
- **Go (Golang):** Core transaction engine and API server
- **Node.js:** Real-time features (e.g., WebSocket for notifications)

**Why:**
- **Go:** High performance for <5-second transaction settlement, ideal for microservices
- **Node.js:** JavaScript ecosystem aligns with Flutter, lightweight for push updates

**Tools:**
- Go: Gorilla Mux for routing, custom SQL queries for Supabase integration
- Node.js: Express.js, Socket.IO

#### 2.1.3 Database, Storage, and Authentication
**Platform:** Supabase  
**Use Case:** Database, file storage, and user authentication for the app

**Why:**
- **Database:** Managed PostgreSQL with real-time capabilities, scalable for 1M+ users
- **Storage:** Secure file storage for KYC documents (e.g., ID scans) and QR codes
- **Authentication:** Built-in support for phone number login with OTP, email/password, and social logins
- **Client:** Supabase JavaScript/TypeScript client for Flutter integration

**Implementation:**
- **Database Schema:**
  - `users`: id, phone_number (unique), hashed_id, kyc_status, wallet_balance
  - `transactions`: id, user_id, recipient_id, amount, fee, type (cash-in/out/P2P/business), timestamp
  - `agents`: id, phone_number, float_balance, commission_balance
  - `businesses`: id, name, wallet_balance

- **Real-Time:** Use Supabase's subscribe feature for live transaction updates
- **Storage:** Upload KYC files to `public/kyc/{user_id}` buckets, secured with Row-Level Security (RLS)
- **Authentication:** Integrate with Flutter using `supabase_flutter` package

#### 2.1.4 Security and Compliance
**Technologies:**
- **TLS 1.3 and AES-256:** Provided by Supabase for data encryption
- **OAuth 2.0:** Supabase's authentication supports secure token-based access
- **TensorFlow Lite:** On-device fraud detection, trainable on transaction patterns

**Compliance:** Supabase's RLS and audit logs meet CBK AML/KYC standards

#### 2.1.5 Additional Tools and Integrations
- **SMS Gateway:** Africa's Talking for OTPs, integrated with Supabase auth
- **KYC:** Use Supabase storage with OpenCV for ID verification
- **CI/CD:** GitHub Actions for automated Flutter builds
- **Testing:** Flutter Test for unit tests, Detox for end-to-end mobile testing

---

## 3. Functional Requirements

### 3.1 User Management

#### 3.1.1 Registration
- Users provide a phone number (e.g., +254712345678) on Cashout
- Supabase Auth sends an OTP via Africa's Talking SMS gateway
- Generate a unique hashed ID (SHA-256 of phone number, truncated to 10 digits + 4-digit random salt) on Cashout
- Store in Supabase users table: id, phone_number, hashed_id, kyc_status (enum: 'pending', 'verified', 'rejected'), wallet_balance (default 0 KES)

#### 3.1.2 Login
- Phone number + OTP authentication using Supabase Auth on Cashout
- Support biometric login (Flutter local_auth package) post-initial login on Cashout

#### 3.1.3 KYC
- Upload ID (e.g., Kenyan National ID) to Supabase storage (public/kyc/{user_id}) on Cashout
- Use OpenCV for ID validation (e.g., detect text, match format), storing kyc_status in users on Cashout
- Admin dashboard (Flutter Web) for manual review if automated fails

### 3.2 P2P Transfers

#### 3.2.1 Flow
- User selects recipient (via hashed ID or QR scan with qr_code_scanner) on Cashout
- Enters amount (e.g., 5,000 KES), fee calculated from Table 3.2.2 on Cashout
- Confirms with PIN/biometrics on Cashout
- Supabase inserts into transactions: id, user_id, recipient_id, amount, fee, type = 'p2p', timestamp on Cashout
- Updates users.wallet_balance for sender (-amount - fee) and recipient (+amount)

#### 3.2.2 Fee Structure
- Competitive fee structure based on transaction amounts on Cashout
- Platform margin = 100% of user fee on Cashout (no agent commission)

**Table 3.2.2: Cashout P2P Fee Structure**
| Transaction Range (KES) | User Fee (KES) | Platform Margin (KES) |
|------------------------|----------------|----------------------|
| 50 - 100 | 8.25 | 8.25 |
| 101 - 500 | 21.75 | 21.75 |
| 501 - 1,000 | 21.75 | 21.75 |
| 1,001 - 1,500 | 21.75 | 21.75 |
| 1,501 - 2,500 | 21.75 | 21.75 |
| 2,501 - 3,500 | 39.00 | 39.00 |
| 3,501 - 5,000 | 51.75 | 51.75 |
| 5,001 - 7,500 | 65.25 | 65.25 |
| 7,501 - 10,000 | 86.25 | 86.25 |
| 10,001 - 15,000 | 125.25 | 125.25 |
| 15,001 - 20,000 | 138.75 | 138.75 |
| 20,001 - 35,000 | 147.75 | 147.75 |
| 35,001 - 50,000 | 208.50 | 208.50 |
| 50,001 - 250,000 | 231.75 | 231.75 |
| 250,001 - 500,000 | 513.00 | 513.00 |
| 500,001 - 1,000,000 | 1,075.50 | 1,075.50 |

#### 3.2.3 Real-Time Updates
- Use Supabase subscribe for instant updates on Cashout (e.g., supabase.channel('transactions').on('INSERT', (payload) => notify(payload)))

### 3.3 Cash-In/Cash-Out

#### 3.3.1 Cash-In
- User visits agent, agent scans QR (generated by qr_code_scanner) or enters hashed ID on Cashout
- Agent inputs amount (e.g., 1,000 KES), confirms with PIN on Cashout
- Supabase updates users.wallet_balance (+amount) and agents.float_balance (-amount) on Cashout
- **No fee or commission on Cashout**

#### 3.3.2 Cash-Out
- User requests withdrawal (e.g., 500 KES) via Cashout, generating QR
- Agent scans QR, confirms with PIN, dispenses cash on Cashout
- Supabase updates users.wallet_balance (-amount - fee), agents.float_balance (+amount + fee), and agents.commission_balance (+agent commission) on Cashout

#### 3.3.3 Fee Structure
- Competitive fee structure based on transaction amounts on Cashout
- Agent commission = 25%, platform margin = 75% on Cashout

**Table 3.3.3: Cashout Cash-Out Fee Structure**
| Transaction Range (KES) | User Fee (KES) | Agent Commission (25%) | Platform Margin (75%) |
|------------------------|----------------|------------------------|----------------------|
| 50 - 100 | 8.25 | 2.0625 | 6.1875 |
| 101 - 500 | 21.75 | 5.4375 | 16.3125 |
| 501 - 1,000 | 21.75 | 5.4375 | 16.3125 |
| 1,001 - 1,500 | 21.75 | 5.4375 | 16.3125 |
| 1,501 - 2,500 | 21.75 | 5.4375 | 16.3125 |
| 2,501 - 3,500 | 39.00 | 9.7500 | 29.2500 |
| 3,501 - 5,000 | 51.75 | 12.9375 | 38.8125 |
| 5,001 - 7,500 | 65.25 | 16.3125 | 48.9375 |
| 7,501 - 10,000 | 86.25 | 21.5625 | 64.6875 |
| 10,001 - 15,000 | 125.25 | 31.3125 | 93.9375 |
| 15,001 - 20,000 | 138.75 | 34.6875 | 104.0625 |
| 20,001 - 35,000 | 147.75 | 36.9375 | 110.8125 |
| 35,001 - 50,000 | 208.50 | 52.1250 | 156.3750 |
| 50,001 - 250,000 | 231.75 | 57.9375 | 173.8125 |
| 250,001 - 500,000 | 513.00 | 128.2500 | 384.7500 |
| 500,001 - 1,000,000 | 1,075.50 | 268.8750 | 806.6250 |

#### 3.3.4 Agent App
- Lightweight Flutter app with QR scanner, PIN input, and offline sync (SQLite), deployed separately for agents

### 3.4 Business Payments

#### 3.4.1 User-Initiated
- User scans business QR (generated by business app), enters amount (e.g., 2,000 KES) on Cashout
- Transaction fee calculated from Table 3.4.3, split 50/50 into business fee and user fee on Cashout
- Confirms with PIN/biometrics on Cashout
- Supabase inserts into transactions: id, user_id, business_id, amount, fee, type = 'business', timestamp on Cashout
- Updates users.wallet_balance (-amount - user fee) and businesses.wallet_balance (+amount - business fee)

#### 3.4.2 Business-Initiated
- Business scans user QR or enters hashed ID, enters amount (e.g., 1,500 KES) on Cashout business app
- Same fee split applies on Cashout
- Confirms with business PIN on Cashout
- Updates businesses.wallet_balance (-amount - business fee) and users.wallet_balance (+amount - user fee)

#### 3.4.3 Fee Structure
- Competitive fee structure based on transaction amounts on Cashout
- Platform margin = 100% of transaction fee, split 50/50 into business fee and user fee on Cashout

**Table 3.4.3: Cashout Business Payment Fee Structure**
| Transaction Range (KES) | Transaction Fee (KES) | Business Fee (50%) | User Fee (50%) |
|------------------------|----------------------|-------------------|----------------|
| 50 - 100 | 8.25 | 4.125 | 4.125 |
| 101 - 500 | 21.75 | 10.875 | 10.875 |
| 501 - 1,000 | 21.75 | 10.875 | 10.875 |
| 1,001 - 1,500 | 21.75 | 10.875 | 10.875 |
| 1,501 - 2,500 | 21.75 | 10.875 | 10.875 |
| 2,501 - 3,500 | 39.00 | 19.500 | 19.500 |
| 3,501 - 5,000 | 51.75 | 25.875 | 25.875 |
| 5,001 - 7,500 | 65.25 | 32.625 | 32.625 |
| 7,501 - 10,000 | 86.25 | 43.125 | 43.125 |
| 10,001 - 15,000 | 125.25 | 62.625 | 62.625 |
| 15,001 - 20,000 | 138.75 | 69.375 | 69.375 |
| 20,001 - 35,000 | 147.75 | 73.875 | 73.875 |
| 35,001 - 50,000 | 208.50 | 104.250 | 104.250 |
| 50,001 - 250,000 | 231.75 | 115.875 | 115.875 |
| 250,001 - 500,000 | 513.00 | 256.500 | 256.500 |
| 500,001 - 1,000,000 | 1,075.50 | 537.750 | 537.750 |

### 3.5 Security Measures

#### 3.5.1 Authentication & Authorization
- **TLS 1.3:** All data in transit encrypted using TLS 1.3
- **AES-256:** Data at rest encrypted using AES-256
- **OAuth 2.0:** Secure token-based access for agents and businesses
- **Biometric Authentication:** Fingerprint/face ID support via Flutter local_auth package
- **PIN Protection:** 4-6 digit PIN for transaction confirmation

#### 3.5.2 On-Device Fraud Detection
- **TensorFlow Lite:** Lightweight AI model for anomaly detection
- **Transaction Pattern Analysis:** Monitor for unusual spending patterns
- **Device Fingerprinting:** Track device characteristics for fraud prevention
- **Offline Capability:** Fraud detection works without internet connection

#### 3.5.3 Cloud-Based Fraud Detection
- **Real-Time Monitoring:** Server-side fraud detection using machine learning
- **Velocity Checks:** Monitor transaction frequency and amounts
- **Geographic Analysis:** Flag transactions from unusual locations
- **Behavioral Analysis:** Track user behavior patterns over time

#### 3.5.4 Data Protection
- **Row-Level Security (RLS):** Supabase RLS policies for data access control
- **Audit Logging:** Complete transaction and access logs for compliance
- **Data Sovereignty:** Store sensitive data in Kenya-compliant regions
- **Regular Security Audits:** Quarterly security assessments

---

## 4. Technical Implementation

### 4.1 Database Schema

#### 4.1.1 Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    hashed_id VARCHAR(14) UNIQUE NOT NULL,
    kyc_status ENUM('pending', 'verified', 'rejected') DEFAULT 'pending',
    wallet_balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### 4.1.2 Transactions Table
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    recipient_id UUID REFERENCES users(id),
    business_id UUID REFERENCES businesses(id),
    amount DECIMAL(15,2) NOT NULL,
    fee DECIMAL(10,4) NOT NULL,
    type ENUM('p2p', 'cash_in', 'cash_out', 'business') NOT NULL,
    status ENUM('pending', 'completed', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### 4.1.3 Agents Table
```sql
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(15) UNIQUE NOT NULL,
    float_balance DECIMAL(15,2) DEFAULT 0.00,
    commission_balance DECIMAL(15,2) DEFAULT 0.00,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### 4.1.4 Businesses Table
```sql
CREATE TABLE businesses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    wallet_balance DECIMAL(15,2) DEFAULT 0.00,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### 4.2 API Endpoints

#### 4.2.1 Authentication
- `POST /auth/register` - User registration with phone number
- `POST /auth/login` - User login with OTP
- `POST /auth/verify-otp` - OTP verification
- `POST /auth/logout` - User logout

#### 4.2.2 Transactions
- `POST /transactions/p2p` - P2P transfer
- `POST /transactions/cash-in` - Cash-in via agent
- `POST /transactions/cash-out` - Cash-out via agent
- `POST /transactions/business` - Business payment
- `GET /transactions/history` - Transaction history
- `GET /transactions/{id}` - Transaction details

#### 4.2.3 User Management
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `POST /users/kyc` - Upload KYC documents
- `GET /users/balance` - Get wallet balance

#### 4.2.4 Agent Management
- `POST /agents/register` - Agent registration
- `GET /agents/transactions` - Agent transaction history
- `POST /agents/float-replenish` - Float replenishment
- `GET /agents/commission` - Commission balance

### 4.3 Flutter Implementation

#### 4.3.1 Main App Structure
```dart
// main.dart
import 'package:supabase_flutter/supabase_flutter.dart';

Future<void> main() async {
  await Supabase.initialize(
    url: 'https://your-project.supabase.co',
    anonKey: 'your-anon-key',
  );
  runApp(CashoutApp());
}

class CashoutApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Cashout',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: AuthWrapper(),
    );
  }
}
```

#### 4.3.2 Authentication Service
```dart
class AuthService {
  final SupabaseClient _supabase = Supabase.instance.client;

  Future<void> signInWithPhone(String phone) async {
    await _supabase.auth.signInWithOtp(
      phone: phone,
      channel: OtpChannel.sms,
    );
  }

  Future<void> verifyOtp(String phone, String token) async {
    await _supabase.auth.verifyOTP(
      phone: phone,
      token: token,
      type: OtpType.sms,
    );
  }
}
```

#### 4.3.3 Transaction Service
```dart
class TransactionService {
  final SupabaseClient _supabase = Supabase.instance.client;

  Future<void> sendMoney(String recipientId, double amount) async {
    final user = _supabase.auth.currentUser;
    final fee = calculateFee(amount);
    
    await _supabase.from('transactions').insert({
      'user_id': user!.id,
      'recipient_id': recipientId,
      'amount': amount,
      'fee': fee,
      'type': 'p2p',
    });
  }

  double calculateFee(double amount) {
    // Implementation based on fee structure tables
    if (amount >= 50 && amount <= 100) return 8.25;
    if (amount >= 101 && amount <= 500) return 21.75;
    // ... continue for all ranges
    return 0.0;
  }
}
```

### 4.4 Go Backend Implementation

#### 4.4.1 Main Server
```go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/supabase-community/supabase-go"
)

func main() {
    supabase := supabase.NewClient("https://your-project.supabase.co", "your-anon-key")
    
    r := mux.NewRouter()
    
    // Authentication routes
    r.HandleFunc("/auth/register", registerHandler).Methods("POST")
    r.HandleFunc("/auth/login", loginHandler).Methods("POST")
    
    // Transaction routes
    r.HandleFunc("/transactions/p2p", p2pHandler).Methods("POST")
    r.HandleFunc("/transactions/cash-in", cashInHandler).Methods("POST")
    r.HandleFunc("/transactions/cash-out", cashOutHandler).Methods("POST")
    
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

#### 4.4.2 Transaction Handler
```go
func p2pHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RecipientID string  `json:"recipient_id"`
        Amount      float64 `json:"amount"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    fee := calculateFee(req.Amount)
    
    // Insert transaction
    _, err := supabase.DB().Exec(`
        INSERT INTO transactions (user_id, recipient_id, amount, fee, type)
        VALUES ($1, $2, $3, $4, 'p2p')
    `, userID, req.RecipientID, req.Amount, fee)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Update balances
    updateBalances(userID, req.RecipientID, req.Amount, fee)
    
    w.WriteHeader(http.StatusOK)
}
```

---

## 5. Deployment and Infrastructure

### 5.1 Supabase Setup
1. **Project Creation:** Create new Supabase project
2. **Database Schema:** Execute SQL schema creation scripts
3. **Authentication:** Configure phone authentication with Africa's Talking
4. **Storage:** Set up KYC document storage with RLS policies
5. **Real-time:** Enable real-time subscriptions for transactions

### 5.2 Flutter App Deployment
1. **iOS:** Deploy to App Store via Xcode
2. **Android:** Deploy to Google Play Store
3. **Web:** Deploy Flutter web app to Supabase hosting or custom server

### 5.3 Go Backend Deployment
1. **Containerization:** Docker container for Go application
2. **Hosting:** Deploy to cloud platform (AWS, Google Cloud, or Azure)
3. **Load Balancing:** Set up load balancer for high availability
4. **Monitoring:** Implement logging and monitoring solutions

---

## 6. Testing Strategy

### 6.1 Unit Testing
- **Flutter:** Widget tests for UI components
- **Go:** Unit tests for business logic
- **Supabase:** Database function tests

### 6.2 Integration Testing
- **API Testing:** Test all endpoints with various scenarios
- **Database Testing:** Test transaction integrity and constraints
- **Authentication Testing:** Test OTP flow and security measures

### 6.3 End-to-End Testing
- **Mobile Testing:** Test complete user flows on iOS and Android
- **Web Testing:** Test web application functionality
- **Cross-Platform Testing:** Ensure consistency across platforms

---

## 7. Security and Compliance

### 7.1 CBK Compliance
- **KYC/AML:** Implement required customer identification procedures
- **Transaction Limits:** Enforce CBK-mandated transaction limits
- **Reporting:** Generate required regulatory reports
- **Audit Trail:** Maintain complete audit logs

### 7.2 Data Protection
- **Encryption:** All data encrypted in transit and at rest
- **Access Control:** Role-based access control for all users
- **Data Retention:** Implement data retention policies
- **Privacy:** Comply with Kenya's data protection laws

### 7.3 Fraud Prevention
- **Real-Time Monitoring:** Monitor transactions for suspicious activity
- **Machine Learning:** Use ML models for fraud detection
- **Manual Review:** Human review for flagged transactions
- **Blocking:** Ability to block suspicious accounts

---

## 8. Monitoring and Analytics

### 8.1 Performance Monitoring
- **Transaction Speed:** Monitor transaction processing times
- **System Uptime:** Track system availability
- **Error Rates:** Monitor application error rates
- **User Experience:** Track app performance metrics

### 8.2 Business Analytics
- **Transaction Volume:** Track daily/monthly transaction volumes
- **Revenue Tracking:** Monitor fee collection and revenue
- **User Growth:** Track user acquisition and retention
- **Agent Performance:** Monitor agent transaction volumes

### 8.3 Security Monitoring
- **Fraud Detection:** Monitor fraud detection system performance
- **Security Events:** Track security-related events
- **Compliance Monitoring:** Ensure regulatory compliance
- **Audit Logs:** Monitor system access and changes

---

## 9. Support and Maintenance

### 9.1 Customer Support
- **Help Desk:** 24/7 customer support system
- **Documentation:** Comprehensive user and agent documentation
- **Training:** Agent training programs
- **Feedback System:** User feedback collection and processing

### 9.2 System Maintenance
- **Regular Updates:** Scheduled system updates and patches
- **Backup Procedures:** Regular data backup and recovery testing
- **Performance Optimization:** Continuous performance improvements
- **Security Updates:** Regular security patches and updates

---

## 10. Success Metrics

### 10.1 User Adoption
- **Target:** 1M users within 18 months
- **Metric:** Monthly active users (MAU)
- **Goal:** 60% user retention rate

### 10.2 Transaction Volume
- **Target:** 10M transactions per month
- **Metric:** Daily transaction volume
- **Goal:** 25% month-over-month growth

### 10.3 Revenue Generation
- **Target:** KES 100M monthly revenue
- **Metric:** Monthly recurring revenue (MRR)
- **Goal:** 30% profit margin

### 10.4 Agent Network
- **Target:** 10,000 active agents
- **Metric:** Active agent count
- **Goal:** 80% agent retention rate

---

This comprehensive PRD provides all the necessary details for implementing the Cashout platform using Cursor. The document covers all technical aspects, business requirements, security measures, and implementation details discussed in the conversation with Grok. 