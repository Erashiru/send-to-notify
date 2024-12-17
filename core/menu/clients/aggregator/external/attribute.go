package external

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	externalModels "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
)

func (m mnm) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	var modifiers = make([]externalModels.Modifier, 0, len(attributes))

	for _, attribute := range attributes {
		modifiers = append(modifiers, m.toModifier(ctx, attribute))
	}

	for _, modifier := range modifiers {
		if err := m.cli.UpdateModifierStopList(ctx, modifier, m.webhookAttributeStoplist); err != nil {
			return "", err
		}
	}

	return "", nil
}
