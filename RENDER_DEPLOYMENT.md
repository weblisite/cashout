# 🚀 **Deploy Cashout Platform on Render**

## 📋 **Overview**

This guide will help you deploy the complete Cashout platform on Render, including:
- **Backend API** (Python Flask)
- **Frontend Web Apps** (Static sites)
- **PWA** (Progressive Web App)

## 🎯 **Deployment Architecture**

```
Render Services:
├── cashout-api.onrender.com     # Backend API
├── cashout-web.onrender.com     # Frontend Web Apps
└── cashout-pwa.onrender.com     # PWA Build
```

---

## 🚀 **Step-by-Step Deployment**

### **Step 1: Prepare Repository**

1. **Ensure all files are committed to GitHub**
   ```bash
   git add .
   git commit -m "Prepare for Render deployment"
   git push origin main
   ```

2. **Verify repository structure**
   ```
   Cashout/
   ├── render.yaml              # Render configuration
   ├── backend/
   │   ├── simple_api.py       # Flask API
   │   ├── wsgi.py             # WSGI entry point
   │   └── requirements.txt    # Python dependencies
   ├── landing-page/           # Frontend files
   └── build/                  # Build outputs
   ```

### **Step 2: Deploy on Render**

#### **Option A: Using render.yaml (Recommended)**

