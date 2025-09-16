package main

// TODO: implement me somehow 🍷🗿

// Student — структуру студента и набор его атрибутов
type Student struct{}

// StudentStorage — хранилище студентов
type StudentStorage struct{}

// NewStudentStorage — создает новое хранилище студентов
func NewStudentStorage(filePath string) *StudentStorage {
	return nil
}

// AddStudent — добавляет нового студента
func (s *StudentStorage) AddStudent(name string, age int, grades []int) error {
	return nil
}

// UpdateStudent — обновляет данные студента
func (s *StudentStorage) UpdateStudent(name string, age int, grades []int) error {
	return nil
}

// GetStudent — возвращает данные студента
func (s *StudentStorage) GetStudent(name string) (Student, error) {
	return Student{}, nil
}

// GetAllStudents — возвращает всех студентов
func (s *StudentStorage) GetAllStudents() []Student {
	return nil
}

// CalculateAverageGrade — вычисляет средний балл студента
func (s *StudentStorage) CalculateAverageGrade(name string) (float64, error) {
	return 0.0, nil
}

// SaveToFile — сохраняет данные в файл JSON
func (s *StudentStorage) SaveToFile() error {
	return nil
}

// LoadFromFile — загружает данные из файла JSON
func (s *StudentStorage) LoadFromFile() error {
	return nil
}
