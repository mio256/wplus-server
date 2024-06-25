-- name: OutputWorkEntriesByWorkplaceAndDate :many
select employees.name as employee_name, work_entries.*
from work_entries
    join employees on work_entries.employee_id = employees.id
    join workplaces on work_entries.workplace_id = workplaces.id
where workplaces.id = $1
    and work_entries.date >= @min_date
    and work_entries.date <= @max_date
    and work_entries.deleted_at is null
order by employee_name;
