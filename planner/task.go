package planner

import (
	"github.com/ivahaev/go-logger"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	id        string
	Name      string
	TypeTask  TypeTask
	DateInput []Date
	DateOut   []Date
	Done      bool
	start     bool
	err       error
}

type TypeTask string

const (
	Dowload = TypeTask("Dowload")
	Resize  = TypeTask("Resize")
)

func (t *Task) Init() {
	t.start = false
	t.Done = false
	t.DateOut = []Date{}
	t.id = uuID()
}

func (t *Task) run(date chan Date) {
	t.start = true

	switch t.TypeTask {
	case Dowload:
		var master DownloadMaster
		logger.Info("Planner -- Task name -- <" + t.Name + "> START")
		imgD := make(chan Date)

		master.init()
		master.date = t.DateInput
		go master.run(imgD)
		for item := range imgD {
			t.DateOut = append(t.DateOut, item)
			date <- item
			logger.Info("Planner --  Task name -- <" + t.Name + "> FINISH  -- " + item.ImgsOriginPath)
		}
		close(date)
		t.complite()

	case Resize:
		var resiz ResizeMaster
		logger.Info("Task name -- <" + t.Name + "> START")
		imgR := make(chan Date)

		resiz.init()
		resiz.date = t.DateInput
		go resiz.run(imgR)
		for item := range imgR {
			t.DateOut = append(t.DateOut, item)
			date <- item
			logger.Info("Task name -- <" + t.Name + "> FINISH  -- " + item.ImgsIconPath)
		}
		close(date)
		t.complite()

	default:
		logger.Error("Planner -- error Task, not found type")
	}
}

func (t *Task) complite() {
	t.start = false
	t.Done = true
}

func uuID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}
