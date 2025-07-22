#!/usr/bin/env python3
"""
WSGI entry point for Cashout API
Production deployment for Render
"""

import os
from simple_api import app

if __name__ == "__main__":
    # Get port from environment variable (Render sets this)
    port = int(os.environ.get("PORT", 10000))
    
    # Run the app
    app.run(host="0.0.0.0", port=port, debug=False) 