document.getElementById('actionSelect').addEventListener('change', function() {
    const action = this.value;
    const contentArea = document.getElementById('contentArea');

    contentArea.innerHTML = ''; // Очищаем содержимое для нового действия

    switch(action) {
        case 'addNote':
            renderAddNoteForm();
            break;
        case 'deleteNote':
            renderDeleteNoteForm();
            break;
        case 'findNote':
            renderFindNoteForm();
            break;
        case 'showPublicNotes':
            fetchNotes('/api/notes/public');
            break;
        case 'showAllNotes':
            fetchNotes('/api/notes');
            break;
        case 'showCollectionNotes':
            renderShowCollectionNotesForm();
            break;
        case 'showTeamNotes':
            fetchNotes('/api/notes/team');
            break;
        case 'showAllTeamNotes':
            fetchNotes('/api/notes/all-teams');
            break;
        case 'addCollection':
            renderAddCollectionForm();
            break;
        case 'deleteCollection':
            renderDeleteCollectionForm();
            break;
        case 'showUserCollections':
            fetchCollections('/api/collections/user');
            break;
        case 'showAllCollections':
            fetchCollections('/api/collections');
            break;
        case 'addNoteToCollection':
            renderAddNoteToCollectionForm();
            break;
        case 'deleteNoteFromCollection':
            renderDeleteNoteFromCollectionForm();
            break;
        case 'deleteUser':
            renderDeleteUserForm();
            break;
        case 'updateUserFio':
            renderUpdateUserFioForm();
            break;
        case 'updateUserRole':
            renderUpdateUserRoleForm();
            break;
        case 'findUser':
            renderFindUserForm();
            break;
        case 'showAllUsers':
            fetchUsers('/api/users');
            break;
        case 'addTeam':
            renderAddTeamForm();
            break;
        case 'deleteTeam':
            renderDeleteTeamForm();
            break;
        case 'findTeam':
            renderFindTeamForm();
            break;
        case 'showAllTeamMembers':
            renderShowTeamMembersForm();
            break;
        case 'showAllTeams':
            fetchTeams('/api/teams');
            break;
        case 'addUserToTeam':
            renderAddUserToTeamForm();
            break;
        case 'deleteUserFromTeam':
            renderDeleteUserFromTeamForm();
            break;
        case 'addSection':
            renderAddSectionForm();
            break;
        case 'deleteSection':
            renderDeleteSectionForm();
            break;
        case 'showAllSections':
            fetchSections('/api/sections');
            break;
        case 'addNoteToSection':
            renderAddNoteToSectionForm();
            break;
        case 'deleteNoteFromSection':
            renderDeleteNoteFromSectionForm();
            break;
        case 'logout':
            logout();
            break;
        default:
            contentArea.innerHTML = '<p>Выберите действие.</p>';
    }
});


// check
function renderAddNoteForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Записку</h3>
    <input class="form-input" type="text" id="noteTitle" placeholder="Введите название записки" required>
    <input class="form-input" type="text" id="noteText" placeholder="Введите текст записки" required>
    <input class="form-input" type="file" id="noteImage" accept="image/png, image/jpeg">
    <input class="form-input" type="file" id="noteRawData" accept=".bin,.dat">
    <button id="submitAddNote">Добавить</button>
    <div id="loadingMessage" style="display: none;">Загрузка...</div>
  `;

    document.getElementById('submitAddNote').addEventListener('click', function() {
        const noteTitle = document.getElementById('noteTitle').value;
        const noteText = document.getElementById('noteText').value;
        const noteImage = document.getElementById('noteImage').files[0];
        const noteRawData = document.getElementById('noteRawData').files[0];

        // Получаем userId из localStorage
        const userId = localStorage.getItem('UserId');
        const userRole = localStorage.getItem('UserRole')

        if (!noteTitle) {
            alert('Пожалуйста, введите название записки.');
            return;
        }

        const formData = new FormData();

        formData.append('title', noteTitle);
        formData.append('text', noteText);
        formData.append('userId', userId);
        formData.append('userRole', userRole);// Передаем ID текущего пользователя
        if (noteImage) {
            formData.append('image', noteImage);
        }
        if (noteRawData) {
            formData.append('rawData', noteRawData);
        }

        // Отображаем сообщение о загрузке
        document.getElementById('loadingMessage').style.display = 'block';

        fetch('/api/notes', {
            method: 'POST',
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('loadingMessage').style.display = 'none';
                alert('Записка добавлена!');
            })
            .catch(error => {
                document.getElementById('loadingMessage').style.display = 'none';
                console.error('Ошибка:', error);
                alert('Ошибка при добавлении записки');
            });
    });
}


// check
function fetchNotes(apiUrl) {
    fetch(apiUrl)
        .then(response => response.json())
        .then(data => {
            const contentArea = document.getElementById('contentArea');
            contentArea.innerHTML = `<h3>Записки</h3>`;
            const list = document.createElement('ul');
            data.forEach(note => {
                const listItem = document.createElement('li');
                listItem.textContent = `${note.id}) ${note.name}`;
                list.appendChild(listItem);
            });
            contentArea.appendChild(list);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}


// check
function renderDeleteNoteFromSectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Записку из Раздела</h3>
    <input class="form-input" type="text" id="noteId" placeholder="ID или название записки">
    <input class="form-input" type="text" id="sectionId" placeholder="ID раздела">
    <button id="submitDeleteNoteFromSection">Удалить</button>
  `;

    document.getElementById('submitDeleteNoteFromSection').addEventListener('click', function() {
        const noteId = document.getElementById('noteId').value;
        const sectionId = document.getElementById('sectionId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(noteId)) {
            apiUrl = `/api/sections/${sectionId}/notes/id/${noteId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/sections/${sectionId}/notes/name/${noteId}`;
        }

        fetch(apiUrl, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Записка удалена');
        });
    });
}


