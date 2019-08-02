CREATE TABLE "person" (
    "person_id" serial,
    "name" text,"email" text,"password" text,"phone" text,"role" text,"picture" text,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    
    "ulid" varchar(26),
    
    PRIMARY KEY ("person_id")
);

CREATE TABLE "usersession_token" (
    "usersession_token_id" serial,
    "cache_token" text,"table_name" text,"record_id" integer,"expiry_date" timestamptz,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    
    
    PRIMARY KEY ("usersession_token_id")
);

CREATE TABLE page (
    "page_id" SERIAL PRIMARY KEY,
    "title" text,
    "slug" text,
    "summary" text,
    "keywords" text,
    "preview_picture" text,
    "kind" text,
    "html" text,
    "html_two" text,
    "html_three" text,
    "html_four" text,
    "html_five" text,
    "html_six" text,
    "html_seven" text,
    "html_eight" text,
    "html_nine" text,
    "html_ten" text,
    "html_eleven" text,
    "html_twelve" text,
    "picture" text,
    "picture_two" text,
    "picture_three" text,
    "picture_four" text,
    "picture_five" text,
    "picture_six" text,
    "misc" text,
    "misc_two" text,
    "misc_three" text,
    "misc_four" text,
    "misc_five" text,
    "misc_six" text,
    "is_locked_slug" boolean,
    "is_special_page" boolean,
    "special_page_for" text,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    "show_in_nav" text,
    "subtitle" text,
    "sort_position" integer,
    "ulid" varchar(26),
    "show_title" boolean,
    "show_subtitle" boolean,
    "seo_title" text,
    "color" text
);

CREATE TABLE block (
    "block_id" SERIAL PRIMARY KEY,
    "picture" text,
    "picture_two" text,
    "picture_three" text,
    "picture_four" text,
    "picture_five" text,
    "picture_six" text,
    "html" text,
    "html_two" text,
    "html_three" text,
    "html_four" text,
    "html_five" text,
    "html_six" text,
    "page_id" integer,
    "date_created" timestamptz,
    "date_modified" timestamptz,
    "content_from_table" text,
    "content_from_table_two" text,
    "content_from_id" integer,
    "content_from_id_two" integer,
    "type" text,
    "sort_position" integer,
    "ulid" varchar(26),
    "additional" text,
    "additional_two" text,
    "additional_three" text,
    "additional_four" text
);