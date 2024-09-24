package bl

import "errors"

// import "errors"

type MyError struct {
	ErrNum   int
	FuncName string
	Err      error
}

func (r *MyError) Error() string {
	return r.Err.Error()
}

const AllIsOk = 0

// Ошибки доступа
const (
	ErrAccessDenied    = iota + 1
	ErrSearchParameter = iota + 2
	ErrNoFile          = iota + 3
)

// Ошибки пользователя
const (
	ErrGetUserByID    = 100 + iota
	ErrGetUserByLogin = 100 + iota
	ErrGetAllUsers    = 100 + iota
	ErrAddUser        = 100 + iota
	ErrDeleteUser     = 100 + iota
	ErrUpdateUser     = 100 + iota
	ErrSignInUser     = 100 + iota
	ErrRegisterUser   = 100 + iota
)

// Ошибки команды
const (
	ErrGetTeamByID        = 200 + iota
	ErrGetTeamByName      = 200 + iota
	ErrGetTeamBySectionID = 200 + iota
	ErrGetAllTeams        = 200 + iota
	ErrAddTeam            = 200 + iota
	ErrDeleteTeam         = 200 + iota
	ErrUpdateTeam         = 200 + iota
	ErrGetTeamMembers     = 200 + iota
	ErrGetUserTeam        = 200 + iota
	ErrAddUserToTeam      = 200 + iota
	ErrDeleteUserFromTeam = 200 + iota
)

// Ошибки заметки
const (
	ErrGetNoteByID              = 300 + iota
	ErrGetNoteByName            = 300 + iota
	ErrGetAllNotes              = 300 + iota
	ErrAddNote                  = 300 + iota
	ErrDeleteNote               = 300 + iota
	ErrUpdateNoteContent        = 300 + iota
	ErrUpdateNoteInfo           = 300 + iota
	ErrAddNoteToCollection      = 300 + iota
	ErrDeleteNoteFromCollection = 300 + iota
)

// Ошибки коллекции
const (
	ErrGetCollectionByID       = 400 + iota
	ErrGetCollectionByName     = 400 + iota
	ErrGetAllCollections       = 400 + iota
	ErrGetAllUserCollections   = 400 + iota
	ErrAddCollection           = 400 + iota
	ErrDeleteCollection        = 400 + iota
	ErrUpdateCollection        = 400 + iota
	ErrGetAllNotesInCollection = 400 + iota
)

// Ошибки раздела
const (
	ErrGetSectionByID       = 500 + iota
	ErrGetSectionByTeamName = 500 + iota
	ErrGetAllSections       = 500 + iota
	ErrAddSection           = 500 + iota
	ErrDeleteSection        = 500 + iota
	ErrUpdateSection        = 500 + iota
	ErrGetAllNotesInSection = 500 + iota
	ErrAddNoteToSection     = 500 + iota
)

const (
	ErrGetFullStat = 600 + iota
)

// Функции для создания ошибок
func ErrAccessDeniedError() error {
	return errors.New("Доступ запрещен")
}

// Ошибки пользователя
func ErrGetUserByIDError() error {
	return errors.New("Ошибка получения пользователя по ID")
}

func ErrGetUserByLoginError() error {
	return errors.New("Ошибка получения пользователя по логину")
}

func ErrGetAllUsersError() error {
	return errors.New("Ошибка получения всех пользователей")
}

func ErrAddUserError() error {
	return errors.New("Ошибка добавления пользователя")
}

func ErrDeleteUserError() error {
	return errors.New("Ошибка удаления пользователя")
}

func ErrUpdateUserError() error {
	return errors.New("Ошибка обновления пользователя")
}

func ErrSignInUserError() error {
	return errors.New("Ошибка авторизации")
}

func ErrRegistrationError() error {
	return errors.New("Ошибка регистрации")
}

// Ошибки команды
func ErrGetTeamByIDError() error {
	return errors.New("Ошибка получения команды по ID")
}

