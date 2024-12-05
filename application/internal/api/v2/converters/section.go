package converters

import (
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/api/v2/transport_models"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"time"
)

func ToSectionInfo(section *models.Section, team *models.Team) transport_models.SectionInfo {
	return transport_models.SectionInfo{
		TeamID: team.Id,
		ID:     section.Id,
	}
}

func ToSectionFullInfo(section *models.Section, team *models.Team) transport_models.SectionFullInfo {
	return transport_models.SectionFullInfo{
		ID:               section.Id,
		TeamID:           team.Id,
		RegistrationDate: section.CreationDate.String(),
	}
}

func FromSectionInfo(sectionInfo transport_models.SectionInfo) (models.Section, int) {
	return models.Section{
		Id: sectionInfo.ID,
	}, sectionInfo.TeamID
}

func FromSectionFullInfo(sectionFullInfo transport_models.SectionFullInfo) (models.Section, int, error) {
	parsedTime, err := time.Parse(sectionFullInfo.RegistrationDate, "2006-01-02 15:04:05-07")
	if err != nil {
		return models.Section{}, -1, err
	}

	return models.Section{
		Id:           sectionFullInfo.ID,
		CreationDate: parsedTime,
	}, sectionFullInfo.TeamID, nil
}
