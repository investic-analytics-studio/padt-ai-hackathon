package checker

import "sync"

var (
	checkerObj *checker
	mu         sync.Mutex
)

func InitChecker() {
	if checkerObj == nil {
		checkerObj = NewChecker()
	}
}

func Struct(val interface{}) error {
	mu.Lock()
	defer mu.Unlock()
	err := checkerObj.Struct(val)
	return err
}

func Var(val interface{}, cond string) error {
	mu.Lock()
	defer mu.Unlock()
	err := checkerObj.Var(val, cond)
	return err
}
