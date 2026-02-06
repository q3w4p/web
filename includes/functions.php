<?php
require_once 'config.php';

/**
 * Validate Discord Token
 */
function validateDiscordToken($token) {
    $ch = curl_init();
    
    curl_setopt_array($ch, [
        CURLOPT_URL => DISCORD_API_BASE . '/users/@me',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_HTTPHEADER => [
            'Authorization: ' . $token,
            'Content-Type: application/json'
        ],
        CURLOPT_SSL_VERIFYPEER => true,
    ]);
    
    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    
    if ($httpCode === 200) {
        return json_decode($response, true);
    }
    
    return false;
}

/**
 * Add User to Database
 */
function addUser($token, $discordUser) {
    try {
        $pdo = getDBConnection();
        
        $avatarUrl = null;
        if ($discordUser['avatar']) {
            $avatarUrl = sprintf(
                'https://cdn.discordapp.com/avatars/%s/%s.png',
                $discordUser['id'],
                $discordUser['avatar']
            );
        }
        
        // Check if user already exists
        $stmt = $pdo->prepare("
            SELECT id FROM users WHERE discord_user_id = ?
        ");
        $stmt->execute([$discordUser['id']]);
        $existingUser = $stmt->fetch();
        
        if ($existingUser) {
            // Update existing user
            $stmt = $pdo->prepare("
                UPDATE users 
                SET discord_token = ?, 
                    discord_username = ?, 
                    avatar_url = ?, 
                    last_active = CURRENT_TIMESTAMP
                WHERE discord_user_id = ?
                RETURNING id
            ");
            $stmt->execute([
                $token,
                $discordUser['username'],
                $avatarUrl,
                $discordUser['id']
            ]);
            $userId = $stmt->fetchColumn();
        } else {
            // Insert new user
            $stmt = $pdo->prepare("
                INSERT INTO users (discord_token, discord_username, discord_user_id, avatar_url)
                VALUES (?, ?, ?, ?)
                RETURNING id
            ");
            $stmt->execute([
                $token,
                $discordUser['username'],
                $discordUser['id'],
                $avatarUrl
            ]);
            $userId = $stmt->fetchColumn();
        }
        
        // Log action
        logAction($userId, 'user_added', "User {$discordUser['username']} added/updated");
        
        return $userId;
    } catch (PDOException $e) {
        error_log("Error adding user: " . $e->getMessage());
        return false;
    }
}

/**
 * Get All Users
 */
function getAllUsers() {
    try {
        $pdo = getDBConnection();
        $stmt = $pdo->query("
            SELECT id, discord_username, discord_user_id, avatar_url, status, 
                   created_at, last_active
            FROM users
            ORDER BY last_active DESC
        ");
        return $stmt->fetchAll();
    } catch (PDOException $e) {
        error_log("Error fetching users: " . $e->getMessage());
        return [];
    }
}

/**
 * Get User by ID
 */
function getUserById($userId) {
    try {
        $pdo = getDBConnection();
        $stmt = $pdo->prepare("SELECT * FROM users WHERE id = ?");
        $stmt->execute([$userId]);
        return $stmt->fetch();
    } catch (PDOException $e) {
        error_log("Error fetching user: " . $e->getMessage());
        return null;
    }
}

/**
 * Delete User
 */
function deleteUser($userId) {
    try {
        $pdo = getDBConnection();
        $stmt = $pdo->prepare("DELETE FROM users WHERE id = ?");
        $stmt->execute([$userId]);
        return true;
    } catch (PDOException $e) {
        error_log("Error deleting user: " . $e->getMessage());
        return false;
    }
}

/**
 * Start Bot Instance
 */
function startBot($userId) {
    try {
        $pdo = getDBConnection();
        $user = getUserById($userId);
        
        if (!$user) {
            return false;
        }
        
        // TODO: Call your Go bot API here
        // Example:
        // $ch = curl_init();
        // curl_setopt_array($ch, [
        //     CURLOPT_URL => 'http://' . GO_BOT_HOST . ':' . GO_BOT_PORT . '/start',
        //     CURLOPT_POST => true,
        //     CURLOPT_POSTFIELDS => json_encode([
        //         'token' => $user['discord_token'],
        //         'userId' => $userId
        //     ]),
        //     CURLOPT_HTTPHEADER => [
        //         'X-API-Key: ' . GO_BOT_API_KEY,
        //         'Content-Type: application/json'
        //     ],
        //     CURLOPT_RETURNTRANSFER => true,
        // ]);
        // $response = curl_exec($ch);
        // curl_close($ch);
        
        // Create bot instance record
        $stmt = $pdo->prepare("
            INSERT INTO bot_instances (user_id, instance_name, status, started_at)
            VALUES (?, ?, 'running', CURRENT_TIMESTAMP)
        ");
        $stmt->execute([$userId, 'bot_' . $user['discord_username']]);
        
        // Update user status
        $stmt = $pdo->prepare("
            UPDATE users SET status = 'active', last_active = CURRENT_TIMESTAMP
            WHERE id = ?
        ");
        $stmt->execute([$userId]);
        
        logAction($userId, 'bot_started', "Bot started for {$user['discord_username']}");
        
        return true;
    } catch (PDOException $e) {
        error_log("Error starting bot: " . $e->getMessage());
        return false;
    }
}

/**
 * Stop Bot Instance
 */
function stopBot($userId) {
    try {
        $pdo = getDBConnection();
        $user = getUserById($userId);
        
        if (!$user) {
            return false;
        }
        
        // TODO: Call your Go bot API here
        
        // Update bot instance
        $stmt = $pdo->prepare("
            UPDATE bot_instances 
            SET status = 'stopped', stopped_at = CURRENT_TIMESTAMP
            WHERE user_id = ? AND status = 'running'
        ");
        $stmt->execute([$userId]);
        
        // Update user status
        $stmt = $pdo->prepare("
            UPDATE users SET status = 'inactive' WHERE id = ?
        ");
        $stmt->execute([$userId]);
        
        logAction($userId, 'bot_stopped', "Bot stopped for {$user['discord_username']}");
        
        return true;
    } catch (PDOException $e) {
        error_log("Error stopping bot: " . $e->getMessage());
        return false;
    }
}

/**
 * Get Bot Instances
 */
function getBotInstances() {
    try {
        $pdo = getDBConnection();
        $stmt = $pdo->query("
            SELECT bi.*, u.discord_username, u.discord_user_id
            FROM bot_instances bi
            JOIN users u ON bi.user_id = u.id
            ORDER BY bi.started_at DESC
            LIMIT 100
        ");
        return $stmt->fetchAll();
    } catch (PDOException $e) {
        error_log("Error fetching instances: " . $e->getMessage());
        return [];
    }
}

/**
 * Log Action
 */
function logAction($userId, $action, $details, $ipAddress = null) {
    try {
        $pdo = getDBConnection();
        $ipAddress = $ipAddress ?? $_SERVER['REMOTE_ADDR'] ?? 'unknown';
        
        $stmt = $pdo->prepare("
            INSERT INTO logs (user_id, action, details, ip_address)
            VALUES (?, ?, ?, ?)
        ");
        $stmt->execute([$userId, $action, $details, $ipAddress]);
        
        return true;
    } catch (PDOException $e) {
        error_log("Error logging action: " . $e->getMessage());
        return false;
    }
}

/**
 * Get Recent Logs
 */
function getRecentLogs($limit = 100) {
    try {
        $pdo = getDBConnection();
        $stmt = $pdo->prepare("
            SELECT l.*, u.discord_username
            FROM logs l
            LEFT JOIN users u ON l.user_id = u.id
            ORDER BY l.created_at DESC
            LIMIT ?
        ");
        $stmt->execute([$limit]);
        return $stmt->fetchAll();
    } catch (PDOException $e) {
        error_log("Error fetching logs: " . $e->getMessage());
        return [];
    }
}

/**
 * Get Statistics
 */
function getStatistics() {
    try {
        $pdo = getDBConnection();
        
        $stats = [
            'total_users' => 0,
            'active_bots' => 0,
            'inactive_bots' => 0,
        ];
        
        $stmt = $pdo->query("SELECT COUNT(*) as total FROM users");
        $stats['total_users'] = $stmt->fetchColumn();
        
        $stmt = $pdo->query("SELECT COUNT(*) as active FROM users WHERE status = 'active'");
        $stats['active_bots'] = $stmt->fetchColumn();
        
        $stmt = $pdo->query("SELECT COUNT(*) as inactive FROM users WHERE status = 'inactive'");
        $stats['inactive_bots'] = $stmt->fetchColumn();
        
        return $stats;
    } catch (PDOException $e) {
        error_log("Error fetching statistics: " . $e->getMessage());
        return ['total_users' => 0, 'active_bots' => 0, 'inactive_bots' => 0];
    }
}

/**
 * Check if user is admin
 */
function isAdmin() {
    return isset($_SESSION['admin']) && $_SESSION['admin'] === true;
}

/**
 * Require admin authentication
 */
function requireAdmin() {
    if (!isAdmin()) {
        header('Location: login.php');
        exit;
    }
}
?>
