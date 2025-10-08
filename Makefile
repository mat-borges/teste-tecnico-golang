# -------------------------
# 🔧 VARIÁVEIS
# -------------------------
GO          := go
PKGS        := ./internal/aggregator/... ./internal/graph/... ./internal/fetcher/...
COVER_FILE  := coverage.out
COVER_HTML  := coverage.html

# Define modo silencioso para testes
export LOG_MODE := silent

# -------------------------
# ⚙️ COMANDOS PRINCIPAIS
# -------------------------

# ✅ Roda testes só nos pacotes principais (sem ruído)
test:
	@echo "Rodando testes silenciosos nos pacotes principais..."
	@$(GO) test $(PKGS) -v -cover -count=1

test-agg:
	@echo "Testando apenas aggregator..."
	@$(GO) test ./internal/aggregator/... -v -cover

# 📊 Gera relatório de cobertura (HTML)
test-report:
	@echo "Gerando relatorio de cobertura..."
	@$(GO) test $(PKGS) -coverprofile=$(COVER_FILE)
	@$(GO) tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
	@echo "✅ Relatório gerado em $(COVER_HTML)"

# 🧹 Limpa arquivos temporários
clean:
	@echo "Limpando arquivos temporarios..."
	@rm -f $(COVER_FILE) $(COVER_HTML)

# 🐳 Sobe o ambiente Docker
docker-up:
	@echo "Subindo o ambiente com Docker Compose..."
	docker compose up --build

# 🛑 Para e remove containers locais
docker-down:
	docker compose down --remove-orphans

# 🧭 Roda o servidor localmente
run:
	@echo "Iniciando servidor local..."
	@$(GO) run cmd/api/main.go
