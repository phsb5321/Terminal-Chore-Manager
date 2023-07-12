package task

type Status int

const (
	ToDo Status = iota

	InProgress

	Done
)

type Task struct {
	Status      Status
	Title       string
	Description string
}

func NewTask(status Status, title, description string) Task {
	return Task{Status: status, Title: title, Description: description}
}

func (s Status) GetNext() Status {
	if s == Done {
		return ToDo
	}
	return s + 1
}

func (s Status) GetPrev() Status {
	if s == ToDo {
		return Done
	}
	return s - 1
}

func (t *Task) Next() {
	t.Status = t.Status.GetNext()
}

func (t Task) FilterValue() string {
	return t.Title
}
