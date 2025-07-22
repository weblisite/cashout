#!/usr/bin/env python3
"""
Simple Cashout API - Flask-based backend for web apps
No complex dependencies, just basic functionality
"""

from flask import Flask, request, jsonify
from flask_cors import CORS
import json
import uuid
from datetime import datetime, timedelta
import os

app = Flask(__name__)
CORS(app)  # Enable CORS for web apps

# In-memory data storage (in production, use a real database)
users_db = {}
agents_db = {}
businesses_db = {}
transactions_db = {}
qr_codes_db = {}

# Sample data initialization
def init_sample_data():
    # Sample users
    users_db['user1'] = {
        'id': 'user1',
        'name': 'John Doe',
        'phone': '+254700123456',
        'email': 'john.doe@example.com',
        'balance': 25000.0,
        'kyc_status': 'verified',
        'created_at': datetime.now().isoformat()
    }
    
    users_db['user2'] = {
        'id': 'user2',
        'name': 'Jane Smith',
        'phone': '+254700789012',
        'email': 'jane.smith@example.com',
        'balance': 15000.0,
        'kyc_status': 'verified',
        'created_at': datetime.now().isoformat()
    }
    
    # Sample agents
    agents_db['agent1'] = {
        'id': 'agent1',
        'name': 'Agent John',
        'phone': '+254700345678',
        'location': 'Nairobi CBD',
        'float_balance': 50000.0,
        'commission_balance': 2500.0,
        'status': 'active',
        'created_at': datetime.now().isoformat()
    }
    
    # Sample businesses
    businesses_db['business1'] = {
        'id': 'business1',
        'name': 'ABC Electronics Store',
        'phone': '+254700567890',
        'location': 'Westlands, Nairobi',
        'balance': 125000.0,
        'status': 'active',
        'created_at': datetime.now().isoformat()
    }
    
    # Sample transactions
    transactions_db['tx1'] = {
        'id': 'tx1',
        'type': 'p2p',
        'from_user': 'user1',
        'to_user': 'user2',
        'amount': 1500.0,
        'fee': 22.0,
        'total': 1522.0,
        'status': 'completed',
        'timestamp': datetime.now().isoformat()
    }

# Fee calculation function
def calculate_fee(amount, transaction_type):
    """Calculate transaction fee based on amount and type"""
    if amount <= 100:
        return 8
    elif amount <= 500:
        return 22
    elif amount <= 1000:
        return 22
    elif amount <= 1500:
        return 22
    elif amount <= 2500:
        return 22
    elif amount <= 3500:
        return 39
    elif amount <= 5000:
        return 52
    elif amount <= 7500:
        return 65
    elif amount <= 10000:
        return 86
    elif amount <= 15000:
        return 125
    elif amount <= 20000:
        return 139
    elif amount <= 35000:
        return 148
    elif amount <= 50000:
        return 209
    elif amount <= 250000:
        return 232
    elif amount <= 500000:
        return 513
    else:
        return 1076

# API Routes

@app.route('/api/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        'status': 'healthy',
        'timestamp': datetime.now().isoformat(),
        'service': 'Cashout API'
    })

@app.route('/api/users/<user_id>', methods=['GET'])
def get_user(user_id):
    """Get user profile"""
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    
    return jsonify(users_db[user_id])

@app.route('/api/users/<user_id>/balance', methods=['GET'])
def get_user_balance(user_id):
    """Get user balance"""
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    
    return jsonify({
        'user_id': user_id,
        'balance': users_db[user_id]['balance'],
        'currency': 'KES',
        'last_updated': datetime.now().isoformat()
    })

@app.route('/api/users/<user_id>/transactions', methods=['GET'])
def get_user_transactions(user_id):
    """Get user transaction history"""
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    
    user_transactions = []
    for tx_id, tx in transactions_db.items():
        if tx['from_user'] == user_id or tx['to_user'] == user_id:
            user_transactions.append(tx)
    
    # Sort by timestamp (newest first)
    user_transactions.sort(key=lambda x: x['timestamp'], reverse=True)
    
    return jsonify({
        'user_id': user_id,
        'transactions': user_transactions,
        'count': len(user_transactions)
    })

