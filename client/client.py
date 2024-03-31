import requests
import json
from datetime import datetime

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
login_res = response.json()

r = requests.get(base_url+'/offices', headers={'Authorization': 'Bearer '+token})  # トークンを使ってリクエストを送信
print(r.json())

# Assuming login_res and token are defined
work_entry_data = {
    "employee_id": login_res['employee_id'],
    "workplace_id": login_res['workplace_id'],
    "date": "2006-01-02T00:00:00.000000Z", # 日付しか使わないので、時間は00:00:00にしておく
    "start_time": "1970-01-01T08:00:00.000000Z", # 時間しか使わないので、日付は1970-01-01にしておく (UNXITime=0からの経過時間
    "end_time": "1970-01-01T17:00:00.000000Z", # 時間しか使わないので、日付は1970-01-01にしておく
}

print(work_entry_data)

# Assuming base_url is defined
r = requests.post(base_url + '/work_entries', headers={'Authorization': 'Bearer ' + token}, json=work_entry_data)
print(r.status_code)
try:
    print(r.json())
except json.JSONDecodeError:
    print("Response is not in JSON format:")
    print(r.text)
