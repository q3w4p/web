<?php
require_once 'config.php';

echo "Initializing Database...\n\n";

try {
    $pdo = getDBConnection();
    echo "✓ Connected to PostgreSQL database\n";
    
    // Create users table
    $pdo->exec("
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            discord_token VARCHAR(255) UNIQUE NOT NULL,
            discord_username VARCHAR(255),
            discord_user_id VARCHAR(255),
            avatar_url TEXT,
            status VARCHAR(50) DEFAULT 'inactive',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ");
    echo "✓ Users table created\n";
    
    // Create bot_instances table
    $pdo->exec("
        CREATE TABLE IF NOT EXISTS bot_instances (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            instance_name VARCHAR(255),
            status VARCHAR(50) DEFAULT 'stopped',
            pid INTEGER,
            started_at TIMESTAMP,
            stopped_at TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ");
    echo "✓ Bot instances table created\n";
    
    // Create admin_sessions table
    $pdo->exec("
        CREATE TABLE IF NOT EXISTS admin_sessions (
            id SERIAL PRIMARY KEY,
            session_token VARCHAR(255) UNIQUE NOT NULL,
            username VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            expires_at TIMESTAMP NOT NULL
        )
    ");
    echo "✓ Admin sessions table created\n";
    
    // Create logs table
    $pdo->exec("
        CREATE TABLE IF NOT EXISTS logs (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            action VARCHAR(255) NOT NULL,
            details TEXT,
            ip_address VARCHAR(45),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ");
    echo "✓ Logs table created\n";
    
    // Create indexes
    $pdo->exec("
        CREATE INDEX IF NOT EXISTS idx_users_discord_user_id ON users(discord_user_id);
        CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
        CREATE INDEX IF NOT EXISTS idx_bot_instances_user_id ON bot_instances(user_id);
        CREATE INDEX IF NOT EXISTS idx_bot_instances_status ON bot_instances(status);
        CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
        CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
    ");
    echo "✓ Indexes created\n";
    
    echo "\n✅ Database initialization completed successfully!\n";
    
} catch (PDOException $e) {
    echo "\n❌ Error: " . $e->getMessage() . "\n";
    exit(1);
}
?>
