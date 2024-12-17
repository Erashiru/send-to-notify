.PHONY: zip docs

mock:
	mockgen --source=pkg/iiko/clients/iiko.go --destination=mocks/iiko.go -package=mocks
	mockgen --source=pkg/menu/client.go --destination=mocks/menu.go -package=mocks
tidy:
	go mod tidy

vet:
	go vet ./...
all: tidy mock vet

docs:
	swag fmt
	swag init --parseDependency --parseInternal -g cmd/integration_api/docs.go -o ./docs -d ./
	aws s3 cp docs/swagger.json s3://kwaaka-files/swagger/doc-integration-api.json

order-retry:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/order_retry/main.go
	zip build/order-retry.zip bin/main
	aws s3 cp build/order-retry.zip s3://kwaaka-files/functions/order-retry.zip
	aws lambda update-function-code --function-name order-retry \
		--s3-bucket kwaaka-files \
		--s3-key functions/order-retry.zip


preorder-on-time:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/preorder_send_on_time/main.go
	zip build/preorder-send-on-time.zip bin/main
		aws s3 cp build/preorder-send-on-time.zip s3://kwaaka-files/functions/preorder-send-on-time.zip
		aws lambda update-function-code --function-name preorder-on-time \
			--s3-bucket kwaaka-files \
			--s3-key functions/preorder-send-on-time.zip

virtual-store-handler:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/virtualstore/main.go
	zip build/virtual-store-handler.zip bin/main
	aws s3 cp build/virtual-store-handler.zip s3://kwaaka-files/functions/virtual-store-handler.zip
	aws lambda update-function-code --function-name virtual-store-handler \
    				--s3-bucket kwaaka-files \
    				--s3-key functions/virtual-store-handler.zip

virtual-store-cron:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/virtualstore/cron/main.go
	zip build/virtual-store-cron.zip bin/main
	aws s3 cp build/virtual-store-cron.zip s3://kwaaka-files/functions/virtual-store-cron.zip
	aws lambda update-function-code --function-name virtual-store-cron \
					--s3-bucket kwaaka-files \
					--s3-key functions/virtual-store-cron.zip

integration-api:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/integration_api/main.go
	zip ./integration-api.zip bin/main fonts/ArialCyrRegular.ttf fonts/kzArialBold.ttf
	aws s3 cp ./integration-api.zip s3://kwaaka-files/functions/integration-api.zip
	aws lambda update-function-code --function-name integration-api \
		--s3-bucket kwaaka-files \
		--s3-key functions/integration-api.zip

prerelease-integration-api:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bootstrap cmd/integration_api/main.go
	zip build/integration-api.zip bootstrap
	aws s3 cp build/integration-api.zip s3://kwaaka-files/functions/integration-api.zip
	aws lambda update-function-code --function-name prerelease-integration-api \
		--s3-bucket kwaaka-files \
		--s3-key functions/integration-api.zip

pos-is-alive:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/integration_api/check-in/main.go
	zip build/pos-is-alive.zip bin/main
	aws s3 cp build/pos-is-alive.zip s3://kwaaka-files/functions/pos-is-alive.zip
	aws lambda update-function-code --function-name pos-is-alive \
			--s3-bucket kwaaka-files \
			--s3-key functions/pos-is-alive.zip

store-schedule-update:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/store_cron/store_schedule_update/main.go
	zip ./store-schedule-update-cron.zip bin/main
		aws s3 cp ./store-schedule-update-cron.zip s3://kwaaka-files/functions/store-schedule-update-cron.zip
		aws lambda update-function-code --function-name store-schedule-update-cron \
			--s3-bucket kwaaka-files \
			--s3-key functions/store-schedule-update-cron.zip


store-status-check:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -tags lambda.norpc -o bootstrap cmd/store_cron/store_status_check/main.go
	zip ./store_status_check_cron.zip bootstrap
		aws s3 cp ./store_status_check_cron.zip s3://kwaaka-files/functions/store_status_check_cron.zip
		aws lambda update-function-code --function-name store_status_check_cron \
			--s3-bucket kwaaka-files \
			--s3-key functions/store_status_check_cron.zip

store-status-report:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -tags lambda.norpc -o bootstrap cmd/store_cron/store_status_report/main.go
	zip ./store_status_report.zip bootstrap
		aws s3 cp ./store_status_report.zip s3://kwaaka-files/functions/store_status_report_cron.zip
		aws lambda update-function-code --function-name store-status-report-cron \
			--s3-bucket kwaaka-files \
			--s3-key functions/store_status_report_cron.zip

offline-iiko-orders:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/offline_iiko_orders/main.go
	zip ./offline_iiko_orders.zip bin/main
		aws s3 cp ./offline_iiko_orders.zip s3://kwaaka-files/functions/offline_iiko_orders.zip
		aws lambda update-function-code --function-name offline_iiko_orders \
			--s3-bucket kwaaka-files \
			--s3-key functions/offline_iiko_orders.zip

