# Makefile pour simplifier les commandes Docker du projet forum

.PHONY: build up down restart logs clean

# Variables
COMPOSE = docker-compose

# Construire les images Docker
build:
	$(COMPOSE) build

# Démarrer tous les services
up:
	$(COMPOSE) up -d

# Arrêter tous les services
down:
	$(COMPOSE) down

# Redémarrer tous les services
restart:
	$(COMPOSE) restart

# Afficher les logs
logs:
	$(COMPOSE) logs -f

# Afficher les logs d'un service spécifique
logs-app:
	$(COMPOSE) logs -f app

logs-db:
	$(COMPOSE) logs -f db

# Nettoyer les ressources Docker
clean:
	$(COMPOSE) down -v --remove-orphans

# Créer une migration pour la base de données
migrate:
	$(COMPOSE) exec app go run /app/internal/database/migrations/migrate.go

# Ouvrir un shell dans un conteneur
shell-app:
	$(COMPOSE) exec app sh

shell-db:
	$(COMPOSE) exec db bash

# Commande pour ajouter un nouvel utilisateur admin
add-admin:
	@echo "Création d'un nouvel administrateur..."
	@read -p "Nom d'utilisateur: " username; \
	read -p "Email: " email; \
	read -s -p "Mot de passe: " password; \
	echo ""; \
	$(COMPOSE) exec db mysql -u forum_user -pforum_password forum_db -e "INSERT INTO users (id, username, email, password, role) VALUES (UUID(), '$$username', '$$email', '$$password', 'admin');"
	@echo "Administrateur créé avec succès!"