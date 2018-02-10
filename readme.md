# Builder - WIP
Builder is designed to quickly scaffold boilerplate code files using templates 

It simply connects to the database defined in .env and fills in the tmpl files in the blueprints folder with varients of field names.

## Available field names

- TableNameSpaces
- TableNameTitle
- TableNameCamel
- TableNameLower
- TableNamePlural
- TableNamePluralTitle
- TableNamePluralCamel
- TableID
- TableIDTitle
- TableIDCamel
- TableIDCamelWithRecord

Currently it generates code for 

- table creation / migrations
- RESTful 
- models
- list pages
- edit pages

Take a look at blueprints for an example. Currently the blueprints are very much tied to my own custom CMS which I hope to get to a point where I can put it on github.

###You can change the blueprints to whatever you want, just follow the naming convention

## Installation
`go get https://github.com/nerdynz/builder.git`

