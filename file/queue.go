package file

import l "github.com/Sirupsen/logrus"

// Queue describes a queue of files ofr the generator workers
type Queue struct {
  files []File
  queue chan File
}

// NewFileQueue returns an initialized file.Queue
func NewFileQueue() *Queue {
  return &Queue{files: []File{}}
}

func (fq *Queue) add(f File) {
  l.WithField("file", f).Debug("Adding file to queue")
  fq.files = append(fq.files, f)
  l.WithField("file", f).Debug("File added")
}

// PopulateQueue feeds parsed files into the underlying channel
func (fq *Queue) PopulateQueue() {
  fq.queue = make(chan File, len(fq.files))
  for _, f := range fq.files {
    fq.queue <- f
  }
  close(fq.queue)
}

// Queue returns the File channel
func (fq *Queue) Queue() chan File {
  return fq.queue
}

// Len returns the length of the queue
func (fq *Queue) Len() int {
  return len(fq.queue)
}

// Files returns the file slice
func (f *Queue) Files() []File {
  return f.files
}
