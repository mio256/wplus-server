import requests

# この2つは固定
base_url = 'http://localhost:8080'
headers = {'Content-Type': 'application/json'}

# ログインフォームから取得してくる
data = {
    'office_id': 1,
    'user_id': 1,
    'password': 'pass'
}

response = requests.post(base_url+'/login', headers=headers, json=data)
print(response.json())  # {'name': 'yamada', 'office_id': 1, 'role': 'admin', 'user_id': 1} → roleを使って画面遷移を制御する
token = response.cookies.get('token')  # JWTトークンを取得
print(token)

r = requests.get(base_url+'/offices', headers={'Authorization': 'Bearer '+token})  # トークンを使ってリクエストを送信
print(r.json())
