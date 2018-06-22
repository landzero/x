dialects=("postgres" "mysql" "mssql" "sqlite")

for dialect in "${dialects[@]}" ; do
    DEBUG=false ORM_DIALECT=${dialect} go test
done
