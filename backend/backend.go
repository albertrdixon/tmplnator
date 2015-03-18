package backend

// Backend types have methods for access key/value stores
type Backend interface {
  Get(key string) (string, error)
  GetAll(key string) ([]string, error)
}

// New returns a Backend. If no peers are passed then returns nil
func New(namespace string, peers []string) Backend {
  if len(peers) < 1 {
    return nil
  }
  return newEtcd(namespace, peers)
}

// NewMock returns a mock Database
func NewMock(vals map[string]string, valslice map[string][]string) Backend {
  return mockBackend{
    vals:     vals,
    valslice: valslice,
  }
}
