// Main JavaScript for Hurry Bot Platform

document.addEventListener('DOMContentLoaded', function() {
    
    // Smooth scroll for internal links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    // Stats counter animation
    const statValues = document.querySelectorAll('.stat-value');
    let animated = false;

    function animateStats() {
        if (animated) return;
        animated = true;

        statValues.forEach(stat => {
            const text = stat.textContent;
            const hasPlus = text.includes('+');
            const hasPercent = text.includes('%');
            
            // Extract number
            let targetValue = parseFloat(text.replace(/[^0-9.]/g, ''));
            let suffix = '';
            
            if (text.includes('k')) {
                targetValue *= 1000;
                suffix = 'k';
            }
            
            if (hasPlus) suffix += '+';
            if (hasPercent) suffix = '%';

            let currentValue = 0;
            const increment = targetValue / 50;
            const duration = 1500;
            const stepTime = duration / 50;

            const counter = setInterval(() => {
                currentValue += increment;
                
                if (currentValue >= targetValue) {
                    currentValue = targetValue;
                    clearInterval(counter);
                }

                let displayValue = currentValue;
                if (text.includes('k')) {
                    displayValue = (currentValue / 1000).toFixed(1);
                } else if (hasPercent) {
                    displayValue = currentValue.toFixed(2);
                } else {
                    displayValue = Math.floor(currentValue);
                }

                stat.textContent = displayValue + suffix;
            }, stepTime);
        });
    }

    // Trigger stats animation when in viewport
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateStats();
            }
        });
    }, { threshold: 0.5 });

    const statsContainer = document.querySelector('.stats-container');
    if (statsContainer) {
        observer.observe(statsContainer);
    }

    // Add active state cycling to stats
    const statItems = document.querySelectorAll('.stat-item');
    let activeIndex = 0;

    setInterval(() => {
        statItems.forEach((item, index) => {
            if (index === activeIndex) {
                item.style.transform = 'scale(1.05)';
                item.style.transition = 'transform 0.3s ease';
            } else {
                item.style.transform = 'scale(1)';
            }
        });
        activeIndex = (activeIndex + 1) % statItems.length;
    }, 2000);

    // Button hover effect enhancement
    const buttons = document.querySelectorAll('.btn');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
        });
        
        button.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });

    // Add parallax effect to gradient orbs
    document.addEventListener('mousemove', (e) => {
        const orbs = document.querySelectorAll('.gradient-orb');
        const x = e.clientX / window.innerWidth;
        const y = e.clientY / window.innerHeight;

        orbs.forEach((orb, index) => {
            const speed = (index + 1) * 0.05;
            const xMove = (x - 0.5) * 100 * speed;
            const yMove = (y - 0.5) * 100 * speed;
            
            orb.style.transform = `translate(${xMove}px, ${yMove}px)`;
        });
    });

    // Navbar scroll effect
    let lastScroll = 0;
    const navbar = document.querySelector('.navbar');

    window.addEventListener('scroll', () => {
        const currentScroll = window.pageYOffset;

        if (currentScroll > lastScroll && currentScroll > 100) {
            // Scrolling down
            navbar.style.transform = 'translateY(-100%)';
        } else {
            // Scrolling up
            navbar.style.transform = 'translateY(0)';
        }

        lastScroll = currentScroll;
    });

    // Add glowing cursor effect
    const cursor = document.createElement('div');
    cursor.className = 'custom-cursor';
    cursor.style.cssText = `
        position: fixed;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background: radial-gradient(circle, rgba(102, 126, 234, 0.8), transparent);
        pointer-events: none;
        z-index: 9999;
        mix-blend-mode: screen;
        transition: transform 0.1s ease;
        display: none;
    `;
    document.body.appendChild(cursor);

    document.addEventListener('mousemove', (e) => {
        cursor.style.left = e.clientX - 10 + 'px';
        cursor.style.top = e.clientY - 10 + 'px';
        cursor.style.display = 'block';
    });

    // Hide cursor when leaving window
    document.addEventListener('mouseleave', () => {
        cursor.style.display = 'none';
    });

    // Add ripple effect on button clicks
    buttons.forEach(button => {
        button.addEventListener('click', function(e) {
            const ripple = document.createElement('span');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = e.clientX - rect.left - size / 2;
            const y = e.clientY - rect.top - size / 2;

            ripple.style.cssText = `
                position: absolute;
                width: ${size}px;
                height: ${size}px;
                left: ${x}px;
                top: ${y}px;
                background: rgba(255, 255, 255, 0.5);
                border-radius: 50%;
                transform: scale(0);
                animation: ripple 0.6s ease-out;
                pointer-events: none;
            `;

            this.appendChild(ripple);

            setTimeout(() => ripple.remove(), 600);
        });
    });

    // Add CSS for ripple animation
    const style = document.createElement('style');
    style.textContent = `
        @keyframes ripple {
            to {
                transform: scale(2);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);

    console.log('🚀 Hurry Bot Platform Loaded');
});
