CREATE TABLE "usersession_token" (
    "usersession_token_id" serial,
    "cache_token" text,"table_name" text,"record_id" integer,"expiry_date" timestamptz,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    
    
    PRIMARY KEY ("usersession_token_id")
);
