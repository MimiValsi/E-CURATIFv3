FROM postgres:16-alpine

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=pass

COPY create_db.sql /docker-entrypoint-initdb.d/

# Les trois étapes pour lancer contenaire
# docker build -t ecuratif_db ./
# docker run -d --name psql_container -p 5432:5432 ecuratif_db:latest
# docker exec -it psql_container bash -c "psql -d ecuratif -U postgres"


