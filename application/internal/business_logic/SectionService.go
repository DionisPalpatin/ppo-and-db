package bl

import "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"

type SectionService struct{}

func (SectionService) GetSection(secID int, name string, user *models.User, searchBy int, isr ISectionRepository) (*models.Section, *MyError) {
	// if user.Role != Admin {
	// 	err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetSection")
	// 	return nil, err
	// }

	var section *models.Section
	var err *MyError

	switch searchBy {
	case SearchByID:
		section, err = isr.GetSectionByID(secID)

	case SearchByString:
		section, err = isr.GetSectionByTeamName(name)

	default:
		section = nil
		err = CreateError(ErrSearchParameter, ErrSearchParameterError(), "GetSection")
	}

	return section, err
}

func (SectionService) GetAllSections(user *models.User, isr ISectionRepository) ([]*models.Section, *MyError) {
	if user.Role != Admin {
		err := CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllSections")
		return nil, err
	}

	return isr.GetAllSections()
}

func (SectionService) GetAllNotesInSection(secID int, user *models.User, isr ISectionRepository, itr ITeamRepository) ([]*models.Note, *MyError) {
	sec, err := isr.GetSectionByID(secID)
	if err.ErrNum != AllIsOk {
		return nil, err
	}

	var team1 *models.Team
	team1, err = itr.GetTeamBySectionID(secID)
	if err.ErrNum != AllIsOk {
		return nil, err
	}

	var team2 *models.Team
	team2, err = itr.GetUserTeam(user)
	if err.ErrNum != AllIsOk {
		return nil, err
	}

	if team1.Id != team2.Id || user.Role != Admin {
		return nil, CreateError(ErrAccessDenied, ErrAccessDeniedError(), "GetAllNotesInSection")
	}

	return isr.GetAllNotesInSection(sec)
}

func (SectionService) AddNoteToSection(section *models.Section, note *models.Note, user *models.User, isr ISectionRepository, itr ITeamRepository) *MyError {
	team1, err := itr.GetTeamBySectionID(section.Id)
	if err.ErrNum != AllIsOk {
		return err
	}

	var team2 *models.Team
	team2, err = itr.GetUserTeam(user)
	if err.ErrNum != AllIsOk {
		return err
	}

	if team1.Id != team2.Id && user.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddNoteToSection")
	}

	return isr.AddNoteToSection(note, section)
}

func (SectionService) DeleteNoteFromSection(section *models.Section, note *models.Note, user *models.User, isr ISectionRepository, itr ITeamRepository) *MyError {
	team1, err := itr.GetTeamBySectionID(section.Id)
	if err.ErrNum != AllIsOk {
		return err
	}

	var team2 *models.Team
	team2, err = itr.GetUserTeam(user)
	if err.ErrNum != AllIsOk {
		return err
	}

	if team1.Id != team2.Id && user.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddNoteToSection")
	}

	return isr.DeleteNoteFromSection(note, section)
}

func (SectionService) AddSection(section *models.Section, team *models.Team, user *models.User, isr ISectionRepository) *MyError {
	if user.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "AddSection")
	}

	return isr.AddSection(section, team)
}

func (SectionService) DeleteSection(id int, user *models.User, isr ISectionRepository) *MyError {
	if user.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "DeleteSection")
	}

	return isr.DeleteSection(id)
}

func (SectionService) UpdateSection(section *models.Section, user *models.User, isr ISectionRepository) *MyError {
	if user.Role != Admin {
		return CreateError(ErrAccessDenied, ErrAccessDeniedError(), "UpdateSection")
	}

	return isr.UpdateSection(section)
}
