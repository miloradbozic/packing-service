# Heroku Deployment Guide

This guide will help you deploy the Packing Service to Heroku.

## Prerequisites

1. **Heroku CLI** - [Install here](https://devcenter.heroku.com/articles/heroku-cli)
2. **Git** - Make sure your project is in a Git repository
3. **Heroku Account** - Sign up at [heroku.com](https://heroku.com)

## Quick Deployment

### Option 1: One-Click Deploy (Recommended)

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/miloradbozic/packing-service)

Click the button above to deploy directly from GitHub with all the necessary configuration.

### Option 2: Manual Deployment

#### Step 1: Login to Heroku
```bash
heroku login
```

#### Step 2: Create a Heroku App
```bash
# Create app (replace 'your-app-name' with your desired name)
heroku create your-app-name

# Or let Heroku generate a random name
heroku create
```

#### Step 3: Add PostgreSQL Database
```bash
# Add the free PostgreSQL addon
heroku addons:create heroku-postgresql:mini
```

#### Step 4: Deploy Your Code
```bash
# Add Heroku remote (if not already added)
git remote add heroku https://git.heroku.com/your-app-name.git

# Deploy to Heroku
git push heroku main
```

#### Step 5: Open Your App
```bash
heroku open
```

## Environment Variables

The app automatically uses Heroku's `DATABASE_URL` environment variable for database connection. No additional configuration is needed.

### Optional Environment Variables

You can set these if you want to override defaults:

```bash
# Set custom host (default: 0.0.0.0)
heroku config:set HOST=0.0.0.0

# Set custom config path (default: config.yaml)
heroku config:set CONFIG_PATH=config.yaml
```

## Database Migrations

The app automatically runs database migrations on startup, so your database will be properly set up with the required tables and initial data.

## Monitoring and Logs

### View Logs
```bash
heroku logs --tail
```

### Check App Status
```bash
heroku ps
```

### Restart App
```bash
heroku restart
```

## API Endpoints

Once deployed, your app will be available at:
- **Web UI**: `https://your-app-name.herokuapp.com/`
- **API**: `https://your-app-name.herokuapp.com/api/v1/`
- **Health Check**: `https://your-app-name.herokuapp.com/health`

### Example API Usage

```bash
# Calculate packing for 501 items
curl -X POST https://your-app-name.herokuapp.com/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"items": 501}'

# Get current pack sizes
curl https://your-app-name.herokuapp.com/api/v1/pack-sizes
```

## Troubleshooting

### Common Issues

1. **Build Fails**: Make sure your `go.mod` file is committed and up to date
2. **Database Connection Issues**: Check that the PostgreSQL addon is properly provisioned
3. **App Crashes**: Check logs with `heroku logs --tail`

### Useful Commands

```bash
# Check app info
heroku info

# Scale the app (if needed)
heroku ps:scale web=1

# Access the database
heroku pg:psql

# View environment variables
heroku config
```

## Updating Your App

To update your deployed app:

```bash
# Make your changes locally
git add .
git commit -m "Your update message"

# Deploy to Heroku
git push heroku main
```

## Cost

- **Free Tier**: Heroku no longer offers a free tier
- **Basic Plan**: $5/month for the web dyno
- **PostgreSQL Mini**: $5/month for the database
- **Total**: ~$10/month for a basic deployment

## Security Notes

- The app uses SSL in production (Heroku handles this automatically)
- Database connections are encrypted
- Environment variables are secure and not exposed in logs
