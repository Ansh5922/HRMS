import type {
    Employee, Department, Designation, AttendanceRecord, LeaveApplication,
    LeaveBalance, PayrollRun, Payslip, Reimbursement, JobPosting, Candidate,
    ReviewCycle, Goal, Course, Notification, Announcement, Role, WorkflowTemplate,
    Organization, User,
} from '../types';

// ── Organization ──
export const mockOrganization: Organization = {
    id: '550e8400-e29b-41d4-a716-446655440001',
    name: 'TechNova Solutions Pvt. Ltd.',
    logo_url: null,
    address: '4th Floor, Tower B, Cyber Hub, Gurugram, HR 122002',
    timezone: 'Asia/Kolkata',
};

// ── Users (different roles) ──
export const mockUsers: User[] = [
    { id: 'u-001', org_id: mockOrganization.id, email: 'admin@technova.com', role_id: 'r-001', is_active: true, is_verified: true, mfa_enabled: false, last_login_at: '2026-08-31T10:00:00Z', created_at: '2026-01-15T00:00:00Z', role: { id: 'r-001', org_id: mockOrganization.id, name: 'Super Admin', description: 'Full access', is_system: true } },
    { id: 'u-002', org_id: mockOrganization.id, email: 'hr@technova.com', role_id: 'r-002', is_active: true, is_verified: true, mfa_enabled: false, last_login_at: '2026-08-31T09:30:00Z', created_at: '2026-02-01T00:00:00Z', role: { id: 'r-002', org_id: mockOrganization.id, name: 'HR Manager', description: 'HR operations', is_system: true } },
    { id: 'u-003', org_id: mockOrganization.id, email: 'manager@technova.com', role_id: 'r-003', is_active: true, is_verified: true, mfa_enabled: false, last_login_at: '2026-08-30T14:00:00Z', created_at: '2026-02-15T00:00:00Z', role: { id: 'r-003', org_id: mockOrganization.id, name: 'Manager', description: 'Team management', is_system: false } },
];

// ── Roles ──
export const mockRoles: Role[] = [
    { id: 'r-001', org_id: mockOrganization.id, name: 'Super Admin', description: 'Full system access with all permissions', is_system: true },
    { id: 'r-002', org_id: mockOrganization.id, name: 'HR Manager', description: 'Manage employees, leaves, payroll, and recruitment', is_system: true },
    { id: 'r-003', org_id: mockOrganization.id, name: 'Manager', description: 'Approve leaves, view team attendance and performance', is_system: false },
    { id: 'r-004', org_id: mockOrganization.id, name: 'Employee', description: 'Self-service: apply leave, view payslips, mark attendance', is_system: true },
    { id: 'r-005', org_id: mockOrganization.id, name: 'Finance', description: 'Payroll processing, reimbursements, tax declarations', is_system: false },
];

// ── Departments ──
export const mockDepartments: Department[] = [
    { id: 'd-001', org_id: mockOrganization.id, name: 'Engineering', code: 'ENG', parent_dept_id: null, head_emp_id: 'e-003', description: 'Software development and architecture', is_active: true },
    { id: 'd-002', org_id: mockOrganization.id, name: 'Human Resources', code: 'HR', parent_dept_id: null, head_emp_id: 'e-002', description: 'People operations and talent management', is_active: true },
    { id: 'd-003', org_id: mockOrganization.id, name: 'Product', code: 'PROD', parent_dept_id: null, head_emp_id: 'e-005', description: 'Product strategy and management', is_active: true },
    { id: 'd-004', org_id: mockOrganization.id, name: 'Design', code: 'DES', parent_dept_id: null, head_emp_id: null, description: 'UI/UX and visual design', is_active: true },
    { id: 'd-005', org_id: mockOrganization.id, name: 'Finance', code: 'FIN', parent_dept_id: null, head_emp_id: null, description: 'Accounting, payroll, and financial planning', is_active: true },
    { id: 'd-006', org_id: mockOrganization.id, name: 'QA', code: 'QA', parent_dept_id: 'd-001', head_emp_id: null, description: 'Quality assurance and testing', is_active: true },
    { id: 'd-007', org_id: mockOrganization.id, name: 'DevOps', code: 'DEVOPS', parent_dept_id: 'd-001', head_emp_id: null, description: 'Infrastructure and CI/CD', is_active: false },
];

