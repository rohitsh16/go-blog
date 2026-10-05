(function () {
    'use strict';

    function getSavedTheme() {
        return localStorage.getItem('theme');
    }

    function getSystemTheme() {
        return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }

    function getCurrentTheme() {
        return getSavedTheme() || getSystemTheme();
    }

    function applyTheme(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        updateToggleButtons(theme);
    }

    function updateToggleButtons(theme) {
        var buttons = document.querySelectorAll('.theme-toggle, #theme-toggle');
        buttons.forEach(function (btn) {
            if (theme === 'dark') {
                btn.innerHTML = '<span class="theme-icon">☀️</span><span class="theme-label">Light</span>';
                btn.setAttribute('aria-label', 'Switch to light mode');
                btn.setAttribute('title', 'Switch to light mode');
            } else {
                btn.innerHTML = '<span class="theme-icon">🌙</span><span class="theme-label">Dark</span>';
                btn.setAttribute('aria-label', 'Switch to dark mode');
                btn.setAttribute('title', 'Switch to dark mode');
            }
        });
    }

    function toggleTheme() {
        var current = document.documentElement.getAttribute('data-theme') || getCurrentTheme();
        var next = current === 'dark' ? 'light' : 'dark';
        localStorage.setItem('theme', next);
        applyTheme(next);
    }

    // Attach click events once DOM is ready
    function init() {
        applyTheme(getCurrentTheme());

        document.addEventListener('click', function (e) {
            var target = e.target.closest('.theme-toggle, #theme-toggle');
            if (target) {
                e.preventDefault();
                toggleTheme();
            }
        });

        // Listen for OS system theme changes if user hasn't manually overridden
        if (window.matchMedia) {
            window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function (e) {
                if (!getSavedTheme()) {
                    applyTheme(e.matches ? 'dark' : 'light');
                }
            });
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // Expose toggleTheme globally if needed
    window.__toggleTheme = toggleTheme;
})();
