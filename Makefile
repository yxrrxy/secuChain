# 定义变量
OPENSCA_SCRIPT_URL := https://raw.githubusercontent.com/XmirrorSecurity/OpenSCA-cli/master/scripts/install.ps1
PROJECT_PATH := E:/Github/blockSBOM/backend
CONFIG_PATH := E:/Github/blockSBOM/backend/app/sbom/config.json
OUTPUT_FILENAME := sbom_output
OUTPUT_SUFFIX := json
TOKEN := your_token_here
RPC_SERVICE := your_rpc_service_url
IMAGE_UPLOAD_URL := http://example.com/upload_image
RESULT_UPLOAD_URL := http://example.com/upload_result

# 默认目标
all: install_opensca generate_database convert_database transmit_image start_opensca_server opensca_scan transmit_results

# 下载并安装 OpenSCA
install_opensca:
	powershell -Command "iex \"&{$(irm $(OPENSCA_SCRIPT_URL) )}\""

# 生成数据库
generate_database:
	echo "Generating database..."
	cd E:/Github/blockSBOM/backend/app/vuln && go run vuln.go

# 将数据库转换为图像
convert_database:
	echo "Converting database to image..."
	cd E:/Github/blockSBOM/backend/output/service && python vuln_chart.py

# 回传图像
transmit_image:
	echo "Transmitting image..."
	curl -F "file=@E:/Github/blockSBOM/backend/output/service/vuln.png" $(IMAGE_UPLOAD_URL)

# 启动 OpenSCA 服务
start_opensca_server:
	echo "Starting OpenSCA server..."
	cd E:/Github/blockSBOM/backend/app/sbom && go run server.go

# 使用 OpenSCA 工具检测项目并生成 SBOM 文件与漏洞文件
opensca_scan:
	echo "Scanning project with OpenSCA..."
	cd E:/Github/blockSBOM/backend && opensca-cli -path $(PROJECT_PATH) -config $(CONFIG_PATH) -out $(OUTPUT_FILENAME).$(OUTPUT_SUFFIX) -token $(TOKEN)

# 回传 SBOM 图表与漏洞图表
transmit_results:
	echo "Transmitting SBOM and vulnerability results..."
	curl -F "file=@E:/Github/blockSBOM/backend/$(OUTPUT_FILENAME).$(OUTPUT_SUFFIX)" $(RESULT_UPLOAD_URL)

# 启动环境（数据库和Fabric网络）
.PHONY: env-up
env-up:
	@echo "Starting environment..."
	docker-compose up -d mysql

# 停止环境
.PHONY: env-down
env-down:
	@echo "Stopping environment..."
	docker-compose down
	cd scripts/fabric/test-network && ./network.sh down

# 启动 API 服务
.PHONY: api
api:
	@echo "Starting API server..."
	go run internal/cmd/api/main.go

# 开发模式启动（带热重载）
.PHONY: dev
dev:
	@echo "Starting API server in development mode..."
	air -c .air.toml

# 清理
clean:
	rm -f E:/Github/blockSBOM/backend/output/service/vuln.png
	rm -f E:/Github/blockSBOM/backend/$(OUTPUT_FILENAME).$(OUTPUT_SUFFIX)