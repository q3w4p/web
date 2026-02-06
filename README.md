# Bot Hosting Platform

A web-based Discord bot hosting platform with token injection, built with Node.js, Express, PostgreSQL (Aiven), and WebSocket for real-time updates.

## Features

- 🚀 **Landing Page**: Professional landing page with animations
- 🔐 **Discord Token Authentication**: Secure token validation and storage
- 📊 **Admin Control Panel (ACCP)**: Real-time monitoring dashboard
- 👥 **Multi-Account Support**: Host multiple Discord accounts simultaneously
- 🔄 **Live Updates**: WebSocket-based real-time status updates
- 📈 **Activity Logging**: Comprehensive action logging
- 🗄️ **PostgreSQL Database**: Secure data storage with Aiven

## Tech Stack

- **Backend**: Node.js, Express
- **Database**: PostgreSQL (Aiven)
- **Real-time**: WebSocket (ws)
- **Frontend**: Vanilla JavaScript, CSS3
- **Deployment**: Compatible with any VPS (Kamatera, DigitalOcean, etc.)

## Prerequisites

- Node.js 20+ installed (LTS recommended)
- Aiven PostgreSQL database account
- Discord API access
- VPS server (Kamatera or similar)

## Installation

### 1. Clone or Upload Files

Upload all files to your VPS server at `/var/www/bot-hosting-platform` or your preferred directory.

### 2. Install Dependencies

```bash
cd /path/to/bot-hosting-platform
npm install
```

### 3. Configure Environment Variables

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Edit `.env` with your actual configuration:

```env
# Server Configuration
PORT=3000
NODE_ENV=production

# Aiven PostgreSQL Database Configuration
DB_HOST=your-aiven-hostname.aivencloud.com
DB_PORT=12345
DB_NAME=defaultdb
DB_USER=avnadmin
DB_PASSWORD=your-secure-password
DB_SSL=true

# Session Secret (generate a random string)
SESSION_SECRET=your-random-secret-key-change-this

# JWT Secret for token generation
JWT_SECRET=your-jwt-secret-key-change-this

# Admin Credentials
ADMIN_USERNAME=admin
ADMIN_PASSWORD=change-this-password

# Discord API
DISCORD_API_BASE=https://discord.com/api/v10

# Go Bot Server Configuration (if applicable)
GO_BOT_HOST=localhost
GO_BOT_PORT=8080
GO_BOT_API_KEY=your-bot-api-key
```

### 4. Initialize Database

Run the database initialization script:

```bash
npm run init-db
```

This will create all necessary tables and indexes in your Aiven PostgreSQL database.

### 5. Start the Server

For production:
```bash
npm start
```

For development (with auto-restart):
```bash
npm run dev
```

The server will start on port 3000 (or your configured PORT).

## Aiven PostgreSQL Setup

### Getting Your Aiven Credentials

