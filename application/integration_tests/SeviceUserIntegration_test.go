package IntegrationTests

// import (
// 	"testing"
// 	"time"
//
// 	"github.com/DionisPalpatin/ppo-and-db/application/config"
// 	"github.com/DionisPalpatin/ppo-and-db/application/internal/business_logic"
// 	"github.com/DionisPalpatin/ppo-and-db/application/internal/database"
// 	"github.com/DionisPalpatin/ppo-and-db/application/internal/models"
// )

// func TestIntegrationGetUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "data_access",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "data_access",
// 		DB:         nil,
// 	}
// 	func(err error) {
// 		if err != nil {
// 			panic(err)
// 		}
// 	}(database.Connect(&dbConfigs))
//
// 	ur := UserRepository{DbConfigs: &dbConfigs}
// 	us := bl.UserService{}
//
// 	user, resState := us.GetUser(1, "", bl.SearchByID, &models.User{Id: 1, Role: bl.Admin}, ur)
//
// 	if resState.ErrNum != bl.AllIsOk {
// 		t.Errorf("GetUser returned an error: %s", resState.Err)
// 	}
// 	if user.Id != 1 {
// 		t.Errorf("GetUser returned wrong user ID, expected: 1, real: %d", user.Id)
// 	}
// }

// func TestIntegrationUpdateUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "data_access",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "data_access",
// 		DB:         nil,
// 	}
// 	func(err error) {
// 		if err != nil {
// 			panic(err)
// 		}
// 	}(database.Connect(&dbConfigs))
//
// 	ur := UserRepository{DbConfigs: &dbConfigs}
// 	us := bl.UserService{}
//
// 	resState := us.UpdateUser(&models.User{Id: 1, Role: bl.Admin}, &models.User{
// 		Id:               1,
// 		Fio:              "stepanov stepan",
// 		RegistrationDate: time.Now(),
// 		Login:            "updateduser",
// 		Password:         "updatedpassword",
// 		Role:             bl.Author,
// 	}, ur)
//
// 	// Проверка результатов
// 	if resState.ErrNum != bl.AllIsOk {
// 		t.Errorf("UpdateUser returned an error: %s", resState.Err)
// 	}
// }

// func TestIntegrationDeleteUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "data_access",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "data_access",
// 		DB:         nil,
// 	}
// 	func(err error) {
// 		if err != nil {
// 			panic(err)
// 		}
// 	}(database.Connect(&dbConfigs))
//
// 	ur := UserRepository{DbConfigs: &dbConfigs}
// 	us := bl.UserService{}
//
// 	resState := us.DeleteUser(&models.User{Id: 1, Role: bl.Admin}, 1, ur)
//
// 	if resState.ErrNum != bl.AllIsOk {
// 		t.Errorf("DeleteUser returned an error: %s", resState.Err)
// 	}
// }
