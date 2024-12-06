package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
	bl "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// Структура для общей конфигурации приложения
type App struct {
	Configs *config.Configs
	Router  *mux.Router

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}

type RegisterRequest struct {
	Fio      string `my_json:"fio"`
	Login    string `my_json:"username"`
	Password string `my_json:"password"`
}

type RegisterResponse struct {
	Success bool   `my_json:"success"`
	Message string `my_json:"message,omitempty"`
	UserID  int    `my_json:"id,omitempty"`
}

// Функция для сохранения файлов на сервере
func saveFile(file multipart.File, path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

func removeFileAfterTime(path string, duration time.Duration, logger *mylogger.MyLogger) {
	time.Sleep(duration)
	err := os.Remove(path)
	if err != nil {
		logger.WriteLog(fmt.Sprintf("Error removing file: %s", err.Error()), slog.LevelError, nil)
	} else {
		logger.WriteLog(fmt.Sprintf("File removed:: %s", path), slog.LevelInfo, nil)
	}
}

func (app *App) AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер загружаемых файлов (до 10MB)
	r.ParseMultipartForm(10 << 20)

	// Получаем данные из формы
	noteTitle := r.FormValue("title")
	noteText := r.FormValue("text")

	// Получаем файл изображения, если он передан
	noteImage, _, err := r.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		http.Error(w, "Error while uploading image", http.StatusInternalServerError)
		return
	}
	defer func() {
		if noteImage != nil {
			noteImage.Close()
		}
	}()

	// Получаем файл бинарных данных, если он передан
	noteRawData, _, err := r.FormFile("rawData")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		http.Error(w, "Error while uploading raw data", http.StatusInternalServerError)
		return
	}
	defer func() {
		if noteRawData != nil {
			noteRawData.Close()
		}
	}()

	// Создаем объект записки
	newNote := &models.Note{
		Name:             noteTitle,
		ContentType:      bl.TextCont,
		RegistrationDate: time.Now(),
		OwnerID:          app.CurUser.Id,
		SectionID:        0,
	}

	// Сохраняем содержимое записки
	filePath := ""
	if noteText != "" {
		newNote.ContentType = bl.TextCont
		err = saveFile(noteImage, "texts/"+newNote.Name+".txt")
		if err != nil {
			http.Error(w, "Error saving image file", http.StatusInternalServerError)
			return
		}
		filePath = "texts/" + newNote.Name + ".txt"

	} else if noteImage != nil {
		newNote.ContentType = bl.ImgCont
		ext := r.FormValue("imageExtension")
		err = saveFile(noteImage, "images/"+newNote.Name+"."+ext)
		if err != nil {
			http.Error(w, "Error saving image file", http.StatusInternalServerError)
			return
		}
		filePath = "images/" + newNote.Name + "." + ext

	} else if noteRawData != nil {
		newNote.ContentType = bl.RawData
		ext := r.FormValue("rawDataExtension")
		err = saveFile(noteRawData, "raw_data/"+newNote.Name+"."+ext)
		if err != nil {
			http.Error(w, "Error saving raw data file", http.StatusInternalServerError)
			return
		}
		filePath = "raw_data/" + newNote.Name + "." + ext

	} else {
		http.Error(w, "No content provided", http.StatusBadRequest)
		return
	}

	// Remove saved file
	removeFileAfterTime(filePath, 10*time.Minute, app.Configs.LogConfigs.Logger)

	if _, err := app.IServices.INoteSvc.AddNote(newNote, app.CurUser, app.IRepos.INoteRepo); err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	note, _ := app.IServices.INoteSvc.GetNote(0, noteTitle, bl.SearchByString, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	if err := app.IServices.INoteSvc.UpdateNoteContent(note.Id, app.CurUser, filePath, app.IRepos.INoteRepo); err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message": "Note added successfully!",
	}
	w.Header().Set("Content-Type", "application/my_json")
	json.NewEncoder(w).Encode(response)
}

