# 🚀 Cashout Platform

A comprehensive digital payment platform with its own banking infrastructure, designed for real money movement through an agent-based cash-in/cash-out model.

## 🌟 **Platform Overview**

Cashout is a complete payment platform that handles the entire payment infrastructure without relying on external payment gateways like M-Pesa or Airtel Money. Users can cash in and cash out through authorized agents, with real money movement and comprehensive transaction management.

## 🏗️ **Architecture**

### **Core Components**
- **Frontend Web Apps**: User, Agent, and Business interfaces
- **Backend API**: Flask-based REST API for transaction processing
- **Database**: Supabase for data storage and user management
- **Authentication**: Supabase Auth for secure user authentication
- **Banking Infrastructure**: Custom settlement engine and liquidity management

### **Key Features**
- ✅ **Real Money Movement**: Complete banking infrastructure
- ✅ **Agent Network**: Cash-in/cash-out through authorized agents
- ✅ **QR Code Payments**: Scan-to-pay functionality
- ✅ **Transaction Management**: Comprehensive transaction history
- ✅ **Fee Structure**: Transparent fee calculation
- ✅ **Commission System**: Agent commission tracking
- ✅ **KYC/AML**: Identity verification and compliance
- ✅ **Real-time Updates**: Live balance and transaction updates

## 📱 **Web Applications**

### **1. User App** (`landing-page/user-app.html`)
- **Features**:
  - Send money to other users
  - QR code generation and scanning
  - Cash-in/cash-out requests
  - Transaction history
  - Real-time balance updates
  - Profile management

### **2. Agent App** (`landing-page/agent-app.html`)
- **Features**:
  - Process cash-in transactions
  - Process cash-out transactions
  - QR code scanning
  - Float balance management
  - Commission tracking
  - Transaction history
  - Daily reports

### **3. Business App** (`landing-page/business-app.html`)
- **Features**:
  - Receive payments from customers
  - QR code generation
  - Transaction history
  - Business analytics
  - Settlement management

## 🔧 **Backend API**

### **API Endpoints**

#### **Health & Status**
- `GET /api/health` - API health check

#### **User Management**
- `GET /api/users/{user_id}` - Get user profile
- `GET /api/users/{user_id}/balance` - Get user balance
- `GET /api/users/{user_id}/transactions` - Get user transactions

#### **Transactions**
- `POST /api/transactions/send` - Send money between users
- `POST /api/transactions/cash-in` - Process cash-in transaction
- `POST /api/transactions/cash-out` - Process cash-out transaction

#### **QR Code Management**
- `POST /api/qr/generate` - Generate QR code for payment
- `POST /api/qr/scan` - Scan and process QR code

#### **Agent Management**
- `GET /api/agents/{agent_id}` - Get agent profile

#### **Business Management**
- `GET /api/businesses/{business_id}` - Get business profile

#### **Fee Calculation**
- `POST /api/fees/calculate` - Calculate transaction fees

### **Fee Structure**
```
Amount Range (KES) | Fee (KES)
-------------------|----------
1 - 100           | 8
101 - 500         | 22
501 - 1,000       | 22
1,001 - 1,500     | 22
1,501 - 2,500     | 22
2,501 - 3,500     | 39
3,501 - 5,000     | 52
5,001 - 7,500     | 65
7,501 - 10,000    | 86
10,001 - 15,000   | 125
15,001 - 20,000   | 139
20,001 - 35,000   | 148
35,001 - 50,000   | 209
50,001 - 250,000  | 232
250,001 - 500,000 | 513
500,001+          | 1,076
```

## 🚀 **Getting Started**

### **Prerequisites**
- Python 3.7+
- Modern web browser
- Supabase account (for production)

### **1. Clone the Repository**
```bash
git clone <repository-url>
cd Cashout
```

### **2. Start the Backend API**
```bash
cd backend
chmod +x start_api.sh
./start_api.sh
```

The API will be available at: `http://localhost:5000`

### **3. Start the Web Apps**
```bash
cd landing-page
python3 -m http.server 8001
```

