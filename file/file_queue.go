package file

import l "github.com/Sirupsen/logrus"

type FileQueue struct {
  files []File
  queue chan File
}

func newFileQueue() *FileQueue {
  return &FileQueue{files: []File{}}
}

func (fq *FileQueue) add(f File) {
  l.WithField("file", f).Debug("Adding file to queue")
  fq.files = append(fq.files, f)
  l.WithField("file", f).Debug("File added")
}

func (fq *FileQueue) populateQueue() {
  fq.queue = make(chan File, len(fq.files))
  for _, f := range fq.files {
    fq.queue <- f
  }
  close(fq.queue)
}

func (fq *FileQueue) Queue() chan File {
  return fq.queue
}

func (fq *FileQueue) Len() int {
  return len(fq.queue)
}
