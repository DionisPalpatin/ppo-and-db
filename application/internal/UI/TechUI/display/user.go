package display

import (
	"fmt"

	"notebook_app/internal/models"
)

func DisplayUser(user *models.User) {
	if user != nil {
		fmt.Println("## Пользователь:")
		fmt.Printf("ID: %d\n", user.Id)
		fmt.Printf("ФИО: %s\n", user.Fio)
		fmt.Printf("Дата регистрации: %s\n", user.RegistrationDate)
		fmt.Printf("Роль: %d\n", user.Role)
		fmt.Printf("login: %s\n", user.Login)
		fmt.Printf("Password: %s\n\n", user.Password)
	} else {
		fmt.Println("## Пользователь не найден:")
	}
}

func DisplayAllUsers(users []*models.User) {
	for i, us := range users {
		fmt.Printf("%d) %s\n", i, us.Fio)
	}
}
