WORKING_DIR := /opt/sale_bot_v2
BINARY_NAME := sale_bot_v2
LOGROTATE_CFG := /etc/logrotate.d/sale_bot_v2
SERVICE_CFG := /etc/systemd/system/sale_bot_v2.service
SERVICE_NAME := sale_bot_v2.service

install: stop_service build
	sudo mkdir -p ${WORKING_DIR}
	sudo cp ${BINARY_NAME} ${WORKING_DIR}/${BINARY_NAME}
	test -f ${WORKING_DIR}/config.json || sudo cp configs/config.json ${WORKING_DIR}/config.json
	sudo chown -R 1000:1000 ${WORKING_DIR}
	sudo cp cmd/sale_bot_v2/sale_bot_v2.logrotate ${LOGROTATE_CFG}
	sudo cp cmd/sale_bot_v2/sale_bot_v2.service ${SERVICE_CFG}
	sudo systemctl daemon-reload
	sudo systemctl restart ${SERVICE_NAME}

clean: stop_service
	sudo rm ${SERVICE_CFG}
	sudo systemctl daemon-reload
	sudo rm ${LOGROTATE_CFG}
	sudo rm ${WORKING_DIR}/${BINARY_NAME}
	sudo rm ${BINARY_NAME}

stop_service:
	sudo systemctl stop ${SERVICE_NAME} || echo "service may not exist or inactive"

build: go.mod go.sum
	go build -o ./${BINARY_NAME} cmd/sale_bot_v2/main.go