package service

import (
	"fmt"
	"log/slog"
	"time"

	"forge.capytal.company/capytalcode/project-comicverse/model"
	"forge.capytal.company/capytalcode/project-comicverse/repository"
	"forge.capytal.company/loreddev/x/tinyssert"
	"github.com/google/uuid"
)

type Project struct {
	projectRepo    *repository.Project
	permissionRepo *repository.Permissions

	log    *slog.Logger
	assert tinyssert.Assertions
}

func NewProject(
	project *repository.Project,
	permissions *repository.Permissions,
	logger *slog.Logger,
	assertions tinyssert.Assertions,
) *Project {
	return &Project{
		projectRepo:    project,
		permissionRepo: permissions,

		log:    logger,
		assert: assertions,
	}
}

func (svc Project) Create(title string, ownerUserID ...uuid.UUID) (model.Project, error) {
	log := svc.log.With(slog.String("title", title))
	log.Info("Creating project")
	defer log.Info("Finished creating project")

	id, err := uuid.NewV7()
	if err != nil {
		return model.Project{}, fmt.Errorf("service: failed to generate id: %w", err)
	}

	now := time.Now()

	p := model.Project{
		ID:          id,
		Title:       title,
		DateCreated: now,
		DateUpdated: now,
	}

	err = svc.projectRepo.Create(p)
	if err != nil {
		return model.Project{}, fmt.Errorf("service: failed to create project: %w", err)
	}

	if len(ownerUserID) > 0 {
		err := svc.SetAuthor(p.ID, ownerUserID[0])
		if err != nil {
			return model.Project{}, err
		}
	}

	return p, nil
}

func (svc Project) SetAuthor(projectID uuid.UUID, userID uuid.UUID) error {
	log := svc.log.With(slog.String("project", projectID.String()), slog.String("user", userID.String()))
	log.Info("Setting project owner")
	defer log.Info("Finished setting project owner")

	if _, err := svc.permissionRepo.GetByID(projectID, userID); err == nil {
		err := svc.permissionRepo.Update(projectID, userID, model.PermissionAuthor)
		if err != nil {
			return fmt.Errorf("service: failed to update project author: %w", err)
		}
	}

	p := model.PermissionAuthor

	err := svc.permissionRepo.Create(projectID, userID, p)
	if err != nil {
		return fmt.Errorf("service: failed to set project owner: %w", err)
	}

	return nil
}

func (svc Project) GetUserProjects(userID uuid.UUID) ([]model.Project, error) {
	perms, err := svc.permissionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user permissions: %w", err)
	}

	ids := []uuid.UUID{}
	for project, permissions := range perms {
		if permissions.Has(model.PermissionRead) {
			ids = append(ids, project)
		}
	}

	if len(ids) == 0 {
		return []model.Project{}, nil
	}

	projects, err := svc.projectRepo.GetByIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get projects: %w", err)
	}

	return projects, nil
}

func (svc Project) GetProject(projectID uuid.UUID) (model.Project, error) {
	p, err := svc.projectRepo.GetByID(projectID)
	if err != nil {
		return model.Project{}, fmt.Errorf("service: failed to get project: %w", err)
	}
	return p, nil
}