// ── Designations ──
export const mockDesignations: Designation[] = [
    { id: 'des-001', org_id: mockOrganization.id, department_id: null, title: 'CEO', grade: 'C-Suite', level: 1 },
    { id: 'des-002', org_id: mockOrganization.id, department_id: 'd-001', title: 'VP Engineering', grade: 'VP', level: 2 },
    { id: 'des-003', org_id: mockOrganization.id, department_id: 'd-001', title: 'Senior Software Engineer', grade: 'IC3', level: 4 },
    { id: 'des-004', org_id: mockOrganization.id, department_id: 'd-001', title: 'Software Engineer', grade: 'IC2', level: 5 },
    { id: 'des-005', org_id: mockOrganization.id, department_id: 'd-002', title: 'HR Manager', grade: 'M1', level: 3 },
    { id: 'des-006', org_id: mockOrganization.id, department_id: 'd-003', title: 'Product Manager', grade: 'M1', level: 3 },
    { id: 'des-007', org_id: mockOrganization.id, department_id: 'd-004', title: 'UI/UX Designer', grade: 'IC2', level: 5 },
    { id: 'des-008', org_id: mockOrganization.id, department_id: 'd-001', title: 'Junior Developer', grade: 'IC1', level: 6 },
    { id: 'des-009', org_id: mockOrganization.id, department_id: 'd-005', title: 'Finance Analyst', grade: 'IC2', level: 5 },
];

