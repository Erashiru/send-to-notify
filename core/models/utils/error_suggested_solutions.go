package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

var suggestedSolutions = map[string]string{
	"Creation timeout expired":                                    "При создании заказа нет ответа на запрос (скорее всего проблема с интернетом на точке)",
	"is excluded from menu for order's table":                     "продукт убрали из меню, необходимо проверить в айко меню, есть ли данный продукт под другим айди, если нет, то удалить из меню, если есть, заменить айди продукта",
	"PRODUCT NOT FOUND IN POS MENU":                               "по данному айди нет продукта в айко меню, необходимо проверить в айко меню, есть ли данный продукт под другим айди, если нет, то удалить из меню, если есть, заменить айди продукта",
	"Cannot find fixed group modifiers item":                      "отправляемый нами атрибут не существует у продукта в айко, необходимо его удалить и добавить валидный атрибут",
	"has invalid group amount":                                    "мы не отправляем атрибут в нужном количестве, необходимо добавить нужный атрибут в продукт",
	"doesn't belong to your api login included organization list": "по данному апи-логину нет организации, по которой отправляем заказ, необходимо проверить валидность параметров key, organization_id, terminal_id в iiko_cloud",
	"is inactive. Only active products can be added to order":     "продукт неактивный. необходимо проверить, не удален ли продукт или не стоит ли на стопе",
}

func GetSuggestedSolution(input, logLink string) string {
	for key := range suggestedSolutions {
		if strings.Contains(input, key) {
			suggestedSolution, ok := suggestedSolutions[key]
			if ok {
				return suggestedSolution
			}
		}
	}

	return "Обратиться в Tech Support" + "\n" + fmt.Sprintf("Ссылка на лог в CloudWatch %s", logLink)
}

type ErrorType string

const (
	Accounting ErrorType = "Бухгалтерия (Iiko Office / Chain)"
	Pos        ErrorType = "Касса (Iiko Front/Api transport)"
	Menu       ErrorType = "Меню (Агрегатора/Выгрузка/Внешнее меню)"
	Settings   ErrorType = "Настройки"
	Aggregator ErrorType = "Агрегатор"
	Kwaaka     ErrorType = "Kwaaka"
)

func (et ErrorType) String() string {
	return string(et)
}

type commonErr struct {
	errType  string
	reason   string
	solution string
}

func GetReasonAndSolution(errInfo string) (string, string, string) {
	ce := &commonErr{}

	checkFuncs := []func(string) bool{
		ce.isAccountingError,
		ce.isMenuError,
		ce.isPosError,
		ce.isSettingsError,
		ce.isAggregatorError,
	}

	for _, f := range checkFuncs {
		if f(errInfo) {
			return ce.errType, ce.reason, ce.solution
		}
	}

	log.Trace().Msgf("UNKNOWN ERROR: %s\n", errInfo)
	ce.errType = Kwaaka.String()
	ce.reason = "Неизвестная ошибка"
	ce.solution = "Для того чтобы понять какая ошибка, нужно\nОбратиться по номеру службы поддержки Kwaaka https://wa.me/+77770224924"

	return ce.errType, ce.reason, ce.solution
}

