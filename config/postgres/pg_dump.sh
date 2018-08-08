pg_dump -U postgres --data-only -t category > init/production/category.psql
pg_dump -U postgres --data-only -t country > init/production/country.psql
pg_dump -U postgres --data-only -t gender > init/production/gender.psql
pg_dump -U postgres --data-only -t language > init/production/language.psql
pg_dump -U postgres --data-only -t post_type > init/production/post_type.psql
pg_dump -U postgres --data-only -t reaction > init/production/reaction.psql