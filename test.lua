local keys = redis.call("ZRANGEBYSCORE", "Comment/user", 1, 1)
local tmpKey = "Comment/user=1"
for i, k in ipairs(keys) do
  redis.call("ZADD", tmpKey, 1, k)
end
redis.call("EXPIRE", tmpKey, 600)
return tmpKey
