<?php
session_start();
require_once 'config.php';
require_once 'includes/functions.php';

$message = '';
$message_type = '';

if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_POST['token'])) {
    $token = trim($_POST['token']);
    
    if (empty($token)) {
        $message = 'Please enter a Discord token';
        $message_type = 'error';
    } elseif (strlen($token) < 50) {
        $message = 'Invalid token format. Discord tokens are typically longer.';
        $message_type = 'error';
    } else {
        // Validate Discord token
        $discord_user = validateDiscordToken($token);
        
        if ($discord_user) {
            $result = addUser($token, $discord_user);
            if ($result) {
                $message = "Successfully added account: " . htmlspecialchars($discord_user['username']);
                $message_type = 'success';
            } else {
                $message = 'Failed to add account. Please try again.';
                $message_type = 'error';
            }
        } else {
            $message = 'Invalid Discord token. Please check and try again.';
            $message_type = 'error';
        }
    }
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Host Bot - Hurry</title>
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/host.css">
</head>
<body>
    <!-- Navigation -->
    <nav class="navbar">
        <div class="nav-container">
            <div class="nav-brand">
                <a href="index.php" style="display: flex; align-items: center; gap: 0.75rem; text-decoration: none;">
                    <img src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg" alt="Hurry" class="brand-logo">
                    <span class="brand-name">Hurry</span>
                </a>
            </div>
            <div class="nav-links">
                <a href="index.php" class="btn-link">← Back to Home</a>
            </div>
        </div>
    </nav>

    <!-- Main Content -->
    <main class="host-main">
        <div class="host-container">
            <div class="host-header">
                <h1 class="gradient-text">Host Your Bot</h1>
                <p>Add your Discord token to start hosting with Hurry</p>
            </div>

            <?php if ($message): ?>
                <div class="alert alert-<?php echo $message_type; ?>">
                    <?php echo htmlspecialchars($message); ?>
                </div>
            <?php endif; ?>

            <div class="content-grid">
                <!-- Token Form -->
                <div class="card card-primary">
                    <h2>Add Discord Account</h2>
                    <form method="POST" action="" class="token-form">
                        <div class="form-group">
                            <label for="token">Discord Token</label>
                            <input 
                                type="password" 
                                id="token" 
                                name="token" 
                                class="form-input" 
                                placeholder="Enter your Discord token..."
                                required
                            >
                            <small class="form-help">Your token is encrypted and stored securely</small>
                        </div>
                        <button type="submit" class="btn btn-primary btn-block">
                            <span>Add Account</span>
                            <div class="btn-border"></div>
                        </button>
                    </form>
                </div>

                <!-- Info Card -->
                <div class="card card-info">
                    <div class="info-icon">ℹ️</div>
                    <h3>How to get your Discord token</h3>
                    <ol class="info-list">
                        <li>Open Discord in your browser (Chrome/Firefox)</li>
                        <li>Press <kbd>F12</kbd> to open Developer Tools</li>
                        <li>Go to the <strong>Console</strong> tab</li>
                        <li>Paste this code and press Enter:</li>
                    </ol>
                    <div class="code-block">
                        <code>window.webpackChunkdiscord_app.push([[''],{},e=>{m=[];for(let c in e.c)m.push(e.c[c])}]);m.find(m=>m?.exports?.default?.getToken!==void 0).exports.default.getToken()</code>
                        <button class="copy-btn" onclick="copyCode()">Copy</button>
                    </div>
                    <p class="warning">⚠️ Never share your token with anyone else!</p>
                </div>
            </div>
        </div>

        <!-- Background Animation -->
        <div class="bg-animation">
            <div class="gradient-orb orb-1"></div>
            <div class="gradient-orb orb-2"></div>
        </div>
    </main>

    <script src="assets/js/main.js"></script>
    <script>
        function copyCode() {
            const code = document.querySelector('.code-block code').textContent;
            navigator.clipboard.writeText(code).then(() => {
                const btn = document.querySelector('.copy-btn');
                btn.textContent = 'Copied!';
                setTimeout(() => {
                    btn.textContent = 'Copy';
                }, 2000);
            });
        }
    </script>
</body>
</html>
