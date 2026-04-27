ALTER TABLE hosts 
ALTER COLUMN address TYPE text[] 
USING ARRAY[address];

ALTER TABLE hosts 
ALTER COLUMN mac_address TYPE text[] 
USING ARRAY[mac_address];

ALTER TABLE hosts 
ALTER COLUMN host_name TYPE text[] 
USING ARRAY[host_name];