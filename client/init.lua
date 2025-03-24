local socket = require("socket.socket")

local ip = "0.0.0.0"
local port = 5000
local player_id = 0
local connected = false

local function join(sock)
	local packet = Encode_packet(0, 0x00, "Player_name")
	sock:send(packet)
	local id = sock:receive()
	if id then
		player_id = id
		connected = true
		print("Connected")
		return true
	end
	return false
end

local function disconnect(sock)
	local packet = Encode_packet(player_id, 0x01, "")
	sock:send(packet)
	local res = sock:receive()
	print("Response: ", res)
end

local sock = Connect(ip, port)

if join(sock) then
	local packet = Encode_packet(player_id, 0x03, "Message")
	sock:send(packet)
	local res = sock:receive()
	print("Response: ", res)
end

disconnect(sock)
