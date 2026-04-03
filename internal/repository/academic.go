package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
	"github.com/fun-dotto/app-bff-api/internal/external"
)

type AcademicRepository struct {
	client *academic_api.ClientWithResponses
}

func NewAcademicRepository(client *academic_api.ClientWithResponses) *AcademicRepository {
	return &AcademicRepository{client: client}
}

// GetFaculties は全教員一覧を取得する
func (r *AcademicRepository) GetFaculties() ([]domain.Faculty, error) {
	response, err := r.client.FacultiesV1ListWithResponse(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculties: status %d", response.StatusCode())
	}

	return external.ToDomainFaculties(response.JSON200.Faculties), nil
}

// GetFaculty は教員詳細を取得する
func (r *AcademicRepository) GetFaculty(id string) (*domain.Faculty, error) {
	response, err := r.client.FacultiesV1DetailWithResponse(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculty: status %d", response.StatusCode())
	}

	faculty := external.ToDomainFaculty(response.JSON200.Faculty)
	return &faculty, nil
}

// GetFacultiesByIDs は指定したIDの教員一覧を取得する
func (r *AcademicRepository) GetFacultiesByIDs(ids []string) (map[string]domain.Faculty, error) {
	if len(ids) == 0 {
		return make(map[string]domain.Faculty), nil
	}

	response, err := r.client.FacultiesV1ListWithResponse(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get faculties by IDs: status %d", response.StatusCode())
	}

	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	result := make(map[string]domain.Faculty, len(ids))
	for _, f := range response.JSON200.Faculties {
		if _, ok := idSet[f.Id]; ok {
			result[f.Id] = external.ToDomainFaculty(f)
		}
	}

	return result, nil
}

// GetSubjects は外部APIから科目一覧を取得する
func (r *AcademicRepository) GetSubjects(query domain.SubjectQuery) ([]domain.Subject, error) {
	params := external.ToExternalSubjectQuery(query)

	response, err := r.client.SubjectsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subjects: status %d", response.StatusCode())
	}

	subjects := response.JSON200.Subjects
	result := make([]domain.Subject, len(subjects))
	for i, s := range subjects {
		result[i] = external.ToDomainSubjectSummary(s)
	}

	return result, nil
}

// GetCourseRegistrations は外部APIから履修登録一覧を取得する
func (r *AcademicRepository) GetCourseRegistrations(userID string, semesters []domain.CourseSemester, year *int) ([]domain.CourseRegistration, error) {
	params := &academic_api.CourseRegistrationsV1ListParams{
		UserId:    userID,
		Semesters: externalToCourseSemesters(semesters),
		Year:      year,
	}

	response, err := r.client.CourseRegistrationsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get course registrations: status %d", response.StatusCode())
	}

	registrations := response.JSON200.CourseRegistrations
	result := make([]domain.CourseRegistration, len(registrations))
	for i, r := range registrations {
		result[i] = external.ToDomainCourseRegistration(r)
	}

	return result, nil
}

// CreateCourseRegistration は外部APIに履修登録を作成する
func (r *AcademicRepository) CreateCourseRegistration(userID string, subjectID string) (*domain.CourseRegistration, error) {
	body := academic_api.CourseRegistrationsV1CreateJSONRequestBody{
		SubjectId: subjectID,
		UserId:    userID,
	}

	response, err := r.client.CourseRegistrationsV1CreateWithResponse(context.Background(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create course registration: status %d", response.StatusCode())
	}

	cr := external.ToDomainCourseRegistration(response.JSON201.CourseRegistration)
	return &cr, nil
}

// DeleteCourseRegistration は外部APIから履修登録を削除する
func (r *AcademicRepository) DeleteCourseRegistration(id string) error {
	response, err := r.client.CourseRegistrationsV1DeleteWithResponse(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete course registration: status %d", response.StatusCode())
	}

	return nil
}

// GetTimetableItems は外部APIから時間割アイテム一覧を取得する
func (r *AcademicRepository) GetTimetableItems(query domain.TimetableItemQuery) ([]domain.TimetableItem, error) {
	params := external.ToExternalTimetableItemQuery(query)

	response, err := r.client.TimetableItemsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get timetable items: status %d", response.StatusCode())
	}

	items := response.JSON200.TimetableItems
	result := make([]domain.TimetableItem, len(items))
	for i, item := range items {
		result[i] = external.ToDomainTimetableItem(item)
	}

	return result, nil
}

// GetPersonalCalendarItems は外部APIから個人カレンダーアイテム一覧を取得する
func (r *AcademicRepository) GetPersonalCalendarItems(userID string, dates []time.Time) ([]domain.PersonalCalendarItem, error) {
	params := external.ToExternalPersonalCalendarItemParams(userID, dates)

	response, err := r.client.PersonalCalendarItemsV1ListWithResponse(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to call academic API: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get personal calendar items: status %d", response.StatusCode())
	}

	items := response.JSON200.PersonalCalendarItems
	result := make([]domain.PersonalCalendarItem, len(items))
	for i, item := range items {
		result[i] = external.ToDomainPersonalCalendarItem(item)
	}

	return result, nil
}

// GetSubject は外部APIから科目詳細を取得する
func (r *AcademicRepository) GetSubject(id string) (*domain.Subject, error) {
	ctx := context.Background()

	// 科目詳細を取得
	subjectResponse, err := r.client.SubjectsV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}
	if subjectResponse.JSON200 == nil {
		return nil, fmt.Errorf("failed to get subject: status %d", subjectResponse.StatusCode())
	}

	// シラバスを取得
	syllabusResponse, err := r.client.SyllabusV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get syllabus: %w", err)
	}
	if syllabusResponse.JSON200 == nil {
		return nil, fmt.Errorf("failed to get syllabus: status %d", syllabusResponse.StatusCode())
	}

	s := external.ToDomainSubject(subjectResponse.JSON200.Subject)

	syllabus := external.ToDomainSyllabus(syllabusResponse.JSON200.Syllabus)
	s.Syllabus = &syllabus

	return &s, nil
}

func externalToCourseSemesters(semesters []domain.CourseSemester) []academic_api.DottoFoundationV1CourseSemester {
	result := make([]academic_api.DottoFoundationV1CourseSemester, len(semesters))
	for i, semester := range semesters {
		result[i] = academic_api.DottoFoundationV1CourseSemester(semester)
	}
	return result
}
