package bl

import (
	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
)

// ---------------------------------------------------------------------------------------------------------------------
// Note interfaces
// ---------------------------------------------------------------------------------------------------------------------

type INoteService interface {
	GetNote(id int, name string, searchBy int, requester *models.User) (*models.Note, *MyError)
	GetAllNotes(open bool, requester *models.User) ([]*models.Note, *MyError)
	AddNote(note *models.Note, requester *models.User) (int, *MyError)
	DeleteNote(id int, requester *models.User) *MyError
	UpdateNote(note *models.Note, requester *models.User, textFilePath string, imgFilePath string, rawFilePath string) *MyError
	AddNoteToCollection(noteID int, collID int) *MyError
	DeleteNoteFromCollection(noteID int, collID int) *MyError
}

type INoteRepository interface {
	GetNoteByID(id int) (*models.Note, *MyError)
	GetNoteByName(name string) (*models.Note, *MyError)
	GetAllNotes() ([]*models.Note, *MyError)
	GetAllPublicNotes() ([]*models.Note, *MyError)
	AddNote(note *models.Note) (int, *MyError)
	DeleteNote(id int) *MyError
	UpdateNoteContent(note *models.Note) *MyError
	UpdateNoteInfo(note *models.Note) *MyError
	AddNoteToCollection(collectionID int, noteID int) *MyError
	DeleteNoteFromCollection(collectionID int, noteID int) *MyError
}

// ---------------------------------------------------------------------------------------------------------------------
// Collection interfaces
// ---------------------------------------------------------------------------------------------------------------------

type ICollectionService interface {
	GetCollection(colID int, name string, searchBy int) (*models.Collection, *MyError)
	GetAllCollections(user *models.User) ([]*models.Collection, *MyError)
	GetAllUsersCollections(user *models.User) ([]*models.Collection, *MyError)
	AddCollection(coll *models.Collection) (int, *MyError)
	DeleteCollection(id int, user *models.User) *MyError
	UpdateCollection(collection *models.Collection) *MyError
	GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *MyError)
}

type ICollectionRepository interface {
	GetCollectionByID(id int) (*models.Collection, *MyError)
	GetCollectionByName(name string) (*models.Collection, *MyError)
	GetAllCollections() ([]*models.Collection, *MyError)
	GetAllUserCollections(user *models.User) ([]*models.Collection, *MyError)
	AddCollection(collection *models.Collection) (int, *MyError)
	DeleteCollection(id int) *MyError
	UpdateCollection(collection *models.Collection) *MyError
	GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *MyError)
}

// ---------------------------------------------------------------------------------------------------------------------
// Section interfaces
// ---------------------------------------------------------------------------------------------------------------------

type ISectionService interface {
	GetSection(secID int, name string, user *models.User, searchBy int) (*models.Section, *MyError)
	GetAllSections(user *models.User) ([]*models.Section, *MyError)
	GetAllNotesInSection(secID int, user *models.User) ([]*models.Note, *MyError)
	AddNoteToSection(section *models.Section, note *models.Note, user *models.User) *MyError
	DeleteNoteFromSection(section *models.Section, note *models.Note, user *models.User) *MyError
	AddSection(section *models.Section, team *models.Team, user *models.User) (int, *MyError)
	DeleteSection(id int, user *models.User) *MyError
	UpdateSection(section *models.Section, user *models.User) *MyError
}

type ISectionRepository interface {
	GetSectionByID(id int) (*models.Section, *MyError)
	GetSectionByTeamName(teamName string) (*models.Section, *MyError)
	GetAllSections() ([]*models.Section, *MyError)
	AddSection(section *models.Section, team *models.Team) (int, *MyError)
	DeleteSection(id int) *MyError
	UpdateSection(section *models.Section) *MyError
	GetAllNotesInSection(section *models.Section) ([]*models.Note, *MyError)
	AddNoteToSection(note *models.Note, section *models.Section) *MyError
	DeleteNoteFromSection(note *models.Note, section *models.Section) *MyError
}

