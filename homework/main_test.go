package main

import (
	"fmt"
	"os"
	"testing"
)

const testFilePath = "test_students.json"

func TestStudentStorage(t *testing.T) {
	// Удаляем тестовый файл, если он существует
	_ = os.Remove(testFilePath)

	storage := NewStudentStorage(testFilePath)

	t.Run("AddStudent", func(t *testing.T) {
		err := storage.AddStudent("Alice", 20, []int{90, 85, 95})
		if err != nil {
			t.Errorf("Failed to add student: %v", err)
		}

		// Попытка добавить того же студента снова
		err = storage.AddStudent("Alice", 20, []int{90, 85, 95})
		if err == nil {
			t.Error("Expected error when adding duplicate student")
		}
	})

	t.Run("GetStudent", func(t *testing.T) {
		student, err := storage.GetStudent("Alice")
		if err != nil {
			t.Errorf("Failed to get student: %v", err)
		}

		if student.Name != "Alice" || student.Age != 20 || len(student.Grades) != 3 {
			t.Error("Student data doesn't match")
		}

		// Попытка получить несуществующего студента
		_, err = storage.GetStudent("Bob")
		if err == nil {
			t.Error("Expected error when getting non-existent student")
		}
	})

	t.Run("UpdateStudent", func(t *testing.T) {
		err := storage.UpdateStudent("Alice", 21, []int{95, 90, 100})
		if err != nil {
			t.Errorf("Failed to update student: %v", err)
		}

		student, err := storage.GetStudent("Alice")
		if err != nil {
			t.Errorf("Failed to get student: %v", err)
		}

		if student.Age != 21 || student.Grades[0] != 95 {
			t.Error("Student data wasn't updated correctly")
		}

		// Попытка обновить несуществующего студента
		err = storage.UpdateStudent("Bob", 21, []int{95, 90, 100})
		if err == nil {
			t.Error("Expected error when updating non-existent student")
		}
	})

	t.Run("CalculateAverageGrade", func(t *testing.T) {
		avg, err := storage.CalculateAverageGrade("Alice")
		if err != nil {
			t.Errorf("Failed to calculate average grade: %v", err)
		}

		expected := (95.0 + 90.0 + 100.0) / 3.0
		if avg != expected {
			t.Errorf("Expected average %.2f, got %.2f", expected, avg)
		}

		// Проверка для студента без оценок
		err = storage.AddStudent("Bob", 22, []int{})
		if err != nil {
			t.Errorf("Failed to add student: %v", err)
		}

		avg, err = storage.CalculateAverageGrade("Bob")
		if err != nil {
			t.Errorf("Failed to calculate average grade: %v", err)
		}
		if avg != 0 {
			t.Errorf("Expected average 0 for student without grades, got %.2f", avg)
		}
	})

	t.Run("SaveAndLoad", func(t *testing.T) {
		err := storage.SaveToFile()
		if err != nil {
			t.Errorf("Failed to save to file: %v", err)
		}

		newStorage := NewStudentStorage(testFilePath)
		err = newStorage.LoadFromFile()
		if err != nil {
			t.Errorf("Failed to load from file: %v", err)
		}

		student, err := newStorage.GetStudent("Alice")
		if err != nil {
			t.Errorf("Failed to get student after load: %v", err)
		}

		if student.Name != "Alice" || student.Age != 21 || len(student.Grades) != 3 {
			t.Error("Loaded student data doesn't match")
		}
	})

	// Очистка после тестов
	_ = os.Remove(testFilePath)
}

func TestConcurrentAccess(t *testing.T) {
	storage := NewStudentStorage(testFilePath)

	// Добавляем начального студента
	err := storage.AddStudent("Initial", 25, []int{80, 85})
	if err != nil {
		t.Fatalf("Failed to add initial student: %v", err)
	}

	// Запускаем несколько горутин для конкурентного доступа
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			// Чтение
			_, err := storage.GetStudent("Initial")
			if err != nil {
				t.Errorf("Failed to get student in goroutine %d: %v", id, err)
			}

			// Запись
			newName := fmt.Sprintf("Student%d", id)
			err = storage.AddStudent(newName, 20+id, []int{70 + id, 75 + id})
			if err != nil {
				t.Errorf("Failed to add student in goroutine %d: %v", id, err)
			}

			done <- true
		}(i)
	}

	// Ожидаем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что все студенты добавлены
	students := storage.GetAllStudents()
	if len(students) != 11 { // Initial + 10 новых
		t.Errorf("Expected 11 students, got %d", len(students))
	}

	// Очистка после тестов
	_ = os.Remove(testFilePath)
}
