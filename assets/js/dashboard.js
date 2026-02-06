// Dashboard functionality

// Start Bot
async function startBot(userId) {
    if (!confirm('Are you sure you want to start this bot?')) {
        return;
    }

    try {
        const formData = new FormData();
        formData.append('user_id', userId);

        const response = await fetch('dashboard.php?action=start_bot', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (data.success) {
            alert('Bot started successfully!');
            location.reload();
        } else {
            alert('Failed to start bot. Please try again.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    }
}

// Stop Bot
async function stopBot(userId) {
    if (!confirm('Are you sure you want to stop this bot?')) {
        return;
    }

    try {
        const formData = new FormData();
        formData.append('user_id', userId);

        const response = await fetch('dashboard.php?action=stop_bot', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (data.success) {
            alert('Bot stopped successfully!');
            location.reload();
        } else {
            alert('Failed to stop bot. Please try again.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    }
}

// Delete User
async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
        return;
    }

    try {
        const formData = new FormData();
        formData.append('user_id', userId);

        const response = await fetch('dashboard.php?action=delete_user', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (data.success) {
            alert('User deleted successfully!');
            location.reload();
        } else {
            alert('Failed to delete user. Please try again.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred. Please try again.');
    }
}

// Refresh Data
function refreshData() {
    location.reload();
}

// Auto-refresh statistics every 30 seconds
setInterval(async () => {
    try {
        const response = await fetch('dashboard.php?action=get_stats');
        const stats = await response.json();

        document.getElementById('total-users').textContent = stats.total_users;
        document.getElementById('active-bots').textContent = stats.active_bots;
        document.getElementById('inactive-bots').textContent = stats.inactive_bots;
    } catch (error) {
        console.error('Error refreshing stats:', error);
    }
}, 30000);

// Auto-scroll logs to bottom on page load
document.addEventListener('DOMContentLoaded', () => {
    const logsContainer = document.querySelector('.logs-container');
    if (logsContainer) {
        logsContainer.scrollTop = logsContainer.scrollHeight;
    }

    console.log('🎯 Dashboard loaded successfully');
});
