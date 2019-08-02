CREATE TABLE "person" (
    "person_id" serial,
    "name" text,"email" text,"password" text,"phone" text,"role" text,"picture" text,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    
    "ulid" varchar(26),
    
    PRIMARY KEY ("person_id")
);
