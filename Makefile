NAME=ecuratif

# Lance la compilation et le web app
# Seulement utilisé pour développement
all: sass build

# Compile le fichier scss vers css
sass:
	@echo "Compiling css file"
	@sassc ui/static/sass/main.scss ui/static/sass/main.css

# Compile seulement côté golang
build:
	@echo "Compiling program"
	@go build -o $(NAME) ./cmd/
	@echo "Start E-Curatif"
	@./$(NAME)


# Lance la création de l'image ecuratif db
# Lance le contenaire
# sleep cmd permet une attente de 10s afin de laisser psql faire son taff
# Lance la compilation du web app (css et go)
install:
	@echo "building docker..."
	@docker buildx build -t ecuratif_db ./database/
	@docker run -d --name ecuratif_psql_container -p 5432:5432 ecuratif_db:latest
	@echo "installing billboard.js via npm..."
	@npm install --prefix=./ui/static/js/ billboard.js
	@echo "generating tls certification(localhost only)..."
	@mkdir tls && cd tls/ && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	@echo "Compiling css files..."
	@sassc ui/static/sass/main.scss ui/static/sass/main.css
	@echo "Compiling program..."
	@go build -o $(NAME) ./cmd/
	@echo "E-Curatif online!"
	@./$(NAME)

test:
	@echo "Compiling test"
	@go build -o test ./test/
	@./test

# Trouver moyen de choper ./launch PID pour envoyer SIGTERM
uninstall:
	@docker stop ecuratif_psql_container
	@docker rm ecuratif_psql_container
	@docker rmi ecuratif_db:latest

clean:
	rm ./$(NAME)
