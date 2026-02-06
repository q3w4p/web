# VPS Deployment Guide (Nginx + PM2)

## 1. Prerequisites
- Node.js (v18+)
- PM2: `npm install -g pm2`
- Nginx
- Go (for the discord bot)

## 2. Setup Application
1. Clone the repository to `/var/www/nazireich.site`
2. Install dependencies: `npm install`
3. Build the project: `npm run build`
4. Start with PM2:
   ```bash
   pm2 start npm --name "web-app" -- run start
   ```

## 3. Nginx Configuration
Since you are using **Cloudflare**, you should configure Nginx to handle the traffic properly. Cloudflare acts as a proxy, so your VPS only needs to listen for incoming requests from Cloudflare.

### Option A: Standard HTTP (Cloudflare handles SSL)
If your Cloudflare SSL/TLS setting is "Flexible", use this:

```nginx
server {
    listen 80;
    server_name nazireich.site;

    location / {
        proxy_pass http://localhost:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### Option B: Full SSL (Recommended)
1. In Cloudflare, set SSL/TLS to **"Full"** or **"Full (strict)"**.
2. Generate an **Origin CA Certificate** in Cloudflare (Websites > nazireich.site > SSL/TLS > Origin Server).
3. Save the certificate to `/etc/ssl/certs/cloudflare.pem` and the private key to `/etc/ssl/private/cloudflare.key`.
4. Use this Nginx config:

```nginx
server {
    listen 80;
    server_name nazireich.site;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name nazireich.site;

    ssl_certificate /etc/ssl/certs/cloudflare.pem;
    ssl_certificate_key /etc/ssl/private/cloudflare.key;

    location / {
        proxy_pass http://localhost:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## 4. Cloudflare DNS Settings
1. Go to your Cloudflare Dashboard.
2. Go to **DNS > Records**.
3. Add an **A Record**:
   - **Name**: `@` (or `nazireich.site`)
   - **IPv4 address**: Your Kamatera VPS IP
   - **Proxy status**: **Proxied** (Orange cloud)

## 5. Troubleshooting Kamatera Firewall
Kamatera has an external firewall. You **MUST** open ports 80 and 443 in the Kamatera console:
1. Log into Kamatera Console.
2. Go to **Network > Firewalls**.
3. Ensure there is a rule allowing **TCP** traffic on ports **80** (HTTP) and **443** (HTTPS) from **Anywhere**.

## 4. Discord Bot (Go)
1. Go to `go-bot/` directory.
2. Build: `go build -o bot`
3. Start: `pm2 start ./bot --name "discord-bot"`
