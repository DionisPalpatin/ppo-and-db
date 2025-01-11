package e2etest

import (
	bl "github.com/DionisPalpatin/ppo-and-db/application/internal/business-logic"
	da "github.com/DionisPalpatin/ppo-and-db/application/internal/data-access"
	mylogger "github.com/DionisPalpatin/ppo-and-db/application/internal/logger"
	"github.com/DionisPalpatin/ppo-and-db/application/internal/config"
	data_builders "github.com/DionisPalpatin/ppo-and-db/application/internal/models/builders"

	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type E2ETestSuite struct {
	suite.Suite
	db          *sql.DB
	authService bl.IOAuthService
	noteService bl.INoteService
	collService bl.ICollectionService
}

func (s *E2ETestSuite) BeforeEach(t provider.T) {
	var err error
	pgInfo := PostgresInfo{
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	s.db, err = InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if s.db == nil {
		return
	}

	dbconfig := config.DBConfigs{
		Host:       pgInfo.Host,
		User:       pgInfo.User,
		Password:   pgInfo.Password,
		Port:       pgInfo.Port,
		Name:       pgInfo.DBName,
		DB:         s.db,
		SchemaName: "interg_tests",
	}

	logger := mylogger.MyLogger{}
	logger.InitLogger("./logs/e2e.log", "debug")

	userRepo := da.NewUserRepository(&dbconfig, &logger)
	s.authService = bl.NewOAuthService(&userRepo, &logger)
	
	noteRepo := da.NewNoteRepository(&dbconfig, &logger)
	s.noteService = bl.NewNoteService(&noteRepo, &logger)

	collRepo := da.NewCollectionRepository(&dbconfig, &logger)
	s.collService = bl.NewCollectionService(&collRepo, &logger)
}

func (s *E2ETestSuite) AfterEach(t provider.T) {
	_ = s.db.Close()
}

func BeforeAll() {
	pgInfo := PostgresInfo{
		Host:     "localhost",
		User:     "postgres",
		Password: "password",
		Port:     15423,
		DBName:   "NotebookAppIntTests",
	}

	db, err := InitDB(pgInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	text, err := os.ReadFile("./init.sql")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec(string(text))
	if err != nil {
		fmt.Println(err)
		return
	}

	db.Close()
	fmt.Println("Соединение с БД в BeforeAll закрыто!")
}

func (s *E2ETestSuite) E2ETest_Login_AddNote_CreateCollection_AddNoteToCollection_Success(t provider.T) {
	t.Title("E2ETest_Login_AddNote_CreateCollection_AddNoteToCollection: Success")
	t.Tags("E2E")

	if IsUnitTestsFailed() || IsIntegrationTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Success", func(sCtx provider.StepCtx) {
		// Step 1
		login := "user1"
		password := "password1"

		token, err := s.authService.SignInUser(login, password)
		
		sCtx.Assert().NotNil(token)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)

		// Step 2
		note := data_builders.NewNoteBuilder().WithName("really unique name").WithOwnerID(1).WithSectionID(1).Build()
		user := data_builders.NewUserBuilder().WithLogin("user1").WithPassword("password1").WithRole(bl.Admin).Build()

		note.Id, err = s.noteService.AddNote(note, user)

		sCtx.Assert().NotEqual(0, note.Id)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)

		// Step 3
		collection := data_builders.NewCollectionBuilder().WithName("Teeeeeest collection").WithOwnerID(user.Id).Build()

		collection.Id, err = s.collService.AddCollection(collection)

		sCtx.Assert().NotEqual(0, collection.Id)
		sCtx.Assert().Equal(bl.Ok, err.ErrNum)

		// Step 4
		err = s.noteService.AddNoteToCollection(note.Id, collection.Id)

		sCtx.Assert().Equal(bl.Ok, err.ErrNum)
	})
}

func TestIntegrationRunner(t *testing.T) {
	BeforeAll()

	suite.RunSuite(t, new(E2ETestSuite))
}