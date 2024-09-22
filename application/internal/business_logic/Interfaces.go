package bl

import (
	"io"

	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
)

// ---------------------------------------------------------------------------------------------------------------------
// Note interfaces
// ---------------------------------------------------------------------------------------------------------------------

type INoteService interface {
	GetNote(id int, name string, searchBy int, requester *models.User, inr INoteRepository, isr ISectionRepository, itr ITeamRepository) (*models.Note, []byte, *MyError)
	GetAllNotes(requester *models.User, inr INoteRepository) ([]*models.Note, *MyError)
	AddNote(note *models.Note, requester *models.User, inr INoteRepository) *MyError
	DeleteNote(id int, requester *models.User, inr INoteRepository) *MyError
	UpdateNoteContent(noteID int, requester *models.User, filePath string, inr INoteRepository) *MyError
	UpdateNoteInfo(requester *models.User, note *models.Note, inr INoteRepository) *MyError
	AddNoteToCollection(noteID int, collID int, inr INoteRepository) *MyError
	DeleteNoteFromCollection(noteID int, collID int, inr INoteRepository) *MyError
}

type INoteRepository interface {
	GetNoteByID(id int) (*models.Note, []byte, *MyError)
	GetNoteByName(name string) (*models.Note, []byte, *MyError)
	GetAllNotes() ([]*models.Note, *MyError)
	AddNote(note *models.Note) *MyError
	DeleteNote(id int) *MyError
	UpdateNoteContentText(reader io.Reader, note *models.Note) *MyError
	UpdateNoteContentImg(reader io.Reader, note *models.Note) *MyError
	UpdateNoteContentRawData(reader io.Reader, note *models.Note) *MyError
	UpdateNoteInfo(note *models.Note) *MyError
	AddNoteToCollection(collectionID int, noteID int) *MyError
	DeleteNoteFromCollection(collectionID int, noteID int) *MyError
}

// ---------------------------------------------------------------------------------------------------------------------
// Collection interfaces
// ---------------------------------------------------------------------------------------------------------------------

type ICollectionService interface {
	GetCollection(colID int, name string, searchBy int, icr ICollectionRepository) (*models.Collection, *MyError)
	GetAllCollections(user *models.User, icr ICollectionRepository) ([]*models.Collection, *MyError)
	GetAllUsersCollections(user *models.User, icr ICollectionRepository) ([]*models.Collection, *MyError)
	AddCollection(coll *models.Collection, icr ICollectionRepository) *MyError
	DeleteCollection(id int, user *models.User, icr ICollectionRepository) *MyError
	UpdateCollection(collection *models.Collection, icr ICollectionRepository) *MyError
	GetAllNotesInCollection(collection *models.Collection, icr ICollectionRepository) ([]*models.Note, *MyError)
}

type ICollectionRepository interface {
	GetCollectionByID(id int) (*models.Collection, *MyError)
	GetCollectionByName(name string) (*models.Collection, *MyError)
	GetAllCollections() ([]*models.Collection, *MyError)
	GetAllUserCollections(user *models.User) ([]*models.Collection, *MyError)
	AddCollection(collection *models.Collection) *MyError
	DeleteCollection(id int) *MyError
	UpdateCollection(collection *models.Collection) *MyError
	GetAllNotesInCollection(collection *models.Collection) ([]*models.Note, *MyError)
}

// ---------------------------------------------------------------------------------------------------------------------
// Section interfaces
// ---------------------------------------------------------------------------------------------------------------------

type ISectionService interface {
	GetSection(secID int, name string, user *models.User, searchBy int, isr ISectionRepository) (*models.Section, *MyError)
	GetAllSections(user *models.User, isr ISectionRepository) ([]*models.Section, *MyError)
	GetAllNotesInSection(secID int, user *models.User, isr ISectionRepository, itr ITeamRepository) ([]*models.Note, *MyError)
	AddNoteToSection(section *models.Section, note *models.Note, user *models.User, isr ISectionRepository, itr ITeamRepository) *MyError
	DeleteNoteFromSection(section *models.Section, note *models.Note, user *models.User, isr ISectionRepository, itr ITeamRepository) *MyError
	AddSection(section *models.Section, team *models.Team, user *models.User, isr ISectionRepository) *MyError
	DeleteSection(id int, user *models.User, isr ISectionRepository) *MyError
	UpdateSection(section *models.Section, user *models.User, isr ISectionRepository) *MyError
}

type ISectionRepository interface {
	GetSectionByID(id int) (*models.Section, *MyError)
	GetSectionByTeamName(teamName string) (*models.Section, *MyError)
	GetAllSections() ([]*models.Section, *MyError)
	AddSection(section *models.Section, team *models.Team) *MyError
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
	GetTeam(id int, name string, searchBy int, requester *models.User, itr ITeamRepository) (*models.Team, *MyError)
	GetAllTeams(requester *models.User, itr ITeamRepository) ([]*models.Team, *MyError)
	UpdateTeam(requester *models.User, team *models.Team, itr ITeamRepository) *MyError
	DeleteTeam(requester *models.User, id int, itr ITeamRepository) *MyError
	AddTeam(requester *models.User, team *models.Team, itr ITeamRepository) *MyError
	AddUserToTeam(requester *models.User, userID int, teamID int, itr ITeamRepository) *MyError
	DeleteUserFromTeam(requester *models.User, userID int, teamID int, itr ITeamRepository) *MyError
	GetTeamMembers(teamID int, requester *models.User, itr ITeamRepository) ([]*models.User, *MyError)
	GetUserTeam(user *models.User, itr ITeamRepository) (*models.Team, *MyError)
}

type ITeamRepository interface {
	GetTeamByID(id int) (*models.Team, *MyError)
	GetTeamByName(name string) (*models.Team, *MyError)
	GetTeamBySectionID(id int) (*models.Team, *MyError)
	GetAllTeams() ([]*models.Team, *MyError)
	AddTeam(team *models.Team) *MyError
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
	GetUser(id int, login string, searchBy int, requester *models.User, iur IUserRepository) (*models.User, *MyError)
	GetAllUsers(requester *models.User, iur IUserRepository) ([]*models.User, *MyError)
	UpdateUser(requester *models.User, user *models.User, iur IUserRepository) *MyError
	DeleteUser(requester *models.User, id int, iur IUserRepository) *MyError
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
	RegisterUser(fio string, login string, password string, iur IUserRepository) (*models.User, *MyError)
	SignInUser(login string, password string, iur IUserRepository) (*models.User, *MyError)
}

type IRepositories struct {
	IUsrRepo  IUserRepository
	ISecRepo  ISectionRepository
	INoteRepo INoteRepository
	IColRepo  ICollectionRepository
	ITeamRepo ITeamRepository
}

type IServices struct {
	IUsrSvc   IUserService
	ISecSvc   ISectionService
	INoteSvc  INoteService
	IColSvc   ICollectionService
	ITeamSvc  ITeamService
	IOAuthSvc IOAuthService
}
