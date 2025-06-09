package services

import (
	"fmt"
)

type RegisterTemplateData struct {
	CsrfToken string
	Errors    []string
	Messages  []string
	Form      *RegisterForm
}

type RegisterForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c RegisterTemplateData) String() string {
	return fmt.Sprintf("CsrfToken: %v, Errors [%v], Messages [%v]",
		c.CsrfToken,
		len(c.Errors),
		len(c.Messages),
	)
}

type RegisterTemplateDataService struct {
	Form *RegisterForm
}

type GetRegisterTemplateContentOptionsFunc func(*RegisterTemplateData)

func NewRegisterTemplateDataService() *RegisterTemplateDataService {
	return &RegisterTemplateDataService{}
}

func (s *RegisterTemplateDataService) WithCsrfToken(csrfToken string) GetRegisterTemplateContentOptionsFunc {
	return func(opts *RegisterTemplateData) {
		opts.CsrfToken = csrfToken
	}
}

func (s *RegisterTemplateDataService) WithErrors(errs []string) GetRegisterTemplateContentOptionsFunc {
	return func(opts *RegisterTemplateData) {
		opts.Errors = append(opts.Errors, errs...)
	}
}

func (s *RegisterTemplateDataService) GetRegisterTemplateContent(opts ...GetRegisterTemplateContentOptionsFunc) (*RegisterTemplateData, error) {

	var templateContent *RegisterTemplateData = &RegisterTemplateData{}

	for _, fn := range opts {
		fn(templateContent)
	}

	var errors []string

	templateContent.Form = &RegisterForm{Username: "", Email: "", Password: ""}
	templateContent.Errors = errors
	templateContent.Messages = []string{}

	return templateContent, nil
}
