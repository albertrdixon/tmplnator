package backend

type Backend interface {
  Get(key string) (string, error)
  GetAll(key string) ([]string, error)
}

func New(namespace string, peers []string) Backend {
  if len(peers) < 1 {
    return nil
  }
  return newEtcd(namespace, peers)
}

func NewMock(vals map[string]string, valslice map[string][]string) Backend {
  return mockBackend{
    vals:     vals,
    valslice: valslice,
  }
}
