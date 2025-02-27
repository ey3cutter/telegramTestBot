package bot

import (
	"database/sql"
	_ "fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramTestBot/db"
)

func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/start" {
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("🗂 Товары"),
						tgbotapi.NewKeyboardButton("👨💼 Написать менеджеру"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Корзина"),
					),
				)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Выберите действие:")
				msg.ReplyMarkup = keyboard
				if _, err := b.Bot.Send(msg); err != nil {
					log.Printf("Ошибка отправки: %v", err) // Не паникуем, а логируем
				}
			} else if update.Message.Text == "🗂 Товары" {
				b.SendCategoryMenu(update.Message.Chat.ID)
			} else if update.Message.Text == "👨💼 Написать менеджеру" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Свяжитесь с менеджером: @MrBuris")
				b.Bot.Send(msg)
			}

		}
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "catalog":
				b.SendCategoryMenu(update.CallbackQuery.Message.Chat.ID)
			case "main_menu":
				// Создаем главное меню
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("🗂 Товары"),
						tgbotapi.NewKeyboardButton("👨💼 Написать менеджеру"),
					),
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("🛒 Корзина"),
						tgbotapi.NewKeyboardButton("ℹ️ Помощь"),
					),
				)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Главное меню:")
				msg.ReplyMarkup = keyboard
				b.Bot.Send(msg)

				// Подтверждаем callback
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				b.Bot.Request(callback)
			default:
				if strings.HasPrefix(update.CallbackQuery.Data, "product_") {
					categoryID, err := strconv.Atoi(strings.TrimPrefix(update.CallbackQuery.Data, "product_"))
					if err != nil {
						log.Panic(err)
					}
					b.SendSubCategoryMenu(update.CallbackQuery.Message.Chat.ID, categoryID)
				} else {
					b.SendProductList(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				}

				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				if _, err := b.Bot.Request(callback); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}

func (b *Bot) SendCategoryMenu(chatID int64) {
	service := db.NewService(db.NewRepository())
	categories, err := service.GetCategories()
	if err != nil {
		log.Panic(err)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, category := range categories {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, "product_"+strconv.Itoa(category.ID)),
		)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "main_menu"),
	))

	msg := tgbotapi.NewMessage(chatID, "Выберите категорию товаров:")
	msg.ReplyMarkup = keyboard
	if _, err := b.Bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *Bot) SendSubCategoryMenu(chatID int64, categoryID int) {
	service := db.NewService(db.NewRepository())
	subCategories, err := service.GetSubCategories(categoryID)
	if err != nil {
		log.Panic(err)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, subCategory := range subCategories {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subCategory.Name, subCategory.Callback),
		)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "catalog"),
	))

	msg := tgbotapi.NewMessage(chatID, "Выберите подкатегорию:")
	msg.ReplyMarkup = keyboard
	if _, err := b.Bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *Bot) SendProductList(chatID int64, callback string) {
	service := db.NewService(db.NewRepository())
	var subCategoryID int
	rows, err := db.Db.Query("SELECT id FROM subcategories WHERE callback = $1", callback)
	if err != nil {
		log.Panic(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)

	for rows.Next() {
		err = rows.Scan(&subCategoryID)
		if err != nil {
			log.Panic(err)
		}
	}

	products, err := service.GetProducts(subCategoryID)
	if err != nil {
		log.Panic(err)
	}

	var msgText string
	for _, product := range products {
		msgText += "- " + product.Name + ": " + product.Description + "\n"
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	if _, err := b.Bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
