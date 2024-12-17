package jowi

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	jowiDto "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/rs/zerolog/log"
	"strconv"
)

func menuFromClient(req jowiDto.ResponseCourse, store storeModels.Settings, sections []models.Section, collections []models.MenuCollection, stopList models.StopListItems) models.Menu {

	menu := models.Menu{
		Products: getProducts(req, store),
	}

	menu.Sections = sections
	menu.Collections = collections
	menu.StopLists = stopList.Products()

	return menu
}

func categoryFromClient(req jowiDto.ResponseCourseCategory) ([]models.Section, []models.MenuCollection, error) {
	var (
		sections    []models.Section
		collections []models.MenuCollection
	)

	var hasCollection bool
	for _, courseCategory := range req.CourseCategories {
		if courseCategory.ParentId != "" {
			hasCollection = true
		}
	}

	for _, courseCategory := range req.CourseCategories {
		if (hasCollection && courseCategory.ParentId != "") || !hasCollection {
			sections = append(sections, models.Section{
				ExtID:      courseCategory.Id,
				Name:       courseCategory.Title,
				Collection: courseCategory.ParentId,
			})
			continue
		}

		collections = append(collections, models.MenuCollection{
			ExtID: courseCategory.Id,
			Name:  courseCategory.Title,
		})
	}

	return sections, collections, nil
}

func stopListFromClient(req jowiDto.ResponseStopList) (models.StopListItems, error) {
	var stopList models.StopListItems

	for _, course := range req.CourseCounts {
		count, err := strconv.ParseFloat(course.Count, 64)
		if err != nil {
			log.Err(err).Msgf("couldn't convert count string into float, %T", course.Count)
			return models.StopListItems{}, err
		}

		stopList = append(stopList, models.StopListItem{
			ProductID: course.Id,
			Balance:   count,
		})
	}

	return stopList, nil
}
