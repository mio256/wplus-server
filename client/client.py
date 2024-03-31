import requests

# この2つは固定
base_url = 'https://wplus-service-7gmfn2s4tq-dt.a.run.app/' # デプロイ先のURL
headers = {'Content-Type': 'application/json'}

# serverの生存確認
response = requests.get(base_url+'/ping', headers=headers)
print(response.json())  # {'message': 'pong'}

# dbの生存確認
response = requests.get(base_url+'/db-ping', headers=headers)
print(response.json())  # {'message': 'db-pong'}

# ログインフォームから取得してくる
data = {
    'office_id': int(input('office_id: ')),
    'user_id': int(input('user_id: ')),
    'password': input('password: ')
}

response = requests.post(base_url+'/login', headers=headers, json=data)
print(response.json())  # {'name': 'yamada', 'office_id': 1, 'role': 'admin', 'user_id': 1} → roleを使って画面遷移を制御する
token = response.cookies.get('token')  # JWTトークンを取得
print(token)

r = requests.get(base_url+'/offices', headers={'Authorization': 'Bearer '+token})  # トークンを使ってリクエストを送信
print(r.json())
