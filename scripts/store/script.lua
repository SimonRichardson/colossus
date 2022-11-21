-- The following code should be treated as a pure function like the following:
-- script(key, field string, score float64, txn, data string)
local key = KEYS[1]
local field = ARGV[1]
local score = tonumber(ARGV[2])
local txn = ARGV[3]
local data = ARGV[4]

local extract = function(value, start)
    local index = string.find(value, 'SEPARATOR', start, true)
    if index and index > start then
        local value = string.sub(value, start, index - 1)
        if value then
            return true, index, value
        end
    end
    return false, index, ''
end

local invalid = function(value, score, txn)
    local ok, scoreOffset, valueScore = extract(value, 1)
    if not ok then
        return true
    end
    if score <= tonumber(valueScore) then
        return true
    end

    local ok, _, valueTxn = extract(value, scoreOffset + 1)
    if not ok then
        return true
    end
    if txn ~= valueTxn then
        return true
    end

    return false
end

-- Check if the score associated with a key is greater than the one already
-- existing in the store for insertions
-- Note: last write wins
local insertion = redis.call('HGET', key .. 'INSERTSUFFIX', field)
if insertion and invalid(insertion, score, txn) then
    return -1
end

-- Check if the score associated with a key is greater than the one already
-- existing in the store for deletions
-- Note: last write wins
local deletion = redis.call('HGET', key .. 'DELETESUFFIX', field)
if deletion and invalid(deletion, score, txn) then
    return -1
end

-- Remove the existing key if it's got a REMSUFFIX.
redis.call('HDEL', key .. 'REMSUFFIX', field)

-- Add the key to the store
return redis.call('HSET', key .. 'ADDSUFFIX', field, data)
