package bl

import "github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"

type StatService struct{}

func (StatService) GetFullStat(req *models.User, isr IStatisticRepository) (*models.Stat, *MyError) {
	if req.Role != Admin {
		myErr := CreateError(ErrGetFullStat, ErrGetFullStatError(), "GetFullStat")
		return nil, myErr
	}

	stat := new(models.Stat)

	stat.TotalMarks = isr.CountMarks()
	stat.TotalNotes = isr.CountNotes()
	stat.TotalCollections = isr.CountCollections()
	stat.TotalTeams = isr.CountTeams()
	stat.TotalSections = isr.CountSections()
	stat.TotalUsers = isr.CountUsers()

	myOk := CreateError(AllIsOk, nil, "")
	return stat, myOk
}
