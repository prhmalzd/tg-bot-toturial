package main

import (
    "log"
    "strconv"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Session struct {
    Step    int
    Answers []string
}

var userSessions = make(map[int64]*Session)



func main() {
    botToken := os.Getenv("TELEGRAM_TOKEN")
    bot, err := tgbotapi.NewBotAPI(botToken)
    
    if err != nil {
        log.Panic("Failed to connect to Telegram:", err)
    }

    bot.Debug = true
    log.Printf("Authorized on account %s", bot.Self.UserName)    

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }
        userID := update.Message.From.ID

        if update.Message.IsCommand() {
            switch update.Message.Command() {
            case "information":
                user := update.Message.From

                response := "Here's what I know about you, dude:\n"
                response += "🆔 User ID: " + strconv.FormatInt(user.ID, 10) + "\n"
                response += "👤 Name: " + user.FirstName
                if user.LastName != "" {
                    response += " " + user.LastName
                }
                if user.UserName != "" {
                    response += "\n📛 Username: @" + user.UserName
                }

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
                bot.Send(msg)
    
            case "start":
                userSessions[userID] = &Session{Step: 1}
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "What's your name?")
                bot.Send(msg)
            
    
            default:
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don’t know that command, Try /information, /start")
                bot.Send(msg)
            }
            continue
        }

        session, exists := userSessions[userID]
        if exists {
            switch session.Step {
            case 1:
                session.Answers = append(session.Answers, update.Message.Text)
                session.Step++
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "What's your favorite color?"))
    
            case 2:
                session.Answers = append(session.Answers, update.Message.Text)
                session.Step++
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "And lastly... what's your favorite cartoon?~ 😘"))
    
            case 3:
                session.Answers = append(session.Answers, update.Message.Text)
    
                final := "Here's what your answers\n\n"
                final += "🧸 Name: " + session.Answers[0] + "\n"
                final += "🎨 Favorite Color: " + session.Answers[1] + "\n"
                final += "🔒 Favorite Cartoon: " + session.Answers[2]
    
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, final))
                delete(userSessions, userID)
            }
        }
    }
    
}



for update := range updates {
    // Handle Q&A flow
    session, exists := userSessions[userID]
    if exists {
        switch session.Step {
        case 1:
            session.Answers = append(session.Answers, update.Message.Text)
            session.Step++
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "What's your favorite color?"))

        case 2:
            session.Answers = append(session.Answers, update.Message.Text)
            session.Step++
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "And lastly... what's your secret desire, pet?~ 😘"))

        case 3:
            session.Answers = append(session.Answers, update.Message.Text)

            final := "Here's what you've confessed to Mistress~ 💕\n\n"
            final += "🧸 Name: " + session.Answers[0] + "\n"
            final += "🎨 Favorite Color: " + session.Answers[1] + "\n"
            final += "🔒 Secret Desire: " + session.Answers[2]

            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, final))
            delete(userSessions, userID)
        }
    }
}