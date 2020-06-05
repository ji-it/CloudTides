# The `validate` target checks for errors and inconsistencies in 
# our specification of an API. This target can check if we're 
# referencing inexistent definitions and gives us hints to where
# to fix problems with our API in a static manner.
validate:
	swagger validate ./swagger/swagger.yml


# The `gen` target depends on the `validate` target as
# it will only succesfully generate the code if the specification
# is valid.
# 
# Here we're specifying some flags:
# --target              the base directory for generating the files;
# --spec                path to the swagger specification;
# --exclude-main        generates only the library code and not a 
#                       sample CLI application;
# --name                the name of the application.
gen: validate
	swagger generate server \
		--target=./pkg \
		--spec=./swagger/swagger.yml \
		--exclude-main \
		--name=cloudTides


# just added `gen` and `validate`
.PHONY: install gen validate
	/Users/szhe/go/src/tides-server/swagger/swagger.yml