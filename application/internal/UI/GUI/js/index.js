// Handle form toggling between login and register
document.getElementById('loginToggle').addEventListener('click', function() {
    document.getElementById('loginForm').classList.remove('hidden');
    document.getElementById('registerForm').classList.add('hidden');
    this.classList.add('active');
    document.getElementById('registerToggle').classList.remove('active');
});

document.getElementById('registerToggle').addEventListener('click', function() {
    document.getElementById('registerForm').classList.remove('hidden');
    document.getElementById('loginForm').classList.add('hidden');
    this.classList.add('active');
    document.getElementById('loginToggle').classList.remove('active');
});

// Handle login submission
document.getElementById('submitLogin').addEventListener('click', function() {
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    // Fetch request for login
    fetch('/api/login', {
        method: 'POST',
        body: JSON.stringify({ username, password }),
        headers: { 'Content-Type': 'application/json' }
    })
        .then(response => response.json())
        .then(user => {
            if (user.success) {
                // Save user ID or role in localStorage
                localStorage.setItem('UserId', user.id);
                localStorage.setItem('UserRole', user.role ? 'admin' : 'user');

                // Redirect based on role
                if (user.role === 2) {
                    window.location.href = '/admin.html'; // Redirect to admin page if admin
                } else if (user.role === 1) {
                    window.location.href = '/author.html'; // Redirect to author dashboard
                } else if (user.role === 0) {
                    window.location.href = '/reader.html'
                } else {
                    alert(`Unknown role: ${user.role}`);
                }
            } else {
                alert('Login failed: ' + user.message);
            }
        })
        .catch(err => {
            console.error('Error during login:', err);
            alert('An error occurred during login.');
        });
});

// Handle registration submission
document.getElementById('submitRegister').addEventListener('click', function() {
    const fio = document.getElementById('registerFio').value;
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    // Fetch request for registration
    fetch('/api/register', {
        method: 'POST',
        body: JSON.stringify({ fio, username, password }),
        headers: { 'Content-Type': 'application/json' }
    })
        .then(response => response.json())
        .then(user => {
            if (user.success) {
                alert('Registration successful');
                // Optionally auto-login the user or redirect to login page
                document.getElementById('loginToggle').click(); // Switch to login form
            } else {
                alert('Registration failed: ' + user.message);
            }
        })
        .catch(err => {
            console.error('Error during registration:', err);
            alert('An error occurred during registration.');
        });
});
