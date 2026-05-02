ALTER TABLE hosts 
ALTER COLUMN address TYPE text[] 
USING ARRAY[address];

ALTER TABLE hosts 
ALTER COLUMN mac_address TYPE text[] 
USING ARRAY[mac_address];

ALTER TABLE hosts 
ALTER COLUMN host_name TYPE text[] 
USING ARRAY[host_name];

update hosts SET address = array[]::text[] WHERE address is NULL;
update hosts SET mac_address = array[]::text[] WHERE mac_address is NULL;
update hosts SET host_name = array[]::text[] WHERE host_name is NULL;

ALTER TABLE IF EXISTS hosts
    ALTER COLUMN address SET DEFAULT array[]::text[],
    ALTER COLUMN mac_address SET DEFAULT array[]::text[],
    ALTER COLUMN host_name SET DEFAULT array[]::text[],
	ALTER COLUMN address SET NOT NULL,
	ALTER COLUMN mac_address SET NOT NULL,
	ALTER COLUMN host_name SET NOT NULL;