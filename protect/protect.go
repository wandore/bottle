package protect

import "sync"

type request struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type Handler struct {
	mu         sync.Mutex
	requestMap map[string]*request
}

func (h *Handler) Query(key string, f func() (interface{}, error)) (interface{}, error) {
	h.mu.Lock()
	if h.requestMap == nil {
		h.requestMap = make(map[string]*request, 0)
	}

	if r, ok := h.requestMap[key]; ok {
		h.mu.Unlock()
		r.wg.Wait()
		return r.value, r.err
	}

	r := &request{}
	r.wg.Add(1)
	h.requestMap[key] = r
	h.mu.Unlock()

	r.value, r.err = f()
	r.wg.Done()

	h.mu.Lock()
	delete(h.requestMap, key)
	h.mu.Unlock()

	return r.value, r.err
}