func ErrGetTeamBySectionIDError() error {
	return errors.New("Ошибка получения команды по ID секции")
}

func ErrGetTeamByNameError() error {
	return errors.New("Ошибка получения команды по ее названию")
}

func ErrGetAllTeamsError() error {
	return errors.New("Ошибка получения всех команд")
}

func ErrAddTeamError() error {
	return errors.New("Ошибка добавления команды")
}

func ErrDeleteTeamError() error {
	return errors.New("Ошибка удаления команды")
}

func ErrUpdateTeamError() error {
	return errors.New("Ошибка обновления команды")
}

func ErrGetTeamMembersError() error {
	return errors.New("Ошибка получения членов команды")
}

func ErrGetUserTeamError() error {
	return errors.New("Ошибка получения команды пользователя")
}

func ErrNoFileError() error {
	return errors.New("Нет такого файла")
}

func ErrSearchParameterError() error {
	return errors.New("Неизвестный параметр поиска")
}

// Ошибки заметки
func ErrGetNoteByIDError() error {
	return errors.New("Ошибка получения заметки по ID")
}

func ErrGetNoteByNameError() error {
	return errors.New("Ошибка получения заметки по имени")
}

func ErrGetAllNotesError() error {
	return errors.New("Ошибка получения всех заметок")
}

func ErrAddNoteError() error {
	return errors.New("Ошибка добавления заметки")
}

func ErrDeleteNoteError() error {
	return errors.New("Ошибка удаления заметки")
}

func ErrUpdateNoteContentError() error {
	return errors.New("Ошибка обновления содержимого заметки")
}

func ErrUpdateNoteInfoError() error {
	return errors.New("Ошибка обновления информации заметки")
}

func ErrAddNoteToCollectionError() error {
	return errors.New("Ошибка добавления заметки в коллекцию")
}

func ErrDeleteNoteFromCollectionError() error {
	return errors.New("Ошибка удаления заметки из коллекции")
}

// Ошибки Подборки
func ErrGetCollectionByIDError() error {
	return errors.New("Ошибка получения коллекции по ID")
}

func ErrGetCollectionByNameError() error {
	return errors.New("Ошибка получения коллекции по имени")
}

func ErrGetAllCollectionsError() error {
	return errors.New("Ошибка получения всех коллекций")
}

func ErrGetAllUserCollectionsError() error {
	return errors.New("Ошибка получения всех коллекций пользователя")
}

func ErrAddCollectionError() error {
	return errors.New("Ошибка добавления коллекции")
}

func ErrDeleteCollectionError() error {
	return errors.New("Ошибка удаления коллекции")
}

func ErrUpdateCollectionError() error {
	return errors.New("Ошибка обновления коллекции")
}

func ErrGetAllNotesInCollectionError() error {
	return errors.New("Ошибка получения всех Записок в коллекции коллекции")
}

// Ошибки раздела
func ErrGetSectionByIDError() error {
	return errors.New("Ошибка получения раздела по ID")
}

func ErrGetSectionByTeamNameError() error {
	return errors.New("Ошибка получения раздела по имени команды")
}

func ErrGetAllSectionsError() error {
	return errors.New("Ошибка получения всех разделов")
}

func ErrGetAllNotesInSectionError() error {
	return errors.New("Ошибка получения всех Записок в разделе")
}

func ErrAddSectionError() error {
	return errors.New("Ошибка добавления раздела")
}

func ErrDeleteSectionError() error {
	return errors.New("Ошибка удаления раздела")
}

func ErrUpdateSectionError() error {
	return errors.New("Ошибка обновления раздела")
}

func ErrAddNoteToSectionError() error {
	return errors.New("Ошибка добавления Записки в раздел")
}

func ErrGetFullStatError() error {
	return errors.New("Ошибка получения статистики")
}

func CreateError(errNum int, err error, funcName string) *MyError {
	myErr := new(MyError)

	myErr.ErrNum = errNum
	myErr.FuncName = funcName
	myErr.Err = err

	return myErr
}
