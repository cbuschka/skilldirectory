#POSTGRES Migration

We are in the process of migrating from Cassandra to Postgres. We have decided to use the GORM (https://github.com/jinzhu/gorm)
to facitlitate this. This is introduces an additional level of abstraction from what we have done before. 

##Roadmap

- [x] Update models to GORM compliant schemas
- [x] Update model tests
- [ ] Update controllers to work with new models
- [ ] Update controller tests
- [ ] Add call to DB.Migrate somewhere it will be called on startup
- [ ] Duplicate all data from the Cassandra database to Postgres

##Tools
- [Branch on skilldirectoryinfra](https://github.com/maryvilledev/skilldirectoryinfra/tree/gorm_crud_demo/postgres)
that prototypes the database and has a script to start Postgres in a Docker container
- [Postman tests](https://github.com/maryvilledev/skilldirectory/tree/master/postman) that can verify correct behavior of the
server endpoints
