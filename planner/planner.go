package planner

import (
	"errors"

	"github.com/ivahaev/go-logger"
)

type Date struct {
	ImgsURL        string
	ImgsOriginPath string
	ImgsIconPath   string
}

type Planner struct {
	TaskCount int
	Tasks     map[int]*Task
	DateOut   []Date
	ImgError  []string
}

func (p *Planner) Init() {
	p.Tasks = make(map[int]*Task)
}

func (p *Planner) GetLength() int {
	var i = 0
	for key, _ := range p.Tasks {
		i = i + 1
		//Заглушка
		key = key + 1
	}
	return i
}

func (p *Planner) AddTask(task *Task) error {
	if &task.id == nil {
		return errors.New("Planner -- Error heap, not value Task")
	}

	p.Tasks[p.GetLength()+1] = task
	return nil
}

/*func (p *Planner) runAllTask() error {
	for _, value := range p.tasks {
		fmt.Println("Planner -- run Task: " + value.Name)
		value.run()
	}
	return nil
}*/

func (p *Planner) RunAllTask(dateOut chan []Date) {
	var dateBuff []Date
	for _, value := range p.Tasks {
		dateTask := make(chan Date)
		logger.Notice("Planner -- run Task: " + value.Name)
		go value.run(dateTask)
		for _, i := range value.DateInput {
			logger.Info("Planner -- date exporting: " + value.Name + i.ImgsOriginPath)
			dateBuff = append(dateBuff, <-dateTask)
		}
	}
	dateOut <- dateBuff
	close(dateOut)
}

func (p *Planner) RunAllTaskType(dateOut chan []Date, taskType TypeTask) {
	var dateBuff []Date
	for _, value := range p.Tasks {
		if value.Done {
			continue
		}
		if value.TypeTask == taskType {
			dateTask := make(chan Date)
			logger.Notice("Planner -- run Task: " + value.Name)
			go value.run(dateTask)
			for _, i := range value.DateInput {
				logger.Info("Planner -- date exporting: " + value.Name + i.ImgsOriginPath)
				var item = <-dateTask
				dateBuff = append(dateBuff, item)

				//Отбраковка
				if value.TypeTask == Dowload {
					if item.ImgsOriginPath == "" {
						p.ImgError = append(p.ImgError, item.ImgsURL)
					}
				} else {
					if item.ImgsIconPath == "" {
						p.ImgError = append(p.ImgError, item.ImgsURL)
					}
				}
			}
		}
	}
	dateOut <- dateBuff
	close(dateOut)
}

func (p *Planner) RunTask(task *Task, dateOut chan []Date) {
	var dateBuff []Date
	if !p.IsTaskPlanner(task) {
		return
	}

	dateTask := make(chan Date)
	logger.Notice("Planner -- run Task: " + task.Name)
	go task.run(dateTask)

	for _, i := range task.DateInput {
		logger.Info("Planner -- date exporting: " + task.Name + i.ImgsOriginPath)
		dateBuff = append(dateBuff, <-dateTask)
	}

	dateOut <- dateBuff
	close(dateOut)
}

func (p *Planner) PoolTaskComplite() (int, *Task) {
	var t *Task
	var i int
	for key, value := range p.Tasks {
		if value.Done == true || value.start == false {
			t = value
			i = key
			break
		}
	}

	if t.id == "" {
		return 0, t
	}

	p.DeleteTask(t)
	return i, t
}

func (p *Planner) PoolsTaskComplite() []*Task {
	var t []*Task
	for _, value := range p.Tasks {
		if value.Done == true || value.start == false {
			t = append(t, value)
			p.DeleteTask(value)
		}
	}
	return t
}

func (p *Planner) GetKeyTask(task *Task) (int, error) {
	if task.id == "" {
		return 0, errors.New("Planner -- Error heap, not value Task")
	}

	for key, value := range p.Tasks {
		if task.id == value.id {
			return key, nil
		}
	}
	return 0, errors.New("Planner -- Error heap, not value Task")
}

func (p *Planner) IsTaskComplite(task *Task) (bool, error) {
	if task.id == "" {
		return false, errors.New("Planner -- Error, not task")
	}
	if task.Done || !task.start {
		return true, nil
	}
	return false, nil
}

func (p *Planner) IsTaskPlanner(task *Task) bool {
	for _, item := range p.Tasks {
		if task.id == item.id {
			return true
		}
	}
	return false
}

func (p *Planner) PoolErrorImages() []string {
	var imgs = p.ImgError
	p.ImgError = []string{}
	return imgs
}

func (p *Planner) DeleteTask(task *Task) error {
	if task.id == "" {
		return errors.New("Planner -- Error heap, not value Task")
	}

	key, err := p.GetKeyTask(task)
	if err != nil {
		return errors.New("Planner -- Error heap, not value Task")
	}
	//task.dateInput = []Date{}
	delete(p.Tasks, key)

	return nil
}
