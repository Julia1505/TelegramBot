package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type User struct {
	ChatID   int64
	Username string
}

type Task struct {
	ID       int
	name     string
	creator  User
	assignee User
}

type TaskStorage struct {
	storage []Task
	count   int
}

func CreateTaskCollection() *TaskStorage {
	return &TaskStorage{
		make([]Task, 0),
		1,
	}
}

func (t *TaskStorage) AddTask(update tgbotapi.Update) string {
	newTask := Task{
		ID:      t.count,
		creator: User{ChatID: update.Message.Chat.ID, Username: update.Message.From.UserName},
	}

	t.count++
	t.storage = append(t.storage, newTask)
	return strconv.Itoa(newTask.ID) + ". " + newTask.name + " by @" + newTask.creator.Username
}

func (t *TaskStorage) DeleteTask(index int) {
	copy(t.storage[index:], t.storage[index+1:])
	t.storage[len(t.storage)-1] = Task{}
	t.storage = t.storage[:len(t.storage)-1]
}

func (t *TaskStorage) ModifyTask(update tgbotapi.Update) {

}

func (t *TaskStorage) Show(update tgbotapi.Update, assignee User, creator User) string {
	var res string
	for i, task := range t.storage {
		res += strconv.Itoa(task.ID) + ". " + task.name + " by @" + task.creator.Username

		if i != len(t.storage)-1 {
			res += "\n"
		}
	}

	if res == "" {
		res = "Нет задач"
	}

	return res
}