cron-update-order-status-by-pos:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/crons/update_order_status_by_pos/main.go
	zip build/update_order_status_by_pos.zip bin/main
	aws s3 cp build/update_order_status_by_pos.zip s3://kwaaka-files/functions/update_order_status_by_pos.zip
	aws lambda update-function-code --function-name update-order-status \
		--s3-bucket kwaaka-files \
		--s3-key functions/update_order_status_by_pos.zip

cron-update-stoplist-by-pos:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/crons/update_stop_list_by_pos/main.go
	zip build/update_stop_list_by_pos.zip bin/main
	aws s3 cp build/update_stop_list_by_pos.zip s3://kwaaka-files/functions/update_stop_list_by_pos.zip
	aws lambda update-function-code --function-name manual-order-update \
		--s3-bucket kwaaka-files \
		--s3-key functions/update_stop_list_by_pos.zip

cron-update-stoplist-by-section:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -tags lambda.norpc -o bootstrap cmd/crons/update_stop_list_by_section/main.go
	zip ./update_stop_list_by_section.zip bootstrap
	aws s3 cp build/update_stop_list_by_section.zip s3://kwaaka-files/functions/update_stop_list_by_section.zip
	aws lambda update-function-code --function-name update-stop-list-by-section \
		--s3-bucket kwaaka-files \
		--s3-key functions/update_stop_list_by_section.zip

validate-store-menus:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o bin/main cmd/menu/validate_store_menus/main.go
	zip ./validate_store_menus.zip bin/main
	aws s3 cp ./validate_store_menus.zip s3://kwaaka-files/functions/validate_store_menus.zip
	aws lambda update-function-code --function-name matching-validate \
    		--s3-bucket kwaaka-files \
    		--s3-key functions/validate_store_menus.zip

order-status-to-closed-cron:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/order_closed_cron/main.go
	zip ./order_closed_cron.zip bootstrap
	aws s3 cp ./order_closed_cron.zip s3://kwaaka-files/functions/order_closed_cron.zip
	aws lambda update-function-code --function-name orderStatusToClose \
		--s3-bucket kwaaka-files \
		--s3-key functions/order_closed_cron.zip

order-stat:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/order_stat/main.go
	zip ./order_stat_cron.zip bootstrap
	aws s3 cp ./order_stat_cron.zip s3://kwaaka-files/functions/order_stat_cron.zip
	aws lambda update-function-code --function-name orderStat \
		--s3-bucket kwaaka-files \
		--s3-key functions/order_stat_cron.zip

migrate-legal-entity-payment-docker:
	migrate -path ./migrations -database 'postgres://localhost:5432/kwaaka?user=kwaaka&password=kwaaka&sslmode=disable' up

accept-telegram-message:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/order_status_change/main.go
	zip ./order_status_change.zip bootstrap
	aws s3 cp ./order_status_change.zip s3://kwaaka-files/functions/order_status_change.zip
	aws lambda update-function-code --function-name accept_telegram_message \
		--s3-bucket kwaaka-files \
		--s3-key functions/order_status_change.zip

s3-wolt-images:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/menu/wolt_images_upload_to_s3/main.go
	zip ./wolt_images_upload_to_s3.zip bootstrap
	aws s3 cp ./wolt_images_upload_to_s3.zip s3://kwaaka-files/functions/wolt_images_upload_to_s3.zip
	aws lambda update-function-code --function-name s3-upload-images \
		--s3-bucket kwaaka-files \
		--s3-key functions/wolt_images_upload_to_s3.zip

whatsapp-messages:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/whatsapp_messages/main.go
	zip ./whatsapp_messages.zip bootstrap
	aws s3 cp ./whatsapp_messages.zip s3://kwaaka-files/functions/whatsapp_messages.zip
	aws lambda update-function-code --function-name whatsapp_messages \
		--s3-bucket kwaaka-files \
		--s3-key functions/whatsapp_messages.zip

indrive-3pl:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/create_delivery/main.go
	zip ./indrive_create_delivery.zip bootstrap
	aws s3 cp ./indrive_create_delivery.zip s3://kwaaka-files/functions/Create3plForIndrive.zip
	aws lambda update-function-code --function-name Create3plForIndrive \
		--s3-bucket kwaaka-files \
		--s3-key functions/Create3plForIndrive.zip

notify-unpaid-customers:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/notify_unpaid_customers/main.go
	zip ./notify_unpaid_customers.zip bootstrap
	aws s3 cp ./notify_unpaid_customers.zip s3://kwaaka-files/functions/notify_unpaid_customers.zip
	aws lambda update-function-code --function-name notify_unpaid_customers \
		--s3-bucket kwaaka-files \
		--s3-key functions/notify_unpaid_customers.zip

send-payment-whatsapp:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/send_payment_in_whatsapp/main.go
	zip ./send_payment_in_whatsapp.zip bootstrap fonts/ArialCyrRegular.ttf fonts/kzArialBold.ttf
	aws s3 cp ./send_payment_in_whatsapp.zip s3://kwaaka-files/functions/send_payment_in_whatsapp.zip
	aws lambda update-function-code --function-name send-payment-whatsapp \
		--s3-bucket kwaaka-files \
		--s3-key functions/send_payment_in_whatsapp.zip

