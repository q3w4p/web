# Deployment Guide for Kamatera VPS

This guide walks you through deploying the Bot Hosting Platform on a Kamatera VPS with Cloudflare.

## Step 1: Kamatera VPS Setup

### 1.1 Create VPS Instance
1. Log in to Kamatera
2. Create new server with:
   - **OS**: Ubuntu 22.04 LTS (recommended)
   - **CPU**: 2 cores minimum
   - **RAM**: 2GB minimum
   - **Storage**: 20GB minimum
3. Note your server IP address

### 1.2 Initial Server Configuration

SSH into your server:
```bash
ssh root@your-server-ip
```

Update system:
```bash
apt update && apt upgrade -y
```

Install Node.js 20.x (LTS):
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt install -y nodejs
```

Verify installation:
```bash
node --version  # Should show v20.x or higher
npm --version   # Should show 10.x or higher
```

Install PM2 (Process Manager):
```bash
npm install -g pm2
```

Install Nginx:
```bash
apt install -y nginx
```

## Step 2: Upload Application Files

### Option A: Using Git (Recommended)
```bash
cd /var/www
git clone your-repository-url bot-hosting-platform
cd bot-hosting-platform
```

### Option B: Using SCP/SFTP
From your local machine:
```bash
scp -r bot-hosting-platform root@your-server-ip:/var/www/
```

### Option C: Manual Upload
Use FileZilla or similar SFTP client to upload files to `/var/www/bot-hosting-platform`

## Step 3: Application Setup

Navigate to application directory:
```bash
cd /var/www/bot-hosting-platform
```

Install dependencies:
```bash
npm install
```

Create .env file:
```bash
cp .env.example .env
nano .env
```

Configure your environment variables (see main README for details):
```env
PORT=3000
NODE_ENV=production
DB_HOST=your-aiven-hostname.aivencloud.com
DB_PORT=12345
DB_NAME=defaultdb
DB_USER=avnadmin
DB_PASSWORD=your-password
DB_SSL=true
SESSION_SECRET=generate-random-string-here
JWT_SECRET=generate-random-string-here
ADMIN_USERNAME=admin
ADMIN_PASSWORD=change-this-secure-password
```

Initialize database:
```bash
npm run init-db
```

## Step 4: Configure Nginx

Create Nginx configuration:
```bash
nano /etc/nginx/sites-available/bot-hosting
```

Add the following configuration:
```nginx
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/bot-hosting-access.log;
    error_log /var/log/nginx/bot-hosting-error.log;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        
        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        
        # Standard proxy headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeout settings
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable the site:
```bash
ln -s /etc/nginx/sites-available/bot-hosting /etc/nginx/sites-enabled/
```

Test Nginx configuration:
```bash
nginx -t
```

Restart Nginx:
```bash
systemctl restart nginx
```

## Step 5: Configure Firewall

