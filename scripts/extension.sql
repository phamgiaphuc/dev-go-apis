CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION uuidv7(timestamptz DEFAULT clock_timestamp())
RETURNS uuid
AS $$
  SELECT encode(
    set_bit(
      set_bit(
        overlay(
          uuid_send(gen_random_uuid())
          placing substring(
            int8send((extract(epoch from $1) * 1000)::bigint) from 3
          )
          from 1 for 6
        ),
        52, 1
      ),
      53, 1
    ),
    'hex'
  )::uuid;
$$ LANGUAGE sql VOLATILE PARALLEL SAFE;

COMMENT ON FUNCTION uuidv7(timestamptz) IS
'Generate a UUIDv7 (48-bit millisecond timestamp + 74 bits randomness)';

CREATE FUNCTION uuidv7_extract_timestamp(uuid) RETURNS timestamptz
AS $$
  SELECT to_timestamp(
    right(substring(uuid_send($1) from 1 for 6)::text, -1)::bit(48)::int8 -- milliseconds
      /1000.0);
$$ LANGUAGE sql immutable strict parallel safe4fr55                   

COMMENT ON FUNCTION uuidv7_extract_timestamp(uuid) IS
'Return the timestamp stored in the first 48 bits of the UUID v7 value';