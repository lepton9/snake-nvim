local socket = require("socket")

function Connect(ip, port)
	local udp = socket.udp()
	udp:settimeout(1)
	udp:setpeername(ip, port)
	return udp
end

function Close(sock)
	sock:close()
end

-- Little Endian
local function uint32_to_bytes(num)
	local b1 = num % 256
	num = (num - b1) / 256
	local b2 = num % 256
	num = (num - b2) / 256
	local b3 = num % 256
	num = (num - b3) / 256
	local b4 = num % 256
	return string.char(b1, b2, b3, b4)
end

-- Encode binary packet
function Encode_packet(player_id, packet_type, data)
	return uint32_to_bytes(player_id) .. string.char(packet_type) .. data
end
