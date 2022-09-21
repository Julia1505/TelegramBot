package main

type User struct {
	ChatID   int64
	Username string
}

type Task struct {
	ID       int
	Name     string
	Creator  User
	Assignee User
}

type TaskStorage struct {
	Storage []Task
	Count   int
}

func CreateTaskCollection() *TaskStorage {
	return &TaskStorage{
		make([]Task, 0),
		1,
	}
}

func (st *TaskStorage) AddTask(newTask *Task) {
	newTask.ID = st.Count
	st.Count++
	st.Storage = append(st.Storage, *newTask)
}

func (st *TaskStorage) DeleteTask(id int, user User) Task {
	for i := range st.Storage {
		if st.Storage[i].ID == id {
			if user == st.Storage[i].Creator || user == st.Storage[i].Assignee {
				task := st.Storage[i]
				copy(st.Storage[i:], st.Storage[i+1:])
				st.Storage[len(st.Storage)-1] = Task{}
				st.Storage = st.Storage[:len(st.Storage)-1]
				return task
			}
		}
	}
	return Task{}

}

func (st *TaskStorage) ModifyTask(Id int, newAssign User, oldAssign User) Task {
	for i := range st.Storage {
		if st.Storage[i].ID == Id {
			if oldAssign.Username == "" {
				st.Storage[i].Assignee = newAssign
				//fmt.Println(st.Storage[i].Assignee)
				return st.Storage[i]
			} else if oldAssign == st.Storage[i].Assignee {
				st.Storage[i].Assignee = User{}
				return st.Storage[i]
			}
		}
	}
	return Task{}
}

func (st *TaskStorage) Get(assignee string, creator string) []Task {
	res := make([]Task, 0)
	for _, task := range st.Storage {
		if (assignee == "" && creator == "") || (assignee == task.Assignee.Username && creator == "") || (creator == task.Creator.Username && assignee == "") {
			res = append(res, task)
		}
	}

	return res
}
