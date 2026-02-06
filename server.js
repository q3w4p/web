require('dotenv').config();
const express = require('express');
const http = require('http');
const WebSocket = require('ws');
const { Pool } = require('pg');
const bodyParser = require('body-parser');
const session = require('express-session');
const axios = require('axios');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const path = require('path');

const app = express();
const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

// PostgreSQL connection pool
const pool = new Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_NAME,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  ssl: process.env.DB_SSL === 'true' ? { rejectUnauthorized: false } : false
});

// Middleware
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(session({
  secret: process.env.SESSION_SECRET,
  resave: false,
  saveUninitialized: false,
  cookie: { secure: false, maxAge: 24 * 60 * 60 * 1000 } // 24 hours
}));

// Serve static files
app.use(express.static('public'));

// Store active WebSocket connections
const clients = new Map();

// WebSocket connection handler
wss.on('connection', (ws) => {
  const clientId = Math.random().toString(36).substring(7);
  clients.set(clientId, ws);
  
  console.log(`Client connected: ${clientId}`);
  
  ws.on('message', async (message) => {
    try {
      const data = JSON.parse(message);
      console.log('Received:', data);
      
      if (data.type === 'ping') {
        ws.send(JSON.stringify({ type: 'pong' }));
      }
    } catch (error) {
      console.error('WebSocket message error:', error);
    }
  });
  
  ws.on('close', () => {
    clients.delete(clientId);
    console.log(`Client disconnected: ${clientId}`);
  });
});

// Broadcast to all connected clients
function broadcast(data) {
  clients.forEach((client) => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify(data));
    }
  });
}

// Validate Discord token
async function validateDiscordToken(token) {
  try {
    const response = await axios.get(`${process.env.DISCORD_API_BASE}/users/@me`, {
      headers: {
        'Authorization': token
      }
    });
    return response.data;
  } catch (error) {
    console.error('Discord token validation error:', error.message);
    return null;
  }
}

// Middleware to check admin authentication
function requireAdmin(req, res, next) {
  if (req.session.isAdmin) {
    next();
  } else {
    res.status(401).json({ error: 'Unauthorized' });
  }
}

// API Routes

// Home route
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'public', 'index.html'));
});

// Admin login
app.post('/api/admin/login', async (req, res) => {
  const { username, password } = req.body;
  
  if (username === process.env.ADMIN_USERNAME && password === process.env.ADMIN_PASSWORD) {
    req.session.isAdmin = true;
    req.session.username = username;
    res.json({ success: true, message: 'Login successful' });
  } else {
    res.status(401).json({ error: 'Invalid credentials' });
  }
});

// Admin logout
app.post('/api/admin/logout', (req, res) => {
  req.session.destroy();
  res.json({ success: true });
});

// Check admin status
app.get('/api/admin/check', (req, res) => {
  res.json({ isAdmin: req.session.isAdmin || false });
});

