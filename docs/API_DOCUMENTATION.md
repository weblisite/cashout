# Cashout API Documentation

## Overview

The Cashout API is a RESTful service built with Go and Gin framework that provides mobile payment functionality including P2P transfers, cash-in/cash-out services, and business payments.

**Base URL**: `http://localhost:8080/api/v1`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### Send OTP
```http
POST /auth/phone/send-otp
```

**Request Body:**
```json
{
  "phone_number": "+254700000000"
}
```

**Response:**
```json
{
  "message": "OTP sent successfully",
  "phone": "+254700000000"
}
```

#### Verify OTP
```http
POST /auth/phone/verify-otp
```

**Request Body:**
```json
{
  "phone_number": "+254700000000",
  "otp": "123456"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "phone_number": "+254700000000",
    "hashed_id": "abc123def456",
    "kyc_status": "not_started",
    "wallet_balance": 0.00,
    "user_status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "is_new_user": true,
  "requires_pin": true,
  "requires_kyc": true
}
```

#### Setup PIN
```http
POST /auth/pin/setup
```

**Request Body:**
```json
{
  "phone_number": "+254700000000",
  "pin": "1234"
}
```

**Response:**
```json
{
  "message": "PIN setup successfully"
}
```

#### Verify PIN
```http
POST /auth/pin/verify
```

**Request Body:**
```json
{
  "phone_number": "+254700000000",
  "pin": "1234"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "phone_number": "+254700000000",
    "hashed_id": "abc123def456",
    "kyc_status": "not_started",
    "wallet_balance": 0.00,
    "user_status": "active",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "is_new_user": false,
  "requires_pin": false,
  "requires_kyc": true
}
```

#### Setup Biometric
```http
POST /auth/biometric/setup
```

**Request Body:**
```json
{
  "phone_number": "+254700000000",
  "biometric_id": "biometric-fingerprint-id"
}
```

**Response:**
```json
{
  "message": "Biometric setup successfully"
}
```

#### Logout
```http
POST /auth/logout
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

### User Management

#### Get Profile
```http
GET /users/profile
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "id": "uuid",
  "phone_number": "+254700000000",
  "hashed_id": "abc123def456",
  "kyc_status": "verified",
  "wallet_balance": 1000.00,
  "user_status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Update Profile
```http
PUT /users/profile
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "kyc_status": "verified"
}
```

**Response:**
```json
{
  "message": "Profile updated successfully"
}
```

#### Get Balance
```http
GET /users/balance
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "balance": 1000.00,
  "currency": "KES"
}
```

### Transactions

#### Send P2P Transfer
```http
POST /transactions/p2p
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "recipient_id": "recipient-hashed-id",
  "amount": 100.00,
  "note": "Payment for lunch"
}
```

