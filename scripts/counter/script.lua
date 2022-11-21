-- The following code should be treated as a pure function like the following:
-- script(key, field string, score, maxSize uint64)
local key = KEYS[1]
local field = ARGV[1]
local score = tonumber(ARGV[2])
local maxSize = tonumber(ARGV[3])

local addKey = key .. 'ADDSUFFIX'
local remKey = key .. 'REMSUFFIX'

-- Make sure that we remain capped to the max size.
local cardinality = tonumber(redis.call('ZCARD', addKey))
if cardinality >= maxSize then
    return -1
end

-- Check to see if an item with in the collection is invalid.
local invalid = function(key, score)
    local valueScore = redis.call('ZSCORE', key, field)
    if valueScore and score <= tonumber(valueScore) then
        return true
    end
    return false
end

-- Check if the key is already invalid and eject out to create a noop.
if invalid(key .. 'INSERTSUFFIX', score) or invalid(key .. 'DELETESUFFIX', score) then
    return -1
end

-- Insert the item after removing any possible trace of the old item.
redis.call('ZREM', remKey, field)
return redis.call('ZADD', addKey, score, field)