func (app *App) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	var notes []*models.Note
	var err *bl.MyError
	if s == "all" {
		notes, err = app.IServices.INoteSvc.GetAllNotes(true, app.CurUser, app.IRepos.INoteRepo)
	} else if s == "public" {
		notes, err = app.IServices.INoteSvc.GetAllNotes(false, app.CurUser, app.IRepos.INoteRepo)
	}

	if err.ErrNum != bl.Ok {
		http.Error(w, "Unable to fetch notes", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/my_json")
	w.WriteHeader(http.StatusOK)

	// Преобразуем список записок в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (app *App) GetAllNotesInTeamSectionHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "team" {
		team, err := app.IServices.ITeamSvc.GetUserTeam(app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sec, err := app.IServices.ISecSvc.GetSection(0, team.Name, app.CurUser, bl.SearchByString, app.IRepos.ISecRepo)
		if err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notes, err := app.IServices.ISecSvc.GetAllNotesInSection(sec.Id, app.CurUser, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
		if err.ErrNum != bl.Ok {
			http.Error(w, "Unable to fetch notes", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовки ответа
		w.Header().Set("Content-Type", "application/my_json")
		w.WriteHeader(http.StatusOK)

		// Преобразуем список записок в JSON и отправляем ответ
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	} else if s == "allteams" {
		return
	}
}

func (app *App) DeleteNoteFromSectionHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем переменные из URL
	vars := mux.Vars(r)
	sectionIDStr := vars["sectionId"]
	noteIdentifier := vars["noteId"] // Может быть ID или название записки

	// Преобразуем sectionID в int
	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	// Получаем NoteRepository
	noteSrv := app.IServices.INoteSvc
	secSrv := app.IServices.ISecSvc
	noteRepo := app.IRepos.INoteRepo
	secRepo := app.IRepos.ISecRepo
	teamRepo := app.IRepos.ITeamRepo

	var note *models.Note
	var myErr *bl.MyError

	// Проверяем, является ли noteIdentifier числом (ID записки) или именем
	noteID, err := strconv.Atoi(noteIdentifier)
	if err == nil {
		// Если это число — получаем записку по ID
		note, myErr = noteSrv.GetNote(noteID, "", bl.SearchByID, app.CurUser, noteRepo, secRepo, teamRepo)
	} else {
		note, myErr = noteSrv.GetNote(0, noteIdentifier, bl.SearchByString, app.CurUser, noteRepo, secRepo, teamRepo)
	}

	// Проверяем, есть ли ошибка при получении записки
	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	sec, myErr := secSrv.GetSection(sectionID, "", app.CurUser, bl.SearchByID, secRepo)
	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	// Удаляем записку из раздела
	myErr = secSrv.DeleteNoteFromSection(sec, note, app.CurUser, secRepo, teamRepo)
	if myErr.ErrNum != bl.Ok {
		http.Error(w, "Error deleting note", http.StatusInternalServerError)
		return
	}

	// Формируем успешный ответ
	w.Header().Set("Content-Type", "application/my_json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note successfully deleted from section"})
}

func (app *App) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectionIDStr := vars["sectionId"]
	noteIdentifier := vars["noteId"]

	// Преобразуем sectionID в int
	sectionID, err := strconv.Atoi(sectionIDStr)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	noteSrv := app.IServices.INoteSvc
	secSrv := app.IServices.ISecSvc

	// Получаем noteID, проверяем, является ли оно числом или строкой
	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	var myErr *bl.MyError
	if err == nil {
		note, myErr = noteSrv.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, myErr = noteSrv.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	section, secErr := secSrv.GetSection(sectionID, "", app.CurUser, 1, app.IRepos.ISecRepo)
	if secErr.ErrNum != bl.Ok {
		http.Error(w, secErr.Error(), http.StatusNotFound)
		return
	}

	// Добавляем записку в раздел
	addErr := secSrv.AddNoteToSection(section, note, app.CurUser, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	if addErr.ErrNum != bl.Ok {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note added to section successfully"})
}

func (app *App) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdentifier := vars["noteId"]

	// Получаем noteID
	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	var myErr *bl.MyError
	if err == nil {
		note, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	// Удаляем записку
	delErr := app.IServices.INoteSvc.DeleteNote(note.Id, app.CurUser, app.IRepos.INoteRepo)
	if delErr.ErrNum != bl.Ok {
		http.Error(w, delErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}

// func (app *handlers) FindNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	noteIdentifier := vars["noteInput"]
//
// 	// Определяем, является ли noteInput ID или именем
// 	noteID, err := strconv.Atoi(noteIdentifier)
// 	var note *models.Note
// 	var myErr *bl.MyError
// 	if err == nil {
// 		note, _, _, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
// 	} else {
// 		note, _, _, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
// 	}
//
// 	if myErr.ErrNum != bl.Ok {
// 		http.Error(w, myErr.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/my_json")
// 	my_json.NewEncoder(w).Encode(note)
// }

func (app *App) FindNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdentifier := vars["noteInput"]

	// Определяем, является ли noteInput ID или именем
	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	var data []byte
	var extension string
	var myErr *bl.MyError
	if err == nil {
		note, myErr = app.IServices.INoteSvc.GetNote(noteID, "", bl.SearchByID, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, bl.SearchByString, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	// В зависимости от расширения файла, выбираем, как его отправить
	switch extension {
	case "txt":
		// Если это текстовая записка
		w.Header().Set("Content-Type", "application/my_json")
		response := map[string]interface{}{
			"note": note,
			"text": string(data), // Преобразуем []byte в строку
		}
		json.NewEncoder(w).Encode(response)
	case "png", "jpg", "jpeg":
		// Если это изображение
		w.Header().Set("Content-Type", "image/"+extension)
		w.Write(data) // Отправляем картинку
	default:
		// Для всех остальных файлов
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\"note."+extension+"\"")
		w.Write(data) // Отправляем файл для скачивания
	}
}

func (app *App) ShowCollectionNotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionIdentifier := vars["collectionId"]

	// Определяем, является ли collectionId ID или названием
	collectionID, err := strconv.Atoi(collectionIdentifier)
	var collection *models.Collection
	var myErr *bl.MyError
	if err == nil {
		collection, myErr = app.IServices.IColSvc.GetCollection(collectionID, "", 1, app.IRepos.IColRepo)
	} else {
		collection, myErr = app.IServices.IColSvc.GetCollection(0, collectionIdentifier, 2, app.IRepos.IColRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	notes, notesErr := app.IServices.IColSvc.GetAllNotesInCollection(collection, app.IRepos.IColRepo)
	if notesErr.ErrNum != bl.Ok {
		http.Error(w, "Error fetching notes from collection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	json.NewEncoder(w).Encode(notes)
}

func (app *App) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionIdentifier := vars["collectionId"]
	noteIdentifier := vars["noteId"]

	// Получаем collectionID и noteID
	collectionID, err := strconv.Atoi(collectionIdentifier)
	var collection *models.Collection
	var myErr *bl.MyError
	if err == nil {
		collection, myErr = app.IServices.IColSvc.GetCollection(collectionID, "", 1, app.IRepos.IColRepo)
	} else {
		collection, myErr = app.IServices.IColSvc.GetCollection(0, collectionIdentifier, 2, app.IRepos.IColRepo)
	}

	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	if err == nil {
		note, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, "Note or Collection not found", http.StatusNotFound)
		return
	}

	// Добавляем записку в подборку
	addErr := app.IServices.INoteSvc.AddNoteToCollection(note.Id, collection.Id, app.IRepos.INoteRepo)
	if addErr.ErrNum != bl.Ok {
		http.Error(w, "Error adding note to collection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note added to collection successfully"})
}

func (app *App) DeleteNoteFromCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteId := vars["noteId"]
	collectionId := vars["collectionId"]

	var noteIDInt, collIDInt int
	var err error

	// Проверяем, является ли noteId числом
	noteIDInt, err = strconv.Atoi(noteId)
	if err != nil {
		http.Error(w, "Некорректный ID записки", http.StatusBadRequest)
		return
	}

	// Проверяем, является ли collectionId числом
	collIDInt, err = strconv.Atoi(collectionId)
	if err != nil {
		http.Error(w, "Некорректный ID подборки", http.StatusBadRequest)
		return
	}

	// Удаление записки из подборки через сервис
	myErr := app.IServices.INoteSvc.DeleteNoteFromCollection(noteIDInt, collIDInt, app.IRepos.INoteRepo)
	if myErr != nil && myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Записка успешно удалена из подборки"})
}

func (app *App) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var newCollection models.Collection

	err := json.NewDecoder(r.Body).Decode(&newCollection)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Устанавливаем владельца
	newCollection.OwnerID = app.CurUser.Id

	_, myErr := app.IServices.IColSvc.AddCollection(&newCollection, app.IRepos.IColRepo)
	if myErr != nil && myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Подборка успешно добавлена"})
}

func (app *App) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collectionId := vars["collectionId"]

	var collIDInt int
	var err error

	// Проверяем, является ли collectionId числом
	collIDInt, err = strconv.Atoi(collectionId)
	if err != nil {
		http.Error(w, "Некорректный ID подборки", http.StatusBadRequest)
		return
	}

	// Удаление подборки через сервис
	myErr := app.IServices.IColSvc.DeleteCollection(collIDInt, app.CurUser, app.IRepos.IColRepo)
	if myErr != nil && myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Подборка успешно удалена"})
}

func (app *App) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	var collections []*models.Collection
	var myErr *bl.MyError
	if s == "all" {
		collections, myErr = app.IServices.IColSvc.GetAllCollections(app.CurUser, app.IRepos.IColRepo)
	} else if s == "user" {
		collections, myErr = app.IServices.IColSvc.GetAllUsersCollections(app.CurUser, app.IRepos.IColRepo)
	}

	if myErr != nil && myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collections)
}

func (app *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	var userIDInt int
	var err error

	// Проверяем, является ли userId числом
	userIDInt, err = strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	// Удаление пользователя через сервис
	myErr := app.IServices.IUsrSvc.DeleteUser(app.CurUser, userIDInt, app.IRepos.IUsrRepo)
	if myErr != nil && myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно удален"})
}

func (app *App) UpdateUserFioHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var newFio struct {
		Fio string `my_json:"fio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newFio); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверим, пользователь запрашивается по ID или по имени
	var user *models.User
	var myErr *bl.MyError
	if id, err := strconv.Atoi(userId); err == nil {
		user, myErr = app.IServices.IUsrSvc.GetUser(id, "", 0, app.CurUser, app.IRepos.IUsrRepo)
	} else {
		user, myErr = app.IServices.IUsrSvc.GetUser(0, userId, 1, app.CurUser, app.IRepos.IUsrRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	user.Fio = newFio.Fio
	updateErr := app.IServices.IUsrSvc.UpdateUser(app.CurUser, user, app.IRepos.IUsrRepo)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "ФИО обновлено"}
	json.NewEncoder(w).Encode(response)
}

func (app *App) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var newRole struct {
		Role int `my_json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newRole); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user *models.User
	var myErr *bl.MyError
	if id, err := strconv.Atoi(userId); err == nil {
		user, myErr = app.IServices.IUsrSvc.GetUser(id, "", 0, app.CurUser, app.IRepos.IUsrRepo)
	} else {
		user, myErr = app.IServices.IUsrSvc.GetUser(0, userId, 1, app.CurUser, app.IRepos.IUsrRepo)
	}

	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	user.Role = newRole.Role
	updateErr := app.IServices.IUsrSvc.UpdateUser(app.CurUser, user, app.IRepos.IUsrRepo)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Роль обновлена"}
	json.NewEncoder(w).Encode(response)
}

func (app *App) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	var newTeam struct {
		Name string `my_json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	team := &models.Team{
		Name:             newTeam.Name,
		RegistrationDate: time.Now(),
	}

	_, addErr := app.IServices.ITeamSvc.AddTeam(app.CurUser, team, app.IRepos.ITeamRepo)
	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Команда добавлена"}
	json.NewEncoder(w).Encode(response)
}

func (app *App) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		app.Configs.LogConfigs.Logger.WriteLog("GetAllUsersHandler is called", slog.LevelInfo, nil)
		users, err := app.IServices.IUsrSvc.GetAllUsers(app.CurUser, app.IRepos.IUsrRepo)
		if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/my_json")
		json.NewEncoder(w).Encode(users)
	}
}

func (app *App) DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var team *models.Team
	var err *bl.MyError

	teamID := vars["teamId"]

	if teamID != "" {
		id, _ := strconv.Atoi(teamID)
		team, err = app.IServices.ITeamSvc.GetTeam(id, "", 1, app.CurUser, app.IRepos.ITeamRepo)
	} else {
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamID, 2, app.CurUser, app.IRepos.ITeamRepo)
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.IServices.ITeamSvc.DeleteTeam(app.CurUser, team.Id, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Команда удалена"})
}

func (app *App) FindTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var team *models.Team
	var err *bl.MyError

	teamID := vars["teamId"]

	if teamID != "" {
		id, _ := strconv.Atoi(teamID)
		team, err = app.IServices.ITeamSvc.GetTeam(id, "", 1, app.CurUser, app.IRepos.ITeamRepo)
	} else {
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamID, 2, app.CurUser, app.IRepos.ITeamRepo)
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func (app *App) ShowTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var team *models.Team
	var err *bl.MyError

	teamID := vars["id"]
	teamName := vars["name"]

	if teamID != "" {
		id, _ := strconv.Atoi(teamID)
		team, err = app.IServices.ITeamSvc.GetTeam(id, "", 1, app.CurUser, app.IRepos.ITeamRepo)
	} else {
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamName, 2, app.CurUser, app.IRepos.ITeamRepo)
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := app.IServices.ITeamSvc.GetTeamMembers(team.Id, app.CurUser, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(members)
}

func (app *App) AddUserToTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var team *models.Team
	var user *models.User
	var err *bl.MyError

	teamID := vars["teamId"]
	userID := vars["userId"]

	if teamID != "" {
		id, _ := strconv.Atoi(teamID)
		team, err = app.IServices.ITeamSvc.GetTeam(id, "", 1, app.CurUser, app.IRepos.ITeamRepo)
	} else {
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamID, 2, app.CurUser, app.IRepos.ITeamRepo)
	}

	if userID != "" {
		id, _ := strconv.Atoi(userID)
		user, err = app.IServices.IUsrSvc.GetUser(id, "", 1, app.CurUser, app.IRepos.IUsrRepo)
	} else {
		user, err = app.IServices.IUsrSvc.GetUser(0, userID, 2, app.CurUser, app.IRepos.IUsrRepo)
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.IServices.ITeamSvc.AddUserToTeam(app.CurUser, user.Id, team.Id, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь добавлен в команду"})
}

func (app *App) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		teams, err := app.IServices.ITeamSvc.GetAllTeams(app.CurUser, app.IRepos.ITeamRepo)
		if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/my_json")
		json.NewEncoder(w).Encode(teams)
	}
}

func (app *App) DeleteUserFromTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID := vars["teamId"]
	userID := vars["userId"]

	var err *bl.MyError

	// Проверка типа поиска команды и пользователя
	idTeam, err1 := strconv.Atoi(teamID)
	idUser, err2 := strconv.Atoi(userID)

	if err1 == nil && err2 == nil {
		err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, idUser, idTeam, app.IRepos.ITeamRepo)
	} else if err1 != nil && err2 != nil {
		team, err := app.IServices.ITeamSvc.GetTeam(0, teamID, bl.SearchByString, app.CurUser, app.IRepos.ITeamRepo)
		if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := app.IServices.IUsrSvc.GetUser(0, userID, bl.SearchByString, app.CurUser, app.IRepos.IUsrRepo)
		if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, user.Id, team.Id, app.IRepos.ITeamRepo)
	} else if err1 != nil {
		team, err := app.IServices.ITeamSvc.GetTeam(0, teamID, bl.SearchByString, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.Ok {
			err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, idUser, team.Id, app.IRepos.ITeamRepo)
		} else if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		user, err := app.IServices.IUsrSvc.GetUser(0, userID, bl.SearchByString, app.CurUser, app.IRepos.IUsrRepo)
		if err.ErrNum == bl.Ok {
			err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, user.Id, idTeam, app.IRepos.ITeamRepo)
		} else if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь удален"})
}

func (app *App) AddSectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamName := vars["teamName"]

	var section models.Section
	_ = json.NewDecoder(r.Body).Decode(&section)

	var err *bl.MyError
	idTeam, err1 := strconv.Atoi(teamName)

	if err1 == nil {
		var team *models.Team
		team, err = app.IServices.ITeamSvc.GetTeam(idTeam, "", bl.SearchByID, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.Ok {
			_, err = app.IServices.ISecSvc.AddSection(&section, team, app.CurUser, app.IRepos.ISecRepo)
		}
	} else {
		var team *models.Team
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamName, bl.SearchByString, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.Ok {
			_, err = app.IServices.ISecSvc.AddSection(&section, team, app.CurUser, app.IRepos.ISecRepo)
		}
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Раздел добавлен"})
}

func (app *App) DeleteSectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamName := vars["teamName"]

	var err *bl.MyError
	idSec, err1 := strconv.Atoi(teamName)

	if err1 == nil {
		err = app.IServices.ISecSvc.DeleteSection(idSec, app.CurUser, app.IRepos.ISecRepo)
	} else {
		var section *models.Section
		section, err = app.IServices.ISecSvc.GetSection(0, teamName, app.CurUser, bl.SearchByString, app.IRepos.ISecRepo)
		if err == nil {
			err = app.IServices.ISecSvc.DeleteSection(section.Id, app.CurUser, app.IRepos.ISecRepo)
		}
	}

	if err != nil && err.ErrNum != bl.Ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Раздел удален"})
}

func (app *App) GetAllSectionsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		sections, err := app.IServices.ISecSvc.GetAllSections(app.CurUser, app.IRepos.ISecRepo)
		if err != nil && err.ErrNum != bl.Ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sections)
	}
}

func (app *App) GetFullStatHendler(w http.ResponseWriter, r *http.Request) {
	stat, myErr := app.IServices.IStatSvc.GetFullStat(app.CurUser, app.IRepos.IStatRepo)

	if myErr.ErrNum != bl.Ok {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/my_json")
	json.NewEncoder(w).Encode(stat)
}