**Response:**
```json
{
  "transaction_id": "uuid",
  "status": "completed",
  "amount": 100.00,
  "fee": 0.50,
  "total": 100.50,
  "recipient": {
    "hashed_id": "recipient-hashed-id",
    "phone_number": "+254700000001"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Process Cash In
```http
POST /transactions/cash-in
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "agent_id": "agent-code",
  "amount": 500.00
}
```

**Response:**
```json
{
  "transaction_id": "uuid",
  "status": "pending",
  "amount": 500.00,
  "fee": 0.00,
  "agent": {
    "agent_code": "agent-code",
    "location": "Nairobi CBD"
  },
  "qr_code": "cashout://cashin?userId=...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Process Cash Out
```http
POST /transactions/cash-out
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "agent_id": "agent-code",
  "amount": 200.00
}
```

**Response:**
```json
{
  "transaction_id": "uuid",
  "status": "pending",
  "amount": 200.00,
  "fee": 4.00,
  "total": 204.00,
  "agent": {
    "agent_code": "agent-code",
    "location": "Nairobi CBD"
  },
  "qr_code": "cashout://cashout?userId=...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Process Business Payment
```http
POST /transactions/business
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "business_id": "business-id",
  "amount": 150.00,
  "note": "Payment for services"
}
```

**Response:**
```json
{
  "transaction_id": "uuid",
  "status": "completed",
  "amount": 150.00,
  "fee": 2.25,
  "total": 152.25,
  "business": {
    "business_id": "business-id",
    "business_name": "Demo Business"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### Get Transaction History
```http
GET /transactions/history?page=1&page_size=10&type=p2p&status=completed
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)
- `type` (optional): Transaction type (p2p, cash_in, cash_out, business)
- `status` (optional): Transaction status (pending, completed, failed)
- `from_date` (optional): Start date (YYYY-MM-DD)
- `to_date` (optional): End date (YYYY-MM-DD)

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "recipient_id": "uuid",
      "amount": 100.00,
      "fee": 0.50,
      "type": "p2p",
      "status": "completed",
      "note": "Payment for lunch",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

#### Get Transaction Details
```http
GET /transactions/{transaction_id}
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "recipient_id": "uuid",
  "amount": 100.00,
  "fee": 0.50,
  "agent_commission": null,
  "platform_margin": null,
  "type": "p2p",
  "status": "completed",
  "note": "Payment for lunch",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Fee Calculation

#### Calculate Fee
```http
GET /fees/calculate?amount=100&type=p2p
```

**Query Parameters:**
- `amount` (required): Transaction amount
- `type` (required): Transaction type (p2p, cash_in, cash_out, business)

**Response:**
```json
{
  "amount": 100.00,
  "fee": 0.50,
  "total": 100.50,
  "fee_percentage": 0.5,
  "type": "p2p"
}
```

#### Get Fee Structure
```http
GET /fees/structure
```

**Response:**
```json
{
  "fee_structure": [
    {
      "min_amount": 50.0,
      "max_amount": 100.0,
      "p2p_fee": 8.0,
      "cash_out_fee": 8.0,
      "business_fee": 8.0
    },
    {
      "min_amount": 101.0,
      "max_amount": 500.0,
      "p2p_fee": 22.0,
      "cash_out_fee": 22.0,
      "business_fee": 22.0
    },
    {
      "min_amount": 501.0,
      "max_amount": 1000.0,
      "p2p_fee": 22.0,
      "cash_out_fee": 22.0,
      "business_fee": 22.0
    },
    {
      "min_amount": 1001.0,
      "max_amount": 1500.0,
      "p2p_fee": 22.0,
      "cash_out_fee": 22.0,
      "business_fee": 22.0
    },
    {
      "min_amount": 1501.0,
      "max_amount": 2500.0,
      "p2p_fee": 22.0,
      "cash_out_fee": 22.0,
      "business_fee": 22.0
    },
    {
      "min_amount": 2501.0,
      "max_amount": 3500.0,
      "p2p_fee": 39.0,
      "cash_out_fee": 39.0,
      "business_fee": 39.0
    },
    {
      "min_amount": 3501.0,
      "max_amount": 5000.0,
      "p2p_fee": 52.0,
      "cash_out_fee": 52.0,
      "business_fee": 52.0
    },
    {
      "min_amount": 5001.0,
      "max_amount": 7500.0,
      "p2p_fee": 65.0,
      "cash_out_fee": 65.0,
      "business_fee": 65.0
    },
    {
      "min_amount": 7501.0,
      "max_amount": 10000.0,
      "p2p_fee": 86.0,
      "cash_out_fee": 86.0,
      "business_fee": 86.0
    },
    {
      "min_amount": 10001.0,
      "max_amount": 15000.0,
      "p2p_fee": 125.0,
      "cash_out_fee": 125.0,
      "business_fee": 125.0
    },
    {
      "min_amount": 15001.0,
      "max_amount": 20000.0,
      "p2p_fee": 139.0,
      "cash_out_fee": 139.0,
      "business_fee": 139.0
    },
    {
      "min_amount": 20001.0,
      "max_amount": 35000.0,
      "p2p_fee": 148.0,
      "cash_out_fee": 148.0,
      "business_fee": 148.0
    },
    {
      "min_amount": 35001.0,
      "max_amount": 50000.0,
      "p2p_fee": 209.0,
      "cash_out_fee": 209.0,
      "business_fee": 209.0
    },
    {
      "min_amount": 50001.0,
      "max_amount": 250000.0,
      "p2p_fee": 232.0,
      "cash_out_fee": 232.0,
      "business_fee": 232.0
    },
    {
      "min_amount": 250001.0,
      "max_amount": 500000.0,
      "p2p_fee": 513.0,
      "cash_out_fee": 513.0,
      "business_fee": 513.0
    },
    {
      "min_amount": 500001.0,
      "max_amount": 1000000.0,
      "p2p_fee": 1076.0,
      "cash_out_fee": 1076.0,
      "business_fee": 1076.0
    }
  ],
  "notes": {
    "p2p": "Platform retains 100% of fee",
    "cash_out": "25% agent commission, 75% platform margin",
    "business": "50/50 split between user and business",
    "cash_in": "Always FREE",
    "rounding": "Fees rounded to nearest whole number (≥0.50 rounds up, <0.50 rounds down)"
  }
}
```

### Agent Management

#### Get Agent Profile
```http
GET /agents/profile
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "agent_code": "AGENT001",
  "float_balance": 5000.00,
  "commission_rate": 0.0200,
  "status": "active",
  "location": "Nairobi CBD",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Get Agent Transactions
```http
GET /agents/transactions?page=1&page_size=10
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "amount": 100.00,
      "fee": 0.00,
      "agent_commission": 2.00,
      "type": "cash_in",
      "status": "completed",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

#### Get Agent Float
```http
GET /agents/float
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "float_balance": 5000.00,
  "currency": "KES"
}
```

#### Update Agent Float
```http
POST /agents/float/update
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "new_balance": 6000.00
}
```

**Response:**
```json
{
  "message": "Float updated successfully",
  "new_balance": 6000.00
}
```

### Business Management

#### Get Business Profile
```http
GET /businesses/profile
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "business_name": "Demo Business",
  "business_type": "Retail",
  "address": "Nairobi, Kenya",
  "phone_number": "+254700000000",
  "email": "business@example.com",
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Get Business Transactions
```http
GET /businesses/transactions?page=1&page_size=10
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "amount": 150.00,
      "fee": 2.25,
      "type": "business",
      "status": "completed",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

#### Generate Business QR
```http
POST /businesses/qr/generate
```

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Request Body:**
```json
{
  "amount": 100.00,
  "note": "Payment for services"
}
```

**Response:**
```json
{
  "qr_code": "cashout://business?businessId=...&amount=100.00",
  "amount": 100.00,
  "business_id": "uuid",
  "expires_at": "2024-01-01T01:00:00Z"
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error type",
  "message": "Human-readable error message",
  "details": "Additional error details (optional)"
}
```

### Common HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `422 Unprocessable Entity`: Validation errors
- `500 Internal Server Error`: Server error

### Validation Errors

```json
{
  "error": "Validation failed",
  "message": "Please check your input",
  "details": {
    "phone_number": ["Phone number is required"],
    "amount": ["Amount must be greater than 0"]
  }
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Authentication endpoints**: 5 requests per minute
- **Transaction endpoints**: 10 requests per minute
- **Other endpoints**: 100 requests per minute

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Webhooks

The API supports webhooks for real-time notifications:

### Webhook Events

- `transaction.completed`: Transaction completed successfully
- `transaction.failed`: Transaction failed
- `user.kyc_verified`: User KYC verified
- `agent.float_updated`: Agent float balance updated

### Webhook Payload

```json
{
  "event": "transaction.completed",
  "timestamp": "2024-01-01T00:00:00Z",
  "data": {
    "transaction_id": "uuid",
    "amount": 100.00,
    "type": "p2p",
    "status": "completed"
  }
}
```

## SDKs and Libraries

### Go Client

```go
package main

import (
    "github.com/cashout/go-sdk"
)

func main() {
    client := cashout.NewClient("your-api-key", "https://api.cashout.com")
    
    // Send P2P transfer
    resp, err := client.Transactions.SendP2P(cashout.P2PRequest{
        RecipientID: "recipient-id",
        Amount:      100.00,
        Note:        "Payment",
    })
}
```

### JavaScript/TypeScript Client

```javascript
import { CashoutClient } from '@cashout/js-sdk';

const client = new CashoutClient('your-api-key', 'https://api.cashout.com');

// Send P2P transfer
const response = await client.transactions.sendP2P({
    recipientId: 'recipient-id',
    amount: 100.00,
    note: 'Payment'
});
```

## Support

For API support and questions:

- **Email**: api-support@cashout.com
- **Documentation**: https://docs.cashout.com
- **Status Page**: https://status.cashout.com 