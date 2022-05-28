package database

import (
	"sync"
)

type MemoryAdapter struct {
	sync.Mutex
}

const (
	Momory = "momory"
)

func init() {
	Register(&MemoryAdapter{}, Momory)
}

func (f *MemoryAdapter) Init() error {
	return nil
}

func (f *MemoryAdapter) UpsetFile(name string, byt []byte) error {

	f.Lock()
	defer f.Unlock()

	return nil
}

func (f *MemoryAdapter) DelFile(name string) error {

	f.Lock()
	defer f.Unlock()

	return nil
}

func (f *MemoryAdapter) FindFile(name string) (BehaviorInfo, error) {

	info := BehaviorInfo{}

	return info, nil
}

func (f *MemoryAdapter) GetAllFiles() ([]BehaviorInfo, error) {

	lst := []BehaviorInfo{}

	return lst, nil
}

func (f *MemoryAdapter) UpdateState(name string, status string) error {
	return nil
}

func (f *MemoryAdapter) UpdateTags(name string, tags []byte) error {
	return nil
}

func (f *MemoryAdapter) FindConfig(name string) (TemplateConfig, error) {
	info := TemplateConfig{}

	return info, nil
}

func (f *MemoryAdapter) UpsetConfig(byt []byte) error {

	f.Lock()
	defer f.Unlock()

	return nil
}

func (f *MemoryAdapter) RemoveReport(id string) error {

	return nil
}

func (f *MemoryAdapter) AppendReport(info ReportInfo) error {

	return nil
}

func (f *MemoryAdapter) GetReport() []ReportInfo {

	lst := []ReportInfo{}

	return lst
}
