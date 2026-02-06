# Hurry - Discord Bot Hosting Platform

A beautiful, modern PHP-based Discord bot hosting platform with animated UI, gradient effects, and real-time management.

## ✨ Features

- 🎨 **Modern UI** - Gradient text, animated borders, and smooth transitions
- 🔐 **Discord Authentication** - Token-based user authentication
- 📊 **Admin Dashboard** - Real-time monitoring and management
- 🤖 **Multi-Bot Support** - Host multiple Discord accounts
- 📝 **Activity Logging** - Track all user actions
- 🗄️ **PostgreSQL (Aiven)** - Secure cloud database
- 🎭 **Animated Effects** - Gradient orbs, button borders, and more

## 🚀 Quick Start

### Prerequisites

- PHP 7.4+ with PDO PostgreSQL extension
- PostgreSQL database (Aiven)
- Apache/Nginx web server
- cURL extension enabled

### Installation

1. **Upload Files**
   ```bash
   # Upload all files to your web server
   /var/www/html/hurry/
   ```

2. **Configure Database**
   
   Edit `config.php` with your Aiven credentials:
   ```php
   define('DB_HOST', 'hurry-hurry.g.aivencloud.com');
   define('DB_PORT', '22637');
   define('DB_NAME', 'defaultdb');
   define('DB_USER', 'avnadmin');
   define('DB_PASSWORD', 'AVNS__v7YcrlbWVN6jtm03JL');
   ```

3. **Set Admin Password**
   
   In `config.php`, change the admin password:
   ```php
   define('ADMIN_PASSWORD', 'YourSecurePassword123!');
   ```

4. **Initialize Database**
   ```bash
   php setup-database.php
   ```

5. **Set Permissions**
   ```bash
   chmod 755 *.php
   chmod 755 assets/ -R
   ```

6. **Access Platform**
   - Main page: `https://yourdomain.com/`
   - Host page: `https://yourdomain.com/host.php`
   - Admin panel: `https://yourdomain.com/login.php`

## 📁 File Structure

```
hurry-bot-platform/
├── index.php              # Landing page
├── host.php              # Discord token submission
├── login.php             # Admin login
├── logout.php            # Logout handler
├── dashboard.php         # Admin dashboard
├── config.php            # Configuration & DB connection
├── setup-database.php    # Database initialization
├── includes/
│   └── functions.php     # Core functions
└── assets/
    ├── css/
    │   ├── style.css     # Main styles
    │   ├── host.css      # Host page styles
    │   └── admin.css     # Admin panel styles
    └── js/
        ├── main.js       # Main JavaScript
        └── dashboard.js  # Dashboard functions
```

## 🎨 Features Breakdown

### Landing Page
- Animated gradient text "Hurry"
- Floating gradient orbs background
- Animated button borders on hover
- Smooth parallax effects
- Responsive design

### Host Page
- Discord token submission
- Token validation via Discord API
- Instructions for obtaining token
- Copy-to-clipboard functionality
- Success/error alerts

### Admin Dashboard
- Real-time statistics (users, active/inactive bots)
- User management (start, stop, delete)
- Bot instance monitoring
- Activity logs
- Auto-refresh every 30 seconds

## 🔧 Configuration

### Environment Settings

Edit `config.php` to customize:

```php
// Admin Credentials
define('ADMIN_USERNAME', 'admin');
define('ADMIN_PASSWORD', 'YourSecurePassword');

// Session Security
define('SESSION_SECRET', 'your-random-secret');

// Discord API
define('DISCORD_API_BASE', 'https://discord.com/api/v10');

// Go Bot Integration (optional)
define('GO_BOT_HOST', 'localhost');
define('GO_BOT_PORT', '8080');
define('GO_BOT_API_KEY', 'your-api-key');
```

### Web Server Configuration

#### Apache (.htaccess)
```apache
RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ index.php [L,QSA]

# Security Headers
Header set X-Frame-Options "SAMEORIGIN"
Header set X-Content-Type-Options "nosniff"
Header set X-XSS-Protection "1; mode=block"
```

