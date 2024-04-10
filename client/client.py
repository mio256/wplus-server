import os
import base64
import json
import requests
import jwt

# 環境変数から基本URLを取得
base_url = os.environ['BASE_URL']
headers = {'Content-Type': 'application/json'}


def ping_server():
    """サーバーの生存確認を行う"""
    response = requests.get(f'{base_url}/ping/', headers=headers)
    if response.ok:
        print('ping: ', response.json())
    else:
        print('Error pinging server')


def ping_database():
    """データベースの生存確認を行う"""
    response = requests.get(f'{base_url}/db-ping/', headers=headers)
    if response.ok:
        print('db-ping: ', response.json())
    else:
        print('Error pinging database')


def login():
    """ログイン処理を行い、トークンとログインレスポンスを取得する"""
    data = {
        'office_id': int(input('office_id: ')),
        'user_id': int(input('user_id: ')),
        'password': input('password: ')
    }
    response = requests.post(f'{base_url}/login/', headers=headers, json=data)
    if response.ok:
        print('login: ', response.json(), response.cookies.get('token'))
        return response.cookies.get('token')
    else:
        print('Login failed')
        return None


def get_workplace(workplace_id):
    r = requests.get(base_url + f'/workplaces/{workplace_id}/', headers=headers)
    return r.json()


def get_employees(role, id):
    # role = 'admin', 'manager', 'employee'
    # id = office_id, workplace_id, employee_id
    if role == 'admin':
        r = requests.get(base_url + f'/employees/', headers=headers)
    elif role == 'manager':
        r = requests.get(base_url + f'/employees/workplace/{id}/', headers=headers)
    elif role == 'employee':
        r = requests.get(base_url + f'/employees/{id}/', headers=headers)
    if r.json() == None:
        print('No employees found')
        return
    return r.json()


def post_work_entry(employee_id, workplace_id):
    work_entry_data = {
        "employee_id": employee_id,
        "workplace_id": workplace_id
    }
    date = input('date (yyyy-mm-dd): ')
    work_entry_data['date'] = date + 'T00:00:00.000Z'  # 日付しか使わないので、時間は00:00:00にしておく

    # workplace_typeを取得
    """ CHECK:
    adminが勤怠登録画面に行く場合、workplace_idがNullなのでworkplace.typeが取得できず、勤怠登録画面を「時間数」「有無」「時間（star,end）」の切り替えを動的にしないといけませんね。
    例えば、adminが勤怠登録画面にいくと「職場選択」→「従業員選択」→入力フォーム のように画面が遷移するなど。
    今はstart,endだけなので、workplace_idがNullだろうと決め打ちでworkplace.type='time'でstart,endの画面固定でよいです。
    """
    workplace_type = get_workplace(workplace_id)['work_type']

    # workplace_typeによって、勤怠登録の方法を変える
    if workplace_type == 'hours':
        work_entry_data['hours'] = int(input('hours: '))
    elif workplace_type == 'time':
        start_time = input('start_time (hh:mm): ')
        end_time = input('end_time (hh:mm): ')
        work_entry_data['start_time'] = '1970-01-01T' + start_time + ':00.000Z'  # 時間しか使わないので、日付はUnixTimeの0秒の日付=1970-01-01にしておく
        work_entry_data['end_time'] = '1970-01-01T' + end_time + ':00.000Z'  # 時間しか使わないので、日付はUnixTimeの0秒の日付=1970-01-01にしておく
    elif workplace_type == 'attendance':
        work_entry_data['attendance'] = bool(input('attendance (y/n): ') == 'y')

    r = requests.post(base_url + '/work_entries', headers=headers, json=work_entry_data)
    return r.json()


def get_work_entries(role, id):
    # role = 'admin', 'manager', 'employee'
    # id = office_id, workplace_id, employee_id
    if role == 'admin':
        r = requests.get(base_url + f'/work_entries/', headers=headers)
    elif role == 'manager':
        r = requests.get(base_url + f'/work_entries/workplace/{id}/', headers=headers)
    elif role == 'employee':
        r = requests.get(base_url + f'/work_entries/employee/{id}/', headers=headers)
    return r.json()


def main():
    ping_server()
    ping_database()

    token = login()

    # tokenが取得できなかった場合、終了
    if not token:
        return

    # Tokenをヘッダーに追加
    headers['Authorization'] = 'Bearer ' + token

    # JWTTokenをParse
    tmp = token.split('.')
    header = json.loads(base64.b64decode(tmp[0]).decode())
    payload = jwt.decode(token, options={"verify_signature": False})
    print(header, payload)

    name = payload['name']
    office_id = payload['office_id']
    user_id = payload['user_id']
    role = payload['role']
    # role == 'admin'の場合、employee_idとworkplace_idが0
    workplace_id = payload['workplace_id']
    employee_id = payload['employee_id']

    # yesならば、勤怠登録画面
    # noならば、勤怠一覧画面
    post = bool(input('post work_entry? (y/n): ') == 'y')

    if role == 'admin':
        # adminはworkplace_idがNullなので注意 (office:workplace = 1:N)
        if post:
            # 管理者はOffice全体から対象のEmployeeを選択して勤怠登録

            # 従業員一覧
            employees = get_employees(role, office_id)
            for _, employee in enumerate(employees):
                print(employee)

            # 従業員のIDを選択
            employee_id = int(input('employee_id: '))
            # 指定させたemployeeからworkplace_idを取得する
            for _, employee in enumerate(employees):
                if employee['id'] == employee_id:
                    workplace_id = employee['workplace_id']

            print(post_work_entry(employee_id, workplace_id))
        else:
            # 管理者はOffice全体の勤怠一覧を表示
            print(get_work_entries(role, office_id))
    elif role == 'manager':
        if post:
            # マネージャーは自分のWorkplaceのEmployeeを選択して勤怠登録

            # 従業員一覧
            employees = get_employees(role, workplace_id)
            print(employees)

            # 従業員のIDを選択
            employee_id = int(input('employee_id: '))

            # managerは自分のworkplace_idを使う
            print(post_work_entry(employee_id, workplace_id))
        else:
            # マネージャーは自分のWorkplaceの勤怠一覧を表示
            print(get_work_entries(role, workplace_id))
    elif role == 'employee':
        if post:
            # 従業員は自分の勤怠登録
            # employee は 画面表示用の情報として取得、employee.nameなど
            employee = get_employees(role, employee_id)
            print(post_work_entry(employee_id, workplace_id))
        else:
            # 従業員は自分の勤怠一覧を表示
            print(get_work_entries(role, employee_id))


if __name__ == "__main__":
    main()