@app.route('/api/transactions/send', methods=['POST'])
def send_money():
    """Send money between users"""
    data = request.get_json()
    
    required_fields = ['from_user', 'to_user', 'amount', 'description']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    from_user = data['from_user']
    to_user = data['to_user']
    amount = float(data['amount'])
    description = data['description']
    
    # Validate users exist
    if from_user not in users_db:
        return jsonify({'error': 'Sender not found'}), 404
    if to_user not in users_db:
        return jsonify({'error': 'Recipient not found'}), 404
    
    # Validate amount
    if amount < 50 or amount > 1000000:
        return jsonify({'error': 'Invalid amount. Must be between 50 and 1,000,000 KES'}), 400
    
    # Check balance
    sender_balance = users_db[from_user]['balance']
    fee = calculate_fee(amount, 'p2p')
    total = amount + fee
    
    if sender_balance < total:
        return jsonify({'error': 'Insufficient balance'}), 400
    
    # Create transaction
    tx_id = str(uuid.uuid4())
    transaction = {
        'id': tx_id,
        'type': 'p2p',
        'from_user': from_user,
        'to_user': to_user,
        'amount': amount,
        'fee': fee,
        'total': total,
        'description': description,
        'status': 'completed',
        'timestamp': datetime.now().isoformat()
    }
    
    # Update balances
    users_db[from_user]['balance'] -= total
    users_db[to_user]['balance'] += amount
    
    # Store transaction
    transactions_db[tx_id] = transaction
    
    return jsonify({
        'success': True,
        'transaction': transaction,
        'message': 'Money sent successfully'
    })

@app.route('/api/transactions/cash-in', methods=['POST'])
def cash_in():
    """Cash in transaction (agent to user)"""
    data = request.get_json()
    
    required_fields = ['user_id', 'agent_id', 'amount']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    user_id = data['user_id']
    agent_id = data['agent_id']
    amount = float(data['amount'])
    
    # Validate user and agent exist
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    if agent_id not in agents_db:
        return jsonify({'error': 'Agent not found'}), 404
    
    # Validate amount
    if amount < 50 or amount > 1000000:
        return jsonify({'error': 'Invalid amount'}), 400
    
    # Check agent float
    agent_float = agents_db[agent_id]['float_balance']
    if agent_float < amount:
        return jsonify({'error': 'Agent has insufficient float'}), 400
    
    # Create transaction
    tx_id = str(uuid.uuid4())
    transaction = {
        'id': tx_id,
        'type': 'cash_in',
        'from_user': agent_id,
        'to_user': user_id,
        'amount': amount,
        'fee': 0.0,  # Cash in is free
        'total': amount,
        'status': 'completed',
        'timestamp': datetime.now().isoformat()
    }
    
    # Update balances
    users_db[user_id]['balance'] += amount
    agents_db[agent_id]['float_balance'] -= amount
    
    # Store transaction
    transactions_db[tx_id] = transaction
    
    return jsonify({
        'success': True,
        'transaction': transaction,
        'message': 'Cash in completed successfully'
    })

@app.route('/api/transactions/cash-out', methods=['POST'])
def cash_out():
    """Cash out transaction (user to agent)"""
    data = request.get_json()
    
    required_fields = ['user_id', 'agent_id', 'amount']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    user_id = data['user_id']
    agent_id = data['agent_id']
    amount = float(data['amount'])
    
    # Validate user and agent exist
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    if agent_id not in agents_db:
        return jsonify({'error': 'Agent not found'}), 404
    
    # Validate amount
    if amount < 50 or amount > 1000000:
        return jsonify({'error': 'Invalid amount'}), 400
    
    # Calculate fee
    fee = calculate_fee(amount, 'cashout')
    total = amount + fee
    
    # Check user balance
    user_balance = users_db[user_id]['balance']
    if user_balance < total:
        return jsonify({'error': 'Insufficient balance'}), 400
    
    # Create transaction
    tx_id = str(uuid.uuid4())
    transaction = {
        'id': tx_id,
        'type': 'cash_out',
        'from_user': user_id,
        'to_user': agent_id,
        'amount': amount,
        'fee': fee,
        'total': total,
        'status': 'completed',
        'timestamp': datetime.now().isoformat()
    }
    
    # Update balances
    users_db[user_id]['balance'] -= total
    agents_db[agent_id]['float_balance'] += amount
    
    # Calculate agent commission (25% of fee)
    agent_commission = fee * 0.25
    agents_db[agent_id]['commission_balance'] += agent_commission
    
    # Store transaction
    transactions_db[tx_id] = transaction
    
    return jsonify({
        'success': True,
        'transaction': transaction,
        'agent_commission': agent_commission,
        'message': 'Cash out completed successfully'
    })

