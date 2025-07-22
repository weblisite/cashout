# Cashout Mobile App Guide

## Overview

The Cashout mobile app is a Flutter-based application that provides secure mobile payment services including P2P transfers, cash-in/cash-out, and business payments. The app is designed for users in Kenya and supports both English and Swahili languages.

## Features

### 🔐 Authentication & Security
- **Phone Number Verification**: OTP-based phone verification
- **PIN Setup**: Secure 4-digit PIN for transactions
- **Biometric Authentication**: Fingerprint and face recognition support
- **JWT Token Management**: Secure session management

### 💰 Wallet Management
- **Real-time Balance**: Live wallet balance display
- **Transaction History**: Complete transaction records with filtering
- **Balance Updates**: Instant balance updates after transactions

### 📱 Core Features
- **P2P Transfers**: Send money to other users instantly
- **QR Code Support**: Scan and generate QR codes for payments
- **Cash-in/Cash-out**: Agent-based deposit and withdrawal
- **Business Payments**: Pay businesses using QR codes
- **Fee Calculation**: Transparent fee structure display

### 🏪 Agent & Business Support
- **Agent Locator**: Find nearby Cashout agents
- **Business Directory**: Discover businesses accepting Cashout
- **Float Management**: Agent float balance tracking

## Installation

### Prerequisites
- Flutter SDK 3.16.0 or higher
- Dart SDK 3.2.0 or higher
- Android Studio / VS Code
- Android SDK (for Android development)
- Xcode (for iOS development, macOS only)

### Setup Instructions

1. **Clone the repository**
   ```bash
   git clone https://github.com/cashout/mobile-app.git
   cd mobile-app
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the app**
   ```bash
   flutter run
   ```

### Environment Configuration

Create a `.env` file in the mobile app root:

```env
# API Configuration
API_BASE_URL=http://localhost:8080/api/v1
API_TIMEOUT=30000

# Supabase Configuration
SUPABASE_URL=your-supabase-url
SUPABASE_ANON_KEY=your-supabase-anon-key

# Intasend Configuration
INTASEND_API_KEY=your-intasend-api-key
INTASEND_PUBLISHABLE_KEY=your-intasend-publishable-key

