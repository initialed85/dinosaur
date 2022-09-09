local socket = require("socket")

PORT = 13337

function sleep(seconds)
    socket.select(nil, nil, seconds)
end

hostname = os.getenv("HOSTNAME")
local_ip = os.getenv("LOCAL_IP")
broadcast_ip = os.getenv("BROADCAST_IP")

local sock = assert(socket.udp4())

assert(sock:setsockname("0.0.0.0", PORT))
assert(sock:setoption('broadcast', true))
assert(sock:settimeout(0))

function receive_callback(data, ip, port)
    if ip == local_ip and port == PORT then
        return
    end

    print(string.format("%s:%d\t%s", ip, port, data))
end

receive_loop = coroutine.create(function()
    while true do
        local data, status_or_ip, port = sock:receivefrom(65507)

        if status_or_ip == "timeout" then
            coroutine.yield()
        else
            receive_callback(data, status_or_ip, port)
        end
    end
end)

send_loop = coroutine.create(function()
    while true do
        assert(sock:sendto(
                string.format("Hello world from Lua @ %s", hostname),
                broadcast_ip,
                PORT
        ))
        coroutine.yield()
    end
end)

while true do
    coroutine.resume(receive_loop)
    coroutine.resume(send_loop)
    sleep(1)
end
