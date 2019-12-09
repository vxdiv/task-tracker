package app

import "time"

const (
	TaskTypeImprovement TaskType = "improvement"
	TaskTypeFuture      TaskType = "feature"
	TaskTypeBug         TaskType = "bug"

	TaskStatusOpen       TaskStatus = "open"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusResolve    TaskStatus = "resolver"
	TaskStatusClose      TaskStatus = "close"
	TaskStatusHold       TaskStatus = "hold"
	TaskStatusReopen     TaskStatus = "reopen"

	TaskResolutionDone            TaskResolution = "done"
	TaskResolutionFixed           TaskResolution = "fixed"
	TaskResolutionDuplicate       TaskResolution = "duplicate"
	TaskResolutionIncomplete      TaskResolution = "incomplete"
	TaskResolutionCannotReproduce TaskResolution = "cannot_reproduce"
	TaskResolutionDoNotNeedToDo   TaskResolution = "do_not_need_to_do"

	TaskPriorityTrivial  TaskPriority = "trivial"
	TaskPriorityMajor    TaskPriority = "major"
	TaskPriorityCritical TaskPriority = "critical"
	TaskPriorityASAP     TaskPriority = "asap"
)

type TaskType string

func (t TaskType) IsImprovement() bool {
	return t == TaskTypeImprovement
}

func (t TaskType) IsFuture() bool {
	return t == TaskTypeFuture
}

func (t TaskType) IsBug() bool {
	return t == TaskTypeBug
}

type TaskResolution string

func (t TaskResolution) IsDone() bool {
	return t == TaskResolutionDone
}

func (t TaskResolution) IsFixed() bool {
	return t == TaskResolutionFixed
}

func (t TaskResolution) IsDuplicate() bool {
	return t == TaskResolutionDuplicate
}

func (t TaskResolution) IsIncomplete() bool {
	return t == TaskResolutionIncomplete
}

func (t TaskResolution) IsCannotReproduce() bool {
	return t == TaskResolutionCannotReproduce
}

func (t TaskResolution) IsDoNotNeedToDo() bool {
	return t == TaskResolutionDoNotNeedToDo
}

type TaskPriority string

func (t TaskPriority) IsTrivial() bool {
	return t == TaskPriorityTrivial
}

func (t TaskPriority) IsMajor() bool {
	return t == TaskPriorityMajor
}

func (t TaskPriority) IsCritical() bool {
	return t == TaskPriorityCritical
}

func (t TaskPriority) IsASAP() bool {
	return t == TaskPriorityASAP
}

type TaskStatus string

type Task struct {
	ID          int64
	Name        string
	Description string
	Type        TaskType
	OwnerID     int64
	AssignedID  int64
	Status      TaskStatus
	DueDate     time.Time
	Resolution  TaskResolution
	Priority    TaskPriority
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) IsOpen() bool {
	return t.Status == TaskStatusOpen
}

func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusInProgress
}

func (t *Task) IsResolve() bool {
	return t.Status == TaskStatusResolve
}

func (t *Task) IsClose() bool {
	return t.Status == TaskStatusClose
}

func (t *Task) IsHold() bool {
	return t.Status == TaskStatusHold
}

func (t *Task) IsReopen() bool {
	return t.Status == TaskStatusReopen
}
