package managers

import (
	"testing"
)

func Test_skipErrors(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "#1",
			args:    "{\"code\":\"Common\",\"message\":{\"orderId\":\"2f0c4821-1aba-488f-bee7-b0cc4e7dd519\",\"terminalGroup\":{\"id\":\"8960d6f2-e5e4-4d76-ba1c-aa7391c95347\",\"name\":\"\u041a\u043e\u0448\u043a\u0430\u0440\u0431\u0430\u0435\u0432\u0430 \u0434\u043e\u0441\u0442\u0430\u0432\u043a\u0430\"},\"timestamp\":1670985096795,\"code\":500,\"message\":\"Resto.Front.Api.Exceptions.EntityNotFoundException: Cannot find fixed simple modifiers item '- \u0441 \u0441\u043e\u0431\u043e\u0439' in order item '\u0417\u0430\u0432\u0442\u0440\u0430\u043a \u0421\u043a\u0440\u044d\u043c\u0431\u043b \u0441 \u043b\u043e\u0441\u043e\u0441\u0435\u043c' (Id = ce392463-61cc-4ee6-bb8f-0aa11cd70c46) Server stack trace: \u0432 Resto.Front.Api.V8Preview3.Editors.Actions.AddOrderModifierItemWithPredefinedId.CreateModifier(Guid id, Int32 amount, Nullable`1 payableAmount, Decimal parentAmount, Nullable`1 predefinedPrice, IProduct modifier, IProductGroup parentGroup, IBaseOrderBuilder order, ICollection`1 simpleBindings, ICollection`1 groupBindings, ICookingOrderItemBuilder cookingItem, Boolean generateItemSaleEventIdSecondary, SessionContext context) \u0432 K: BuildAgent work master-installer dev iikoFront.NetApiResto.Front.ApiV8Preview3Editors\\\\\\\\Actions\\\\\\\\AddOrderModifierItem.cs:\u0441\u0442\u0440\u043e\u043a\u0430 152\\\\r\\\\n   \u0432 Resto.Front.Api.V8Preview3.Editors.Actions.AddOrderModifierItemWithPredefinedId.Apply(Guid id, Int32 amount, Nullable`1 payableAmount, Nullable`1 predefinedPrice, IProduct modifier, IProductGroup parentGroup, IBaseOrderBuilder order, IProductOrderItemBuilder productItem, SessionContext context) \u0432 K:\\\\\\\\BuildAgent\\\\\\\\work\\\\\\\\master-installer\\\\\\\\dev\\\\\\\\iikoFront.Net\\\\\\\\Api\\\\\\\\Resto.Front.Api\\\\\\\\V8Preview3\\\\\\\\Editors\\\\\\\\Actions\\\\\\\\AddOrderModifierItem.cs:\u0441\u0442\u0440\u043e\u043a\u0430 120\\\\r\\\\n   \u0432 Resto.Front.Api.V8Preview3.Editors.Actions.AddOrderModifierItemWithPredefinedId.Resto.Front.Api.V8Preview3.Editors.IEditActionBase.WriteChanges(SessionContext context, IDictionary`2 createdEntities, IDictionary`2 actionNumbersToCreatedApiEntities, Int32 i) \u0432 K:\\\\\\\\BuildAgent\\\\\\\\work\\\\\\\\master-installer\\\\\\\\dev\\\\\\\\iikoFront.Net\\\\\\\\Api\\\\\\\\Resto.Front.Api\\\\\\\\V8Preview3\\\\\\\\Editors\\\\\\\\Editors.g.cs:\u0441\u0442\u0440\u043e\u043a\u0430 1110\\\\r\\\\n   \u0432 Resto.Front.Api.V8Preview3.Editors.EditSessionWriter.WriteChanges() \u0432 K:\\\\\\\\BuildAgent\\\\\\\\work\\\\\\\\master-installer\\\\\\\\dev\\\\\\\\iikoFront.Net\\\\\\\\Api\\\\\\\\Resto.Front.Api\\\\\\\\V8Preview3\\\\\\\\Editors\\\\\\\\EditSessionWriter.cs:\u0441\u0442\u0440\u043e\u043a\u0430 202\\\\r\\\\n   \u0432 Resto.Front.Api.V8Preview3.OperationServiceInternal.SubmitChanges(IUser user, IEditSession editSession) \u0432 K:\\\\\\\\BuildAgent\\\\\\\\work\\\\\\\\master-installer\\\\\\\\dev\\\\\\\\iikoFront.Net\\\\\\\\Api\\\\\\\\Resto.Front.Api\\\\\\\\V8Preview3\\\\\\\\OperationServiceInternal.cs:\u0441\u0442\u0440\u043e\u043a\u0430 1967\\\\r\\\\n   \u0432 Resto.Front.Api.V8Preview3.OperationServiceInternal.Resto.Front.Api.IOperationServiceInternal.SubmitChanges(ICredentials credentials, IEditSession editSession) \u0432 K:\\\\\\\\BuildAgent\\\\\\\\work\\\\\\\\master-installer\\\\\\\\dev\\\\\\\\iikoFront.Net\\\\\\\\Api\\\\\\\\Resto.Front.Api\\\\\\\\V8Preview3\\\\\\\\Operations.g.cs:\u0441\u0442\u0440\u043e\u043a\u0430 2557\\\\r\\\\n   \u0432 System.Runtime.Remoting.Messaging.StackBuilderSink._PrivateProcessMessage(IntPtr md, Object[] args, Object server, Object[]& outArgs)\\\\r\\\\n   \u0432 System.Runtime.Remoting.Messaging.StackBuilderSink.SyncProcessMessage(IMessage msg)\\\\r\\\\n\\\\r\\\\nException rethrown at [0]: \\\\r\\\\n   \u0432 Resto.Front.Api.iikoTransport.Extensions.OperationServiceExtensions.ExecuteContinuousOperationWithCorrectErrorHandling(IOperationService operationService, Action`1 action)\\\\r\\\\n   \u0432 Resto.Front.Api.iikoTransport.CreateDeliveryOrder.OrderConverter.Convert(CreateOrderRequest request, IOrderType orderType) \u0432 Resto.Front.Api.iikoTransport.CreateDeliveryOrder.CreateDeliveryOrderProcessor.Process(CreateOrderRequest request, CommandContextLogger log)\",\"description\":\"Cannot find fixed simple modifiers item '- \u0441 \u0441\u043e\u0431\u043e\u0439' in order item '\u0417\u0430\u0432\u0442\u0440\u0430\u043a \u0421\u043a\u0440\u044d\u043c\u0431\u043b \u0441 \u043b\u043e\u0441\u043e\u0441\u0435\u043c' (Id = ce392463-61cc-4ee6-bb8f-0aa11cd70c46)\",\"additionalData\":null},\"description\":null,\"additionalData\":null}",
			wantErr: false,
		},
		{
			name:    "#2",
			args:    "{\"code\":\"Common\",\"message\":{\n \"orderId\": \"9cff39b3-d517-4a9c-9595-c730c48f80d7\",\n \"terminalGroup\": {\n   \"id\": \"9c03ddbc-70fa-4020-a787-d29d5f3adb4b\",\n   \"name\": \"Free Dog Pizza Фартуна\"\n },\n \"timestamp\": 1672725972351,\n \"code\": 500,\n \"message\": \"Resto.Front.Api.Exceptions.ConstraintViolationException: Order item modifier “Колбаски 1шт” (39013999-951a-44ce-b633-642fcf3a7a84) has invalid group amount: min = 1, max = 1, actual = 0. Ensure that interconnected product and modifier changes are in the same edit session.\\r\\n\\r\\nServer stack trace: \\r\\n   в Resto.Front.Api.V8Preview1.Editors.EditSessionWriter.CheckMandatoryGroupModifiers(ICollection`1 modifiers, IEnumerable`1 groupBindings, Decimal parentAmount) в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Editors\\\\EditSessionWriter.Checks.cs:строка 0\\r\\n   в Resto.Front.Api.V8Preview1.Editors.EditSessionWriter.CheckMandatoryGroupModifiers(ICookingOrderItemBuilder item) в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Editors\\\\EditSessionWriter.Checks.cs:строка 344\\r\\n   в Resto.Front.Api.V8Preview1.Editors.EditSessionWriter.CheckMandatoryModifiers(IBaseOrderBuilder order) в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Editors\\\\EditSessionWriter.Checks.cs:строка 282\\r\\n   в Resto.Front.Api.V8Preview1.Editors.EditSessionWriter.CheckModifiers() в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Editors\\\\EditSessionWriter.Checks.cs:строка 276\\r\\n   в Resto.Front.Api.V8Preview1.Editors.EditSessionWriter.PostCheckEditSession() в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Editors\\\\EditSessionWriter.cs:строка 71\\r\\n   в Resto.Front.Api.V8Preview1.OperationServiceInternal.SubmitChanges(IUser user, IEditSession editSession) в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\OperationServiceInternal.cs:строка 1965\\r\\n   в Resto.Front.Api.V8Preview1.OperationServiceInternal.Resto.Front.Api.IOperationServiceInternal.SubmitChanges(ICredentials credentials, IEditSession editSession) в M:\\\\BuildAgent\\\\work\\\\master-installer\\\\dev\\\\iikoFront.Net\\\\Api\\\\Resto.Front.Api\\\\V8Preview1\\\\Operations.g.cs:строка 2474\\r\\n   в System.Runtime.Remoting.Messaging.StackBuilderSink._PrivateProcessMessage(IntPtr md, Object[] args, Object server, Object[]& outArgs)\\r\\n   в System.Runtime.Remoting.Messaging.StackBuilderSink.SyncProcessMessage(IMessage msg)\\r\\n\\r\\nException rethrown at [0]: \\r\\n   в Resto.Front.Api.iikoTransport.Extensions.OperationServiceExtensions.ExecuteContinuousOperationWithCorrectErrorHandling(IOperationService operationService, Action`1 action)\\r\\n   в Resto.Front.Api.iikoTransport.CreateDeliveryOrder.OrderConverter.Convert(CreateOrderRequest request, IOrderType orderType)\\r\\n   в Resto.Front.Api.iikoTransport.CreateDeliveryOrder.CreateDeliveryOrderProcessor.Process(CreateOrderRequest request, CommandContextLogger log)\",\n \"description\": \"Order item modifier “Колбаски 1шт” (39013999-951a-44ce-b633-642fcf3a7a84) has invalid group amount: min = 1, max = 1, actual = 0. Ensure that interconnected product and modifier changes are in the same edit session.\",\n \"additionalData\": null\n}}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := skipErrors(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("skipErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
