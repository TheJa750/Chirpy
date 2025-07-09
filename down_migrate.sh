cd sql/schema
echo Peforming migration down...
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" down
cd ../..
