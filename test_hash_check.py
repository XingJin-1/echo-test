import bcrypt

passwd = b'Jx15528250227!'

salt = bcrypt.gensalt()
hashed = bcrypt.hashpw(passwd, salt)

if bcrypt.checkpw(passwd, hashed):
    print("match")
else:
    print("does not match")