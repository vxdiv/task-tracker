package model

import "time"

const (
	TaskTypeImprovement = "improvement"
	TaskTypeFuture      = "future"
	TaskTypeBug         = "bug"

	TaskStatusOpen       = "open"
	TaskStatusInProgress = "in_progress"
	TaskStatusResolve    = "resolver"
	TaskStatusClose      = "close"
	TaskStatusHold       = "hold"
	TaskStatusReopen     = "reopen"

	TaskResolutionDone            = "done"
	TaskResolutionFixed           = "fixed"
	TaskResolutionDuplicate       = "duplicate"
	TaskResolutionIncomplete      = "incomplete"
	TaskResolutionCannotReproduce = "cannot_reproduce"
	TaskResolutionDoNotNeedToDo   = "do_no_need_to_do"

	TaskPriorityTrivial  = "trivial"
	TaskPriorityMajor    = "major"
	TaskPriorityCritical = "critical"
	TaskPriorityASAP     = "asap"
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

type Task struct {
	ID          int64
	Name        string
	Description string
	Type        TaskType
	OwnerID     int64
	AssignedID  int64
	Status      string
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
