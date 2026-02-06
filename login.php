<?php
session_start();
require_once 'config.php';

$error = '';

// Redirect if already logged in
if (isset($_SESSION['admin']) && $_SESSION['admin'] === true) {
    header('Location: dashboard.php');
    exit;
}

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'] ?? '';
    $password = $_POST['password'] ?? '';
    
    if ($username === ADMIN_USERNAME && $password === ADMIN_PASSWORD) {
        $_SESSION['admin'] = true;
        $_SESSION['username'] = $username;
        header('Location: dashboard.php');
        exit;
    } else {
        $error = 'Invalid credentials';
    }
}
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Login - Hurry</title>
    <link rel="stylesheet" href="assets/css/style.css">
    <link rel="stylesheet" href="assets/css/admin.css">
</head>
<body>
    <div class="login-page">
        <div class="login-container">
            <div class="login-header">
                <div class="logo-small">
                    <img src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg" alt="Hurry">
                </div>
                <h2>Admin Login</h2>
                <p>Access Control Panel</p>
            </div>

            <?php if ($error): ?>
                <div class="alert alert-error">
                    <?php echo htmlspecialchars($error); ?>
                </div>
            <?php endif; ?>

            <form method="POST" action="" class="login-form">
                <div class="form-group">
                    <label for="username">Username</label>
                    <input 
                        type="text" 
                        id="username" 
                        name="username" 
                        class="form-input" 
                        required
                        autocomplete="username"
                    >
                </div>

                <div class="form-group">
                    <label for="password">Password</label>
                    <input 
                        type="password" 
                        id="password" 
                        name="password" 
                        class="form-input" 
                        required
                        autocomplete="current-password"
                    >
                </div>

                <button type="submit" class="btn btn-primary btn-block">
                    <span>Login</span>
                    <div class="btn-border"></div>
                </button>
            </form>

            <div class="login-footer">
                <a href="index.php" class="btn-link">← Back to Home</a>
            </div>
        </div>

        <!-- Background Animation -->
        <div class="bg-animation">
            <div class="gradient-orb orb-1"></div>
            <div class="gradient-orb orb-2"></div>
        </div>
    </div>

    <script src="assets/js/main.js"></script>
</body>
</html>