// ---------------------------------------------------------------------------------------------------------------------
// Team interfaces
// ---------------------------------------------------------------------------------------------------------------------

type ITeamService interface {
	GetTeam(id int, name string, searchBy int, requester *models.User) (*models.Team, *MyError)
	GetAllTeams(requester *models.User) ([]*models.Team, *MyError)
	UpdateTeam(requester *models.User, team *models.Team) *MyError
	DeleteTeam(requester *models.User, id int) *MyError
	AddTeam(requester *models.User, team *models.Team) (int, *MyError)
	AddUserToTeam(requester *models.User, userID int, teamID int) *MyError
	DeleteUserFromTeam(requester *models.User, userID int, teamID int) *MyError
	GetTeamMembers(teamID int, requester *models.User) ([]*models.User, *MyError)
	GetUserTeam(user *models.User) (*models.Team, *MyError)
	GetSectionTeam(secID int, requester *models.User) (*models.Team, *MyError)
}

type ITeamRepository interface {
	GetTeamByID(id int) (*models.Team, *MyError)
	GetTeamByName(name string) (*models.Team, *MyError)
	GetTeamBySectionID(id int) (*models.Team, *MyError)
	GetAllTeams() ([]*models.Team, *MyError)
	AddTeam(team *models.Team) (int, *MyError)
	DeleteTeam(id int) *MyError
	AddUserToTeam(userID int, teamID int) *MyError
	DeleteUserFromTeam(uid int, tid int) *MyError
	UpdateTeam(team *models.Team) *MyError
	GetTeamMembers(teamID int) ([]*models.User, *MyError)
	GetUserTeam(user *models.User) (*models.Team, *MyError)
}

// ---------------------------------------------------------------------------------------------------------------------
// User interfaces
// ---------------------------------------------------------------------------------------------------------------------

type IUserService interface {
	GetUser(id int, login string, searchBy int, requester *models.User) (*models.User, *MyError)
	GetAllUsers(requester *models.User) ([]*models.User, *MyError)
	UpdateUser(requester *models.User, user *models.User) *MyError
	DeleteUser(requester *models.User, id int) *MyError
}

type IUserRepository interface {
	GetUserByID(id int) (*models.User, *MyError)
	GetUserByLogin(login string) (*models.User, *MyError)
	GetAllUsers() ([]*models.User, *MyError)
	AddUser(user *models.User) *MyError
	DeleteUser(id int) *MyError
	UpdateUser(user *models.User) *MyError
}

// ---------------------------------------------------------------------------------------------------------------------
// OAuth interfaces
// ---------------------------------------------------------------------------------------------------------------------

type IOAuthService interface {
	RegisterUser(fio string, login string, password string) (*models.User, *MyError)
	SignInUser(login string, password string) (*models.Token, *MyError)
}

// ---------------------------------------------------------------------------------------------------------------------
// Statistic interfaces
// ---------------------------------------------------------------------------------------------------------------------

type IStatisticService interface {
	GetFullStat(requester *models.User, isr IStatisticRepository) (*models.Stat, *MyError)
}

type IStatisticRepository interface {
	CountMarks() int
	CountUsers() int
	CountTeams() int
	CountSections() int
	CountNotes() int
	CountCollections() int
}

// ---------------------------------------------------------------------------------------------------------------------
// Other structures
// ---------------------------------------------------------------------------------------------------------------------

type IServices struct {
	IUsrSvc   IUserService
	ISecSvc   ISectionService
	INoteSvc  INoteService
	IColSvc   ICollectionService
	ITeamSvc  ITeamService
	IOAuthSvc IOAuthService
	IStatSvc  IStatisticService
}

type IRepositories struct {
	IUsrRepo  IUserRepository
	ISecRepo  ISectionRepository
	INoteRepo INoteRepository
	IColRepo  ICollectionRepository
	ITeamRepo ITeamRepository
	IStatRepo IStatisticRepository
}