// ── Employees (12 employees, various roles/departments) ──
export const mockEmployees: Employee[] = [
    { id: 'e-001', org_id: mockOrganization.id, user_id: 'u-001', emp_code: 'EMPTN0001', first_name: 'Rajesh', last_name: 'Sharma', middle_name: null, dob: '1985-03-15', gender: 'Male', phone: '+91-98765-43210', personal_email: 'rajesh.s@gmail.com', address: '12, Sector 45', city: 'Gurugram', state: 'Haryana', country: 'India', pincode: '122003', dept_id: 'd-001', designation_id: 'des-001', manager_id: null, employment_type: 'full_time', work_location: 'Gurugram HQ', joining_date: '2020-01-15', status: 'active', photo_url: null, department: 'Engineering', designation: 'CEO', manager_name: undefined },
    { id: 'e-002', org_id: mockOrganization.id, user_id: 'u-002', emp_code: 'EMPTN0002', first_name: 'Priya', last_name: 'Patel', middle_name: null, dob: '1990-07-22', gender: 'Female', phone: '+91-98765-43211', personal_email: 'priya.p@gmail.com', address: '34, DLF Phase 3', city: 'Gurugram', state: 'Haryana', country: 'India', pincode: '122002', dept_id: 'd-002', designation_id: 'des-005', manager_id: 'e-001', employment_type: 'full_time', work_location: 'Gurugram HQ', joining_date: '2021-04-01', status: 'active', photo_url: null, department: 'Human Resources', designation: 'HR Manager', manager_name: 'Rajesh Sharma' },
    { id: 'e-003', org_id: mockOrganization.id, user_id: 'u-003', emp_code: 'EMPTN0003', first_name: 'Vikram', last_name: 'Singh', middle_name: 'Kumar', dob: '1988-11-05', gender: 'Male', phone: '+91-98765-43212', personal_email: 'vikram.s@gmail.com', address: '56, Cyber City', city: 'Gurugram', state: 'Haryana', country: 'India', pincode: '122001', dept_id: 'd-001', designation_id: 'des-002', manager_id: 'e-001', employment_type: 'full_time', work_location: 'Gurugram HQ', joining_date: '2020-06-15', status: 'active', photo_url: null, department: 'Engineering', designation: 'VP Engineering', manager_name: 'Rajesh Sharma' },
    { id: 'e-004', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0004', first_name: 'Ananya', last_name: 'Gupta', middle_name: null, dob: '1993-02-18', gender: 'Female', phone: '+91-98765-43213', personal_email: 'ananya.g@gmail.com', address: '78, MG Road', city: 'Bengaluru', state: 'Karnataka', country: 'India', pincode: '560001', dept_id: 'd-001', designation_id: 'des-003', manager_id: 'e-003', employment_type: 'full_time', work_location: 'Remote', joining_date: '2022-01-10', status: 'active', photo_url: null, department: 'Engineering', designation: 'Senior Software Engineer', manager_name: 'Vikram Singh' },
    { id: 'e-005', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0005', first_name: 'Arjun', last_name: 'Mehta', middle_name: null, dob: '1991-09-30', gender: 'Male', phone: '+91-98765-43214', personal_email: 'arjun.m@gmail.com', address: '90, Connaught Place', city: 'New Delhi', state: 'Delhi', country: 'India', pincode: '110001', dept_id: 'd-003', designation_id: 'des-006', manager_id: 'e-001', employment_type: 'full_time', work_location: 'Gurugram HQ', joining_date: '2021-08-01', status: 'active', photo_url: null, department: 'Product', designation: 'Product Manager', manager_name: 'Rajesh Sharma' },
    { id: 'e-006', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0006', first_name: 'Neha', last_name: 'Reddy', middle_name: null, dob: '1995-12-10', gender: 'Female', phone: '+91-98765-43215', personal_email: 'neha.r@gmail.com', address: '23, Banjara Hills', city: 'Hyderabad', state: 'Telangana', country: 'India', pincode: '500034', dept_id: 'd-004', designation_id: 'des-007', manager_id: 'e-005', employment_type: 'full_time', work_location: 'Remote', joining_date: '2023-03-20', status: 'active', photo_url: null, department: 'Design', designation: 'UI/UX Designer', manager_name: 'Arjun Mehta' },
    { id: 'e-007', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0007', first_name: 'Rohit', last_name: 'Joshi', middle_name: null, dob: '1997-06-25', gender: 'Male', phone: '+91-98765-43216', personal_email: 'rohit.j@gmail.com', address: '45, Aundh', city: 'Pune', state: 'Maharashtra', country: 'India', pincode: '411007', dept_id: 'd-001', designation_id: 'des-004', manager_id: 'e-003', employment_type: 'full_time', work_location: 'Remote', joining_date: '2023-07-01', status: 'active', photo_url: null, department: 'Engineering', designation: 'Software Engineer', manager_name: 'Vikram Singh' },
    { id: 'e-008', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0008', first_name: 'Kavya', last_name: 'Nair', middle_name: null, dob: '1999-01-12', gender: 'Female', phone: '+91-98765-43217', personal_email: 'kavya.n@gmail.com', address: '67, Indiranagar', city: 'Bengaluru', state: 'Karnataka', country: 'India', pincode: '560038', dept_id: 'd-001', designation_id: 'des-008', manager_id: 'e-004', employment_type: 'intern', work_location: 'Gurugram HQ', joining_date: '2026-06-01', status: 'probation', photo_url: null, department: 'Engineering', designation: 'Junior Developer', manager_name: 'Ananya Gupta' },
    { id: 'e-009', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0009', first_name: 'Amit', last_name: 'Verma', middle_name: null, dob: '1992-04-08', gender: 'Male', phone: '+91-98765-43218', personal_email: 'amit.v@gmail.com', address: '89, Salt Lake', city: 'Kolkata', state: 'West Bengal', country: 'India', pincode: '700091', dept_id: 'd-005', designation_id: 'des-009', manager_id: 'e-001', employment_type: 'full_time', work_location: 'Remote', joining_date: '2022-09-15', status: 'active', photo_url: null, department: 'Finance', designation: 'Finance Analyst', manager_name: 'Rajesh Sharma' },
    { id: 'e-010', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0010', first_name: 'Sneha', last_name: 'Iyer', middle_name: null, dob: '1994-08-20', gender: 'Female', phone: '+91-98765-43219', personal_email: 'sneha.i@gmail.com', address: '12, T Nagar', city: 'Chennai', state: 'Tamil Nadu', country: 'India', pincode: '600017', dept_id: 'd-001', designation_id: 'des-004', manager_id: 'e-003', employment_type: 'contract', work_location: 'Remote', joining_date: '2024-01-10', status: 'notice_period', photo_url: null, department: 'Engineering', designation: 'Software Engineer', manager_name: 'Vikram Singh' },
    { id: 'e-011', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0011', first_name: 'Deepak', last_name: 'Rao', middle_name: null, dob: '1996-05-14', gender: 'Male', phone: '+91-98765-43220', personal_email: 'deepak.r@gmail.com', address: '34, Koramangala', city: 'Bengaluru', state: 'Karnataka', country: 'India', pincode: '560034', dept_id: 'd-006', designation_id: 'des-004', manager_id: 'e-003', employment_type: 'full_time', work_location: 'Gurugram HQ', joining_date: '2023-11-01', status: 'active', photo_url: null, department: 'QA', designation: 'Software Engineer', manager_name: 'Vikram Singh' },
    { id: 'e-012', org_id: mockOrganization.id, user_id: null, emp_code: 'EMPTN0012', first_name: 'Meera', last_name: 'Kapoor', middle_name: null, dob: '1998-10-03', gender: 'Female', phone: '+91-98765-43221', personal_email: 'meera.k@gmail.com', address: '56, Whitefield', city: 'Bengaluru', state: 'Karnataka', country: 'India', pincode: '560066', dept_id: 'd-001', designation_id: 'des-004', manager_id: 'e-004', employment_type: 'full_time', work_location: 'Remote', joining_date: '2025-02-01', status: 'active', photo_url: null, department: 'Engineering', designation: 'Software Engineer', manager_name: 'Ananya Gupta' },
];

// ── Attendance (today's records) ──
export const mockAttendance: AttendanceRecord[] = [
    { id: 'a-001', emp_id: 'e-001', date: '2026-08-31', check_in: '2026-08-31T09:02:00+05:30', check_out: '2026-08-31T18:15:00+05:30', method: 'face', status: 'present', duration_mins: 553, overtime_mins: 15 },
    { id: 'a-002', emp_id: 'e-002', date: '2026-08-31', check_in: '2026-08-31T09:00:00+05:30', check_out: '2026-08-31T18:00:00+05:30', method: 'face', status: 'present', duration_mins: 540, overtime_mins: 0 },
    { id: 'a-003', emp_id: 'e-003', date: '2026-08-31', check_in: '2026-08-31T09:25:00+05:30', check_out: null, method: 'geo', status: 'late', duration_mins: null, overtime_mins: 0 },
    { id: 'a-004', emp_id: 'e-004', date: '2026-08-31', check_in: '2026-08-31T09:10:00+05:30', check_out: '2026-08-31T13:00:00+05:30', method: 'manual', status: 'half_day', duration_mins: 230, overtime_mins: 0 },
    { id: 'a-005', emp_id: 'e-005', date: '2026-08-31', check_in: null, check_out: null, method: 'manual', status: 'on_leave', duration_mins: null, overtime_mins: 0 },
    { id: 'a-006', emp_id: 'e-006', date: '2026-08-31', check_in: '2026-08-31T10:00:00+05:30', check_out: null, method: 'geo', status: 'wfh', duration_mins: null, overtime_mins: 0 },
    { id: 'a-007', emp_id: 'e-007', date: '2026-08-31', check_in: '2026-08-31T08:55:00+05:30', check_out: '2026-08-31T18:30:00+05:30', method: 'face', status: 'present', duration_mins: 575, overtime_mins: 30 },
    { id: 'a-008', emp_id: 'e-008', date: '2026-08-31', check_in: null, check_out: null, method: 'manual', status: 'absent', duration_mins: null, overtime_mins: 0 },
    { id: 'a-009', emp_id: 'e-009', date: '2026-08-31', check_in: '2026-08-31T09:05:00+05:30', check_out: '2026-08-31T18:00:00+05:30', method: 'biometric', status: 'present', duration_mins: 535, overtime_mins: 0 },
    { id: 'a-010', emp_id: 'e-010', date: '2026-08-31', check_in: '2026-08-31T09:00:00+05:30', check_out: '2026-08-31T18:00:00+05:30', method: 'geo', status: 'present', duration_mins: 540, overtime_mins: 0 },
    { id: 'a-011', emp_id: 'e-011', date: '2026-08-31', check_in: '2026-08-31T09:30:00+05:30', check_out: null, method: 'face', status: 'late', duration_mins: null, overtime_mins: 0 },
    { id: 'a-012', emp_id: 'e-012', date: '2026-08-31', check_in: '2026-08-31T09:00:00+05:30', check_out: '2026-08-31T18:10:00+05:30', method: 'face', status: 'present', duration_mins: 550, overtime_mins: 10 },
];

// ── Leave Applications ──
export const mockLeaveApplications: LeaveApplication[] = [
    { id: 'la-001', emp_id: 'e-005', leave_type_id: 'lt-001', from_date: '2026-08-31', to_date: '2026-09-02', days: 2, session: 'full', reason: 'Family function in hometown', status: 'approved', created_at: '2026-08-25T10:00:00Z', leave_type: 'Casual Leave', employee_name: 'Arjun Mehta' },
    { id: 'la-002', emp_id: 'e-007', leave_type_id: 'lt-002', from_date: '2026-09-05', to_date: '2026-09-08', days: 3, session: 'full', reason: 'Fever and cold', status: 'pending', created_at: '2026-08-30T14:00:00Z', leave_type: 'Sick Leave', employee_name: 'Rohit Joshi' },
    { id: 'la-003', emp_id: 'e-004', leave_type_id: 'lt-001', from_date: '2026-09-15', to_date: '2026-09-15', days: 0.5, session: 'first_half', reason: 'Doctor appointment', status: 'pending', created_at: '2026-08-31T09:00:00Z', leave_type: 'Casual Leave', employee_name: 'Ananya Gupta' },
    { id: 'la-004', emp_id: 'e-006', leave_type_id: 'lt-003', from_date: '2026-09-10', to_date: '2026-09-12', days: 3, session: 'full', reason: 'Personal work', status: 'rejected', created_at: '2026-08-28T11:00:00Z', leave_type: 'Earned Leave', employee_name: 'Neha Reddy' },
    { id: 'la-005', emp_id: 'e-012', leave_type_id: 'lt-001', from_date: '2026-09-20', to_date: '2026-09-22', days: 2, session: 'full', reason: 'Wedding at home', status: 'pending', created_at: '2026-08-31T08:00:00Z', leave_type: 'Casual Leave', employee_name: 'Meera Kapoor' },
    { id: 'la-006', emp_id: 'e-011', leave_type_id: 'lt-002', from_date: '2026-08-28', to_date: '2026-08-29', days: 2, session: 'full', reason: 'Food poisoning', status: 'approved', created_at: '2026-08-28T07:00:00Z', leave_type: 'Sick Leave', employee_name: 'Deepak Rao' },
];

// ── Leave Balances ──
export const mockLeaveBalances: LeaveBalance[] = [
    { id: 'lb-001', emp_id: 'e-004', leave_type_id: 'lt-001', year: 2026, total: 12, used: 4, pending: 0.5, available: 7.5, leave_type: 'Casual Leave' },
    { id: 'lb-002', emp_id: 'e-004', leave_type_id: 'lt-002', year: 2026, total: 8, used: 2, pending: 0, available: 6, leave_type: 'Sick Leave' },
    { id: 'lb-003', emp_id: 'e-004', leave_type_id: 'lt-003', year: 2026, total: 15, used: 5, pending: 0, available: 10, leave_type: 'Earned Leave' },
];

// ── Payroll Runs ──
export const mockPayrollRuns: PayrollRun[] = [
    { id: 'pr-001', org_id: mockOrganization.id, month: 8, year: 2026, status: 'disbursed', total_gross: 2450000, total_deductions: 612500, total_net: 1837500 },
    { id: 'pr-002', org_id: mockOrganization.id, month: 7, year: 2026, status: 'disbursed', total_gross: 2380000, total_deductions: 595000, total_net: 1785000 },
    { id: 'pr-003', org_id: mockOrganization.id, month: 6, year: 2026, status: 'disbursed', total_gross: 2380000, total_deductions: 595000, total_net: 1785000 },
    { id: 'pr-004', org_id: mockOrganization.id, month: 9, year: 2026, status: 'draft', total_gross: null, total_deductions: null, total_net: null },
];

// ── Payslips (Aug 2026) ──
export const mockPayslips: Payslip[] = [
    { id: 'ps-001', emp_id: 'e-001', payroll_run_id: 'pr-001', gross: 450000, total_deductions: 112500, net_pay: 337500, working_days: 22, lop_days: 0, breakdown: { basic: 180000, hra: 90000, special: 90000, conveyance: 19200, medical: 15000, pf: 21600, professional_tax: 2500, income_tax: 88400 } },
    { id: 'ps-002', emp_id: 'e-002', payroll_run_id: 'pr-001', gross: 180000, total_deductions: 45000, net_pay: 135000, working_days: 22, lop_days: 0, breakdown: { basic: 72000, hra: 36000, special: 36000, pf: 8640, income_tax: 36360 } },
    { id: 'ps-003', emp_id: 'e-003', payroll_run_id: 'pr-001', gross: 320000, total_deductions: 80000, net_pay: 240000, working_days: 22, lop_days: 0, breakdown: { basic: 128000, hra: 64000, special: 64000, pf: 15360, income_tax: 64640 } },
    { id: 'ps-004', emp_id: 'e-004', payroll_run_id: 'pr-001', gross: 200000, total_deductions: 50000, net_pay: 150000, working_days: 22, lop_days: 0, breakdown: { basic: 80000, hra: 40000, special: 40000, pf: 9600, income_tax: 40400 } },
    { id: 'ps-005', emp_id: 'e-007', payroll_run_id: 'pr-001', gross: 120000, total_deductions: 30000, net_pay: 90000, working_days: 22, lop_days: 0, breakdown: { basic: 48000, hra: 24000, special: 24000, pf: 5760, income_tax: 24240 } },
];

// ── Reimbursements ──
export const mockReimbursements: Reimbursement[] = [
    { id: 'rb-001', emp_id: 'e-003', category: 'travel', amount: 12500, description: 'Client visit to Mumbai - flight + cab', status: 'approved', created_at: '2026-08-20T10:00:00Z' },
    { id: 'rb-002', emp_id: 'e-005', category: 'food', amount: 3200, description: 'Team lunch for product sprint completion', status: 'pending', created_at: '2026-08-28T14:00:00Z' },
    { id: 'rb-003', emp_id: 'e-004', category: 'medical', amount: 8500, description: 'Eye checkup and prescription glasses', status: 'pending', created_at: '2026-08-30T09:00:00Z' },
    { id: 'rb-004', emp_id: 'e-007', category: 'travel', amount: 1800, description: 'Cab to office for weekend deployment', status: 'rejected', created_at: '2026-08-15T11:00:00Z' },
    { id: 'rb-005', emp_id: 'e-009', category: 'medical', amount: 15000, description: 'Annual health checkup', status: 'paid', created_at: '2026-08-10T10:00:00Z' },
];

// ── Job Postings ──
export const mockJobPostings: JobPosting[] = [
    { id: 'jp-001', title: 'Senior Backend Developer (Go)', description: 'Build microservices for HRMS platform', type: 'full_time', openings: 2, location: 'Gurugram', is_remote: true, status: 'open', created_at: '2026-08-01T00:00:00Z' },
    { id: 'jp-002', title: 'React Frontend Developer', description: 'Build responsive dashboards and admin panels', type: 'full_time', openings: 1, location: 'Remote', is_remote: true, status: 'open', created_at: '2026-08-10T00:00:00Z' },
    { id: 'jp-003', title: 'DevOps Engineer', description: 'Manage K8s clusters and CI/CD pipelines', type: 'full_time', openings: 1, location: 'Bengaluru', is_remote: false, status: 'paused', created_at: '2026-07-15T00:00:00Z' },
    { id: 'jp-004', title: 'Product Design Intern', description: 'Assist in designing user flows for HRMS modules', type: 'internship', openings: 2, location: 'Gurugram', is_remote: false, status: 'open', created_at: '2026-08-20T00:00:00Z' },
    { id: 'jp-005', title: 'QA Automation Engineer', description: 'Build E2E test suites with Playwright', type: 'contract', openings: 1, location: 'Remote', is_remote: true, status: 'closed', created_at: '2026-06-01T00:00:00Z' },
];

// ── Candidates ──
export const mockCandidates: Candidate[] = [
    { id: 'c-001', first_name: 'Saurabh', last_name: 'Tiwari', email: 'saurabh.t@outlook.com', phone: '+91-98123-45678', source: 'LinkedIn', created_at: '2026-08-05T00:00:00Z' },
    { id: 'c-002', first_name: 'Ritu', last_name: 'Dubey', email: 'ritu.d@gmail.com', phone: '+91-99876-54321', source: 'Naukri', created_at: '2026-08-12T00:00:00Z' },
    { id: 'c-003', first_name: 'Karan', last_name: 'Malhotra', email: 'karan.m@proton.me', phone: '+91-97654-32100', source: 'Referral', created_at: '2026-08-18T00:00:00Z' },
    { id: 'c-004', first_name: 'Divya', last_name: 'Krishnan', email: 'divya.k@yahoo.com', phone: '+91-96543-21098', source: 'LinkedIn', created_at: '2026-08-25T00:00:00Z' },
];

// ── Review Cycles ──
export const mockReviewCycles: ReviewCycle[] = [
    { id: 'rc-001', name: 'H1 2026 Performance Review', type: 'semi_annual', year: 2026, start_date: '2026-07-01', end_date: '2026-07-31', status: 'completed' },
    { id: 'rc-002', name: 'H2 2026 Performance Review', type: 'semi_annual', year: 2026, start_date: '2027-01-01', end_date: '2027-01-31', status: 'upcoming' },
    { id: 'rc-003', name: 'Q3 2026 OKR Review', type: 'quarterly', year: 2026, start_date: '2026-10-01', end_date: '2026-10-15', status: 'upcoming' },
];

// ── Goals ──
export const mockGoals: Goal[] = [
    { id: 'g-001', emp_id: 'e-004', title: 'Complete authentication module', description: 'Implement login, register, MFA, and RBAC', target_value: 100, current_value: 85, weight: 40, type: 'kpi', status: 'active', due_date: '2026-09-30' },
    { id: 'g-002', emp_id: 'e-004', title: 'Improve code coverage to 80%', description: 'Write unit and integration tests', target_value: 80, current_value: 62, weight: 30, type: 'kpi', status: 'active', due_date: '2026-10-15' },
    { id: 'g-003', emp_id: 'e-004', title: 'Learn Kubernetes', description: 'Complete CKA certification prep', target_value: 100, current_value: 40, weight: 15, type: 'learning', status: 'active', due_date: '2026-12-31' },
    { id: 'g-004', emp_id: 'e-007', title: 'Build payroll processing engine', description: 'Batch processing for 500+ employees', target_value: 100, current_value: 100, weight: 50, type: 'kpi', status: 'completed', due_date: '2026-08-01' },
    { id: 'g-005', emp_id: 'e-006', title: 'Redesign leave management UI', description: 'Mobile-first responsive dashboard', target_value: 100, current_value: 70, weight: 40, type: 'okr', status: 'active', due_date: '2026-09-15' },
];

// ── Courses ──
export const mockCourses: Course[] = [
    { id: 'cr-001', title: 'React Advanced Patterns', description: 'Compound components, render props, custom hooks', category: 'Frontend', type: 'video', duration_hrs: 8, is_mandatory: false },
    { id: 'cr-002', title: 'Go Concurrency Deep Dive', description: 'Goroutines, channels, select, and worker pools', category: 'Backend', type: 'video', duration_hrs: 6, is_mandatory: false },
    { id: 'cr-003', title: 'Information Security Awareness', description: 'Phishing, password hygiene, data handling', category: 'Compliance', type: 'pdf', duration_hrs: 1.5, is_mandatory: true },
    { id: 'cr-004', title: 'POSH Training', description: 'Prevention of Sexual Harassment at workplace', category: 'Compliance', type: 'classroom', duration_hrs: 2, is_mandatory: true },
    { id: 'cr-005', title: 'System Design Fundamentals', description: 'Load balancers, caching, databases, message queues', category: 'Engineering', type: 'video', duration_hrs: 12, is_mandatory: false },
];

// ── Notifications ──
export const mockNotifications: Notification[] = [
    { id: 'n-001', type: 'leave', title: 'Leave request from Rohit Joshi', message: 'Rohit has applied for 3 days sick leave (Sep 5-8)', action_url: '/leaves', is_read: false, created_at: '2026-08-31T14:00:00Z' },
    { id: 'n-002', type: 'leave', title: 'Leave approved', message: 'Your casual leave for Aug 31 - Sep 2 has been approved', action_url: '/leaves', is_read: false, created_at: '2026-08-30T16:00:00Z' },
    { id: 'n-003', type: 'payroll', title: 'August payroll disbursed', message: 'Payslips for August 2026 are now available', action_url: '/payroll', is_read: true, created_at: '2026-08-29T10:00:00Z' },
    { id: 'n-004', type: 'attendance', title: 'Late check-in alert', message: 'Vikram Singh checked in at 9:25 AM (25 min late)', action_url: '/attendance', is_read: true, created_at: '2026-08-31T09:30:00Z' },
    { id: 'n-005', type: 'recruitment', title: 'New candidate application', message: 'Divya Krishnan applied for React Frontend Developer', action_url: '/recruitment', is_read: false, created_at: '2026-08-31T11:00:00Z' },
    { id: 'n-006', type: 'performance', title: 'Goal deadline approaching', message: 'Your goal "Redesign leave management UI" is due in 15 days', action_url: '/performance/goals', is_read: true, created_at: '2026-08-31T08:00:00Z' },
];

// ── Announcements ──
export const mockAnnouncements: Announcement[] = [
    { id: 'ann-001', title: 'Independence Day Holiday', content: 'The office will remain closed on August 15th (Friday) for Independence Day. Wishing everyone a happy Independence Day!', published_at: '2026-08-12T10:00:00Z', expires_at: '2026-08-16T00:00:00Z' },
    { id: 'ann-002', title: 'Q3 Town Hall - September 5th', content: 'Join the all-hands town hall on September 5th at 3 PM IST. CEO Rajesh Sharma will share the Q2 performance review and Q3 roadmap. Link will be shared via email.', published_at: '2026-08-28T09:00:00Z', expires_at: '2026-09-06T00:00:00Z' },
    { id: 'ann-003', title: 'New Health Insurance Partner', content: 'We have partnered with Star Health Insurance for enhanced health coverage starting October 1st. Details will be shared by HR team next week.', published_at: '2026-08-30T11:00:00Z', expires_at: null },
];

// ── Workflow Templates ──
export const mockWorkflows: WorkflowTemplate[] = [
    { id: 'wf-001', name: 'Leave Approval (2-Level)', module: 'leave', steps: [{ step: 1, approver_role: 'Manager', order: 1 }, { step: 2, approver_role: 'HR Manager', order: 2 }], is_active: true },
    { id: 'wf-002', name: 'Expense Approval', module: 'payroll', steps: [{ step: 1, approver_role: 'Manager', order: 1 }, { step: 2, approver_role: 'Finance', order: 2 }], is_active: true },
    { id: 'wf-003', name: 'Document Verification', module: 'employee', steps: [{ step: 1, approver_role: 'HR Manager', order: 1 }], is_active: false },
];

// ── Helper: get employee name by ID ──
export function getEmployeeName(empId: string): string {
    const emp = mockEmployees.find(e => e.id === empId);
    return emp ? `${emp.first_name} ${emp.last_name}` : empId;
}

// ── Dashboard stats ──
export const mockDashboardStats = {
    totalEmployees: mockEmployees.length,
    departments: mockDepartments.filter(d => d.is_active).length,
    presentToday: mockAttendance.filter(a => a.status === 'present').length,
    pendingLeaves: mockLeaveApplications.filter(l => l.status === 'pending').length,
    payrollThisMonth: mockPayrollRuns.find(p => p.month === 8 && p.year === 2026)?.total_net ?? 0,
    openPositions: mockJobPostings.filter(j => j.status === 'open').reduce((sum, j) => sum + j.openings, 0),
};
