package mail

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"fmt"
	"strings"

	"github.com/oasdiff/oasdiff/checker"
	"github.com/oasdiff/oasdiff/formatters"
	"github.com/wneessen/go-mail"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/mail Service

type Service interface {
	SendMail(to string, subject string, body string) error
	SendOpenAPIDifferencesNotification(ctx context.Context, updatedApplication *model.Application, applicationDependencies []*model.AppEndpointDependencies, changes checker.Changes)
}

type mailService struct {
	client         *mail.Client
	senderMail     string
	storageService storage.Service
	logger         log.Logger
}

func NewMailService(mailConfig config.MailConfig, storageService storage.Service, logger log.Logger) (Service, error) {
	client, err := mail.NewClient(
		mailConfig.SMTPHost,
		mail.WithPort(mailConfig.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(mailConfig.SenderEmail),
		mail.WithPassword(mailConfig.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %v", err)
	}

	return &mailService{
		client:         client,
		senderMail:     mailConfig.SenderEmail,
		storageService: storageService,
		logger:         logger,
	}, nil
}

func (ms *mailService) SendMail(to string, subject string, body string) error {
	m := mail.NewMsg()
	err := m.From(ms.senderMail)
	if err != nil {
		return fmt.Errorf("failed to set sender email: %v", err)
	}
	err = m.To(to)
	if err != nil {
		return fmt.Errorf("failed to set recipient email: %v", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	err = ms.client.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func (ms *mailService) SendOpenAPIDifferencesNotification(ctx context.Context, updatedApplication *model.Application, applicationDependencies []*model.AppEndpointDependencies, changes checker.Changes) {
	teamMembersEmails := make(map[string][]string)

	for _, appDep := range applicationDependencies {
		if appDep.Application.Team == nil {
			continue
		}

		relevantChanges := ms.filterRelevantChanges(changes, appDep.Endpoints)
		if len(relevantChanges) == 0 {
			ms.logger.Infof("No relevant changes for application %s depending on %s, skipping email notification", appDep.Application.Name, updatedApplication.Name)
			continue
		}

		if _, exists := teamMembersEmails[appDep.Application.Team.Name]; !exists {
			members, err := ms.storageService.GetTeamMembers(ctx, appDep.Application.Team.Name)
			if err != nil {
				ms.logger.Errorf("Failed to retrieve team members for team %s: %v", appDep.Application.Team.Name, err)
				continue
			}

			emails := make([]string, 0)
			for _, member := range members {
				emails = append(emails, member.Email)
			}
			teamMembersEmails[appDep.Application.Team.Name] = emails
		}

		err := ms.sendEmailToTeamMembers(teamMembersEmails[appDep.Application.Team.Name], relevantChanges, updatedApplication, appDep.Application)
		if err != nil {
			ms.logger.Errorf("Failed to send email notification to team %s: %v", appDep.Application.Team.Name, err)
		}
	}
}

func (ms *mailService) filterRelevantChanges(changes checker.Changes, appEndpoints map[string]bool) checker.Changes {
	relevantChanges := make(checker.Changes, 0)

	for _, change := range changes {
		if ms.isChangeRelevant(change, appEndpoints) {
			relevantChanges = append(relevantChanges, change)
		}
	}

	return relevantChanges
}

func (ms *mailService) isChangeRelevant(change checker.Change, appEndpoints map[string]bool) bool {
	endpoint := change.GetPath()
	method := change.GetOperation()

	if endpoint == "" || method == "" {
		return false
	}

	key := strings.ToLower(method) + " " + strings.ToLower(endpoint)
	_, exists := appEndpoints[key]
	return exists
}

func (ms *mailService) sendEmailToTeamMembers(teamMembersEmails []string, changes checker.Changes, updatedApplication *model.Application, dependingApplication *model.Application) error {
	subject := fmt.Sprintf("[%s application dependency change] Detected changes in used %s endpoints", dependingApplication.Name, updatedApplication.Name)
	body, err := ms.formatChangesAsHTML(changes)
	if err != nil {
		return fmt.Errorf("failed to format changes as HTML: %s", err.Error())
	}

	for _, email := range teamMembersEmails {
		err := ms.SendMail(email, subject, body)
		if err != nil {
			ms.logger.Errorf("Failed to send email to %s: %v", email, err)
		}
	}

	return nil
}

func (ms *mailService) formatChangesAsHTML(changes checker.Changes) (string, error) {
	htmlFormatter, err := formatters.Lookup(string(formatters.FormatHTML), formatters.FormatterOpts{
		Language: "en",
	})

	if err != nil {
		return "", fmt.Errorf("failed to create HTML formatter: %s", err.Error())
	}

	htmlBytes, err := htmlFormatter.RenderChangelog(changes, formatters.NewRenderOpts(), "old", "new")
	if err != nil {
		return "", fmt.Errorf("failed to render HTML changelog: %s", err.Error())
	}

	return string(htmlBytes), nil
}
