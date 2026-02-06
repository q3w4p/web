require('dotenv').config();
const { Pool } = require('pg');

const pool = new Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_NAME,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  ssl: process.env.DB_SSL === 'true' ? { rejectUnauthorized: false } : false
});

async function initDatabase() {
  const client = await pool.connect();
  
  try {
    console.log('Connected to Aiven PostgreSQL database');
    
    // Create users table
    await client.query(`
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
    `);
    console.log('✓ Users table created');

    // Create bot_instances table
    await client.query(`
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
    `);
    console.log('✓ Bot instances table created');

    // Create admin_sessions table
    await client.query(`
      CREATE TABLE IF NOT EXISTS admin_sessions (
        id SERIAL PRIMARY KEY,
        session_token VARCHAR(255) UNIQUE NOT NULL,
        username VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL
      )
    `);
    console.log('✓ Admin sessions table created');

    // Create logs table
    await client.query(`
      CREATE TABLE IF NOT EXISTS logs (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
        action VARCHAR(255) NOT NULL,
        details TEXT,
        ip_address VARCHAR(45),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
      )
    `);
    console.log('✓ Logs table created');

    // Create indexes
    await client.query(`
      CREATE INDEX IF NOT EXISTS idx_users_discord_user_id ON users(discord_user_id);
      CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
      CREATE INDEX IF NOT EXISTS idx_bot_instances_user_id ON bot_instances(user_id);
      CREATE INDEX IF NOT EXISTS idx_bot_instances_status ON bot_instances(status);
      CREATE INDEX IF NOT EXISTS idx_logs_user_id ON logs(user_id);
      CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs(created_at);
    `);
    console.log('✓ Indexes created');

    console.log('\n✅ Database initialization completed successfully!');
  } catch (error) {
    console.error('❌ Error initializing database:', error);
    throw error;
  } finally {
    client.release();
    await pool.end();
  }
}

initDatabase().catch(console.error);
