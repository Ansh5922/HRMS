CREATE OR REPLACE VIEW v_active_employees AS
SELECT
    e.id, e.emp_code, e.first_name, e.last_name,
    d.name   AS department,
    des.title AS designation,
    m.first_name || ' ' || m.last_name AS manager_name,
    e.joining_date, e.employment_type, e.work_location, e.org_id
FROM employees e
LEFT JOIN departments  d   ON d.id = e.dept_id
LEFT JOIN designations des ON des.id = e.designation_id
LEFT JOIN employees    m   ON m.id = e.manager_id
WHERE e.status = 'active';

CREATE OR REPLACE VIEW v_monthly_attendance AS
SELECT
    emp_id,
    DATE_TRUNC('month', date) AS month,
    COUNT(*) FILTER (WHERE status = 'present')  AS present_days,
    COUNT(*) FILTER (WHERE status = 'absent')   AS absent_days,
    COUNT(*) FILTER (WHERE status = 'half_day') AS half_days,
    COUNT(*) FILTER (WHERE status = 'wfh')      AS wfh_days,
    COUNT(*) FILTER (WHERE status = 'late')     AS late_days,
    SUM(overtime_mins)                          AS total_overtime_mins
FROM attendance_records
GROUP BY emp_id, DATE_TRUNC('month', date);

CREATE OR REPLACE VIEW v_leave_summary AS
SELECT
    lb.emp_id,
    lt.name AS leave_type,
    lb.year,
    lb.total, lb.used, lb.pending, lb.available
FROM leave_balances lb
JOIN leave_types lt ON lt.id = lb.leave_type_id;
