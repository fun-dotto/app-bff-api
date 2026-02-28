package handler

import (
	api "github.com/fun-dotto/app-bff-api/generated"
	"github.com/fun-dotto/app-bff-api/internal/domain"
)

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:    announcement.ID,
		Title: announcement.Title,
		Date:  announcement.AvailableFrom,
		Url:   announcement.URL,
	}
}

func toApiSubjectDetail(subject domain.Subject) api.SubjectDetail {
	var syllabus api.SubjectServiceSyllabus
	if subject.Syllabus != nil {
		syllabus = toApiSyllabus(*subject.Syllabus)
	}

	return api.SubjectDetail{
		Id:                 subject.ID,
		Name:               subject.Name,
		Credit:             subject.Credit,
		Semester:           api.DottoFoundationV1CourseSemester(subject.Semester),
		Faculties:          toApiFaculties(subject.Faculties),
		Requirements:       toApiRequirements(subject.Requirements),
		EligibleAttributes: toApiTargetClasses(subject.EligibleAttributes),
		Syllabus:           syllabus,
	}
}

func toApiSyllabus(syllabus domain.Syllabus) api.SubjectServiceSyllabus {
	return api.SubjectServiceSyllabus{
		Id:                         syllabus.ID,
		Name:                       syllabus.Name,
		EnName:                     syllabus.EnName,
		Summary:                    syllabus.Summary,
		LearningOutcomes:           syllabus.LearningOutcomes,
		ContentsAndSchedule:        syllabus.ContentsAndSchedule,
		PreLearning:                syllabus.PreLearning,
		PostLearning:               syllabus.PostLearning,
		Assignments:                syllabus.Assignments,
		EvaluationMethod:           syllabus.EvaluationMethod,
		Textbooks:                  syllabus.Textbooks,
		ReferenceBooks:             syllabus.ReferenceBooks,
		Prerequisites:              syllabus.Prerequisites,
		Notes:                      syllabus.Notes,
		Keywords:                   syllabus.Keywords,
		Classifications:            syllabus.Classifications,
		Grades:                     syllabus.Grades,
		Credit:                     syllabus.Credit,
		FacultyNames:               syllabus.FacultyNames,
		TeachingForm:               syllabus.TeachingForm,
		TeachingAndExamForm:        syllabus.TeachingAndExamForm,
		TeachingLanguage:           syllabus.TeachingLanguage,
		MultiplePersonTeachingForm: syllabus.MultiplePersonTeachingForm,
		PracticalHomeFacultyCategory: syllabus.PracticalHomeFacultyCategory,
		DspoSubject:                syllabus.DspoSubject,
		TargetAreas:                syllabus.TargetAreas,
		TargetCourses:              syllabus.TargetCourses,
	}
}

func toApiFaculties(faculties []domain.SubjectFaculty) []api.SubjectServiceSubjectFaculty {
	result := make([]api.SubjectServiceSubjectFaculty, len(faculties))
	for i, f := range faculties {
		result[i] = api.SubjectServiceSubjectFaculty{
			Faculty: api.DottoFoundationV1Faculty{
				Id:    f.Faculty.ID,
				Name:  f.Faculty.Name,
				Email: f.Faculty.Email,
			},
			IsPrimary: f.IsPrimary,
		}
	}
	return result
}

func toApiRequirements(requirements []domain.SubjectRequirement) []api.SubjectServiceSubjectRequirement {
	result := make([]api.SubjectServiceSubjectRequirement, len(requirements))
	for i, r := range requirements {
		result[i] = api.SubjectServiceSubjectRequirement{
			Course:          api.DottoFoundationV1Course(r.Course),
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}
	return result
}

func toApiTargetClasses(targetClasses []domain.SubjectTargetClass) []api.SubjectServiceSubjectTargetClass {
	result := make([]api.SubjectServiceSubjectTargetClass, len(targetClasses))
	for i, tc := range targetClasses {
		var class *api.DottoFoundationV1Class
		if tc.Class != nil {
			c := api.DottoFoundationV1Class(*tc.Class)
			class = &c
		}
		result[i] = api.SubjectServiceSubjectTargetClass{
			Grade: api.DottoFoundationV1Grade(tc.Grade),
			Class: class,
		}
	}
	return result
}
