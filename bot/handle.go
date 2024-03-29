package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Answer struct {
	Task
	User
}

var helloMessage = `Приветики-пистолетики, я бот Apopope ;) 
Вот что я умею:
/tasks - вывод всех задач
/new задача - создание новой задачи
/assign_$ID - делает пользователя исполнителем задачи
/unassign_$ID - снимает задачу с текущего исполнителя
/resolve_$ID - выполняет задачу
/my - показывает задачи, которые назначены на меня
/owner - показывает задачи, которые были созданы мной`

var LIST = `{{if .}}
{{range .}}{{.Task.ID}}. {{.Task.Name}} by @{{.Task.Creator.Username}}{{if .Task.Assignee.Username}}
assignee: {{if (eq .User.Username .Task.Assignee.Username)}}я
/unassign_{{.Task.ID}} /resolve_{{.Task.ID}}{{else}}@{{.Assignee.Username}}{{end}}{{else}}
/assign_{{.Task.ID}}{{end}}
{{end}}
{{else}}Нет задач{{end}}`

var SPECLIST = `{{if .}}
{{range .}}{{.ID}}. {{.Name}} by @{{.Creator.Username}}{{if .Assignee.Username}}
/unassign_{{.ID}} /resolve_{{.ID}}{{else}}
/assign_{{.ID}}{{end}}
{{end}}
{{else}}Нет задач{{end}}`

var TASK = `Задача "{{.Name}}" создана, id={{.ID}}`

var ASSIGN = `Задача "{{.Task.Name}}" назначена на {{if (ne .User.Username .Task.Assignee.Username)}}@{{.Task.Assignee.Username}}{{else}}вас{{end}}`

var UNASSIGN = `{{if .Task.Name}}{{if .User.Username}}Задача "{{.Task.Name}}" осталась без исполнителя{{else}}Принято{{end}}{{else}}Задача не на вас{{end}}`

var RESOLVE = `Задача "{{.Task.Name}}" выполнена {{if (ne .User.Username .Task.Assignee.Username)}}@{{.Task.Assignee.Username}}{{end}}`

var patternAssign = `^/assign_\d+$`
var patternUnassign = `^/unassign_\d+$`
var patternResolve = `^/resolve_\d+$`

func (t *TelegramBot) HandleUpdates(st *TaskStorage, update tgbotapi.Update) {
	//fmt.Println(update.Message.Text)
	if message := update.Message.Text; message != "" {
		parseMessage := strings.SplitN(message, " ", 2)

		//fmt.Println(parseMessage[0])
		switch parseMessage[0] {
		case "hello":
			fallthrough
		case "/help":
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage))
		case "/apopope":
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "А по попе ШЛЁП!!!!\U0001FAF1🏻🍑"))
		case "/tasks":
			ShowAll(t, *st, update)
		case "/new":
			NewMessage(t, st, update)
		case "/my":
			ShowMy(t, *st, update)
		case "/owner":
			ShowMyCreate(t, *st, update)
		default:
			fmt.Println(parseMessage[0])
			if is, _ := regexp.MatchString(patternAssign, parseMessage[0]); is {
				re, _ := regexp.Compile(`\d+`)
				taskId, _ := strconv.Atoi(re.FindAllString(parseMessage[0], 1)[0])
				AssignUser(t, st, update, taskId)
			}

			if is, _ := regexp.MatchString(patternUnassign, parseMessage[0]); is {
				re, _ := regexp.Compile(`\d+`)
				taskId, _ := strconv.Atoi(re.FindAllString(parseMessage[0], 1)[0])
				UnssignUser(t, st, update, taskId)
			}

			if is, _ := regexp.MatchString(patternResolve, parseMessage[0]); is {
				re, _ := regexp.Compile(`\d+`)
				taskId, _ := strconv.Atoi(re.FindAllString(parseMessage[0], 1)[0])
				ResolveUser(t, st, update, taskId)
			}

		}

	}
}

func ResolveUser(t *TelegramBot, st *TaskStorage, update tgbotapi.Update, Id int) {
	currentUser := User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName}
	task := st.DeleteTask(Id, currentUser)
	if task != (Task{}) {
		var tmpl = template.New("resolve")
		tmpl, _ = tmpl.Parse(RESOLVE)
		buf := bytes.NewBufferString("")
		answer := Answer{
			Task: task,
			User: currentUser,
		}
		tmpl.Execute(buf, answer)
		t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
		buf.Reset()

		if currentUser.Username != task.Creator.Username {
			answer.User.Username = task.Creator.Username
			tmpl.Execute(buf, answer)
			t.bot.Send(tgbotapi.NewMessage(task.Creator.ChatID, buf.String()))
		}
	}
}

func UnssignUser(t *TelegramBot, st *TaskStorage, update tgbotapi.Update, Id int) {
	currentUser := User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName}
	task := st.ModifyTask(Id, User{}, currentUser)
	var tmpl = template.New("unassign")
	tmpl, _ = tmpl.Parse(UNASSIGN)
	buf := bytes.NewBufferString("")
	answer := Answer{
		Task: task,
	}
	tmpl.Execute(buf, answer)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
	buf.Reset()

	if currentUser.Username != task.Creator.Username && task != (Task{}) {
		answer.User.Username = task.Creator.Username
		tmpl.Execute(buf, answer)
		t.bot.Send(tgbotapi.NewMessage(task.Creator.ChatID, buf.String()))
	}
}

func AssignUser(t *TelegramBot, st *TaskStorage, update tgbotapi.Update, Id int) {
	newAssign := User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName}
	task := st.ModifyTask(Id, newAssign, User{})
	var tmpl = template.New("assign")
	tmpl, _ = tmpl.Parse(ASSIGN)
	buf := bytes.NewBufferString("")
	answer := Answer{
		Task: task,
		User: newAssign,
	}

	tmpl.Execute(buf, answer)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
	buf.Reset()

	if newAssign.Username != task.Creator.Username {
		answer.User.Username = task.Creator.Username
		tmpl.Execute(buf, answer)
		t.bot.Send(tgbotapi.NewMessage(task.Creator.ChatID, buf.String()))
	}

}

func NewMessage(t *TelegramBot, st *TaskStorage, update tgbotapi.Update) {
	name := strings.Trim(update.Message.Text, "/new ")
	newTask := Task{Name: name, Creator: User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName}}
	//fmt.Println(newTask)
	st.AddTask(&newTask)
	//fmt.Println(newTask)

	var tmpl = template.New("new")
	tmpl, _ = tmpl.Parse(TASK)
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, newTask)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
}

func ShowAll(t *TelegramBot, st TaskStorage, update tgbotapi.Update) {
	dataToShow := st.Get("", "")
	currentUser := User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName}
	answer := make([]Answer, len(dataToShow))
	for i, data := range dataToShow {
		answer[i] = Answer{Task: data, User: currentUser}
	}
	var tmpl = template.New("list")
	tmpl, _ = tmpl.Parse(LIST)
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, answer)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
}

func ShowMy(t *TelegramBot, st TaskStorage, update tgbotapi.Update) {
	dataToShow := st.Get(update.Message.From.UserName, "")
	var tmpl = template.New("list")
	tmpl, _ = tmpl.Parse(SPECLIST)
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, dataToShow)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
}

func ShowMyCreate(t *TelegramBot, st TaskStorage, update tgbotapi.Update) {
	dataToShow := st.Get("", update.Message.From.UserName)
	var tmpl = template.New("list")
	tmpl, _ = tmpl.Parse(SPECLIST)
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, dataToShow)
	t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, buf.String()))
}