The web apps will be available at: `http://localhost:8001`

### **4. Access the Applications**
- **User App**: http://localhost:8001/user-app.html
- **Agent App**: http://localhost:8001/agent-app.html
- **Business App**: http://localhost:8001/business-app.html

## 🧪 **Testing the Platform**

### **Demo Mode**
The applications work in demo mode when the API is not available, using sample data for testing.

### **API Mode**
When the backend API is running, the applications connect to real endpoints for:
- User authentication
- Transaction processing
- Balance updates
- QR code generation

### **Test Scenarios**

#### **User App Testing**
1. **Send Money**: Enter recipient phone and amount
2. **QR Code**: Generate and scan QR codes
3. **Cash In/Out**: Request agent services
4. **Transaction History**: View past transactions

#### **Agent App Testing**
1. **Cash In**: Process customer deposits
2. **Cash Out**: Process customer withdrawals
3. **QR Scanning**: Scan customer QR codes
4. **Commission Tracking**: Monitor earnings

#### **Business App Testing**
1. **Receive Payments**: Accept customer payments
2. **QR Generation**: Create payment QR codes
3. **Settlement**: Manage business funds

## 🔒 **Security Features**

### **Authentication & Authorization**
- Supabase Auth integration
- JWT token management
- Role-based access control

### **Data Protection**
- Encrypted data transmission
- Secure API endpoints
- Input validation and sanitization

### **Compliance**
- KYC/AML integration
- Transaction monitoring
- Regulatory compliance features

## 📊 **Business Model**

### **Revenue Streams**
1. **Transaction Fees**: Per-transaction charges
2. **Agent Commissions**: 25% of fees to agents
3. **Business Services**: Premium business features
4. **Settlement Services**: Inter-bank transfers

### **Agent Network**
- **Float Management**: Agents maintain cash float
- **Commission Structure**: 25% of transaction fees
- **Training & Support**: Comprehensive agent support
- **Geographic Coverage**: Nationwide agent network

## 🔮 **Future Enhancements**

### **Phase 2: Advanced Features**
- **Settlement Engine**: Automated inter-bank settlements
- **Liquidity Management**: Dynamic liquidity allocation
- **Risk Management**: Fraud detection and prevention

### **Phase 3: Compliance & Security**
- **KYC/AML**: Advanced identity verification
- **Fraud Detection**: AI-powered fraud prevention
- **Regulatory Compliance**: Full regulatory adherence

### **Phase 4: Scale & Performance**
- **High Availability**: 99.9% uptime guarantee
- **Performance Optimization**: Sub-second response times
- **Analytics Dashboard**: Comprehensive business intelligence

## 🛠️ **Development**

### **Project Structure**
```
Cashout/
├── landing-page/          # Web applications
│   ├── user-app.html     # User interface
│   ├── agent-app.html    # Agent interface
│   ├── business-app.html # Business interface
│   ├── index.html        # Landing page
│   └── styles.css        # Shared styles
├── backend/              # Backend API
│   ├── simple_api.py     # Flask API
│   ├── requirements.txt  # Python dependencies
│   └── start_api.sh      # Startup script
├── .cursor/              # Cursor IDE configuration
│   └── mcp.json         # MCP server configuration
└── README.md            # This file
```

### **API Development**
The backend API is built with Flask and includes:
- RESTful endpoints
- JSON response format
- Error handling
- CORS support
- Sample data for testing

### **Frontend Development**
The web apps are built with:
- Vanilla JavaScript
- HTML5/CSS3
- Responsive design
- Progressive Web App features

## 📞 **Support**

### **Documentation**
- API documentation available at `/api/health`
- Code comments for all major functions
- Sample data for testing

### **Troubleshooting**
1. **API Not Starting**: Check Python installation and dependencies
2. **Web Apps Not Loading**: Ensure HTTP server is running
3. **CORS Issues**: Verify API CORS configuration
4. **Demo Mode**: Check API connectivity

## 📄 **License**

This project is proprietary software. All rights reserved.

---

**Built with ❤️ for the Cashout Platform** 