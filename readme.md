# Builder - WIP
Builder is designed to quickly scaffold boilerplate code files using templates. 

It simply connects to the database defined in .env and fills in the tmpl files in the blueprints folder with varients of field names.

## Available field names

- TableName `person` `image_meta`
- TableNameSpaces `person` `image meta`
- TableNamePascal `Person` `ImageMeta`
- TableNameCamel `person` `imageMeta`
- TableNameLower `person` `imagemeta`
- TableNamePlural `people` `image_metas`
- TableNamePluralPascal `people` `ImageMetas`
- TableNamePluralCamel `people` `imageMetas`
- TableID `person_ulid`
- TableIDTitle `PersonULID`
- TableIDCamel `personULID`

Currently it generates code for 

- table creation / migrations
- RESTful 
- models
- list pages (blueprint is for vue 3)
- edit pages (blueprint is for vue 3)
- API definition

Take a look at blueprints for an example. Currently the blueprints are very much tied to my take on a web framework in Golang. I hope to get to a point where I can put it on github in an understandable manor.

__ You can change the blueprints to whatever you want, just follow the naming convention __

## Installation
`go get github.com/nerdynz/builder`

add a `rest/migrations` folder to test migrations


## NEEDS CLEANUP
- case statements and logic a rebit all over the place, first use of bubble tea, so it was a bit of a peacemeal process