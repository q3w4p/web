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
Create a file at `/etc/nginx/sites-available/nazireich.site`:

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
        proxy_cache_bypass $http_upgrade;
    }
}
```

Enable the site:
```bash
ln -s /etc/nginx/sites-available/nazireich.site /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

## 4. Discord Bot (Go)
1. Go to `go-bot/` directory.
2. Build: `go build -o bot`
3. Start: `pm2 start ./bot --name "discord-bot"`
