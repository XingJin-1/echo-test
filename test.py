import base64

message = "XingJin:Jx15528250227!"
message_bytes = message.encode('ascii')
base64_bytes = base64.b64encode(message_bytes)
base64_message = base64_bytes.decode('ascii')

print("message_bytes: ", message_bytes)

print("base64_bytes: ", base64_bytes)

print("base64_message: ", base64_message)