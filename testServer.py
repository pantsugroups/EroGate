from flask import Flask,request
from json import dumps,loads
app = Flask(__name__)

@app.route('/login',methods=["POST"])
def login():
    return dumps({
        "secret":"this is a secret",
        "u":{
        "ID":123,
        "username":"aaa"
        }
        })
@app.route("/website",methods=["POST"])
def website():
    return "this is a test page"
if __name__ == '__main__':
    app.run()