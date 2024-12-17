package error_solutions

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	models3 "github.com/kwaaka-team/orders-core/core/models"
	storecoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/error_solutions/models"
	"github.com/kwaaka-team/orders-core/service/error_solutions/repository"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetAllErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error)
	GetErrorSolutionByCode(ctx context.Context, store storecoreModels.Store, code string) (models.ErrorSolution, bool, error)
	GetTimeoutErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error)
	SetFailReason(ctx context.Context, store storecoreModels.Store, errorMessage, code, errorType string) (models3.FailReason, bool, error)
}

type ServiceImpl struct {
	repo repository.Repository
}

func NewErrorSolutionService(repo repository.Repository) (*ServiceImpl, error) {
	return &ServiceImpl{
		repo: repo,
	}, nil
}

func (s *ServiceImpl) GetAllErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error) {

	errorSolutions, err := s.repo.GetErrorSolutions(ctx)
	if err != nil {
		log.Err(err).Msgf("error GetAllErrorSolutions")
		return nil, err
	}

	return errorSolutions, nil
}

func (s *ServiceImpl) GetErrorSolutionByCode(ctx context.Context, store storecoreModels.Store, code string) (models.ErrorSolution, bool, error) {

	log.Info().Msgf("GetErrorSolutionByCode for code: %s", code)

	if len(code) == 0 {
		return models.ErrorSolution{}, false, nil
	}

	errSolutionByCode, err := s.repo.GetErrorSolutionByCode(ctx, code)
	if err != nil {
		log.Err(err).Msgf("ServiceImpl GetErrorSolutionByCode error for code: %s", code)
		return models.ErrorSolution{}, false, err
	}
	if len(store.ID) == 0 {
		return errSolutionByCode, false, nil
	}
	sendToStopListStatus := s.checkProductAddToStopListStatus(store, errSolutionByCode)

	return errSolutionByCode, sendToStopListStatus, nil
}

func (s *ServiceImpl) GetErrorSolutionByType(ctx context.Context, errorType string) (models.ErrorSolution, error) {

	if len(errorType) == 0 {
		return models.ErrorSolution{}, nil
	}

	errSolutionByCode, err := s.repo.GetErrorSolutionByType(ctx, errorType)
	if err != nil {
		return models.ErrorSolution{}, err
	}

	return errSolutionByCode, nil
}

func (s *ServiceImpl) checkProductAddToStopListStatus(store storecoreModels.Store, errorSolution models.ErrorSolution) bool {

	var storePosType string

	if store.PosType == models2.IIKO.String() && store.IikoCloud.IsExternalMenu {
		storePosType = models2.IIKOWEB.String()
	} else {
		storePosType = store.PosType
	}

	for _, posType := range errorSolution.StopListPosTypes {
		if posType == storePosType {
			return true
		}
	}

	return false
}

func (s *ServiceImpl) GetTimeoutErrorSolutions(ctx context.Context) ([]models.ErrorSolution, error) {
	return s.repo.GetTimeoutErrSolutions(ctx)
}

func (s *ServiceImpl) SetFailReason(ctx context.Context, store storecoreModels.Store, errorMessage, code, errorType string) (models3.FailReason, bool, error) {

	log.Info().Msgf("set fail reason: errorMessage:%s\n code:%s\n type:%s\n", errorMessage, code, errorType)

	var (
		errorSolution        models.ErrorSolution
		sendToStoplistStatus bool
		err                  error
	)
	if code != "" {
		errorSolution, sendToStoplistStatus, err = s.GetErrorSolutionByCode(ctx, store, code)
		if err != nil {
			return models3.FailReason{}, false, err
		}
	}

	if errorType != "" {
		errorSolution, err = s.GetErrorSolutionByType(ctx, errorType)
		if err != nil {
			return models3.FailReason{}, false, err
		}
	}

	return models3.FailReason{
		Code:         errorSolution.Code,
		Message:      errorMessage,
		BusinessName: errorSolution.BusinessName,
		Reason:       errorSolution.Reason,
		Solution:     errorSolution.Solution,
	}, sendToStoplistStatus, nil
}