// Add/Update Discord user
app.post('/api/users/add', async (req, res) => {
  const { token } = req.body;
  
  if (!token) {
    return res.status(400).json({ error: 'Discord token is required' });
  }
  
  try {
    // Validate Discord token
    const discordUser = await validateDiscordToken(token);
    
    if (!discordUser) {
      return res.status(400).json({ error: 'Invalid Discord token' });
    }
    
    const client = await pool.connect();
    
    try {
      // Check if user exists
      const existingUser = await client.query(
        'SELECT * FROM users WHERE discord_user_id = $1',
        [discordUser.id]
      );
      
      let userId;
      
      if (existingUser.rows.length > 0) {
        // Update existing user
        const result = await client.query(
          `UPDATE users 
           SET discord_token = $1, discord_username = $2, avatar_url = $3, last_active = CURRENT_TIMESTAMP
           WHERE discord_user_id = $4
           RETURNING id`,
          [token, discordUser.username, discordUser.avatar ? `https://cdn.discordapp.com/avatars/${discordUser.id}/${discordUser.avatar}.png` : null, discordUser.id]
        );
        userId = result.rows[0].id;
      } else {
        // Insert new user
        const result = await client.query(
          `INSERT INTO users (discord_token, discord_username, discord_user_id, avatar_url)
           VALUES ($1, $2, $3, $4)
           RETURNING id`,
          [token, discordUser.username, discordUser.id, discordUser.avatar ? `https://cdn.discordapp.com/avatars/${discordUser.id}/${discordUser.avatar}.png` : null]
        );
        userId = result.rows[0].id;
      }
      
      // Log action
      await client.query(
        'INSERT INTO logs (user_id, action, details, ip_address) VALUES ($1, $2, $3, $4)',
        [userId, 'user_added', `User ${discordUser.username} added/updated`, req.ip]
      );
      
      // Broadcast update to all connected clients
      broadcast({
        type: 'user_update',
        data: {
          id: userId,
          username: discordUser.username,
          discord_id: discordUser.id,
          avatar: discordUser.avatar ? `https://cdn.discordapp.com/avatars/${discordUser.id}/${discordUser.avatar}.png` : null,
          status: 'active'
        }
      });
      
      res.json({
        success: true,
        user: {
          id: userId,
          username: discordUser.username,
          discord_id: discordUser.id
        }
      });
    } finally {
      client.release();
    }
  } catch (error) {
    console.error('Error adding user:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Get all users
app.get('/api/users', requireAdmin, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT id, discord_username, discord_user_id, avatar_url, status, created_at, last_active
       FROM users
       ORDER BY last_active DESC`
    );
    
    res.json({ users: result.rows });
  } catch (error) {
    console.error('Error fetching users:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Start bot instance
app.post('/api/bot/start', requireAdmin, async (req, res) => {
  const { userId } = req.body;
  
  try {
    const client = await pool.connect();
    
    try {
      // Get user
      const userResult = await pool.query('SELECT * FROM users WHERE id = $1', [userId]);
      
      if (userResult.rows.length === 0) {
        return res.status(404).json({ error: 'User not found' });
      }
      
      const user = userResult.rows[0];
      
      // TODO: Call your Go bot API to start the bot with the user's token
      // Example:
      // await axios.post(`http://${process.env.GO_BOT_HOST}:${process.env.GO_BOT_PORT}/start`, {
      //   token: user.discord_token,
      //   userId: userId
      // }, {
      //   headers: { 'X-API-Key': process.env.GO_BOT_API_KEY }
      // });
      
      // Create bot instance record
      await client.query(
        `INSERT INTO bot_instances (user_id, instance_name, status, started_at)
         VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`,
        [userId, `bot_${user.discord_username}`, 'running']
      );
      
      // Update user status
      await client.query(
        'UPDATE users SET status = $1, last_active = CURRENT_TIMESTAMP WHERE id = $2',
        ['active', userId]
      );
      
      // Log action
      await client.query(
        'INSERT INTO logs (user_id, action, details, ip_address) VALUES ($1, $2, $3, $4)',
        [userId, 'bot_started', `Bot started for ${user.discord_username}`, req.ip]
      );
      
      broadcast({
        type: 'bot_status',
        data: { userId, status: 'running', username: user.discord_username }
      });
      
      res.json({ success: true, message: 'Bot started successfully' });
    } finally {
      client.release();
    }
  } catch (error) {
    console.error('Error starting bot:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Stop bot instance
app.post('/api/bot/stop', requireAdmin, async (req, res) => {
  const { userId } = req.body;
  
  try {
    const client = await pool.connect();
    
    try {
      const userResult = await pool.query('SELECT * FROM users WHERE id = $1', [userId]);
      
      if (userResult.rows.length === 0) {
        return res.status(404).json({ error: 'User not found' });
      }
      
      const user = userResult.rows[0];
      
      // TODO: Call your Go bot API to stop the bot
      
      // Update bot instance
      await client.query(
        `UPDATE bot_instances 
         SET status = $1, stopped_at = CURRENT_TIMESTAMP
         WHERE user_id = $2 AND status = $3`,
        ['stopped', userId, 'running']
      );
      
      // Update user status
      await client.query(
        'UPDATE users SET status = $1 WHERE id = $2',
        ['inactive', userId]
      );
      
      // Log action
      await client.query(
        'INSERT INTO logs (user_id, action, details, ip_address) VALUES ($1, $2, $3, $4)',
        [userId, 'bot_stopped', `Bot stopped for ${user.discord_username}`, req.ip]
      );
      
      broadcast({
        type: 'bot_status',
        data: { userId, status: 'stopped', username: user.discord_username }
      });
      
      res.json({ success: true, message: 'Bot stopped successfully' });
    } finally {
      client.release();
    }
  } catch (error) {
    console.error('Error stopping bot:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Get bot instances
app.get('/api/bot/instances', requireAdmin, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT bi.*, u.discord_username, u.discord_user_id
       FROM bot_instances bi
       JOIN users u ON bi.user_id = u.id
       ORDER BY bi.started_at DESC
       LIMIT 100`
    );
    
    res.json({ instances: result.rows });
  } catch (error) {
    console.error('Error fetching instances:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Delete user
app.delete('/api/users/:id', requireAdmin, async (req, res) => {
  const { id } = req.params;
  
  try {
    await pool.query('DELETE FROM users WHERE id = $1', [id]);
    
    broadcast({
      type: 'user_deleted',
      data: { userId: id }
    });
    
    res.json({ success: true });
  } catch (error) {
    console.error('Error deleting user:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Get logs
app.get('/api/logs', requireAdmin, async (req, res) => {
  try {
    const result = await pool.query(
      `SELECT l.*, u.discord_username
       FROM logs l
       LEFT JOIN users u ON l.user_id = u.id
       ORDER BY l.created_at DESC
       LIMIT 100`
    );
    
    res.json({ logs: result.rows });
  } catch (error) {
    console.error('Error fetching logs:', error);
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Health check
app.get('/api/health', (req, res) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

// Start server
const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log(`🚀 Server running on port ${PORT}`);
  console.log(`📡 WebSocket server ready`);
  console.log(`🔗 Admin panel: http://localhost:${PORT}/admin.html`);
});
