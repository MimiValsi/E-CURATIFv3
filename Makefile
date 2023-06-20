NAME=launch

# Lauch everything
all: sass build start

# Build sass file to css
sass:
	sass --no-source-map ui/static/sass/main.scss:ui/static/sass/main.css

# Compile entirely the program
build:
	echo "Compiling program"
	go build -o $(NAME) ./cmd/

# start program (if exists) else run make all
start:
	if [ -f ./$(NAME) ]; then \
		./$(NAME) \
	else \
		make all; \
	fi

# run temporarely the program (dev only)
temp: 
	go run ./cmd/

clean:
	rm ./$(NAME)