// check
function renderAddNoteToSectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Записку в Раздел</h3>
    <input class="form-input" type="text" id="noteId" placeholder="ID или название записки">
    <input class="form-input" type="text" id="sectionId" placeholder="ID раздела">
    <button id="submitAddNoteToSection">Добавить</button>
  `;

    document.getElementById('submitAddNoteToSection').addEventListener('click', function() {
        const noteId = document.getElementById('noteId').value;
        const sectionId = document.getElementById('sectionId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(noteId)) {
            apiUrl = `/api/sections/${sectionId}/notes/id/${noteId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/sections/${sectionId}/notes/name/${noteId}`;
        }

        fetch(apiUrl, {
            method: 'POST',
            body: JSON.stringify({ noteId }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Записка добавлена');
        });
    });
}


// check
function renderDeleteNoteForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Записку</h3>
    <input class="form-input" type="text" id="noteId" placeholder="Введите ID или имя записки">
    <button id="submitDeleteNote">Удалить</button>
  `;

    document.getElementById('submitDeleteNote').addEventListener('click', function() {
        const noteId = document.getElementById('noteId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(noteId)) {
            apiUrl = `/api/notes/id/${noteId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/notes/name/${noteId}`;
        }

        fetch(apiUrl, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                alert('Записка удалена!');
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function renderFindNoteForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Найти Записку</h3>
    <input class="form-input" type="text" id="noteInput" placeholder="Введите ID или имя записки">
    <button id="submitFindNote">Найти</button>
  `;

    document.getElementById('submitFindNote').addEventListener('click', function() {
        const noteInput = document.getElementById('noteInput').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(noteInput)) {
            apiUrl = `/api/notes/id/${noteInput}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/notes/name/${noteInput}`;
        }

        fetch(apiUrl)
            .then(response => response.json())
            .then(note => {
                contentArea.innerHTML = `
        <h3>Записка</h3>
        <p>ID: ${note.id}</p>
        <p>Текст: ${note.text}</p>
        ${note.image ? `<img src="/uploads/${note.image}" alt="Записка Картинка" />` : ''}
        ${note.rawData ? `<a href="/uploads/${note.rawData}" download>Скачать RAW данные</a>` : ''}
      `;
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function renderShowCollectionNotesForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Показать Записки из Подборки</h3>
    <input class="form-input" type="text" id="collectionId" placeholder="ID или название подборки">
    <button id="submitShowCollectionNotes">Показать</button>
  `;

    document.getElementById('submitShowCollectionNotes').addEventListener('click', function() {
        const collectionId = document.getElementById('collectionId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(collectionId)) {
            apiUrl = `/api/collections/id/${collectionId}/notes`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/collections/name/${collectionId}/notes`;
        }

        fetch(apiUrl)
            .then(response => response.json())
            .then(notes => {
                const contentArea = document.getElementById('contentArea');
                contentArea.innerHTML = notes.map(note => `<div>${note.id}) ${note.name}</div>`).join('');
            });
    });
}


// check
function renderAddCollectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Подборку</h3>
    <input class="form-input" type="text" id="collectionName" placeholder="Введите название подборки">
    <button id="submitAddCollection">Добавить</button>
  `;

    document.getElementById('submitAddCollection').addEventListener('click', function() {
        const collectionName = document.getElementById('collectionName').value;

        fetch('/api/collections', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: collectionName, owner_id: UserID })
        })
            .then(response => response.json())
            .then(data => {
                alert('Подборка добавлена!');
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function renderDeleteCollectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Подборку</h3>
    <input class="form-input" type="text" id="collectionId" placeholder="Введите ID или имя подборки">
    <button id="submitDeleteCollection">Удалить</button>
  `;

    document.getElementById('submitDeleteCollection').addEventListener('click', function() {
        const collectionId = document.getElementById('collectionId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(collectionId)) {
            apiUrl = `/api/collections/id/${collectionId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/collections/name/${collectionId}`;
        }

        fetch(apiUrl, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Подборка удалена');
        });
    });
}


// check
function fetchCollections(s) {
    fetch(`/api/collections?s=${s}`)
        .then(response => response.json())
        .then(collections => {
            const contentArea = document.getElementById('contentArea');
            contentArea.innerHTML = collections.map(coll => `<div>${coll.id}) ${coll.name}</div>`).join('');
        });
}


// check
function renderAddNoteToCollectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Записку в Подборку</h3>
    <input class="form-input" type="text" id="noteId" placeholder="ID или название записки">
    <input class="form-input" type="text" id="collectionId" placeholder="ID или название подборки">
    <button id="submitAddNoteToCollection">Добавить</button>
  `;

    document.getElementById('submitAddNoteToCollection').addEventListener('click', function() {
        const noteId = document.getElementById('noteId').value;
        const collectionId = document.getElementById('collectionId').value;

        let apiUrl;
        if (!isNaN(noteId)) {
            if (!isNaN(collectionId)) {
                apiUrl = `/api/collections/id/${collectionId}/notes/id/${noteId}`;
            } else {
                apiUrl = `/api/collections/name/${collectionId}/notes/id/${noteId}`;
            }
        } else {
            if (!isNaN(collectionId)) {
                apiUrl = `/api/collections/id/${collectionId}/notes/name/${noteId}`;
            } else {
                apiUrl = `/api/collections/name/${collectionId}/notes/name/${noteId}`;
            }
        }

        fetch(apiUrl, {
            method: 'POST',
            body: JSON.stringify({ noteId }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Записка добавлена');
        });
    });
}


// check
function renderDeleteNoteFromCollectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Записку из Подборки</h3>
    <input class="form-input" type="text" id="noteId" placeholder="ID или название записки">
    <input class="form-input" type="text" id="collectionId" placeholder="ID или название подборки">
    <button id="submitDeleteNoteFromCollection">Удалить</button>
  `;

    document.getElementById('submitDeleteNoteFromCollection').addEventListener('click', function() {
        const noteId = document.getElementById('noteId').value;
        const collectionId = document.getElementById('collectionId').value;

        let apiUrl;
        if (!isNaN(noteId)) {
            if (!isNaN(collectionId)) {
                apiUrl = `/api/collections/id/${collectionId}/notes/id/${noteId}`;
            } else {
                apiUrl = `/api/collections/name/${collectionId}/notes/id/${noteId}`;
            }
        } else {
            if (!isNaN(collectionId)) {
                apiUrl = `/api/collections/id/${collectionId}/notes/name/${noteId}`;
            } else {
                apiUrl = `/api/collections/name/${collectionId}/notes/name/${noteId}`;
            }
        }

        fetch(apiUrl, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Записка удалена');
        });
    });
}


// check
function renderDeleteUserForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Пользователя</h3>
    <input class="form-input" type="text" id="userId" placeholder="Введите ID или имя пользователя">
    <button id="submitDeleteUser">Удалить</button>
  `;

    document.getElementById('submitDeleteUser').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;

        let apiUrl;
        if (!isNaN(userId)) {
            apiUrl = `/api/users/id/${userId}`;
        } else {
            apiUrl = `/api/users/name/${userId}`;
        }

        fetch(apiUrl, {
            method: 'DELETE'
        })
            .then(response => response.json())
            .then(data => {
                alert('Пользователь удален!');
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function renderUpdateUserFioForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Изменить ФИО пользователя</h3>
    <input class="form-input" type="text" id="userId" placeholder="ID или имя пользователя">
    <input class="form-input" type="text" id="newFio" placeholder="Новое ФИО">
    <button id="submitUpdateUserFio">Изменить</button>
  `;

    document.getElementById('submitUpdateUserFio').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;
        const newFio = document.getElementById('newFio').value;

        let apiUrl;
        if (!isNaN(userId)) {
            apiUrl = `/api/users/id/${userId}/fio`;
        } else {
            apiUrl = `/api/users/name/${userId}/fio`;
        }

        fetch(apiUrl, {
            method: 'PATCH',
            body: JSON.stringify({ fio: newFio }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'ФИО обновлено');
        });
    });
}


// check
function renderUpdateUserRoleForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Изменить Роль пользователя</h3>
    <input class="form-input" type="text" id="userId" placeholder="ID или имя пользователя">
    <input class="form-input" type="number" id="newRole" placeholder="Новая роль (число)">
    <button id="submitUpdateUserRole">Изменить</button>
  `;

    document.getElementById('submitUpdateUserRole').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;
        const newRole = document.getElementById('newRole').value;

        let apiUrl;
        if (!isNaN(userId)) {
            apiUrl = `/api/users/id/${userId}/role`;
        } else {
            apiUrl = `/api/users/name/${userId}/role`;
        }

        fetch(apiUrl, {
            method: 'PATCH',
            body: JSON.stringify({ role: newRole }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Роль обновлена');
        });
    });
}


// check
function renderFindUserForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Найти Пользователя</h3>
    <input class="form-input" type="text" id="userId" placeholder="Введите ID или ФИО пользователя">
    <button id="submitFindUser">Найти</button>
  `;

    document.getElementById('submitFindUser').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;

        let apiUrl;
        if (!isNaN(userId)) {
            apiUrl = `/api/users/id/${userId}`;
        } else {
            apiUrl = `/api/users/name/${userId}`;
        }

        fetch(apiUrl)
            .then(response => response.json())
            .then(user => {
                contentArea.innerHTML = `
        <h3>Пользователь</h3>
        <p>ID: ${user.id}</p>
        <p>ФИО: ${user.fio}</p>
        <p>Роль: ${user.role}</p>
      `;
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function fetchUsers(apiUrl) {
    fetch(apiUrl)
        .then(response => response.json())
        .then(users => {
            const contentArea = document.getElementById('contentArea');
            contentArea.innerHTML = '<h3>Все Пользователи</h3>';
            const list = document.createElement('ul');
            users.forEach(user => {
                const listItem = document.createElement('li');
                listItem.textContent = `${user.id}) ${user.fio}`;
                list.appendChild(listItem);
            });
            contentArea.appendChild(list);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}


// check
function renderAddTeamForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Команду</h3>
    <input class="form-input" type="text" id="teamName" placeholder="Введите название команды">
    <button id="submitAddTeam">Добавить</button>
  `;

    document.getElementById('submitAddTeam').addEventListener('click', function() {
        const teamName = document.getElementById('teamName').value;

        fetch('/api/teams', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: teamName })
        })
            .then(response => response.json())
            .then(data => {
                alert('Команда добавлена!');
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
}


// check
function renderDeleteTeamForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Команду</h3>
    <input class="form-input" type="text" id="teamId" placeholder="ID или название команды">
    <button id="submitDeleteTeam">Удалить</button>
  `;

    document.getElementById('submitDeleteTeam').addEventListener('click', function() {
        const teamId = document.getElementById('teamId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(teamId)) {
            apiUrl = `/api/teams/id/${teamId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/teams/name/${teamId}`;
        }

        fetch(apiUrl, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Команда удалена');
        });
    });
}


// check
function renderFindTeamForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Найти Команду</h3>
    <input class="form-input" type="text" id="teamIdOrName" placeholder="ID или название команды">
    <button id="submitFindTeam">Найти</button>
  `;

    document.getElementById('submitFindTeam').addEventListener('click', function() {
        const teamIdOrName = document.getElementById('teamIdOrName').value;

        let apiUrl;
        if (!isNaN(teamIdOrName)) {
            apiUrl = `/api/teams/id/${teamIdOrName}`;
        } else {
            apiUrl = `/api/teams/name/${teamIdOrName}`;
        }

        fetch(apiUrl)
            .then(response => response.json())
            .then(team => {
                const contentArea = document.getElementById('contentArea');
                contentArea.innerHTML = `<div>${team.id}) ${team.name}</div>`;
            });
    });
}


// check
function renderShowTeamMembersForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Показать членов команды</h3>
    <input class="form-input" type="text" id="teamId" placeholder="ID или название команды">
    <button id="submitShowTeamMembers">Показать</button>
  `;

    document.getElementById('submitShowTeamMembers').addEventListener('click', function() {
        const teamId = document.getElementById('teamId').value;

        let apiUrl;
        if (!isNaN(teamId)) {
            apiUrl = `/api/teams/id/${teamId}/members`;
        } else {
            apiUrl = `/api/teams/name/${teamId}/members`;
        }

        fetch(apiUrl)
            .then(response => response.json())
            .then(members => {
                const contentArea = document.getElementById('contentArea');
                contentArea.innerHTML = members.map(member => `<div>${member.id}) ${member.fio}</div>`).join('');
            });
    });
}


// check
function fetchTeams(s) {
    fetch(`/api/teams?s=${s}`)
        .then(response => response.json())
        .then(teams => {
            const contentArea = document.getElementById('contentArea');
            contentArea.innerHTML = teams.map(team => `<div>${team.id}) ${team.name}</div>`).join('');
        });
}


// check
function renderAddUserToTeamForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить пользователя в команду</h3>
    <input class="form-input" type="text" id="userId" placeholder="ID или ФИО пользователя">
    <input class="form-input" type="text" id="teamId" placeholder="ID или название команды">
    <button id="submitAddUserToTeam">Добавить</button>
  `;

    document.getElementById('submitAddUserToTeam').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;
        const teamId = document.getElementById('teamId').value;

        let apiUrl;
        if (!isNaN(teamId)) {
            if (!isNaN(userId)) {
                apiUrl = `/api/teams/id/${teamId}/members/id/${userId}`;
            } else {
                apiUrl = `/api/teams/id/${teamId}/members/name/${userId}`;
            }
        } else {
            if (!isNaN(userId)) {
                apiUrl = `/api/teams/name/${teamId}/members/id/${userId}`;
            } else {
                apiUrl = `/api/teams/name/${teamId}/members/name/${userId}`;
            }
        }

        fetch(apiUrl, {
            method: 'POST',
            body: JSON.stringify({ userId }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Пользователь добавлен');
        });
    });
}


// check
function renderDeleteUserFromTeamForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить пользователя из команды</h3>
    <input class="form-input" type="text" id="userId" placeholder="ID или ФИО пользователя">
    <input class="form-input" type="text" id="teamId" placeholder="ID или команды">
    <button id="submitDeleteUserFromTeam">Удалить</button>
  `;

    document.getElementById('submitDeleteUserFromTeam').addEventListener('click', function() {
        const userId = document.getElementById('userId').value;
        const teamId = document.getElementById('teamId').value;

        let apiUrl;
        if (!isNaN(teamId)) {
            if (!isNaN(userId)) {
                apiUrl = `/api/teams/id/${teamId}/members/id/${userId}`;
            } else {
                apiUrl = `/api/teams/id/${teamId}/members/name/${userId}`;
            }
        } else {
            if (!isNaN(userId)) {
                apiUrl = `/api/teams/name/${teamId}/members/id/${userId}`;
            } else {
                apiUrl = `/api/teams/name/${teamId}/members/name/${userId}`;
            }
        }

        fetch(apiUrl, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Пользователь удален');
        });
    });
}


// check
function renderAddSectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Добавить Раздел</h3>
    <input class="form-input" type="text" id="sectionName" placeholder="ID или название команды, для которой создается раздел">
    <button id="submitAddSection">Добавить</button>
  `;

    document.getElementById('submitAddSection').addEventListener('click', function() {
        const teamName = document.getElementById('sectionName').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(teamName)) {
            apiUrl = `/api/sections/id/${teamName}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/sections/name/${teamName}`;
        }

        fetch(apiUrl, {
            method: 'POST',
            body: JSON.stringify({ name: teamName }),
            headers: { 'Content-Type': 'application/json' }
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Раздел добавлен');
        });
    });
}


function renderDeleteSectionForm() {
    const contentArea = document.getElementById('contentArea');
    contentArea.innerHTML = `
    <h3>Удалить Раздел</h3>
    <input class="form-input" type="text" id="sectionId" placeholder="ID раздела или имя команды, к которой он относится">
    <button id="submitDeleteSection">Удалить</button>
  `;

    document.getElementById('submitDeleteSection').addEventListener('click', function() {
        const sectionId = document.getElementById('sectionId').value;

        let apiUrl;
        // Проверяем, является ли ввод числом, если да — это ID
        if (!isNaN(sectionId)) {
            apiUrl = `/api/sections/id/${sectionId}`;
        } else {
            // Иначе будем искать по имени
            apiUrl = `/api/sections/name/${sectionId}`;
        }

        fetch(`/api/sections/${sectionId}`, {
            method: 'DELETE'
        }).then(response => response.json()).then(data => {
            alert(data.message || 'Раздел удален');
        });
    });
}


// check
function fetchSections(s) {
    fetch(`/api/sections?s=${s}`)
        .then(response => response.json())
        .then(sections => {
            const contentArea = document.getElementById('contentArea');
            contentArea.innerHTML = sections.map(section => `<div>${section.id}) ${section.name}</div>`).join('');
        });
}

