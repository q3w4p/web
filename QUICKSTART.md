# Quick Start Guide

Get your bot hosting platform running in 5 minutes!

## Prerequisites Checklist
- [ ] Node.js 20+ installed (LTS recommended)
- [ ] Aiven PostgreSQL database created
- [ ] VPS server access (Kamatera)
- [ ] Domain pointed to Cloudflare (optional)

## Quick Setup

### 1. Install Dependencies (30 seconds)
```bash
npm install
```

### 2. Configure Environment (2 minutes)
```bash
cp .env.example .env
nano .env
```

**Minimum Required Settings:**
```env
PORT=3000
DB_HOST=your-aiven-hostname.aivencloud.com
DB_PORT=12345
DB_NAME=defaultdb
DB_USER=avnadmin
DB_PASSWORD=your-password
DB_SSL=true
SESSION_SECRET=any-random-string-here
JWT_SECRET=another-random-string
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password
```

### 3. Initialize Database (30 seconds)
```bash
npm run init-db
```

You should see:
```
✓ Users table created
✓ Bot instances table created
✓ Admin sessions table created
✓ Logs table created
✓ Indexes created
✅ Database initialization completed successfully!
```

### 4. Start Server (10 seconds)
```bash
npm start
```

You should see:
```
🚀 Server running on port 3000
📡 WebSocket server ready
🔗 Admin panel: http://localhost:3000/admin.html
```

### 5. Access Your Platform
- **Landing Page**: http://localhost:3000
- **Host Page**: http://localhost:3000/host.html
- **Admin Panel**: http://localhost:3000/admin.html

## Testing It Works

### Test 1: Landing Page
1. Visit http://localhost:3000
2. Should see animated landing page
3. Click "Get Started" → Should go to host page

### Test 2: Host Page
1. Visit http://localhost:3000/host.html
2. Enter a test Discord token
3. Should validate and add user (or show error)

### Test 3: Admin Panel
1. Visit http://localhost:3000/admin.html
2. Login with credentials from .env
3. Should see dashboard with stats

## Common Issues

### "Cannot connect to database"
- Check Aiven credentials in .env
- Verify DB_SSL=true for Aiven
- Test connection: `psql "postgresql://user:password@host:port/database?sslmode=require"`

### "Port 3000 already in use"
- Change PORT in .env
- Or kill process: `lsof -ti:3000 | xargs kill`

### "Module not found"
- Run `npm install` again
- Check Node.js version: `node --version` (should be 20+)

## Next Steps

### For Development
```bash
npm run dev  # Auto-restart on file changes
```

### For Production
See DEPLOYMENT.md for full Kamatera + Cloudflare setup

### Adding Your Go Bot
1. Update `server.js` with Go bot API calls
2. Set GO_BOT_HOST, GO_BOT_PORT in .env
3. Restart server

## Getting Discord Token (For Testing)

**Desktop App Method:**
1. Open Discord
2. Press `Ctrl+Shift+I` (Windows/Linux) or `Cmd+Option+I` (Mac)
3. Go to Console tab
4. Paste and press Enter:
```javascript
window.webpackChunkdiscord_app.push([[''],{},e=>{m=[];for(let c in e.c)m.push(e.c[c])}]);m.find(m=>m?.exports?.default?.getToken!==void 0).exports.default.getToken()
```
5. Copy token (without quotes)

**Browser Method (Chrome/Firefox):**
1. Open Discord in browser: https://discord.com/app
2. Press F12 to open Developer Tools
3. Go to Console tab
4. Paste same code above
5. Copy token

⚠️ **Security Warning**: Never share your token with anyone!

## File Structure Overview

```
bot-hosting-platform/
├── server.js              # Main application
├── init-database.js       # DB setup script
├── package.json           # Dependencies
├── .env                   # Your configuration
├── README.md              # Full documentation
├── DEPLOYMENT.md          # VPS deployment guide
├── QUICKSTART.md          # This file
└── public/                # Web files
    ├── index.html         # Landing page
    ├── host.html          # Token submission
    ├── admin.html         # Control panel
    ├── css/
    │   └── styles.css     # All styles
    └── js/
        ├── landing.js     # Landing logic
        ├── host.js        # Host logic
        └── admin.js       # Admin logic
```

## Useful Commands

```bash
# Development
npm run dev              # Run with auto-restart

# Production
npm start                # Start server
pm2 start server.js      # Start with PM2 (recommended)
pm2 logs                 # View logs
pm2 restart all          # Restart
pm2 stop all             # Stop

# Database
npm run init-db          # Initialize/reset database

# Maintenance
pm2 monit                # Monitor resources
pm2 save                 # Save PM2 config
pm2 startup              # Enable auto-start on boot
```

## Support

Need help?
1. Check main README.md for detailed info
2. Check DEPLOYMENT.md for VPS setup
3. Review logs: `pm2 logs` or `npm start` output
4. Verify .env configuration
5. Test database connection

## Security Checklist

Before going live:
- [ ] Change ADMIN_USERNAME and ADMIN_PASSWORD
- [ ] Use strong SESSION_SECRET and JWT_SECRET
- [ ] Enable HTTPS (Cloudflare + Let's Encrypt)
- [ ] Configure firewall (ufw)
- [ ] Set DB_SSL=true
- [ ] Regular backups enabled

## What's Next?

1. **Customize**: Edit HTML/CSS in public/ folder
2. **Integrate**: Connect your Go bot in server.js
3. **Deploy**: Follow DEPLOYMENT.md for production
4. **Monitor**: Set up logging and monitoring
5. **Scale**: Use PM2 cluster mode for multiple cores

---

**Ready to deploy to production?** See DEPLOYMENT.md

**Need more details?** See README.md

**Having issues?** Check the troubleshooting sections!
