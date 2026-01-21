package userhandler

import (
	"log"

	userservice "github.com/Ferari430/tg_sender/internal/service/userService"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter interface {
	Successes(id int64, text string) error
	Error(id int64, text string) error
	Welcome(id int64, name string) error
	Files(id int64, text string) error
}

const (
	startCommand = "start"
	helpCommand  = "help"
	myFiles      = "MyFiles"
)

type UserHandler struct {
	P Presenter
	s *userservice.UserService
}

func NewUserHandler(userService *userservice.UserService, p Presenter) *UserHandler {
	return &UserHandler{
		P: p,
		s: userService,
	}
}

func (u *UserHandler) HandleMessage(msg *tgbotapi.Message) {
	id := msg.Chat.ID
	if msg.IsCommand() {
		switch msg.Command() {

		case startCommand:

			dto := userservice.UserDTO{
				ChatID:   id,
				Username: msg.From.UserName,
			}

			err := u.s.Start(dto)
			if err != nil {
				log.Println(err)
				return
			}

			err = u.P.Welcome(id, dto.Username)
			if err != nil {
				log.Println("Error:", err)
				return
			}
		//todo
		case helpCommand:
			u.s.Help()

		case myFiles:
			u.s.Files()

			err := u.P.Files(msg.Chat.ID, msg.Text)

			if err != nil {
				log.Println(err)
				return
			}

		default:
			log.Println("unknown command: ", msg.Command())
		}
	}

}
