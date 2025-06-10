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
	Username string `json:"username" form:"username" validate:"required,min=4,max=255"`
	Email    string `json:"email" form:"email" validate:"required,email,min=6,max=255"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=255"`
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
