<?php
session_start();
require_once 'config.php';
require_once 'includes/functions.php';

requireAdmin();

$stats = getStatistics();
$users = getAllUsers();
$instances = getBotInstances();
$logs = getRecentLogs(50);

// Handle AJAX requests
if (isset($_GET['action'])) {
    header('Content-Type: application/json');
    
    switch ($_GET['action']) {
        case 'start_bot':
            $userId = $_POST['user_id'] ?? 0;
            $result = startBot($userId);
            echo json_encode(['success' => $result]);
            exit;
            
        case 'stop_bot':
            $userId = $_POST['user_id'] ?? 0;
            $result = stopBot($userId);
            echo json_encode(['success' => $result]);
            exit;
            
        case 'delete_user':
            $userId = $_POST['user_id'] ?? 0;
            $result = deleteUser($userId);
            echo json_encode(['success' => $result]);
            exit;
            
        case 'get_stats':
            echo json_encode(getStatistics());
            exit;
    }
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard - Hurry Admin</title>
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/admin.css">
</head>
<body>
    <!-- Navigation -->
    <nav class="navbar">
        <div class="nav-container">
            <div class="nav-brand">
                <img src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg" alt="Hurry" class="brand-logo">
                <span class="brand-name">Hurry Admin</span>
            </div>
            <div class="nav-links">
                <div class="status-indicator">
                    <span class="status-dot active"></span>
                    <span>Connected</span>
                </div>
                <span class="user-info">👤 <?php echo htmlspecialchars($_SESSION['username']); ?></span>
                <a href="logout.php" class="btn-link">Logout</a>
            </div>
        </div>
    </nav>

    <!-- Main Dashboard -->
    <main class="dashboard-main">
        <div class="dashboard-container">
            
            <!-- Stats Grid -->
            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-icon">👥</div>
                    <div class="stat-content">
                        <div class="stat-value" id="total-users"><?php echo $stats['total_users']; ?></div>
                        <div class="stat-label">Total Users</div>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">▶️</div>
                    <div class="stat-content">
                        <div class="stat-value" id="active-bots"><?php echo $stats['active_bots']; ?></div>
                        <div class="stat-label">Active Bots</div>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">⏸️</div>
                    <div class="stat-content">
                        <div class="stat-value" id="inactive-bots"><?php echo $stats['inactive_bots']; ?></div>
                        <div class="stat-label">Inactive Bots</div>
                    </div>
                </div>
                <div class="stat-card">
                    <div class="stat-icon">⚡</div>
                    <div class="stat-content">
                        <div class="stat-value">99.99%</div>
                        <div class="stat-label">Uptime</div>
                    </div>
                </div>
            </div>

            <!-- Users Section -->
            <div class="panel-section">
                <div class="section-header">
                    <h2>User Accounts</h2>
                    <button class="btn btn-small" onclick="refreshData()">
                        <span>🔄 Refresh</span>
                    </button>
                </div>

                <div class="table-container">
                    <table class="data-table">
                        <thead>
                            <tr>
                                <th>Avatar</th>
                                <th>Username</th>
                                <th>Discord ID</th>
                                <th>Status</th>
                                <th>Last Active</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <?php if (empty($users)): ?>
                                <tr>
                                    <td colspan="6" class="text-center">No users found</td>
                                </tr>
                            <?php else: ?>
                                <?php foreach ($users as $user): ?>
                                    <tr>
                                        <td>
                                            <?php if ($user['avatar_url']): ?>
                                                <img src="<?php echo htmlspecialchars($user['avatar_url']); ?>" alt="Avatar" class="avatar">
                                            <?php else: ?>
                                                <div class="avatar avatar-placeholder"></div>
                                            <?php endif; ?>
                                        </td>
                                        <td><?php echo htmlspecialchars($user['discord_username']); ?></td>
                                        <td><?php echo htmlspecialchars($user['discord_user_id']); ?></td>
                                        <td>
                                            <span class="status-badge status-<?php echo $user['status']; ?>">
                                                <?php echo ucfirst($user['status']); ?>
                                            </span>
                                        </td>
                                        <td><?php echo date('M d, Y H:i', strtotime($user['last_active'])); ?></td>
                                        <td>
                                            <div class="action-buttons">
                                                <?php if ($user['status'] === 'active'): ?>
                                                    <button class="btn btn-small btn-danger" onclick="stopBot(<?php echo $user['id']; ?>)">Stop</button>
                                                <?php else: ?>
                                                    <button class="btn btn-small btn-success" onclick="startBot(<?php echo $user['id']; ?>)">Start</button>
                                                <?php endif; ?>
                                                <button class="btn btn-small btn-danger" onclick="deleteUser(<?php echo $user['id']; ?>)">Delete</button>
                                            </div>
                                        </td>
                                    </tr>
                                <?php endforeach; ?>
                            <?php endif; ?>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Bot Instances Section -->
            <div class="panel-section">
                <div class="section-header">
                    <h2>Bot Instances</h2>
                </div>

                <div class="table-container">
                    <table class="data-table">
                        <thead>
                            <tr>
                                <th>Instance Name</th>
                                <th>Username</th>
                                <th>Status</th>
                                <th>Started At</th>
                            </tr>
                        </thead>
                        <tbody>
                            <?php if (empty($instances)): ?>
                                <tr>
                                    <td colspan="4" class="text-center">No instances running</td>
                                </tr>
                            <?php else: ?>
                                <?php foreach ($instances as $instance): ?>
                                    <tr>
                                        <td><?php echo htmlspecialchars($instance['instance_name']); ?></td>
                                        <td><?php echo htmlspecialchars($instance['discord_username']); ?></td>
                                        <td>
                                            <span class="status-badge status-<?php echo $instance['status']; ?>">
                                                <?php echo ucfirst($instance['status']); ?>
                                            </span>
                                        </td>
                                        <td><?php echo $instance['started_at'] ? date('M d, Y H:i', strtotime($instance['started_at'])) : 'N/A'; ?></td>
                                    </tr>
                                <?php endforeach; ?>
                            <?php endif; ?>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Activity Logs -->
            <div class="panel-section">
                <div class="section-header">
                    <h2>Activity Logs</h2>
                </div>

                <div class="logs-container">
                    <?php if (empty($logs)): ?>
                        <div class="log-entry">
                            <span class="log-message">No logs available</span>
                        </div>
                    <?php else: ?>
                        <?php foreach ($logs as $log): ?>
                            <div class="log-entry">
                                <span class="log-time"><?php echo date('H:i:s', strtotime($log['created_at'])); ?></span>
                                <span class="log-message">
                                    [<?php echo htmlspecialchars($log['action']); ?>] 
                                    <?php echo htmlspecialchars($log['details']); ?>
                                    <?php if ($log['discord_username']): ?>
                                        (<?php echo htmlspecialchars($log['discord_username']); ?>)
                                    <?php endif; ?>
                                </span>
                            </div>
                        <?php endforeach; ?>
                    <?php endif; ?>
                </div>
            </div>

        </div>
    </main>

    <script src="assets/js/main.js"></script>
    <script src="assets/js/dashboard.js"></script>
</body>
</html>
