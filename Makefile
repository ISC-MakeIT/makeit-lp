COMPOSE := docker compose

.PHONY: help up stop down logs ps build db reset env

help:
	@echo "使い方: make <target>"
	@echo ""
	@echo "  up      起動"
	@echo "  stop    停止"
	@echo "  down    停止してコンテナを削除 (DBデータは残る)"
	@echo "  logs    ログを追う (例: make logs SERVICE=backend)"
	@echo "  ps      状態を表示"
	@echo "  build   イメージを作り直して起動"
	@echo "  db      MySQL に入る"
	@echo "  reset   ボリュームごと削除 (DBデータも消える)"
	@echo "  env     .env が無ければ .env.example から作る"

up:
	$(COMPOSE) up -d

stop:
	$(COMPOSE) stop

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f $(SERVICE)

ps:
	$(COMPOSE) ps

build:
	$(COMPOSE) up -d --build

db:
	$(COMPOSE) exec db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE'

reset:
	$(COMPOSE) down -v

env:
	@test -f .env || cp .env.example .env
