package bl

import "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"

type SectionService struct {
	isr ISectionRepository
	itr ITeamRepository
}

func (ss *SectionService) GetSection(secID int, name string, requester *models.User, searchBy int) (*models.Section, *MyError) {
	var section *models.Section
	var err *MyError

	if requester.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "DeleteNote", "bl")
	}

	switch searchBy {
	case SearchByID:
		section, err = ss.isr.GetSectionByID(secID)

	case SearchByString:
		section, err = ss.isr.GetSectionByTeamName(name)

	default:
		section = nil
		err = CreateError(ErrSearchParameter, "GetSection", "bl")
	}

	return section, err
}

func (ss *SectionService) GetAllSections(user *models.User) ([]*models.Section, *MyError) {
	if user.Role != Admin {
		err := CreateError(ErrAccessDenied, "GetAllSections", "bl")
		return nil, err
	}

	return ss.isr.GetAllSections()
}

func (ss *SectionService) GetAllNotesInSection(secID int, user *models.User) ([]*models.Note, *MyError) {
	sec, err := ss.isr.GetSectionByID(secID)
	if err.ErrNum != Ok {
		return nil, err
	}

	var team1 *models.Team
	team1, err = ss.itr.GetTeamBySectionID(secID)
	if err.ErrNum != Ok {
		return nil, err
	}

	var team2 *models.Team
	team2, err = ss.itr.GetUserTeam(user)
	if err.ErrNum != Ok {
		return nil, err
	}

	if team1.Id != team2.Id || user.Role != Admin {
		return nil, CreateError(ErrAccessDenied, "GetAllNotesInSection", "bl")
	}

	return ss.isr.GetAllNotesInSection(sec)
}

func (ss *SectionService) AddNoteToSection(section *models.Section, note *models.Note, user *models.User) *MyError {
	team1, err := ss.itr.GetTeamBySectionID(section.Id)
	if err.ErrNum != Ok {
		return err
	}

	var team2 *models.Team
	team2, err = ss.itr.GetUserTeam(user)
	if err.ErrNum != Ok {
		return err
	}

	if team1.Id != team2.Id && user.Role != Admin {
		return CreateError(ErrAccessDenied, "AddNoteToSection", "bl")
	}

	return ss.isr.AddNoteToSection(note, section)
}

func (ss *SectionService) DeleteNoteFromSection(section *models.Section, note *models.Note, user *models.User) *MyError {
	team1, err := ss.itr.GetTeamBySectionID(section.Id)
	if err.ErrNum != Ok {
		return err
	}

	var team2 *models.Team
	team2, err = ss.itr.GetUserTeam(user)
	if err.ErrNum != Ok {
		return err
	}

	if team1.Id != team2.Id && user.Role != Admin {
		return CreateError(ErrAccessDenied, "AddNoteToSection", "bl")
	}

	return ss.isr.DeleteNoteFromSection(note, section)
}

func (ss *SectionService) AddSection(section *models.Section, team *models.Team, user *models.User) (int, *MyError) {
	if user.Role != Admin {
		return 0, CreateError(ErrAccessDenied, "AddSection", "bl")
	}

	return ss.isr.AddSection(section, team)
}

func (ss *SectionService) DeleteSection(id int, user *models.User) *MyError {
	if user.Role != Admin {
		return CreateError(ErrAccessDenied, "DeleteSection", "bl")
	}

	return ss.isr.DeleteSection(id)
}

func (ss *SectionService) UpdateSection(section *models.Section, user *models.User) *MyError {
	if user.Role != Admin {
		return CreateError(ErrAccessDenied, "UpdateSection", "bl")
	}

	return ss.isr.UpdateSection(section)
}
