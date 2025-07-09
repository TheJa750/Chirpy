cd sql/schema
echo Peforming migration up...
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" up
cd ../..
