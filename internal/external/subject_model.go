package external

import (
	"github.com/fun-dotto/app-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

// ToDomainSubject は外部APIのSubjectをDomainのSubjectに変換する
func ToDomainSubject(m academic_api.Subject) domain.Subject {
	return domain.Subject{
		ID:                 m.Id,
		Name:               m.Name,
		Year:               m.Year,
		Credit:             m.Credit,
		Semester:           domain.CourseSemester(m.Semester),
		Faculties:          toDomainSubjectFaculties(m.Faculties),
		Requirements:       toDomainSubjectRequirements(m.Requirements),
		EligibleAttributes: toDomainSubjectTargetClasses(m.EligibleAttributes),
		Syllabus:           nil,
	}
}

func ToDomainSubjectSummary(m academic_api.SubjectSummary) domain.Subject {
	return domain.Subject{
		ID:        m.Id,
		Name:      m.Name,
		Faculties: toDomainSubjectFaculties(m.Faculties),
	}
}

func ToDomainSyllabus(m academic_api.Syllabus) domain.Syllabus {
	return domain.Syllabus{
		ID:                           m.Id,
		Name:                         m.Name,
		EnName:                       m.EnName,
		Summary:                      m.Summary,
		LearningOutcomes:             m.LearningOutcomes,
		ContentsAndSchedule:          m.ContentsAndSchedule,
		PreLearning:                  m.PreLearning,
		PostLearning:                 m.PostLearning,
		Assignments:                  m.Assignments,
		EvaluationMethod:             m.EvaluationMethod,
		Textbooks:                    m.Textbooks,
		ReferenceBooks:               m.ReferenceBooks,
		Prerequisites:                m.Prerequisites,
		Notes:                        m.Notes,
		Keywords:                     m.Keywords,
		Classifications:              m.Classifications,
		Grades:                       m.Grades,
		Credit:                       m.Credit,
		FacultyNames:                 m.FacultyNames,
		TeachingForm:                 m.TeachingForm,
		TeachingAndExamForm:          m.TeachingAndExamForm,
		TeachingLanguage:             m.TeachingLanguage,
		MultiplePersonTeachingForm:   m.MultiplePersonTeachingForm,
		PracticalHomeFacultyCategory: m.PracticalHomeFacultyCategory,
		DsopSubject:                  m.DsopSubject,
		TargetAreas:                  m.TargetAreas,
		TargetCourses:                m.TargetCourses,
	}
}

func toDomainSubjectFaculties(faculties []academic_api.SubjectFaculty) []domain.SubjectFaculty {
	result := make([]domain.SubjectFaculty, len(faculties))
	for i, f := range faculties {
		result[i] = domain.SubjectFaculty{
			Faculty: domain.Faculty{
				ID: f.FacultyId,
			},
			IsPrimary: f.IsPrimary,
		}
	}
	return result
}

func toDomainSubjectRequirements(requirements []academic_api.SubjectRequirement) []domain.SubjectRequirement {
	result := make([]domain.SubjectRequirement, len(requirements))
	for i, r := range requirements {
		result[i] = domain.SubjectRequirement{
			Course:          domain.Course(r.Course),
			RequirementType: domain.SubjectRequirementType(r.RequirementType),
		}
	}
	return result
}

func toDomainSubjectTargetClasses(targetClasses []academic_api.SubjectTargetClass) []domain.SubjectTargetClass {
	result := make([]domain.SubjectTargetClass, len(targetClasses))
	for i, tc := range targetClasses {
		var class *domain.Class
		if tc.Class != nil {
			c := domain.Class(*tc.Class)
			class = &c
		}
		result[i] = domain.SubjectTargetClass{
			Grade: domain.Grade(tc.Grade),
			Class: class,
		}
	}
	return result
}

// ToExternalSubjectQuery はDomainのSubjectQueryを外部APIのSubjectsV1ListParamsに変換する
func ToExternalSubjectQuery(q domain.SubjectQuery) *academic_api.SubjectsV1ListParams {
	params := &academic_api.SubjectsV1ListParams{}

	if q.Q != "" {
		params.Q = &q.Q
	}
	if len(q.Grade) > 0 {
		grades := toExternalGrades(q.Grade)
		params.Grade = &grades
	}
	if len(q.Courses) > 0 {
		courses := toExternalCourses(q.Courses)
		params.Courses = &courses
	}
	if len(q.Class) > 0 {
		classes := toExternalClasses(q.Class)
		params.Class = &classes
	}
	if len(q.Classification) > 0 {
		classifications := toExternalClassifications(q.Classification)
		params.Classification = &classifications
	}
	if len(q.Semester) > 0 {
		semesters := toExternalSemesters(q.Semester)
		params.Semester = &semesters
	}
	if len(q.RequirementType) > 0 {
		requirementTypes := toExternalRequirementTypes(q.RequirementType)
		params.RequirementType = &requirementTypes
	}
	if len(q.CulturalSubjectCategory) > 0 {
		culturalSubjectCategories := toExternalCulturalSubjectCategories(q.CulturalSubjectCategory)
		params.CulturalSubjectCategory = &culturalSubjectCategories
	}

	return params
}

func toExternalGrades(grades []domain.Grade) []academic_api.DottoFoundationV1Grade {
	result := make([]academic_api.DottoFoundationV1Grade, len(grades))
	for i, g := range grades {
		result[i] = academic_api.DottoFoundationV1Grade(g)
	}
	return result
}

func toExternalCourses(courses []domain.Course) []academic_api.DottoFoundationV1Course {
	result := make([]academic_api.DottoFoundationV1Course, len(courses))
	for i, c := range courses {
		result[i] = academic_api.DottoFoundationV1Course(c)
	}
	return result
}

func toExternalClasses(classes []domain.Class) []academic_api.DottoFoundationV1Class {
	result := make([]academic_api.DottoFoundationV1Class, len(classes))
	for i, c := range classes {
		result[i] = academic_api.DottoFoundationV1Class(c)
	}
	return result
}

func toExternalClassifications(classifications []domain.SubjectClassification) []academic_api.DottoFoundationV1SubjectClassification {
	result := make([]academic_api.DottoFoundationV1SubjectClassification, len(classifications))
	for i, c := range classifications {
		result[i] = academic_api.DottoFoundationV1SubjectClassification(c)
	}
	return result
}

func toExternalSemesters(semesters []domain.CourseSemester) []academic_api.DottoFoundationV1CourseSemester {
	result := make([]academic_api.DottoFoundationV1CourseSemester, len(semesters))
	for i, s := range semesters {
		result[i] = academic_api.DottoFoundationV1CourseSemester(s)
	}
	return result
}

func toExternalRequirementTypes(types []domain.SubjectRequirementType) []academic_api.DottoFoundationV1SubjectRequirementType {
	result := make([]academic_api.DottoFoundationV1SubjectRequirementType, len(types))
	for i, t := range types {
		result[i] = academic_api.DottoFoundationV1SubjectRequirementType(t)
	}
	return result
}

func toExternalCulturalSubjectCategories(categories []domain.CulturalSubjectCategory) []academic_api.DottoFoundationV1CulturalSubjectCategory {
	result := make([]academic_api.DottoFoundationV1CulturalSubjectCategory, len(categories))
	for i, c := range categories {
		result[i] = academic_api.DottoFoundationV1CulturalSubjectCategory(c)
	}
	return result
}