1. **Go to Render Dashboard**
   - Visit [https://dashboard.render.com](https://dashboard.render.com)
   - Sign in or create account

2. **Create New Blueprint**
   - Click "New +"
   - Select "Blueprint"
   - Connect your GitHub repository: `weblisite/cashout`

3. **Deploy Services**
   - Render will automatically detect `render.yaml`
   - Click "Apply" to deploy all services
   - Wait for deployment to complete

#### **Option B: Manual Deployment**

**Deploy Backend API:**
1. Click "New +" → "Web Service"
2. Connect GitHub repository
3. Configure:
   - **Name**: `cashout-api`
   - **Environment**: `Python 3`
   - **Build Command**: `cd backend && pip install -r requirements.txt`
   - **Start Command**: `cd backend && gunicorn wsgi:app`
   - **Plan**: Free

**Deploy Frontend:**
1. Click "New +" → "Static Site"
2. Connect GitHub repository
3. Configure:
   - **Name**: `cashout-web`
   - **Build Command**: `echo "Frontend ready"`
   - **Publish Directory**: `landing-page`
   - **Plan**: Free

---

## 🔧 **Configuration Files**

### **render.yaml**
```yaml
services:
  # Backend API Service
  - type: web
    name: cashout-api
    env: python
    plan: free
    buildCommand: |
      cd backend
      python -m pip install --upgrade pip
      pip install -r requirements.txt
    startCommand: |
      cd backend
      gunicorn wsgi:app
    envVars:
      - key: PYTHON_VERSION
        value: 3.9.0
      - key: PORT
        value: 10000
      - key: FLASK_ENV
        value: production
      - key: FLASK_DEBUG
        value: false

  # Frontend Web App Service
  - type: web
    name: cashout-web
    env: static
    plan: free
    buildCommand: |
      echo "Building frontend..."
      mkdir -p build/web
      cp -r landing-page/* build/web/
      echo "Frontend build complete"
    staticPublishPath: ./build/web
    envVars:
      - key: NODE_VERSION
        value: 18.0.0
```

### **backend/wsgi.py**
```python
#!/usr/bin/env python3
"""
WSGI entry point for Cashout API
Production deployment for Render
"""

import os
from simple_api import app

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 10000))
    app.run(host="0.0.0.0", port=port, debug=False)
```

### **backend/requirements.txt**
```
Flask==2.3.3
Flask-CORS==4.0.0
gunicorn==21.2.0
```

---

## 🌐 **Deployment URLs**

After deployment, you'll get these URLs:

### **Production URLs**
- **Backend API**: `https://cashout-api.onrender.com`
- **Frontend Web**: `https://cashout-web.onrender.com`
- **PWA**: `https://cashout-pwa.onrender.com`

### **API Endpoints**
- **Health Check**: `https://cashout-api.onrender.com/api/health`
- **User Data**: `https://cashout-api.onrender.com/api/users/{user_id}`
- **Transactions**: `https://cashout-api.onrender.com/api/transactions/*`

---

## 📱 **Mobile App Deployment**

### **PWA Deployment**
1. **Deploy PWA Service**
   - Use the PWA service in `render.yaml`
   - PWA will be available at `https://cashout-pwa.onrender.com`

2. **Test PWA Installation**
   - Open on mobile device
   - Tap "Add to Home Screen"
   - App works like native app

### **Android APK Generation**
1. **Use PWA Builder**
   - Go to [https://www.pwabuilder.com](https://www.pwabuilder.com)
   - Enter: `https://cashout-pwa.onrender.com`
   - Generate Android package
   - Download APK

---

## 🔍 **Testing Deployment**

### **1. Test Backend API**
```bash
# Health check
curl https://cashout-api.onrender.com/api/health

# Get user data
curl https://cashout-api.onrender.com/api/users/user1

# Calculate fees
curl -X POST https://cashout-api.onrender.com/api/fees/calculate \
  -H "Content-Type: application/json" \
  -d '{"amount": 1000, "transaction_type": "p2p"}'
```

### **2. Test Frontend Apps**
- **User App**: `https://cashout-web.onrender.com/user-app.html`
- **Agent App**: `https://cashout-web.onrender.com/agent-app.html`
- **Business App**: `https://cashout-web.onrender.com/business-app.html`

### **3. Test PWA**
- Open `https://cashout-pwa.onrender.com` on mobile
- Install as app
- Test offline functionality

---

## 🔧 **Environment Variables**

### **Backend Environment Variables**
```bash
FLASK_ENV=production
FLASK_DEBUG=false
PORT=10000
PYTHON_VERSION=3.9.0
```

### **Frontend Environment Variables**
```bash
NODE_VERSION=18.0.0
```

---

## 📊 **Monitoring & Logs**

### **View Logs**
1. Go to Render Dashboard
2. Select your service
3. Click "Logs" tab
4. Monitor real-time logs

### **Health Monitoring**
- **API Health**: `https://cashout-api.onrender.com/api/health`
- **Uptime**: Render provides 99.9% uptime
- **Performance**: Automatic scaling

---

## 🚀 **Production Features**

### **Automatic Scaling**
- Render automatically scales based on traffic
- Free tier includes basic scaling
- Paid plans for advanced scaling

### **SSL Certificates**
- Automatic HTTPS certificates
- Secure connections by default
- No additional configuration needed

### **CDN**
- Global content delivery network
- Fast loading worldwide
- Automatic caching

---

## 🔒 **Security**

### **HTTPS**
- All connections use HTTPS
- Automatic SSL certificates
- Secure API endpoints

### **CORS Configuration**
- Properly configured for production
- Secure cross-origin requests
- API protection

---

## 📈 **Performance Optimization**

### **Backend Optimization**
- Gunicorn for production WSGI server
- Optimized for Render environment
- Efficient request handling

### **Frontend Optimization**
- Static file serving
- CDN delivery
- Caching headers

---

## 🛠️ **Troubleshooting**

### **Common Issues**

**1. Build Failures**
```bash
# Check build logs in Render dashboard
# Verify requirements.txt is correct
# Ensure Python version compatibility
```

**2. API Connection Issues**
```bash
# Verify API URL is correct
# Check CORS configuration
# Test API endpoints directly
```

**3. PWA Issues**
```bash
# Verify manifest.json
# Check service worker
# Test on HTTPS only
```

### **Debug Commands**
```bash
# Test API locally
cd backend
python simple_api.py

# Test frontend locally
cd landing-page
python3 -m http.server 8000
```

---

## 📞 **Support**

### **Render Support**
- [Render Documentation](https://render.com/docs)
- [Render Community](https://community.render.com)
- [Render Status](https://status.render.com)

### **Cashout Platform Support**
- Check logs in Render dashboard
- Review API documentation
- Test endpoints manually

---

## 🎉 **Deployment Complete**

After successful deployment:

### **✅ What's Live**
- **Backend API**: Full transaction processing
- **Frontend Apps**: User, Agent, Business interfaces
- **PWA**: Mobile app capabilities
- **Documentation**: Complete API docs

### **🚀 Next Steps**
1. **Test all endpoints**
2. **Verify mobile app installation**
3. **Monitor performance**
4. **Scale as needed**

### **📱 Mobile App Distribution**
- **PWA**: Share URL for instant installation
- **Android**: Generate APK via PWA Builder
- **iOS**: Use Safari "Add to Home Screen"

---

**🎯 Your Cashout platform is now live on Render! 🎯**

**Production URLs:**
- **API**: `https://cashout-api.onrender.com`
- **Web**: `https://cashout-web.onrender.com`
- **PWA**: `https://cashout-pwa.onrender.com` 