func (ce *commonErr) isAccountingError(errInfo string) bool {
	ce.errType = Accounting.String()

	productNameRx := regexp.MustCompile(`Product\s+“([^”]+)”`)
	matchesProduct := productNameRx.FindStringSubmatch(errInfo)
	if len(matchesProduct) > 1 {
		if strings.Contains(errInfo, "CannotAddInactiveProductException") {
			ce.reason = fmt.Sprintf("блюдо %s не добавлена в прейскурант или не включена в прейскурант данного ресторана", matchesProduct[1])
			ce.solution = "Включить данное блюдо в приказ внутри Прейскуранта, в данном ресторане"
		} else if strings.Contains(errInfo, "ConstraintViolationException") {
			tableNameRx := regexp.MustCompile(`order's table\s+([0-9]+)\s+\(([a-zA-Z0-9-]+)\)`)
			matchesTable := tableNameRx.FindStringSubmatch(errInfo)
			if len(matchesTable) > 2 {
				ce.reason = fmt.Sprintf("блюдо %s удалено из меню в таблице %s с id %s.", matchesProduct[1], matchesTable[1], matchesTable[2])
			} else {
				ce.reason = fmt.Sprintf("блюдо %s удалено из меню", matchesProduct[1])
			}
			ce.solution = "1) Восстановить данное блюдо,\n2) Если создано другое блюдо, необходимо его синхронизировать в Личном кабинете Kwaaka, по инструкции \"Добавление блюда\""
		}

		return ce.check()
	}

	desc := extractDescription(errInfo)
	if desc != "" {
		reName := regexp.MustCompile(`Order item modifier\s+“([^”]+)”`)
		reMin := regexp.MustCompile(`min\s*=\s*(\d+)`)
		reMax := regexp.MustCompile(`max\s*=\s*(\d+)`)
		reActual := regexp.MustCompile(`actual\s*=\s*(\d+)`)

		nameMatches := reName.FindStringSubmatch(desc)
		minMatches := reMin.FindStringSubmatch(desc)
		maxMatches := reMax.FindStringSubmatch(desc)
		actualMatches := reActual.FindStringSubmatch(desc)

		if len(nameMatches) > 1 && len(minMatches) > 1 && len(maxMatches) > 1 && len(actualMatches) > 1 {
			ce.reason = fmt.Sprintf("У блюда %s, параметры минимального и максимального выбора группы модификаторов не совподают с параметрами установленными в Iiko Office/Chain\nу \"%s” стоит min = %s, max = %s, нужно отправлять = %s", nameMatches[1], nameMatches[1], minMatches[1], maxMatches[1], actualMatches[1])
			ce.solution = "1) Скорректировать, в личном кабинете Kwaaka параметры выбора группы модификатора, внутри блюда\n2) Если это новая группа мрдификаторов, необходимо сделать привязку или загрузить меню из внешнего меню Iikoweb"
			return ce.check()
		}
	}

	if strings.Contains(desc, "Comments are not allowed") {
		ce.reason = "Запрещена передача комментариев к позициям заказа"
		ce.solution = "Нужно включить комментарий внутри Iiko Office/Chain - Администрирование - Настройки торгового предприятия - Выбрать соответствующу кассу - Нажать на галочку \"Разрешить текстовые комментарии к позициям заказа\""
	} else if strings.Contains(desc, "Apilogin's license for using the Cloud API has expired") {
		ce.reason = "Ваша лицензия Api Iiko истекла"
		ce.solution = "Необходимо оплатить за лицензию Диллеру. Если уже оплатили, то необходимо уточнить статус у диллеров Iiko"
	} else if strings.Contains(errInfo, "Payment item can not be externally processed due to payment type settings") {
		ce.reason = fmt.Sprintf("Тип оплаты который используется, запрещен для ввода из вне")
		ce.solution = "Необходимо поменять его в Iiko Office/Chain, внутри розничных продаж - Типы оплат -  зайти в тип оплат и поменять его настрйоку на \"Внешний\" или \"Как на стороне ресторана и вне\""
	}

	return ce.check()
}

func (ce *commonErr) isMenuError(errInfo string) bool {
	ce.errType = Menu.String()

	if strings.Contains(errInfo, "ATTRIBUTE NOT FOUND IN POS MENU") {
		ce.reason = "Данный модификатор удален в выгрузке или во внешнем меню\n"
		ce.solution = "Если модификатор не актуальный:\n1) Удалить блюдо в меню агрегаторов (Glovo/Wolt/Yandex)\nЕсли модификатор актуальный:\n1) Проверить наличие модификатора в выгрузке меню/внешнем меню\n2) Восстановить/создать модификатор\n3) Обновить меню в личном кабинете Kwaaka\n4) Синхронизировать его в агрегатор меню (Glovo/Wolt/Yandex)\n5) Опубликовать меню в личном кабинете Kwaakan"
	}

	return ce.check()
}