# App Configuration
APP_NAME=Cashout
APP_VERSION=1.0.0
ENVIRONMENT=development
```

## User Guide

### Getting Started

#### 1. App Installation
1. Download Cashout from Google Play Store or App Store
2. Open the app and tap "Get Started"
3. Accept terms and conditions

#### 2. Phone Verification
1. Enter your phone number (Kenyan format: +254...)
2. Tap "Send OTP"
3. Enter the 6-digit OTP received via SMS
4. Tap "Verify"

#### 3. PIN Setup
1. Create a 4-digit PIN
2. Confirm the PIN
3. Tap "Set PIN"

#### 4. Biometric Setup (Optional)
1. Tap "Setup Biometric"
2. Follow device prompts for fingerprint/face setup
3. Tap "Done"

### Using the App

#### Home Screen
- **Wallet Balance**: Displays current balance prominently
- **Quick Actions**: Send, Receive, Cash-in, Cash-out buttons
- **Recent Transactions**: Last 5 transactions
- **Statistics**: Monthly transaction summary

#### Send Money
1. Tap "Send" on home screen
2. Choose recipient method:
   - **Phone Number**: Enter recipient's phone
   - **User ID**: Enter recipient's hashed ID
   - **QR Code**: Scan recipient's QR code
3. Enter amount
4. Add optional note
5. Review transaction details
6. Enter PIN to confirm
7. Transaction completed

#### Receive Money
1. Tap "Receive" on home screen
2. Choose amount option:
   - **Any Amount**: Generate QR for any amount
   - **Specific Amount**: Set fixed amount
3. Share QR code or payment details
4. Wait for sender to complete transaction

#### Cash In (Deposit)
1. Tap "Cash In" on home screen
2. Enter amount to deposit
3. Choose agent:
   - **Scan Agent QR**: Scan agent's QR code
   - **Enter Agent ID**: Manually enter agent code
4. Show your QR code to agent
5. Agent processes transaction
6. Balance updated instantly

#### Cash Out (Withdraw)
1. Tap "Cash Out" on home screen
2. Enter amount to withdraw
3. Choose agent (scan QR or enter ID)
4. Show QR code to agent
5. Agent provides cash
6. Transaction completed

#### Business Payments
1. Tap "Business" on home screen
2. Choose business:
   - **Scan Business QR**: Scan business QR code
   - **Enter Business ID**: Manually enter business ID
3. Enter payment amount
4. Add optional note
5. Confirm payment with PIN
6. Payment processed instantly

#### Transaction History
1. Tap "History" in bottom navigation
2. View all transactions with filters:
   - **All**: All transaction types
   - **P2P**: Person-to-person transfers
   - **Cash In**: Deposit transactions
   - **Cash Out**: Withdrawal transactions
   - **Business**: Business payments
3. Tap transaction for details
4. Use search and date filters

#### Profile Management
1. Tap "Profile" in bottom navigation
2. View account information:
   - Phone number
   - User ID
   - KYC status
   - Account status
3. Access settings:
   - **Notifications**: Configure push notifications
   - **Language**: Switch between English/Swahili
   - **Security**: Change PIN, biometric settings
   - **Support**: Help center, contact support

## Security Features

### PIN Protection
- 4-digit PIN required for all transactions
- PIN is hashed and stored securely
- Failed attempts are limited
- PIN can be changed in settings

### Biometric Authentication
- Fingerprint recognition
- Face recognition (iOS)
- Optional setup for convenience
- Can be disabled in settings

### Transaction Security
- All transactions require PIN confirmation
- Real-time fraud detection
- Transaction limits and monitoring
- Secure QR code generation

### Data Protection
- End-to-end encryption
- Secure API communication
- Local data encryption
- GDPR compliance

## Troubleshooting

### Common Issues

#### OTP Not Received
1. Check phone number format (+254...)
2. Check SMS permissions
3. Wait 60 seconds before resending
4. Contact support if issue persists

#### Transaction Failed
1. Check internet connection
2. Verify recipient details
3. Ensure sufficient balance
4. Check transaction limits
5. Contact support with transaction ID

#### App Crashes
1. Restart the app
2. Clear app cache
3. Update to latest version
4. Reinstall if necessary

#### QR Code Issues
1. Ensure good lighting
2. Hold device steady
3. Clean camera lens
4. Try manual entry instead

### Error Messages

| Error | Solution |
|-------|----------|
| "Insufficient Balance" | Add money to wallet via cash-in |
| "Invalid PIN" | Check PIN and try again |
| "Network Error" | Check internet connection |
| "Transaction Limit Exceeded" | Wait or contact support |
| "Invalid Recipient" | Verify recipient details |

## Support

### In-App Support
1. Go to Profile → Support
2. Choose support option:
   - **Help Center**: FAQ and guides
   - **Contact Support**: Live chat/email
   - **Report Problem**: Bug reporting

### Contact Information
- **Email**: support@cashout.com
- **Phone**: +254 700 000 000
- **WhatsApp**: +254 700 000 000
- **Hours**: 24/7 support

### Social Media
- **Twitter**: @CashoutApp
- **Facebook**: /CashoutApp
- **Instagram**: @CashoutApp

## Privacy & Terms

### Privacy Policy
- Personal data protection
- Transaction data security
- Third-party sharing policies
- User rights and controls

### Terms of Service
- Account usage rules
- Transaction limits
- Dispute resolution
- Liability limitations

### Data Usage
- Transaction history retention
- Analytics and improvements
- Marketing communications
- Legal compliance

## Updates & Maintenance

### App Updates
- Automatic updates enabled
- Manual updates available
- Release notes provided
- Backward compatibility maintained

### Maintenance Windows
- Scheduled maintenance: 2-4 AM EAT
- Emergency maintenance: As needed
- Advance notification provided
- Minimal service disruption

### Feature Roadmap
- **Q1 2024**: International transfers
- **Q2 2024**: Bill payments
- **Q3 2024**: Savings features
- **Q4 2024**: Investment options

## Technical Specifications

### Supported Platforms
- **Android**: 6.0 (API 23) and higher
- **iOS**: 12.0 and higher
- **Web**: Chrome, Safari, Firefox

### Device Requirements
- **RAM**: 2GB minimum, 4GB recommended
- **Storage**: 100MB free space
- **Camera**: For QR code scanning
- **Biometric**: For biometric authentication

### Network Requirements
- **Internet**: 3G/4G/5G or WiFi
- **Bandwidth**: 1MB/s minimum
- **Latency**: <500ms recommended

## Development

### Architecture
- **Framework**: Flutter 3.16.0
- **State Management**: Provider
- **Architecture**: MVVM pattern
- **Testing**: Unit, widget, integration tests

### Key Dependencies
```yaml
dependencies:
  flutter: ^3.16.0
  provider: ^6.0.5
  http: ^1.1.0
  shared_preferences: ^2.2.2
  qr_flutter: ^4.1.0
  qr_code_scanner: ^1.0.1
  local_auth: ^2.1.6
  supabase_flutter: ^1.10.25
```

### Build Commands
```bash
# Debug build
flutter run

# Release build (Android)
flutter build apk --release

# Release build (iOS)
flutter build ios --release

# Web build
flutter build web --release
```

### Testing
```bash
# Unit tests
flutter test

# Widget tests
flutter test test/widget_test.dart

# Integration tests
flutter test integration_test/
```

## Contributing

### Development Setup
1. Fork the repository
2. Create feature branch
3. Make changes
4. Add tests
5. Submit pull request

### Code Standards
- Follow Flutter style guide
- Use meaningful variable names
- Add comments for complex logic
- Maintain test coverage >80%

### Review Process
1. Automated tests must pass
2. Code review required
3. Documentation updated
4. Security review for sensitive changes

---

**Version**: 1.0.0  
**Last Updated**: January 2024  
**Support**: support@cashout.com 