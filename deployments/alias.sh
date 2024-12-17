if [ $1 = "Stage" ]; then
  aws lambda update-function-configuration --function-name auto-update-aggregator-menu-rkeeper7xml --environment '{"Variables":$(CONFIG)}'
  	aws lambda update-alias --function-name auto-update-aggregator-menu-rkeeper7xml --name Stage --function-version $(VERSION)
else
  echo Try again
fi