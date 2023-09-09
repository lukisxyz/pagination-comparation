-- Open a file for writing the SQL DML script
local file = assert(io.open("db/migrations/000002_seed_product_table.up.sql", "w"))

local ulid_characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
local ulid_length = 26

function generate_ulid()
    local ulid = ""
    
    -- Add the time-based component (10 characters)
    local time = tostring(os.time()):reverse()
    ulid = ulid .. time:sub(1, 10):reverse()

    -- Add the random component (16 characters)
    for i = 1, 16 do
        local random_index = math.random(1, ulid_length)
        ulid = ulid .. ulid_characters:sub(random_index, random_index)
    end

    return ulid
end

-- Initialize a starting timestamp (adjust as needed)
local start_timestamp = os.time({year = 2023, month = 1, day = 1, hour = 0, min = 0, sec = 0})
local current_timestamp = start_timestamp

file:write("BEGIN;\n")
for i = 1, 1000000 do
    local ulid = i
    local name = "Product " .. i
    local amount = math.random(1, 1000)
    local description = "Description for Product " .. i

    local query = string.format([[
        INSERT INTO products (id, name, amount, description, created_at)
        VALUES ('%s', '%s', %f, '%s', TIMESTAMP 'epoch' + INTERVAL '%d seconds');
    ]], ulid, name, amount, description, current_timestamp)

    file:write(query .. "\n")

    -- Increment the timestamp randomly (adjust the range as needed)
    current_timestamp = current_timestamp + math.random(0, 60)  -- Random increment between 1 minute and 1 hour
end
file:write("COMMIT;\n")
file:close()

print("SQL DML script generated to db/seeder/seed_data.sql")