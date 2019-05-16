from flask import Flask,request,make_response
from json import dumps,loads
from requests import get,post
app = Flask(__name__)

@app.route('/login2',methods=["POST"])
def login():
    res = post("http://127.0.0.1/login", json ={
        "secret":"this is a secret",
        "u":{
        "ID":123,
        "username":"aaa"
        }
        })
    response = make_response(res.text)
    response.set_cookie("test","yes")
    response.set_cookie("tes2","no")
    return response
@app.route("/website",methods=["POST"])
def website():
    return "this is a test page"
if __name__ == '__main__':
    app.run()