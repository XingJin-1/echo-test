# Importing flask module in the project is mandatory
# An object of Flask class is our WSGI application.
from flask import Flask, request
import jwt
import cryptography

  
# Flask constructor takes the name of 
# current module (__name__) as argument.
app = Flask(__name__)
  
# The route() function of the Flask class is a decorator, 
# which tells the application which URL should call 
# the associated function.
@app.route('/')
# ‘/’ URL is bound with hello_world() function.
def hello_world():
    #encoded_jwt = request.headers.get('Authorization')
    #encoded_jwt = encoded_jwt.replace('Bearer ', "")

    #jwt.register_algorithm('RS256', RSAAlgorithm(RSAAlgorithm.SHA256))
    #secrete = "oKlbVzzdP7qxOCDDKmdzUJD1vNMwReR4p2KLLEMg8zRQKSmvQZcBT3naIdYW2rYV"
    #decoded_jwt = jwt.decode(encoded_jwt, secrete, algorithms=["RS256"])

    return 'hello from xing'
  
# main driver function
if __name__ == '__main__':
  
    # run() method of Flask class runs the application 
    # on the local development server.
    #app.run(port=4200)
    app.run(host="0.0.0.0", port=4200)