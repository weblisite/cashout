#!/bin/bash

echo "🚀 Starting Cashout Backend API..."

# Check if Python 3 is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed. Please install Python 3 first."
    exit 1
fi

# Check if virtual environment exists, create if not
if [ ! -d "venv" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "📥 Installing dependencies..."
pip install -r requirements.txt

# Start the API
echo "🌟 Starting API server..."
echo "📱 API will be available at: http://localhost:5000"
echo "🔍 Health check: http://localhost:5000/api/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

python3 simple_api.py 