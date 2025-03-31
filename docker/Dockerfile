# Utiliser l'image officielle de Go comme base
FROM golang:1.23.5-alpine 

# Installation des dépendances nécessaires
RUN apk add --no-cache gcc musl-dev

# Définir le répertoire de travail dans le conteneur
WORKDIR /app

# Copier les fichiers go.mod et go.sum d'abord pour tirer parti du cache Docker
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o forum main.go


# Installer les dépendances nécessaires
RUN apk add --no-cache ca-certificates

# Créer un utilisateur non-root pour des raisons de sécurité
RUN adduser -D -g '' appuser


# Définir l'utilisateur non-root
USER appuser

# Exposer le port sur lequel l'application va écouter
EXPOSE 8080

# Commande pour démarrer l'application
CMD ["./forum"]