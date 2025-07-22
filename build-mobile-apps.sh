#!/bin/bash

echo "🚀 Building Cashout Mobile Apps..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}Node.js not found. Installing via Homebrew...${NC}"
    if command -v brew &> /dev/null; then
        brew install node
    else
        echo -e "${RED}Homebrew not found. Please install Node.js manually.${NC}"
        exit 1
    fi
fi

# Create build directory
mkdir -p build
cd build

echo -e "${BLUE}📱 Building Progressive Web App (PWA)...${NC}"

# Create PWA build
mkdir -p pwa
cp -r ../landing-page/* pwa/

# Create icons directory
mkdir -p pwa/icons

echo -e "${GREEN}✅ PWA build created in build/pwa/${NC}"
echo -e "${BLUE}🌐 To test PWA:${NC}"
echo -e "   1. cd build/pwa"
echo -e "   2. python3 -m http.server 8000"
echo -e "   3. Open http://localhost:8000 in mobile browser"
echo -e "   4. Tap 'Add to Home Screen'"

echo ""
echo -e "${BLUE}📦 Building with Capacitor (Android/iOS)...${NC}"

# Check if Capacitor is available
if command -v npx &> /dev/null; then
    echo -e "${YELLOW}Installing Capacitor...${NC}"
    npx @capacitor/cli create cashout-app com.cashout.app Cashout
    
    if [ $? -eq 0 ]; then
        cd cashout-app
        
        # Copy web assets
        rm -rf www
        cp -r ../../landing-page www
        
        # Install Capacitor
        npm install @capacitor/core @capacitor/cli
        npm install @capacitor/android @capacitor/ios
        
        # Build for platforms
        npx cap add android
        npx cap add ios
        
        echo -e "${GREEN}✅ Capacitor app created in build/cashout-app/${NC}"
        echo -e "${BLUE}📱 To build Android:${NC}"
        echo -e "   1. cd build/cashout-app"
        echo -e "   2. npx cap sync"
        echo -e "   3. npx cap open android"
        echo -e "${BLUE}🍎 To build iOS:${NC}"
        echo -e "   1. cd build/cashout-app"
        echo -e "   2. npx cap sync"
        echo -e "   3. npx cap open ios"
    else
        echo -e "${RED}Failed to create Capacitor app${NC}"
    fi
else
    echo -e "${YELLOW}Node.js/npm not available. Skipping Capacitor build.${NC}"
fi

echo ""
echo -e "${BLUE}📱 Building with TWA (Trusted Web Activity)...${NC}"

# Create TWA build
mkdir -p twa
cp -r ../landing-page/* twa/

# Create TWA manifest
cat > twa/manifest.json << 'EOF'
{
  "name": "Cashout",
  "short_name": "Cashout",
  "start_url": "/user-app.html",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#667eea",
  "orientation": "portrait",
  "scope": "/",
  "icons": [
    {
      "src": "icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
EOF

echo -e "${GREEN}✅ TWA build created in build/twa/${NC}"

echo ""
echo -e "${BLUE}📋 Build Summary:${NC}"
echo -e "${GREEN}✅ PWA: build/pwa/ (Ready to deploy)${NC}"
echo -e "${GREEN}✅ TWA: build/twa/ (Ready for Android)${NC}"

if [ -d "cashout-app" ]; then
    echo -e "${GREEN}✅ Capacitor: build/cashout-app/ (Ready for native builds)${NC}"
fi

echo ""
echo -e "${BLUE}🚀 Next Steps:${NC}"
echo -e "1. ${YELLOW}For PWA:${NC} Deploy build/pwa/ to any web server"
echo -e "2. ${YELLOW}For Android TWA:${NC} Use Bubblewrap or PWA Builder"
echo -e "3. ${YELLOW}For Native Apps:${NC} Use Android Studio or Xcode with Capacitor"

echo ""
echo -e "${GREEN}🎉 Build completed successfully!${NC}" 