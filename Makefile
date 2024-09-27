NAME=launch

# Lauch everything
all: sass build start

# Build sass file to css
sass:
	@echo "Compiling css file"
	@sassc ui/static/sass/main.scss ui/static/sass/main.css

# @sass --no-source-map ui/static/sass/main.scss:ui/static/sass/main.css
# Compile entirely the program
build:
	@echo "Compiling program"
	@go build -o $(NAME) ./cmd/
	@echo "Start E-Curatif"
	@make start

# start program (if exists) else run make all
start:
	@if [ -f ./$(NAME) ]; then \
		./$(NAME) \
	else \
		make all; \
	fi

# run temporarely the program (dev only)
temp: 
	go run ./cmd/

test:
	@echo "Compiling test"
	@go build -o test ./test/
	@./test

clean:
	rm ./$(NAME)
