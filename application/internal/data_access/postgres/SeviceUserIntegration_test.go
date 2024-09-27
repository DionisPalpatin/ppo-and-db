package postgres

// import (
// 	"testing"
// 	"time"
//
// 	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/config"
// 	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/business_logic"
// 	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/database"
// 	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
// )

// func TestIntegrationGetUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "postgres",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "postgres",
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
// 	user, myErr := us.GetUser(1, "", bl.SearchByID, &models.User{Id: 1, Role: bl.Admin}, ur)
//
// 	if myErr.ErrNum != bl.AllIsOk {
// 		t.Errorf("GetUser returned an error: %s", myErr.Err)
// 	}
// 	if user.Id != 1 {
// 		t.Errorf("GetUser returned wrong user ID, expected: 1, real: %d", user.Id)
// 	}
// }

// func TestIntegrationUpdateUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "postgres",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "postgres",
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
// 	myErr := us.UpdateUser(&models.User{Id: 1, Role: bl.Admin}, &models.User{
// 		Id:               1,
// 		Fio:              "stepanov stepan",
// 		RegistrationDate: time.Now(),
// 		Login:            "updateduser",
// 		Password:         "updatedpassword",
// 		Role:             bl.Author,
// 	}, ur)
//
// 	// Проверка результатов
// 	if myErr.ErrNum != bl.AllIsOk {
// 		t.Errorf("UpdateUser returned an error: %s", myErr.Err)
// 	}
// }

// func TestIntegrationDeleteUser(t *testing.T) {
// 	dbConfigs := config.DBConfigs{
// 		Host:       "localhost",
// 		Port:       5432,
// 		User:       "postgres",
// 		Password:   "postgrespassword",
// 		Name:       "NotebookApp",
// 		SchemaName: "test",
// 		DriverName: "postgres",
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
// 	myErr := us.DeleteUser(&models.User{Id: 1, Role: bl.Admin}, 1, ur)
//
// 	if myErr.ErrNum != bl.AllIsOk {
// 		t.Errorf("DeleteUser returned an error: %s", myErr.Err)
// 	}
// }
