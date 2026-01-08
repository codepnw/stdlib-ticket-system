ALTER TABLE bookings 
ALTER COLUMN id SET DEFAULT gen_random_uuid();