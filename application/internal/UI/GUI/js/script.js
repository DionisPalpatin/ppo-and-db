// document.addEventListener("DOMContentLoaded", function() {
//     const loginBtn = document.getElementById('login-btn');
//     const errorMessage = document.getElementById('error-message');
//     const authForm = document.getElementById('auth-form');
//     const readerMenu = document.getElementById('reader-menu');
//     const output = document.getElementById('output');
//
//     // Авторизация пользователя
//     loginBtn.addEventListener('click', async function() {
//         const login = document.getElementById('login').value;
//         const password = document.getElementById('password').value;
//
//         // Пример авторизации через API
//         const response = await fetch('/api/login', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json',
//             },
//             body: JSON.stringify({ login, password })
//         });
//
//         if (response.ok) {
//             const data = await response.json();
//             // Если авторизация успешна, показываем меню для Читателя
//             authForm.classList.add('hidden');
//             readerMenu.classList.remove('hidden');
//         } else {
//             errorMessage.textContent = 'Invalid login or password';
//         }
//     });
//
//     // Функция для отображения данных
//     function displayList(items, type) {
//         output.innerHTML = ''; // Очистка контейнера перед выводом новых данных
//
//         const ul = document.createElement('ul');
//         items.forEach(item => {
//             const li = document.createElement('li');
//             li.textContent = `${item.Id}) ${item.Name}`;
//             ul.appendChild(li);
//         });
//         output.appendChild(ul);
//     }
//
//     // Получение всех общедоступных заметок
//     document.getElementById('view-public-notes').addEventListener('click', async function() {
//         const response = await fetch('/api/notes/public', { method: 'GET' });
//         if (response.ok) {
//             const notes = await response.json();
//             displayList(notes, 'note');
//         } else {
//             console.error('Failed to fetch notes');
//         }
//     });
//
//     // Получение всех коллекций пользователя
//     document.getElementById('view-collections').addEventListener('click', async function() {
//         const response = await fetch('/api/collections', { method: 'GET' });
//         if (response.ok) {
//             const collections = await response.json();
//             displayList(collections, 'collection');
//         } else {
//             console.error('Failed to fetch collections');
//         }
//     });
//
//     // Пример выхода из аккаунта
//     document.getElementById('logout').addEventListener('click', function() {
//         // Обработать логику выхода, например очистку токена
//         window.location.reload(); // Перезагружаем страницу для возврата к форме логина
//     });
// });


document.addEventListener("DOMContentLoaded", function() {
    const loginBtn = document.getElementById('login-btn');
    const registerBtn = document.getElementById('register-btn');
    const loginErrorMessage = document.getElementById('login-error-message');
    const registerErrorMessage = document.getElementById('register-error-message');
    const authForm = document.getElementById('auth-form');
    const registerForm = document.getElementById('register-form');
    const readerMenu = document.getElementById('reader-menu');
    const output = document.getElementById('output');
    const showRegisterFormLink = document.getElementById('show-register-form');
    const showLoginFormLink = document.getElementById('show-login-form');

    showRegisterFormLink.addEventListener('click', function(event) {
        event.preventDefault();
        authForm.classList.add('hidden');
        registerForm.classList.remove('hidden');
    });

    showLoginFormLink.addEventListener('click', function(event) {
        event.preventDefault();
        registerForm.classList.add('hidden');
        authForm.classList.remove('hidden');
    });

    // Авторизация пользователя
    loginBtn.addEventListener('click', async function() {
        const login = document.getElementById('login').value;
        const password = document.getElementById('password').value;

        // Запрос на сервер для входа
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ login, password })
        });

        if (response.ok) {
            const data = await response.json();
            // Проверяем роль пользователя (1 - Читатель)
            if (data.role === 1) {
                // Если авторизация успешна и роль - Читатель, показываем меню
                authForm.classList.add('hidden');
                registerForm.classList.add('hidden');
                readerMenu.classList.remove('hidden');
            } else {
                // Обработка для других ролей, если необходимо
                alert('Эта роль еще не реализована в интерфейсе.');
            }
        } else {
            loginErrorMessage.textContent = 'Неверный логин или пароль';
        }
    });

    // Регистрация пользователя
    registerBtn.addEventListener('click', async function() {
        const fio = document.getElementById('register-fio').value;
        const login = document.getElementById('register-login').value;
        const password = document.getElementById('register-password').value;

        // Запрос на сервер для регистрации
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ fio, login, password })
        });

        if (response.ok) {
            const data = await response.json();
            alert('Регистрация прошла успешно! Теперь вы можете войти в систему.');
            // Очищаем поля регистрации
            document.getElementById('register-fio').value = '';
            document.getElementById('register-login').value = '';
            document.getElementById('register-password').value = '';
        } else {
            const errorData = await response.json();
            registerErrorMessage.textContent = errorData.message || 'Ошибка регистрации';
        }
    });

    // Функция для отображения данных
    function displayList(items, type) {
        output.innerHTML = ''; // Очистка контейнера перед выводом новых данных

        const ul = document.createElement('ul');
        items.forEach(item => {
            const li = document.createElement('li');
            li.textContent = `${item.Id}) ${item.Name}`;
            ul.appendChild(li);
        });
        output.appendChild(ul);
    }


    // Получение всех общедоступных заметок
    document.getElementById('view-public-notes').addEventListener('click', async function() {
        const response = await fetch('/api/notes/public', { method: 'GET' });
        if (response.ok) {
            const notes = await response.json();
            displayList(notes, 'note');
        } else {
            console.error('Не удалось получить заметки');
        }
    });

    // Получение всех коллекций пользователя
    document.getElementById('view-collections').addEventListener('click', async function() {
        const response = await fetch('/api/collections', { method: 'GET' });
        if (response.ok) {
            const collections = await response.json();
            displayList(collections, 'collection');
        } else {
            console.error('Не удалось получить коллекции');
        }
    });

    // Выход из аккаунта
    document.getElementById('logout').addEventListener('click', function() {
        // Обработка выхода из системы, например, очистка токена
        window.location.reload(); // Перезагружаем страницу для возврата к форме авторизации
    });
});

