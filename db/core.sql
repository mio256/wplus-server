-- 事業所テーブル
create table offices (
    id bigserial primary key,
    name varchar(255) not null,
    deleted_at timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- 勤務種類
create type work_type as enum ('hours', 'time', 'attendance');

-- 職場テーブル
create table workplaces (
    id bigserial primary key,
    name varchar(255) not null,
    office_id bigint not null,
    work_type work_type not null,
    deleted_at timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- 従業員テーブル
create table employees (
    id bigserial primary key,
    name varchar(255) not null,
    workplace_id bigint not null,
    deleted_at timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- 勤務テーブル
create table work_entries (
    id bigserial primary key,
    employee_id bigint not null,
    workplace_id bigint not null,
    date date not null,
    hours smallint,
    start_time time,
    end_time time,
    attendance boolean,
    constraint chk_work_entries_check check (
        (hours is not null and start_time is null and end_time is null and attendance is null) or
        (hours is null and start_time is not null and end_time is not null and attendance is null) or
        (hours is null and start_time is null and end_time is null and attendance is not null)
    ),
    comment varchar(255),
    deleted_at timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- 外部キー制約
alter table workplaces add constraint fk_workplaces_offices foreign key (office_id) references offices(id);
alter table employees add constraint fk_employees_workplaces foreign key (workplace_id) references workplaces(id);
alter table work_entries add constraint fk_work_hours_entries_employees foreign key (employee_id) references employees(id);
alter table work_entries add constraint fk_work_hours_entries_workplaces foreign key (workplace_id) references workplaces(id);