1. Log in to [Aiven Console](https://console.aiven.io/)
2. Create or select your PostgreSQL service
3. From the service page, copy:
   - **Host**: Found under "Connection Information"
   - **Port**: Usually 12345 or similar
   - **Database**: Default is `defaultdb`
   - **User**: Usually `avnadmin`
   - **Password**: Set during service creation
4. SSL is automatically enabled for Aiven

### Connection String Example

```
postgresql://avnadmin:password@hostname.aivencloud.com:12345/defaultdb?sslmode=require
```

## Cloudflare Setup

### 1. Point Domain to Your Server

In Cloudflare DNS settings:
- Add an A record pointing to your VPS IP address
- Set Proxy status to "Proxied" (orange cloud)

### 2. SSL/TLS Settings

- Set SSL/TLS encryption mode to "Full" or "Full (strict)"
- Enable "Always Use HTTPS"

### 3. Reverse Proxy (Nginx)

Install and configure Nginx on your VPS:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Enable and restart Nginx:
```bash
sudo ln -s /etc/nginx/sites-available/bot-hosting /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## Pages Overview

### 1. Landing Page (`/`)
- Main entry point with animated hero section
- "Get Started" button → Host page
- "Admin Panel" button → Admin login

### 2. Host Page (`/host.html`)
- Discord token submission form
- Token validation and account addition
- Instructions for obtaining Discord token

### 3. Admin Panel (`/admin.html`)
- Login required (credentials from .env)
- Real-time user monitoring
- Bot instance management
- Activity logs
- Statistics dashboard

## API Endpoints

### Public Endpoints
- `POST /api/users/add` - Add Discord user with token

### Admin Endpoints (Authentication Required)
- `POST /api/admin/login` - Admin login
- `POST /api/admin/logout` - Admin logout
- `GET /api/admin/check` - Check admin status
- `GET /api/users` - Get all users
- `POST /api/bot/start` - Start bot instance
- `POST /api/bot/stop` - Stop bot instance
- `GET /api/bot/instances` - Get all bot instances
- `DELETE /api/users/:id` - Delete user
- `GET /api/logs` - Get activity logs

## WebSocket Events

### Client → Server
- `ping` - Keep-alive ping

### Server → Client
- `pong` - Keep-alive response
- `user_update` - User added/updated
- `bot_status` - Bot started/stopped
- `user_deleted` - User deleted

## Integration with Go Bot

To integrate with your existing Go bot:

1. Update `server.js` in the bot start/stop functions
2. Add API calls to your Go bot server:

```javascript
// Example in server.js
await axios.post(`http://${process.env.GO_BOT_HOST}:${process.env.GO_BOT_PORT}/start`, {
    token: user.discord_token,
    userId: userId
}, {
    headers: { 'X-API-Key': process.env.GO_BOT_API_KEY }
});
```

## Security Best Practices

1. **Change Default Credentials**: Update admin username/password in .env
2. **Use Strong Secrets**: Generate random strings for SESSION_SECRET and JWT_SECRET
3. **Enable HTTPS**: Always use SSL/TLS in production
4. **Database Security**: Use Aiven's built-in security features
5. **Token Storage**: Tokens are stored encrypted in PostgreSQL
6. **Rate Limiting**: Consider adding rate limiting for production

## Process Management (PM2)

For production deployment, use PM2:

```bash
# Install PM2
npm install -g pm2

# Start application
pm2 start server.js --name "bot-hosting"

# Enable auto-restart on server reboot
pm2 startup
pm2 save

# Monitor
pm2 monit

# View logs
pm2 logs bot-hosting
```

## Troubleshooting

### Database Connection Issues
- Verify Aiven credentials in .env
- Check if SSL is enabled (should be true for Aiven)
- Ensure your VPS IP is whitelisted in Aiven (usually automatic)

### WebSocket Connection Issues
- Check if firewall allows WebSocket connections
- Verify Cloudflare WebSocket support is enabled
- Ensure proxy configuration passes Upgrade headers

### Token Validation Fails
- Verify Discord API base URL is correct
- Check if token format is valid
- Ensure network can reach Discord API

## File Structure

```
bot-hosting-platform/
├── server.js                 # Main Express server
├── init-database.js          # Database initialization
├── package.json              # Dependencies
├── .env                      # Environment variables
├── .env.example             # Environment template
└── public/
    ├── index.html           # Landing page
    ├── host.html            # Host page
    ├── admin.html           # Admin panel
    ├── css/
    │   └── styles.css       # All styles
    └── js/
        ├── landing.js       # Landing page logic
        ├── host.js          # Host page logic
        └── admin.js         # Admin panel logic
```

## Support

For issues or questions:
1. Check the troubleshooting section
2. Review server logs: `pm2 logs bot-hosting`
3. Check database connection: `npm run init-db`

## License

MIT License - Feel free to modify and use for your needs.

## Contributing

Contributions are welcome! Please ensure:
- Code follows existing style
- Test thoroughly before submitting
- Update documentation as needed