func (ce *commonErr) isPosError(errInfo string) bool {
	ce.errType = Pos.String()

	if strings.Contains(errInfo, "Creation timeout expired, order automatically transited to error creation") {
		ce.reason = "В момент поступления заказа, касса: \n1) Отключена\n2) Не подключена к интернету/электроснабжению\n3) Не обменивается данными с Api"
		ce.solution = "В 1 и 2 случаях надо проверит: подключение интернета и электроснабжения\nВ 3 случае нужно чтобы ваш диллер проверил обмен транспорта кассы, и перезагрухил транспорт на кассе"
	} else if strings.Contains(errInfo, "OUT OF MEMORY") {
		ce.reason = "На кассе в момент поступления заказа, перегуржена оперативная память\nВ момент поступления заказа у кассы перегружена оперативная память"
		ce.solution = "Открыто слишком много приложений на точке, память на диске в ПК перегружена программами, недостаточно ОЗУ\nСвязаться с Iiko диллерами, проверить состояние памяти ОЗУ"
	} else if strings.Contains(errInfo, "System.ArgumentOutOfRangeException") && strings.Contains(errInfo, "but it cannot be greater than") {
		ce.reason = "Комментарий клиента, превышает ограничение Iiko в 255 символов"
		ce.solution = "Обратиться к Iiko по поводу увеличения ограничения в 255 символов"
	}

	return ce.check()
}

func (ce *commonErr) isSettingsError(errInfo string) bool {
	ce.errType = Settings.String()

	if strings.Contains(errInfo, "too small delivery date") {
		ce.reason = "Дата и время которое указано в заказе, не совподает с датой и временем в кассе"
		ce.solution = "Проверить акутальное время в заказе (агрегаторе) и в кассе\nЕсли время в агрегаторе некорректное, обратиться в СП агрегатора\nЕсли время в кассе не правильное:\n1) Поменять дату и время в моноблоке (Виндоусе)\n2) Поменять дату и время в Iiko Front"
	} else if strings.Contains(errInfo, "it's top parent group isn't allowed for current department") {
		ce.reason = "Тип заказа который передается в кассу, не соотвествует вашему сервису"
		ce.solution = "Нужно поменять доставку на корректную, выбор между \"Доставка самовывоз\" и \"Доставка курьером\""
	} else if strings.Contains(errInfo, "Code:DuplicatedOrderId Message:Order already exists") {
		ce.reason = "Данный заказ уже существует в кассе"
		ce.solution = "Удалить второй дублированный заказ\nОтключить настройку retry в личном кабинете Kwaaka"
	} else if strings.Contains(errInfo, "Entity not found: id=IOrderType, targetType") {
		ce.reason = "В выгрузке меню/внешнем меню у блюда нет группы к которому оно привязано"
		ce.solution = "Нужно привязать данное блюдо к группе вунтри Товары - Блюда или во внешнем меню"
	} else if strings.Contains(errInfo, "Cannot find fixed simple modifiers item") ||
		strings.Contains(errInfo, "Cannot find fixed group modifiers") {

		var modifier string
		var item string
		modifierRx := regexp.MustCompile(`'[^']+'`)
		modifierMatches := modifierRx.FindStringSubmatch(errInfo)
		if len(modifierMatches) > 1 {
			modifier = modifierMatches[1]
			if len(modifierMatches) > 2 {
				item = fmt.Sprintf(" в блюде %s", modifierMatches[2])
			}
		}
		ce.reason = fmt.Sprintf("Модификатор %s%s заведен в меню Iiko, но не добавлен в меню Агрегатора (Glovo, Wolt, Yandex)", modifier, item)
		ce.solution = "Необходимо Нажать обновить меню в личном кабинете Kwaaka"
	}

	if ce.check() {
		return true
	}

	var targetProduct string
	var errScaleProduct string
	scaleRx := regexp.MustCompile(`Product\s+(“[^”]+”\s+\([^)]+\)).*scale\s+(“[^”]+”\s+\([^)]+\))`)
	scaleMatches := scaleRx.FindStringSubmatch(errInfo)
	if len(scaleMatches) > 2 {
		targetProduct, errScaleProduct = scaleMatches[1], scaleMatches[2]
		ce.reason = fmt.Sprintf("Продукт %s имеет scale следующего продукта: %s", targetProduct, errScaleProduct)
		ce.solution = "Необходимо обратиться в техподержку квааки по этому поводу"
	}

	return ce.check()
}

func (ce *commonErr) isAggregatorError(errInfo string) bool {
	ce.errType = Aggregator.String()

	if strings.Contains(errInfo, "create order error: response status: 401 Unauthorized") {
		ce.reason = "Интеграция не включена на стороне Агрегатора"
		ce.solution = "Обратиться к агрегатору с запросом включить интеграцию по вашему ресторану"
	}

	return ce.check()
}

func (ce *commonErr) check() bool {
	return ce.reason != "" && ce.solution != ""
}

func extractDescription(input string) string {
	re := regexp.MustCompile(`"description":\s*"([^"]+)"`)

	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
