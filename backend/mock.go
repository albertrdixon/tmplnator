package backend

type mockBackend struct {
  vals     map[string]string
  valslice map[string][]string
}

func (m mockBackend) Get(k string) (string, error) {
  if v, ok := m.vals[k]; ok {
    return v, nil
  }
  return "", KeyNotFound(k)
}

func (m mockBackend) GetAll(k string) ([]string, error) {
  if v, ok := m.valslice[k]; ok {
    return v, nil
  }
  return []string{}, KeyNotFound(k)
}
