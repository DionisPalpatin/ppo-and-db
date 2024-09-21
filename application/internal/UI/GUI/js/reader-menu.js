// reader-menu.js

async function findNote() {
    const noteId = prompt("Введите ID заметки:");
    if (noteId) {
        const response = await fetch(`/api/note/${noteId}`, {
            method: 'GET'
        });
        const note = await response.json();
        alert(`Заметка: ${note.name}\nСодержание: ${note.content}`);
    }
}

async function getAllNotes() {
    const response = await fetch('/api/notes', {
        method: 'GET'
    });
    const notes = await response.json();
    alert(`Заметки: ${notes.map(note => note.name).join(', ')}`);
}

async function getUserCollections() {
    const response = await fetch('/api/user/collections', {
        method: 'GET'
    });
    const collections = await response.json();
    alert(`Коллекции: ${collections.map(collection => collection.name).join(', ')}`);
}

async function getNotesInCollection() {
    const collectionId = prompt("Введите ID коллекции:");
    if (collectionId) {
        const response = await fetch(`/api/collection/${collectionId}/notes`, {
            method: 'GET'
        });
        const notes = await response.json();
        alert(`Заметки в коллекции: ${notes.map(note => note.name).join(', ')}`);
    }
}

async function getTeamNotes() {
    const response = await fetch('/api/team/notes', {
        method: 'GET'
    });
    const notes = await response.json();
    alert(`Заметки команды: ${notes.map(note => note.name).join(', ')}`);
}

function logout() {
    window.location.href = '/';
}
