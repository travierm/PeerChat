package lib

type DataPayload struct {
	Hash string
	Data interface{}
}

type SignalServer struct {
	Cache map[string]interface{}
}

type AnswerServer struct {
	Cache map[string][]interface{}
}

func (s SignalServer) Store(hash string, signal interface{}) {
	s.Cache[hash] = signal
}

func (s SignalServer) GetByHash(hash string) interface{} {
	return s.Cache[hash]
}

func (s AnswerServer) Push(hash string, answer interface{}) {
	if s.Cache[hash] == nil {
		s.Cache[hash] = make([]interface{}, 0)
	}

	s.Cache[hash] = append(s.Cache[hash], answer)
}

func (s AnswerServer) GetByHash(hash string) []interface{} {
	return s.Cache[hash]
}