kaspi-salescout:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/kaspi_salescout/main.go
	zip ./kaspi_salescout.zip bootstrap
	aws s3 cp ./kaspi_salescout.zip s3://kwaaka-files/functions/kaspi_salescout.zip
	aws lambda update-function-code --function-name kaspi-salescout \
		--s3-bucket kwaaka-files \
		--s3-key functions/kaspi_salescout.zip


delivery-dispatcher-price:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/delivery_dispatcher_price/main.go
	zip ./delivery_dispatcher_price.zip bootstrap
	aws s3 cp ./delivery_dispatcher_price.zip s3://kwaaka-files/functions/delivery_dispatcher_price.zip
	aws lambda update-function-code --function-name delivery-dispatcher-price \
		--s3-bucket kwaaka-files \
		--s3-key functions/delivery_dispatcher_price.zip

deferred-status-send:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/send_defer_statuses/main.go
	zip ./send_defer_statuses.zip bootstrap
	aws s3 cp ./send_defer_statuses.zip s3://kwaaka-files/functions/send_defer_statuses.zip
	aws lambda update-function-code --function-name DeferredStatusSend \
		--s3-bucket kwaaka-files \
		--s3-key functions/send_defer_statuses.zip

no-dispatcher-message:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/no_dispatcher_message/main.go
	zip ./no_dispatcher_message.zip bootstrap
	aws s3 cp ./no_dispatcher_message.zip s3://kwaaka-files/functions/no_dispatcher_message.zip
	aws lambda update-function-code --function-name no-dispatcher-message \
		--s3-bucket kwaaka-files \
		--s3-key functions/no_dispatcher_message.zip

multricard-refresh-token:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/refresh_multicard_token/main.go
	zip ./refresh_multicard_token.zip bootstrap
	aws s3 cp ./refresh_multicard_token.zip s3://kwaaka-files/functions/multricard-refresh-token.zip
	aws lambda update-function-code --function-name multricard-refresh-token \
		--s3-bucket kwaaka-files \
		--s3-key functions/multricard-refresh-token.zip

performer-lookup-time-more-15-minute:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/performer_lookup_time_more_15_minute/main.go
	zip ./performer_lookup_time_more_15_minute.zip bootstrap
	aws s3 cp ./performer_lookup_time_more_15_minute.zip s3://kwaaka-files/functions/performer_lookup_time_more_15_minute.zip
	aws lambda update-function-code --function-name performer-lookup-time-more-15-minute \
		--s3-bucket kwaaka-files \
		--s3-key functions/performer_lookup_time_more_15_minute.zip

order-auto-close:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/crons/order_auto_close/main.go
	zip ./order_auto_close.zip bootstrap
	aws s3 cp ./order_auto_close.zip s3://kwaaka-files/functions/order_auto_close.zip
	aws lambda update-function-code --function-name order-auto-close \
		--s3-bucket kwaaka-files \
		--s3-key functions/order_auto_close.zip

wolt-discount-run:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/wolt_discount_run/main.go
	zip ./wolt_discount_run.zip bootstrap
	aws s3 cp ./wolt_discount_run.zip s3://kwaaka-files/functions/wolt_discount_run.zip
	aws lambda update-function-code --function-name wolt-discount-run \
		--s3-bucket kwaaka-files \
		--s3-key functions/wolt_discount_run.zip

CONFIG = $(shell echo | sh ./deployments/.sh)

VERSION = $(shell echo | aws lambda publish-version --function-name auto-update-aggregator-menu-rkeeper7xml | jq -r .Version)


alias-name: conf
ifdef DESCRIPTION
	aws lambda update-function-configuration --function-name auto-update-aggregator-menu-rkeeper7xml --environment '{"Variables":$(CONFIG)}'
	aws lambda update-alias --function-name auto-update-aggregator-menu-rkeeper7xml --name Stage --function-version $(VERSION)
endif

conf:
	$(bash if [ $(foo $arg) -eq 0 ] ; then echo Hello)


auto-update-aggregator-menu-rkeeper7xml:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/rkeeper7xml/auto_update_menu_cron/main.go
	zip ./auto_update_menu_cron.zip bootstrap
	aws s3 cp ./auto_update_menu_cron.zip s3://kwaaka-files/functions/auto_update_menu_cron.zip
	aws lambda update-function-code --function-name auto-update-aggregator-menu-rkeeper7xml \
		--s3-bucket kwaaka-files \
		--s3-key functions/auto_update_menu_cron.zip
	$(shell sh ./deployments/alias.sh)

salescout-proxy:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -o bootstrap cmd/starter_app_salescout_proxy/main.go
	zip ./kaspi-salescout-payment-api.zip bootstrap
	aws s3 cp ./kaspi-salescout-payment-api.zip s3://kwaaka-files/functions/kaspi-salescout-payment-api.zip
	aws lambda update-function-code --function-name kaspi-salescout-payment-api \
		--s3-bucket kwaaka-files \
		--s3-key functions/kaspi-salescout-payment-api.zip