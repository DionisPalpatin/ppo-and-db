package appconfigs

import (
	"encoding/json"
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
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
	mylogger "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// Структура для общей конфигурации приложения
type AppConfigs struct {
	Configs *config.Configs
	Router  *mux.Router
	CurUser *models.User

	IServices *bl.IServices
	IRepos    *bl.IRepositories
}

type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UserID  int    `json:"id,omitempty"`
	Role    int    `json:"role,omitempty"`
}

type RegisterRequest struct {
	Fio      string `json:"fio"`
	Login    string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	UserID  int    `json:"id,omitempty"`
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

func (app *AppConfigs) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid request payload"}`, http.StatusBadRequest)
		app.Configs.LogConfigs.Logger.WriteLog("Invalid request payload", slog.LevelError, nil)
		return
	}

	// Вызов SignInUser через сервис
	user, myErr := app.IServices.IOAuthSvc.SignInUser(loginReq.Login, loginReq.Password, app.IRepos.IUsrRepo)
	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, `{"success": false, "message": "`+myErr.Error()+`"}`, http.StatusUnauthorized)
		app.Configs.LogConfigs.Logger.WriteLog(myErr.Error(), slog.LevelError, nil)
		return
	}

	// Успешная аутентификация
	app.CurUser = user
	response := LoginResponse{
		Success: true,
		UserID:  user.Id,
		Role:    user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *AppConfigs) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Вызов RegisterUser через сервис
	user, myErr := app.IServices.IOAuthSvc.RegisterUser(registerReq.Fio, registerReq.Login, registerReq.Password, app.IRepos.IUsrRepo)
	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusBadRequest)
		return
	}

	// Успешная регистрация
	response := RegisterResponse{
		Success: true,
		UserID:  user.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *AppConfigs) AddNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Ограничиваем размер загружаемых файлов (до 10MB)
	r.ParseMultipartForm(10 << 20)

	// Получаем данные из формы
	noteTitle := r.FormValue("title")
	noteText := r.FormValue("text")

	// Получаем файл изображения, если он передан
	noteImage, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
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
	if err != nil && err != http.ErrMissingFile {
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
		RegistrationDate: time.Now().Format("2006-01-02 15:04:05"),
		OwnerID:          app.CurUser.Id,
		SectionID:        0,
	}

	// Сохраняем содержимое записки
	filePath := ""
	if noteText != "" {
		newNote.ContentType = bl.TextCont
		err := saveFile(noteImage, "texts/"+newNote.Name+".txt")
		if err != nil {
			http.Error(w, "Error saving image file", http.StatusInternalServerError)
			return
		}
		filePath = "texts/" + newNote.Name + ".txt"
	} else if noteImage != nil {
		newNote.ContentType = bl.ImgCont
		err := saveFile(noteImage, "images/"+newNote.Name+".jpg")
		if err != nil {
			http.Error(w, "Error saving image file", http.StatusInternalServerError)
			return
		}
		filePath = "images/" + newNote.Name + ".jpg"
	} else if noteRawData != nil {
		newNote.ContentType = bl.RawData
		err := saveFile(noteRawData, "raw_data/"+newNote.Name+".dat")
		if err != nil {
			http.Error(w, "Error saving raw data file", http.StatusInternalServerError)
			return
		}
		filePath = "raw_data/" + newNote.Name + ".dat"
	} else {
		http.Error(w, "No content provided", http.StatusBadRequest)
		return
	}

	// Remove saved file
	removeFileAfterTime(filePath, 10*time.Minute, app.Configs.LogConfigs.Logger)

	if err := app.IServices.INoteSvc.AddNote(newNote, app.CurUser, app.IRepos.INoteRepo); err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	note, _, _ := app.IServices.INoteSvc.GetNote(0, noteTitle, bl.SearchByString, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	if err := app.IServices.INoteSvc.UpdateNoteContent(note.Id, app.CurUser, filePath, app.IRepos.INoteRepo); err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{
		"message": "Note added successfully!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *AppConfigs) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	var notes []*models.Note
	var err *bl.MyError
	if s == "all" {
		notes, err = app.IServices.INoteSvc.GetAllNotes(true, app.CurUser, app.IRepos.INoteRepo)
	} else if s == "public" {
		notes, err = app.IServices.INoteSvc.GetAllNotes(false, app.CurUser, app.IRepos.INoteRepo)
	}

	if err.ErrNum != bl.AllIsOk {
		http.Error(w, "Unable to fetch notes", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Преобразуем список записок в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (app *AppConfigs) GetAllNotesInTeamSectionHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "team" {
		team, err := app.IServices.ITeamSvc.GetUserTeam(app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sec, err := app.IServices.ISecSvc.GetSection(0, team.Name, app.CurUser, bl.SearchByString, app.IRepos.ISecRepo)
		if err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notes, err := app.IServices.ISecSvc.GetAllNotesInSection(sec.Id, app.CurUser, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
		if err.ErrNum != bl.AllIsOk {
			http.Error(w, "Unable to fetch notes", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовки ответа
		w.Header().Set("Content-Type", "application/json")
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

func (app *AppConfigs) DeleteNoteFromSectionHandler(w http.ResponseWriter, r *http.Request) {
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
		note, _, myErr = noteSrv.GetNote(noteID, "", bl.SearchByID, app.CurUser, noteRepo, secRepo, teamRepo)
	} else {
		note, _, myErr = noteSrv.GetNote(0, noteIdentifier, bl.SearchByString, app.CurUser, noteRepo, secRepo, teamRepo)
	}

	// Проверяем, есть ли ошибка при получении записки
	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	sec, myErr := secSrv.GetSection(sectionID, "", app.CurUser, bl.SearchByID, secRepo)
	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	// Удаляем записку из раздела
	myErr = secSrv.DeleteNoteFromSection(sec, note, app.CurUser, secRepo, teamRepo)
	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Error deleting note", http.StatusInternalServerError)
		return
	}

	// Формируем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note successfully deleted from section"})
}

func (app *AppConfigs) AddNoteToSectionHandler(w http.ResponseWriter, r *http.Request) {
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
		note, _, myErr = noteSrv.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, _, myErr = noteSrv.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	section, secErr := secSrv.GetSection(sectionID, "", app.CurUser, 1, app.IRepos.ISecRepo)
	if secErr.ErrNum != bl.AllIsOk {
		http.Error(w, secErr.Error(), http.StatusNotFound)
		return
	}

	// Добавляем записку в раздел
	addErr := secSrv.AddNoteToSection(section, note, app.CurUser, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	if addErr.ErrNum != bl.AllIsOk {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note added to section successfully"})
}

func (app *AppConfigs) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdentifier := vars["noteId"]

	// Получаем noteID
	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	var myErr *bl.MyError
	if err == nil {
		note, _, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, _, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	// Удаляем записку
	delErr := app.IServices.INoteSvc.DeleteNote(note.Id, app.CurUser, app.IRepos.INoteRepo)
	if delErr.ErrNum != bl.AllIsOk {
		http.Error(w, delErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note deleted successfully"))
}

func (app *AppConfigs) FindNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdentifier := vars["noteInput"]

	// Определяем, является ли noteInput ID или именем
	noteID, err := strconv.Atoi(noteIdentifier)
	var note *models.Note
	var myErr *bl.MyError
	if err == nil {
		note, _, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, _, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (app *AppConfigs) ShowCollectionNotesHandler(w http.ResponseWriter, r *http.Request) {
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

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	notes, notesErr := app.IServices.IColSvc.GetAllNotesInCollection(collection, app.IRepos.IColRepo)
	if notesErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Error fetching notes from collection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (app *AppConfigs) AddNoteToCollectionHandler(w http.ResponseWriter, r *http.Request) {
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
		note, _, myErr = app.IServices.INoteSvc.GetNote(noteID, "", 1, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	} else {
		note, _, myErr = app.IServices.INoteSvc.GetNote(0, noteIdentifier, 2, app.CurUser, app.IRepos.INoteRepo, app.IRepos.ISecRepo, app.IRepos.ITeamRepo)
	}

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Note or Collection not found", http.StatusNotFound)
		return
	}

	// Добавляем записку в подборку
	addErr := app.IServices.INoteSvc.AddNoteToCollection(note.Id, collection.Id, app.IRepos.INoteRepo)
	if addErr.ErrNum != bl.AllIsOk {
		http.Error(w, "Error adding note to collection", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note added to collection successfully"})
}

func (app *AppConfigs) DeleteNoteFromCollectionHandler(w http.ResponseWriter, r *http.Request) {
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
	if myErr != nil && myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Записка успешно удалена из подборки"})
}

func (app *AppConfigs) AddCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var newCollection models.Collection

	err := json.NewDecoder(r.Body).Decode(&newCollection)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Устанавливаем владельца
	newCollection.OwnerID = app.CurUser.Id

	myErr := app.IServices.IColSvc.AddCollection(&newCollection, app.IRepos.IColRepo)
	if myErr != nil && myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Подборка успешно добавлена"})
}

func (app *AppConfigs) DeleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
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
	if myErr != nil && myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Подборка успешно удалена"})
}

func (app *AppConfigs) GetAllCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	var collections []*models.Collection
	var myErr *bl.MyError
	if s == "all" {
		collections, myErr = app.IServices.IColSvc.GetAllCollections(app.CurUser, app.IRepos.IColRepo)
	} else if s == "user" {
		collections, myErr = app.IServices.IColSvc.GetAllUsersCollections(app.CurUser, app.IRepos.IColRepo)
	}

	if myErr != nil && myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collections)
}

func (app *AppConfigs) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
	if myErr != nil && myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно удален"})
}

func (app *AppConfigs) UpdateUserFioHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var newFio struct {
		Fio string `json:"fio"`
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

	if myErr.ErrNum != bl.AllIsOk {
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

func (app *AppConfigs) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var newRole struct {
		Role int `json:"role"`
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

	if myErr.ErrNum != bl.AllIsOk {
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

func (app *AppConfigs) FindUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	app.Configs.LogConfigs.Logger.WriteLog("FindUserHandle get strings: "+userId, slog.LevelInfo, nil)

	var user *models.User
	var myErr *bl.MyError
	if id, err := strconv.Atoi(userId); err == nil {
		user, myErr = app.IServices.IUsrSvc.GetUser(id, "", 0, app.CurUser, app.IRepos.IUsrRepo)
	} else {
		user, myErr = app.IServices.IUsrSvc.GetUser(0, userId, 1, app.CurUser, app.IRepos.IUsrRepo)
	}

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (app *AppConfigs) AddTeamHandler(w http.ResponseWriter, r *http.Request) {
	var newTeam struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newTeam); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	team := &models.Team{
		Name:             newTeam.Name,
		RegistrationDate: time.Now().Format("2006-01-02"),
	}

	addErr := app.IServices.ITeamSvc.AddTeam(app.CurUser, team, app.IRepos.ITeamRepo)
	if addErr != nil {
		http.Error(w, addErr.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Команда добавлена"}
	json.NewEncoder(w).Encode(response)
}

func (app *AppConfigs) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		app.Configs.LogConfigs.Logger.WriteLog("GetAllUsersHandler is called", slog.LevelInfo, nil)
		users, err := app.IServices.IUsrSvc.GetAllUsers(app.CurUser, app.IRepos.IUsrRepo)
		if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func (app *AppConfigs) DeleteTeamHandler(w http.ResponseWriter, r *http.Request) {
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

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.IServices.ITeamSvc.DeleteTeam(app.CurUser, team.Id, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Команда удалена"})
}

func (app *AppConfigs) FindTeamHandler(w http.ResponseWriter, r *http.Request) {
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

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func (app *AppConfigs) ShowTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
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

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := app.IServices.ITeamSvc.GetTeamMembers(team.Id, app.CurUser, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(members)
}

func (app *AppConfigs) AddUserToTeamHandler(w http.ResponseWriter, r *http.Request) {
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

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.IServices.ITeamSvc.AddUserToTeam(app.CurUser, user.Id, team.Id, app.IRepos.ITeamRepo)
	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь добавлен в команду"})
}

func (app *AppConfigs) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		teams, err := app.IServices.ITeamSvc.GetAllTeams(app.CurUser, app.IRepos.ITeamRepo)
		if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teams)
	}
}

func (app *AppConfigs) DeleteUserFromTeamHandler(w http.ResponseWriter, r *http.Request) {
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
		if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := app.IServices.IUsrSvc.GetUser(0, userID, bl.SearchByString, app.CurUser, app.IRepos.IUsrRepo)
		if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, user.Id, team.Id, app.IRepos.ITeamRepo)
	} else if err1 != nil {
		team, err := app.IServices.ITeamSvc.GetTeam(0, teamID, bl.SearchByString, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.AllIsOk {
			err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, idUser, team.Id, app.IRepos.ITeamRepo)
		} else if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		user, err := app.IServices.IUsrSvc.GetUser(0, userID, bl.SearchByString, app.CurUser, app.IRepos.IUsrRepo)
		if err.ErrNum == bl.AllIsOk {
			err = app.IServices.ITeamSvc.DeleteUserFromTeam(app.CurUser, user.Id, idTeam, app.IRepos.ITeamRepo)
		} else if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь удален"})
}

func (app *AppConfigs) AddSectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamName := vars["teamName"]

	var section models.Section
	_ = json.NewDecoder(r.Body).Decode(&section)

	var err *bl.MyError
	idTeam, err1 := strconv.Atoi(teamName)

	if err1 == nil {
		var team *models.Team
		team, err = app.IServices.ITeamSvc.GetTeam(idTeam, "", bl.SearchByID, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.AllIsOk {
			err = app.IServices.ISecSvc.AddSection(&section, team, app.CurUser, app.IRepos.ISecRepo)
		}
	} else {
		var team *models.Team
		team, err = app.IServices.ITeamSvc.GetTeam(0, teamName, bl.SearchByString, app.CurUser, app.IRepos.ITeamRepo)
		if err.ErrNum == bl.AllIsOk {
			err = app.IServices.ISecSvc.AddSection(&section, team, app.CurUser, app.IRepos.ISecRepo)
		}
	}

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Раздел добавлен"})
}

func (app *AppConfigs) DeleteSectionHandler(w http.ResponseWriter, r *http.Request) {
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

	if err != nil && err.ErrNum != bl.AllIsOk {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Раздел удален"})
}

func (app *AppConfigs) GetAllSectionsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("s")

	if s == "all" {
		sections, err := app.IServices.ISecSvc.GetAllSections(app.CurUser, app.IRepos.ISecRepo)
		if err != nil && err.ErrNum != bl.AllIsOk {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sections)
	}
}

func (app *AppConfigs) GetFullStatHendler(w http.ResponseWriter, r *http.Request) {
	stat, myErr := app.IServices.IStatSvc.GetFullStat(app.CurUser, app.IRepos.IStatRepo)

	if myErr.ErrNum != bl.AllIsOk {
		http.Error(w, myErr.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}
