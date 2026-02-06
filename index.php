<?php
session_start();
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hurry - Discord Bot Hosting</title>
    <link rel="stylesheet" href="assets/css/style.css">
</head>
<body>
    <!-- Navigation -->
    <nav class="navbar">
        <div class="nav-container">
            <div class="nav-brand">
                <img src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg" alt="Hurry" class="brand-logo">
                <span class="brand-name">Hurry</span>
            </div>
            <div class="nav-links">
                <?php if(isset($_SESSION['user'])): ?>
                    <a href="dashboard.php" class="btn-link">Dashboard</a>
                    <a href="logout.php" class="btn-link">Logout</a>
                <?php else: ?>
                    <a href="login.php" class="btn-link">Login</a>
                <?php endif; ?>
            </div>
        </div>
    </nav>

    <!-- Hero Section -->
    <main class="hero-section">
        <div class="hero-container">
            <div class="hero-content">
                <!-- Logo -->
                <div class="hero-logo">
                    <div class="logo-wrapper">
                        <img src="https://cdn.discordapp.com/attachments/1465030716118139095/1469105444948807882/1.jpeg" alt="Hurry Logo" class="logo-image">
                    </div>
                </div>

                <!-- Title with Gradient -->
                <h1 class="hero-title">
                    <span class="gradient-text">Hurry</span>
                </h1>

                <!-- CTA Buttons -->
                <div class="cta-container">
                    <?php if(isset($_SESSION['user'])): ?>
                        <a href="dashboard.php" class="btn btn-primary">
                            <span>Dashboard</span>
                            <div class="btn-border"></div>
                        </a>
                    <?php else: ?>
                        <a href="host.php" class="btn btn-primary">
                            <span>Get Started</span>
                            <div class="btn-border"></div>
                        </a>
                        <a href="login.php" class="btn btn-secondary">
                            <span>Admin Panel</span>
                            <div class="btn-border"></div>
                        </a>
                    <?php endif; ?>
                </div>

                <!-- Stats -->
                <div class="stats-container">
                    <div class="stat-item">
                        <div class="stat-value">2,400+</div>
                        <div class="stat-label">Active Bots</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-value">150k+</div>
                        <div class="stat-label">Total Users</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-value">99.99%</div>
                        <div class="stat-label">Uptime</div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Animated Background -->
        <div class="bg-animation">
            <div class="gradient-orb orb-1"></div>
            <div class="gradient-orb orb-2"></div>
            <div class="gradient-orb orb-3"></div>
        </div>
    </main>

    <!-- Footer -->
    <footer class="footer">
        <div class="footer-container">
            <p>&copy; <?php echo date('Y'); ?> Hurry. All rights reserved.</p>
        </div>
    </footer>

    <script src="assets/js/main.js"></script>
</body>
</html>
