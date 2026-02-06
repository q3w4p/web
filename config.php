<?php
// Database Configuration (Aiven PostgreSQL)
define('DB_HOST', 'hurry-hurry.g.aivencloud.com');
define('DB_PORT', '22637');
define('DB_NAME', 'defaultdb');
define('DB_USER', 'avnadmin');
define('DB_PASSWORD', 'AVNS__v7YcrlbWVN6jtm03JL');

// Session Configuration
define('SESSION_SECRET', 'a8f3h9s7d6k2m4n5b7v9c1x3z5a7s9d1f3g5h7j9k1m3n5b7v9');

// Admin Credentials
define('ADMIN_USERNAME', 'admin');
define('ADMIN_PASSWORD', 'HurryAdmin2024!'); // Change this!

// Discord API
define('DISCORD_API_BASE', 'https://discord.com/api/v10');

// Go Bot Configuration (if using)
define('GO_BOT_HOST', 'localhost');
define('GO_BOT_PORT', '8080');
define('GO_BOT_API_KEY', 'your-bot-api-key');

// Timezone
date_default_timezone_set('UTC');

// Error Reporting (disable in production)
error_reporting(E_ALL);
ini_set('display_errors', 1);

// Database Connection Function
function getDBConnection() {
    try {
        $dsn = sprintf(
            "pgsql:host=%s;port=%s;dbname=%s;sslmode=require",
            DB_HOST,
            DB_PORT,
            DB_NAME
        );
        
        $pdo = new PDO($dsn, DB_USER, DB_PASSWORD, [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
            PDO::ATTR_EMULATE_PREPARES => false,
        ]);
        
        return $pdo;
    } catch (PDOException $e) {
        error_log("Database connection failed: " . $e->getMessage());
        die("Database connection failed. Please check your configuration.");
    }
}
?>
