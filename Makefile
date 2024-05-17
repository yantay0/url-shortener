# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = '174.138.89.148'

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh url_shortener@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api url_shortener@${production_host_ip}:~
	rsync -rP --delete ./migrations url_shortener@${production_host_ip}:~
	ssh -t url_shortener@${production_host_ip} 'migrate -path ~/migrations -database $$URL_SHORTENER_DB_DSN up'