Allow HTTP, HTTPS, and SSH:
```bash
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

Check status:
```bash
ufw status
```

## Step 6: Start Application with PM2

Start the application:
```bash
cd /var/www/bot-hosting-platform
pm2 start server.js --name bot-hosting
```

Configure PM2 to start on boot:
```bash
pm2 startup
pm2 save
```

Check application status:
```bash
pm2 status
pm2 logs bot-hosting
```

## Step 7: Cloudflare Configuration

### 7.1 Add Domain to Cloudflare
1. Log in to Cloudflare
2. Add your domain
3. Update nameservers at your domain registrar

### 7.2 DNS Configuration
In Cloudflare DNS settings:
- Add A record:
  - **Type**: A
  - **Name**: @ (or subdomain)
  - **Content**: Your Kamatera VPS IP
  - **Proxy status**: Proxied (orange cloud)
  - **TTL**: Auto

- Add CNAME for www (optional):
  - **Type**: CNAME
  - **Name**: www
  - **Content**: yourdomain.com
  - **Proxy status**: Proxied

### 7.3 SSL/TLS Settings
1. Go to SSL/TLS → Overview
2. Set encryption mode to **Full**
3. Enable "Always Use HTTPS"
4. Enable "Automatic HTTPS Rewrites"

### 7.4 Security Settings
1. Go to Security → Settings
2. Set Security Level to **Medium** or **High**
3. Enable "Browser Integrity Check"
4. Configure Challenge Passage as needed

### 7.5 Speed Optimization
1. Go to Speed → Optimization
2. Enable "Auto Minify" (CSS, JavaScript, HTML)
3. Enable "Brotli" compression
4. Enable "Rocket Loader" (optional)

### 7.6 Page Rules (Optional but Recommended)
Create page rules for:
- Cache everything: `yourdomain.com/css/*`, `yourdomain.com/js/*`
- Browser Cache TTL: 1 year for static assets

## Step 8: Install SSL Certificate (Let's Encrypt)

Install Certbot:
```bash
apt install -y certbot python3-certbot-nginx
```

Temporarily disable Cloudflare proxy (set to DNS only):
1. Go to Cloudflare DNS
2. Click on orange cloud icon to make it gray
3. Wait 2-3 minutes for DNS to propagate

Obtain SSL certificate:
```bash
certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

Follow prompts:
- Enter email address
- Agree to terms
- Choose to redirect HTTP to HTTPS (option 2)

Re-enable Cloudflare proxy:
1. Go back to Cloudflare DNS
2. Click gray cloud to make it orange again

## Step 9: Verify Deployment

Test the application:
```bash
# Check if server is running
pm2 status

# Check logs
pm2 logs bot-hosting --lines 100

# Test locally
curl http://localhost:3000

# Test Nginx
curl http://your-server-ip
```

Access from browser:
- Visit: https://yourdomain.com
- Should see landing page
- Try: https://yourdomain.com/admin.html
- Try: https://yourdomain.com/host.html

## Step 10: Monitoring and Maintenance

### Monitor Application
```bash
# View real-time logs
pm2 logs bot-hosting

# View monitoring dashboard
pm2 monit

# Check resource usage
pm2 show bot-hosting
```

### Monitor Server Resources
```bash
# CPU and memory
htop

# Disk usage
df -h

# Network connections
netstat -tulpn
```

### Application Updates
```bash
cd /var/www/bot-hosting-platform
git pull  # If using git
npm install  # Update dependencies
pm2 restart bot-hosting
```

### Database Backup (Recommended)
Create automated backups using Aiven's built-in backup feature:
1. Log in to Aiven Console
2. Select your PostgreSQL service
3. Go to "Backups" tab
4. Configure backup schedule

## Step 11: Go Bot Integration

If you have your Go bot on the same server:

### Configure Go Bot
Set up your Go bot to listen on a specific port (e.g., 8080):
```bash
# In your Go bot directory
./your-bot-binary --port 8080
```

Or use PM2 to manage Go bot:
```bash
pm2 start your-bot-binary --name go-bot -- --port 8080
```

### Update Environment Variables
Edit `/var/www/bot-hosting-platform/.env`:
```env
GO_BOT_HOST=localhost
GO_BOT_PORT=8080
GO_BOT_API_KEY=your-secure-api-key
```

Restart Node.js application:
```bash
pm2 restart bot-hosting
```

## Troubleshooting

### Application won't start
```bash
# Check logs
pm2 logs bot-hosting

# Check if port is in use
netstat -tulpn | grep 3000

# Check Node.js version
node --version
```

### Database connection issues
```bash
# Test connection from server
psql "postgresql://user:password@hostname:port/database?sslmode=require"

# Check environment variables
cat /var/www/bot-hosting-platform/.env
```

### Nginx issues
```bash
# Test configuration
nginx -t

# Check error logs
tail -f /var/log/nginx/error.log

# Restart Nginx
systemctl restart nginx
```

### Cloudflare issues
- Verify DNS propagation: `dig yourdomain.com`
- Check SSL/TLS mode is "Full"
- Disable rocket loader if JavaScript issues occur
- Check Cloudflare firewall rules

### WebSocket connection issues
1. Ensure Cloudflare WebSocket is enabled (automatic in most plans)
2. Check Nginx WebSocket proxy configuration
3. Verify firewall allows WebSocket connections

## Security Hardening (Recommended)

### 1. Create Non-Root User
```bash
adduser nodeuser
usermod -aG sudo nodeuser
```

### 2. Configure SSH
```bash
nano /etc/ssh/sshd_config
```
Set:
```
PermitRootLogin no
PasswordAuthentication no  # If using SSH keys
```

Restart SSH:
```bash
systemctl restart ssh
```

### 3. Install Fail2Ban
```bash
apt install -y fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### 4. Regular Updates
```bash
# Create update script
nano /root/update.sh
```

Add:
```bash
#!/bin/bash
apt update
apt upgrade -y
apt autoremove -y
```

Make executable:
```bash
chmod +x /root/update.sh
```

Schedule with cron:
```bash
crontab -e
```

Add:
```
0 2 * * 0 /root/update.sh
```

## Performance Optimization

### Enable Nginx Caching
```nginx
# Add to nginx config
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m use_temp_path=off;

location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
    proxy_cache my_cache;
    proxy_cache_valid 200 30d;
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

### PM2 Cluster Mode (Optional)
For better performance with multiple CPU cores:
```bash
pm2 delete bot-hosting
pm2 start server.js --name bot-hosting -i max
```

## Backup Strategy

### 1. Database Backups
Handled by Aiven automatically

### 2. Application Backups
```bash
# Create backup script
nano /root/backup-app.sh
```

Add:
```bash
#!/bin/bash
BACKUP_DIR="/backup"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
tar -czf $BACKUP_DIR/bot-hosting-$DATE.tar.gz /var/www/bot-hosting-platform
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

Make executable and schedule:
```bash
chmod +x /root/backup-app.sh
```

Add to cron:
```bash
0 3 * * * /root/backup-app.sh
```

## Conclusion

Your bot hosting platform should now be fully deployed and accessible at your domain. Monitor logs regularly and keep everything updated for optimal performance and security.

For additional help, refer to the main README.md file.
