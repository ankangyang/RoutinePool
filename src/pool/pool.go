package pool

import (
	"fmt"
	)


//任务
type Task struct {
	f func() error
}

func (t Task) Execute() error{
	return t.f()
}

func NewTask(newFun func() error) *Task{
	t := Task{f:newFun}
	return &t
}


//池
type WorkPool struct {
	workNum int //池大小
	taskChan chan *Task
	isTaskChanClose bool //是否关闭
	//taskChanMutex sync.Mutex

	EntryChan chan *Task
	isEntryChanClose bool //是否关闭
	//entryChanMutex sync.Mutex

}

func(pool *WorkPool) Init(taskSize , chanSize int){
	pool.workNum = taskSize

	pool.taskChan = make(chan *Task, chanSize)
	pool.isTaskChanClose = false
	pool.EntryChan = make(chan *Task, chanSize)
	pool.isEntryChanClose =  false
}

func (pool *WorkPool) work(workId int ){
	for task := range pool.taskChan {
		task.Execute()
		fmt.Printf("workID:%v\n", workId)
	}
}

func(pool *WorkPool) Fini(){
	if(pool.isTaskChanClose){
		//pool.taskChanMutex.Lock()
		close(pool.taskChan)
		pool.isTaskChanClose = true
		//pool.taskChanMutex.Unlock()
	}
	if(pool.isEntryChanClose ){
		//pool.entryChanMutex.Lock()
		close(pool.EntryChan)
		pool.isTaskChanClose = false
		//pool.entryChanMutex.Unlock()
	}

}
func (pool *WorkPool) AddTask(task *Task) bool{
	//pool.entryChanMutex.Lock()
	//defer pool.entryChanMutex.Unlock()
	if(false == pool.isEntryChanClose ){
		pool.EntryChan <- task
		return true
	}
	return false
}
func(pool *WorkPool) Run(){
	//任务转移到内部chan
	go func() {
		for{
			task := <- pool.EntryChan
			//pool.taskChanMutex.Lock()
			if false == pool.isTaskChanClose{
				pool.taskChan <- task
			}
			//pool.taskChanMutex.Unlock()
		}
	}()

	//初始化工作协程
	for i := 0; i < pool.workNum; i++ {
		go pool.work(i+1)
	}
}



