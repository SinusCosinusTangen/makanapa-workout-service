POSTGRES_ADMIN_DB="makanapa"
POSTGRES_ADMIN_USER="postgres"
POSTGRES_ADMIN_PASS="postgres"

docker run -d \
    --name golang-postgres \
    -p 5432:5432 \
    -e POSTGRES_DB=$POSTGRES_ADMIN_DB \
    -e POSTGRES_USER=$POSTGRES_ADMIN_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_ADMIN_PASS \
    postgres