package services

import (
	"fmt"
)

type LoginTemplateData struct {
	CsrfToken string
	Errors    []string
	Messages  []string
	Form      *LoginForm
}

type LoginForm struct {
	Contact  string `json:"contact" form:"contact" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (c LoginTemplateData) String() string {
	return fmt.Sprintf("CsrfToken: %v, Errors [%v], Messages [%v]",
		c.CsrfToken,
		len(c.Errors),
		len(c.Messages),
	)
}

type LoginTemplateDataService struct {
	Form *LoginForm
}

type GetLoginTemplateContentOptionsFunc func(*LoginTemplateData)

func NewLoginTemplateDataService() *LoginTemplateDataService {
	return &LoginTemplateDataService{}
}

func (s *LoginTemplateDataService) WithCsrfToken(csrfToken string) GetLoginTemplateContentOptionsFunc {
	return func(opts *LoginTemplateData) {
		opts.CsrfToken = csrfToken
	}
}

func (s *LoginTemplateDataService) WithErrors(errs []string) GetLoginTemplateContentOptionsFunc {
	return func(opts *LoginTemplateData) {
		opts.Errors = append(opts.Errors, errs...)
	}
}

func (s *LoginTemplateDataService) GetLoginTemplateContent(opts ...GetLoginTemplateContentOptionsFunc) (*LoginTemplateData, error) {

	var templateContent *LoginTemplateData = &LoginTemplateData{}

	for _, fn := range opts {
		fn(templateContent)
	}

	var errors []string

	templateContent.Form = &LoginForm{Contact: "", Password: ""}
	templateContent.Errors = errors
	templateContent.Messages = []string{}

	return templateContent, nil
}
