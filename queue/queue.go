package queue

import "mime/multipart"

type File struct {
	File   multipart.File
	Header *multipart.FileHeader
}

type Queue struct {
	taskCh chan File
}

func NewQueue() *Queue {
	return &Queue{
		taskCh: make(chan File, 1024),
	}
}

func (q *Queue) AddTask(file multipart.File, header *multipart.FileHeader) {
	q.taskCh <- File{File: file, Header: header}
}

func (q *Queue) ConsumeTasks() <-chan File {
	return q.taskCh
}