#### Nginx
```nginx
location / {
    try_files $uri $uri/ /index.php?$query_string;
}

location ~ \.php$ {
    fastcgi_pass unix:/var/run/php/php7.4-fpm.sock;
    fastcgi_index index.php;
    include fastcgi_params;
    fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
}

# Security Headers
add_header X-Frame-Options "SAMEORIGIN";
add_header X-Content-Type-Options "nosniff";
add_header X-XSS-Protection "1; mode=block";
```

## 🔌 Go Bot Integration

To integrate with your Go bot, edit `includes/functions.php`:

### Start Bot Function
```php
function startBot($userId) {
    $user = getUserById($userId);
    
    // Call your Go bot API
    $ch = curl_init();
    curl_setopt_array($ch, [
        CURLOPT_URL => 'http://' . GO_BOT_HOST . ':' . GO_BOT_PORT . '/start',
        CURLOPT_POST => true,
        CURLOPT_POSTFIELDS => json_encode([
            'token' => $user['discord_token'],
            'userId' => $userId
        ]),
        CURLOPT_HTTPHEADER => [
            'X-API-Key: ' . GO_BOT_API_KEY,
            'Content-Type: application/json'
        ],
        CURLOPT_RETURNTRANSFER => true,
    ]);
    
    $response = curl_exec($ch);
    curl_close($ch);
    
    // Rest of the function...
}
```

## 🎯 Usage

### For Users

1. Visit the main page
2. Click "Get Started"
3. Enter Discord token
4. Wait for validation
5. Account added successfully!

### For Admins

1. Click "Admin Panel" on main page
2. Login with credentials
3. View all users and statistics
4. Start/stop bots for any user
5. Delete users if needed
6. Monitor activity logs

## 🛡️ Security

- ✅ Password-protected admin panel
- ✅ Session management with secrets
- ✅ SQL injection prevention (PDO prepared statements)
- ✅ XSS protection (htmlspecialchars on all outputs)
- ✅ CSRF protection (session tokens)
- ✅ SSL/TLS encryption (Aiven)
- ✅ Input validation

### Recommendations

1. **Change default admin password** immediately
2. **Use HTTPS** - Configure SSL certificate
3. **Regular backups** - Backup database regularly
4. **Update PHP** - Keep PHP version current
5. **Monitor logs** - Check activity logs regularly

## 📊 Database Schema

### Users Table
```sql
id              SERIAL PRIMARY KEY
discord_token   VARCHAR(255) UNIQUE NOT NULL
discord_username VARCHAR(255)
discord_user_id VARCHAR(255)
avatar_url      TEXT
status          VARCHAR(50) DEFAULT 'inactive'
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
last_active     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Bot Instances Table
```sql
id              SERIAL PRIMARY KEY
user_id         INTEGER REFERENCES users(id)
instance_name   VARCHAR(255)
status          VARCHAR(50) DEFAULT 'stopped'
pid             INTEGER
started_at      TIMESTAMP
stopped_at      TIMESTAMP
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Logs Table
```sql
id              SERIAL PRIMARY KEY
user_id         INTEGER REFERENCES users(id)
action          VARCHAR(255) NOT NULL
details         TEXT
ip_address      VARCHAR(45)
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Check PHP PDO PostgreSQL extension
php -m | grep pdo_pgsql

# Install if missing (Ubuntu/Debian)
sudo apt-get install php-pgsql
sudo systemctl restart apache2
```

### Permission Denied
```bash
# Set correct permissions
chmod 755 *.php
chmod 755 assets/ -R
chown www-data:www-data -R .
```

### Discord API Not Responding
- Check if cURL is enabled: `php -m | grep curl`
- Verify Discord API base URL in config.php
- Check firewall settings

### Styles Not Loading
- Verify file paths in HTML
- Check web server configuration
- Clear browser cache

## 📝 Changelog

### Version 1.0.0
- Initial release
- Landing page with animated effects
- Discord token authentication
- Admin dashboard
- User management
- Bot instance control
- Activity logging

## 🤝 Support

For issues or questions:
1. Check troubleshooting section
2. Review configuration settings
3. Check server logs
4. Verify database connection

## 📄 License

MIT License - Free to use and modify

## 🙏 Credits

- Design inspired by modern SaaS platforms
- Icons: Unicode emojis
- Fonts: System fonts
- Database: Aiven PostgreSQL

---

**Made with ❤️ for the Discord bot hosting community**
