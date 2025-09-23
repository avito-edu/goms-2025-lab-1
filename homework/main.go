package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Student struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Grades []int  `json:"grades"`
}

type StudentStorage struct {
	mu       sync.RWMutex
	students map[string]Student
	filePath string
}

func NewStudentStorage(filePath string) *StudentStorage {
	return &StudentStorage{
		students: make(map[string]Student),
		filePath: filePath,
	}
}

func (s *StudentStorage) AddStudent(name string, age int, grades []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if name == "" {
		return errors.New("empty name")
	}
	if _, exists := s.students[name]; exists {
		return errors.New("student already exists")
	}
	s.students[name] = Student{Name: name, Age: age, Grades: append([]int(nil), grades...)}
	return nil
}

func (s *StudentStorage) GetStudent(name string) (Student, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.students[name]
	if !ok {
		return Student{}, errors.New("student not found")
	}
	return st, nil
}

func (s *StudentStorage) UpdateStudent(name string, age int, grades []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.students[name]; !ok {
		return errors.New("student not found")
	}
	s.students[name] = Student{Name: name, Age: age, Grades: append([]int(nil), grades...)}
	return nil
}

func (s *StudentStorage) CalculateAverageGrade(name string) (float64, error) {
	s.mu.RLock()
	st, ok := s.students[name]
	s.mu.RUnlock()
	if !ok {
		return 0, errors.New("student not found")
	}
	if len(st.Grades) == 0 {
		return 0, nil
	}
	sum := 0
	for _, g := range st.Grades {
		sum += g
	}
	return float64(sum) / float64(len(st.Grades)), nil
}

func (s *StudentStorage) AddGrades(name string, grades []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.students[name]
	if !ok {
		return errors.New("student not found")
	}
	st.Grades = append(st.Grades, grades...)
	s.students[name] = st
	return nil
}

func (s *StudentStorage) GetAllStudents() []Student {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Student, 0, len(s.students))
	for _, st := range s.students {
		res = append(res, st)
	}
	return res
}

func (s *StudentStorage) SaveToFile() error {
	s.mu.RLock()
	data := make([]Student, 0, len(s.students))
	for _, st := range s.students {
		data = append(data, st)
	}
	s.mu.RUnlock()

	tmpPath := s.filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, s.filePath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func (s *StudentStorage) LoadFromFile() error {
	f, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	var data []Student
	dec := json.NewDecoder(f)
	if err := dec.Decode(&data); err != nil {
		return err
	}

	s.mu.Lock()
	s.students = make(map[string]Student, len(data))
	for _, st := range data {
		s.students[st.Name] = st
	}
	s.mu.Unlock()
	return nil
}