@app.route('/api/qr/generate', methods=['POST'])
def generate_qr():
    """Generate QR code for payment"""
    data = request.get_json()
    
    required_fields = ['user_id', 'amount']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    user_id = data['user_id']
    amount = float(data['amount'])
    
    if user_id not in users_db:
        return jsonify({'error': 'User not found'}), 404
    
    # Generate QR code data
    qr_id = str(uuid.uuid4())
    qr_data = {
        'id': qr_id,
        'user_id': user_id,
        'amount': amount,
        'type': 'payment_request',
        'expires_at': (datetime.now() + timedelta(hours=1)).isoformat(),
        'created_at': datetime.now().isoformat()
    }
    
    qr_codes_db[qr_id] = qr_data
    
    return jsonify({
        'success': True,
        'qr_code': qr_data,
        'qr_url': f'cashout://pay?qr_id={qr_id}&amount={amount}'
    })

@app.route('/api/qr/scan', methods=['POST'])
def scan_qr():
    """Scan and process QR code"""
    data = request.get_json()
    
    required_fields = ['qr_id', 'scanner_user_id']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    qr_id = data['qr_id']
    scanner_user_id = data['scanner_user_id']
    
    if qr_id not in qr_codes_db:
        return jsonify({'error': 'Invalid QR code'}), 404
    
    qr_data = qr_codes_db[qr_id]
    
    # Check if QR code is expired
    expires_at = datetime.fromisoformat(qr_data['expires_at'])
    if datetime.now() > expires_at:
        return jsonify({'error': 'QR code has expired'}), 400
    
    # Check if scanner has sufficient balance
    scanner_balance = users_db[scanner_user_id]['balance']
    amount = qr_data['amount']
    fee = calculate_fee(amount, 'p2p')
    total = amount + fee
    
    if scanner_balance < total:
        return jsonify({'error': 'Insufficient balance'}), 400
    
    # Process payment
    tx_id = str(uuid.uuid4())
    transaction = {
        'id': tx_id,
        'type': 'qr_payment',
        'from_user': scanner_user_id,
        'to_user': qr_data['user_id'],
        'amount': amount,
        'fee': fee,
        'total': total,
        'qr_id': qr_id,
        'status': 'completed',
        'timestamp': datetime.now().isoformat()
    }
    
    # Update balances
    users_db[scanner_user_id]['balance'] -= total
    users_db[qr_data['user_id']]['balance'] += amount
    
    # Store transaction
    transactions_db[tx_id] = transaction
    
    # Remove QR code (one-time use)
    del qr_codes_db[qr_id]
    
    return jsonify({
        'success': True,
        'transaction': transaction,
        'message': 'Payment completed successfully'
    })

@app.route('/api/agents/<agent_id>', methods=['GET'])
def get_agent(agent_id):
    """Get agent profile"""
    if agent_id not in agents_db:
        return jsonify({'error': 'Agent not found'}), 404
    
    return jsonify(agents_db[agent_id])

@app.route('/api/businesses/<business_id>', methods=['GET'])
def get_business(business_id):
    """Get business profile"""
    if business_id not in businesses_db:
        return jsonify({'error': 'Business not found'}), 404
    
    return jsonify(businesses_db[business_id])

@app.route('/api/fees/calculate', methods=['POST'])
def calculate_fees():
    """Calculate transaction fees"""
    data = request.get_json()
    
    required_fields = ['amount', 'transaction_type']
    for field in required_fields:
        if field not in data:
            return jsonify({'error': f'Missing required field: {field}'}), 400
    
    amount = float(data['amount'])
    transaction_type = data['transaction_type']
    
    fee = calculate_fee(amount, transaction_type)
    total = amount + fee
    
    return jsonify({
        'amount': amount,
        'fee': fee,
        'total': total,
        'transaction_type': transaction_type
    })

if __name__ == '__main__':
    # Initialize sample data
    init_sample_data()
    
    # Get port from environment or use default
    port = int(os.environ.get('PORT', 5000))
    
    print(f"🚀 Starting Cashout API on port {port}")
    print(f"📱 Health check: http://localhost:{port}/api/health")
    print(f"📚 API Documentation:")
    print(f"   - GET  /api/health")
    print(f"   - GET  /api/users/<user_id>")
    print(f"   - GET  /api/users/<user_id>/balance")
    print(f"   - GET  /api/users/<user_id>/transactions")
    print(f"   - POST /api/transactions/send")
    print(f"   - POST /api/transactions/cash-in")
    print(f"   - POST /api/transactions/cash-out")
    print(f"   - POST /api/qr/generate")
    print(f"   - POST /api/qr/scan")
    print(f"   - POST /api/fees/calculate")
    
    app.run(host='0.0.0.0', port=port, debug=True) 