package repository

import (
	"fmt"
	"time"

	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// MockAcademicRepository はテスト用のモック
type MockAcademicRepository struct {
	subjects                    []domain.Subject
	courseRegistrations         []domain.CourseRegistration
	personalCalendarItems       []domain.PersonalCalendarItem
	faculties                   []domain.Faculty
	createError                 error
	deleteError                 error
	getSubjectsError            error
	getSubjectError             error
	getRegistrationsErr         error
	getPersonalCalendarItemsErr error
}

func NewMockAcademicRepository() *MockAcademicRepository {
	faculty := domain.Faculty{
		ID:    "f1",
		Name:  "田中太郎",
		Email: "tanaka@example.com",
	}

	subject := domain.Subject{
		ID:       "s1",
		Name:     "プログラミング基礎",
		Year:     2026,
		Credit:   2,
		Semester: domain.CourseSemesterQ1,
		Faculties: []domain.SubjectFaculty{
			{Faculty: faculty, IsPrimary: true},
		},
		Requirements: []domain.SubjectRequirement{
			{Course: domain.CourseInformationSystem, RequirementType: domain.SubjectRequirementTypeRequired},
		},
		EligibleAttributes: []domain.SubjectTargetClass{
			{Grade: domain.GradeB1, Class: nil},
		},
	}

	slot := domain.TimetableSlot{
		DayOfWeek: domain.DayOfWeekMonday,
		Period:    domain.PeriodPeriod1,
	}

	room := domain.Room{
		ID:    "r1",
		Name:  "101講義室",
		Floor: domain.Floor1,
	}

	return &MockAcademicRepository{
		subjects: []domain.Subject{subject},
		courseRegistrations: []domain.CourseRegistration{
			{ID: "cr1", Subject: subject},
		},
		personalCalendarItems: []domain.PersonalCalendarItem{
			{
				Date:    time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),
				Period:  slot.Period,
				Rooms:   []domain.Room{room},
				Status:  domain.PersonalCalendarItemStatusNormal,
				Subject: subject,
			},
		},
		faculties: []domain.Faculty{faculty},
	}
}

func NewMockAcademicRepositoryWithError(field string, err error) *MockAcademicRepository {
	m := NewMockAcademicRepository()
	switch field {
	case "create":
		m.createError = err
	case "delete":
		m.deleteError = err
	case "getSubjects":
		m.getSubjectsError = err
	case "getSubject":
		m.getSubjectError = err
	case "getRegistrations":
		m.getRegistrationsErr = err
	case "getPersonalCalendarItems":
		m.getPersonalCalendarItemsErr = err
	}
	return m
}

func (m *MockAcademicRepository) GetFaculties() ([]domain.Faculty, error) {
	return m.faculties, nil
}

func (m *MockAcademicRepository) GetFaculty(id string) (*domain.Faculty, error) {
	for _, f := range m.faculties {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("faculty not found: %s", id)
}

func (m *MockAcademicRepository) GetSubjects(_ domain.SubjectQuery) ([]domain.Subject, error) {
	if m.getSubjectsError != nil {
		return nil, m.getSubjectsError
	}
	return m.subjects, nil
}

func (m *MockAcademicRepository) GetSubject(id string) (*domain.Subject, error) {
	if m.getSubjectError != nil {
		return nil, m.getSubjectError
	}
	for _, s := range m.subjects {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("subject not found: %s", id)
}

func (m *MockAcademicRepository) GetCourseRegistrations(_ string, _ []domain.CourseSemester, _ *int) ([]domain.CourseRegistration, error) {
	if m.getRegistrationsErr != nil {
		return nil, m.getRegistrationsErr
	}
	return m.courseRegistrations, nil
}

func (m *MockAcademicRepository) CreateCourseRegistration(_ string, _ string) (*domain.CourseRegistration, error) {
	if m.createError != nil {
		return nil, m.createError
	}
	reg := domain.CourseRegistration{
		ID:      "cr-new",
		Subject: m.subjects[0],
	}
	return &reg, nil
}

func (m *MockAcademicRepository) DeleteCourseRegistration(_ string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	return nil
}

func (m *MockAcademicRepository) GetTimetableItems(_ domain.TimetableItemQuery) ([]domain.TimetableItem, error) {
	return []domain.TimetableItem{
		{
			ID:      "t1",
			Slot:    &domain.TimetableSlot{DayOfWeek: domain.DayOfWeekMonday, Period: m.personalCalendarItems[0].Period},
			Rooms:   m.personalCalendarItems[0].Rooms,
			Subject: m.personalCalendarItems[0].Subject,
		},
	}, nil
}

func (m *MockAcademicRepository) GetPersonalCalendarItems(_ string, _ []time.Time) ([]domain.PersonalCalendarItem, error) {
	if m.getPersonalCalendarItemsErr != nil {
		return nil, m.getPersonalCalendarItemsErr
	}
	return m.personalCalendarItems, nil
}

func (m *MockAcademicRepository) GetCancelledClasses(_ domain.CancelledClassQuery) ([]domain.CancelledClass, error) {
	return []domain.CancelledClass{}, nil
}

func (m *MockAcademicRepository) GetMakeupClasses(_ domain.MakeupClassQuery) ([]domain.MakeupClass, error) {
	return []domain.MakeupClass{}, nil
}

func (m *MockAcademicRepository) GetRoomChanges(_ domain.RoomChangeQuery) ([]domain.RoomChange, error) {
	return []domain.RoomChange{}, nil
}
