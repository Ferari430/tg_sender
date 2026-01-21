package out

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgPresenter struct {
	bot *tgbotapi.BotAPI
}

func NewTgPresenter(b *tgbotapi.BotAPI) *TgPresenter {
	return &TgPresenter{bot: b}
}

func (tg *TgPresenter) Successes(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Error(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *TgPresenter) Welcome(id int64, name string) error {

	msg := tgbotapi.NewMessage(id, StartMessage)
	_, err := tg.bot.Send(msg)

	return err
}

func (tg *TgPresenter) Files(id int64, fileNames []string) error {
	m := fmt.Sprintf("Твои файлы: %s", strings.Join(fileNames, "\n "))

	msg := tgbotapi.NewMessage(id, m)
	_, err := tg.bot.Send(msg)

	return err
}

func (tg *TgPresenter) Message(id int64, text string) error {

	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)

	return err
}

func (tg *TgPresenter) Help(id int64) error {
	text := helpMessage
	msg := tgbotapi.NewMessage(id, text)
	_, err := tg.bot.Send(msg)

	return err
}

const (
	helpMessage = `📚 *Помощь по использованию бота*

*📦 Поддерживаемые форматы архивов:*
• ZIP (.zip) — основной формат
• TAR.GZ (.tar.gz, .tgz) — архивы Linux
• 7Z (.7z) — высокое сжатие  
• RAR (.rar) — архивы WinRAR

*⚠️ Ограничения:*
• Максимальный размер архива: 10 МБ
• В архиве должно быть не более 50 файлов
• Имена файлов должны быть в кодировке *UTF-8*
• Запрещены архивы с вложенными архивами

*📋 Основные команды:*
/start — Начать работу с ботом
/help — Эта справка
/files — Получить PDF файл из базы

*🔒 Безопасность:*
• Все файлы проверяются антивирусом
• Архивы распаковываются в изолированном окружении
• Ваши данные не передаются третьим лицам

*💡 Советы:*
1. Убедитесь, что архив не защищен паролем
2. Проверьте, что все файлы в архиве — нужные
3. Если архив большой, используйте .7z для лучшего сжатия

По всем вопросам обращайтесь к @Airfool`
				
	StartMessage = "Привет! 👋  \nЯ бот для работы с файлами 📦\n\nЧто я умею:\n• принимаю ZIP-архивы" +
		"  \n• сохраняю их на сервере  " +
		"\n• запоминаю, кто и что загрузил  " +
		"\n\nПросто отправь ZIP-файл в этот чат.\n\n" +
		"Доступные команды:\n" +
		"/help — помощь  \n/myfiles — посмотреть загруженные файлы"
)
