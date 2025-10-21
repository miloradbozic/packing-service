#!/bin/bash

# Packing Service Heroku Deployment Script
# This script automates the deployment process to Heroku

set -e  # Exit on any error

echo "🚀 Packing Service Heroku Deployment Script"
echo "=========================================="

# Check if Heroku CLI is installed
if ! command -v heroku &> /dev/null; then
    echo "❌ Heroku CLI is not installed. Please install it first:"
    echo "   https://devcenter.heroku.com/articles/heroku-cli"
    exit 1
fi

# Check if user is logged in to Heroku
if ! heroku auth:whoami &> /dev/null; then
    echo "🔐 Please login to Heroku first:"
    heroku login
fi

# Get app name from user or use default
if [ -z "$1" ]; then
    echo "📝 Enter your Heroku app name (or press Enter for auto-generated name):"
    read -r APP_NAME
else
    APP_NAME="$1"
fi

# Create Heroku app
if [ -n "$APP_NAME" ]; then
    echo "🏗️  Creating Heroku app: $APP_NAME"
    heroku create "$APP_NAME" || {
        echo "⚠️  App might already exist, continuing..."
    }
else
    echo "🏗️  Creating Heroku app with auto-generated name"
    heroku create || {
        echo "⚠️  App creation failed, continuing..."
    }
fi

# Get the actual app name (in case it was auto-generated)
ACTUAL_APP_NAME=$(heroku apps:info --json | jq -r '.app.name')
echo "📱 Using app: $ACTUAL_APP_NAME"

# Add PostgreSQL addon
echo "🗄️  Adding PostgreSQL database..."
heroku addons:create heroku-postgresql:mini -a "$ACTUAL_APP_NAME" || {
    echo "⚠️  PostgreSQL addon might already exist, continuing..."
}

# Set environment variables
echo "⚙️  Setting environment variables..."
heroku config:set CONFIG_PATH=config.yaml -a "$ACTUAL_APP_NAME"

# Deploy the app
echo "🚀 Deploying to Heroku..."
git push heroku main

# Wait for deployment to complete
echo "⏳ Waiting for deployment to complete..."
sleep 10

# Check if the app is running
echo "🔍 Checking app status..."
heroku ps -a "$ACTUAL_APP_NAME"

# Open the app
echo "🌐 Opening your app in the browser..."
heroku open -a "$ACTUAL_APP_NAME"

echo ""
echo "✅ Deployment completed successfully!"
echo "📱 App URL: https://$ACTUAL_APP_NAME.herokuapp.com"
echo "🔧 API URL: https://$ACTUAL_APP_NAME.herokuapp.com/api/v1"
echo "❤️  Health Check: https://$ACTUAL_APP_NAME.herokuapp.com/health"
echo ""
echo "📊 Useful commands:"
echo "   heroku logs --tail -a $ACTUAL_APP_NAME    # View logs"
echo "   heroku ps -a $ACTUAL_APP_NAME             # Check status"
echo "   heroku restart -a $ACTUAL_APP_NAME        # Restart app"
echo "   heroku config -a $ACTUAL_APP_NAME         # View config"
echo ""
echo "🎉 Your Packing Service is now live on Heroku!"
