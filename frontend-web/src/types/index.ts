// ── Auth & Users ──
export interface User {
    id: string;
    org_id: string;
    email: string;
    role_id: string;
    is_active: boolean;
    is_verified: boolean;
    mfa_enabled: boolean;
    last_login_at: string | null;
    created_at: string;
    role?: Role;
    permissions?: string[];
}

export interface Role {
    id: string;
    org_id: string;
    name: string;
    description: string;
    is_system: boolean;
}

export interface Permission {
    id: string;
    resource: string;
    action: string;
}

export interface Organization {
    id: string;
    name: string;
    logo_url: string | null;
    address: string | null;
    timezone: string;
}

// ── Employee ──
export interface Employee {
    id: string;
    org_id: string;
    user_id: string | null;
    emp_code: string;
    first_name: string;
    last_name: string;
    middle_name: string | null;
    dob: string | null;
    gender: string | null;
    phone: string | null;
    personal_email: string | null;
    address: string | null;
    city: string | null;
    state: string | null;
    country: string | null;
    pincode: string | null;
    dept_id: string | null;
    designation_id: string | null;
    manager_id: string | null;
    employment_type: 'full_time' | 'part_time' | 'contract' | 'intern';
    work_location: string | null;
    joining_date: string;
    status: 'active' | 'probation' | 'notice_period' | 'resigned' | 'terminated' | 'absconding';
    photo_url: string | null;
    department?: string;
    designation?: string;
    manager_name?: string;
}

export interface Department {
    id: string;
    org_id: string;
    name: string;
    code: string | null;
    parent_dept_id: string | null;
    head_emp_id: string | null;
    description: string | null;
    is_active: boolean;
}

export interface Designation {
    id: string;
    org_id: string;
    department_id: string | null;
    title: string;
    grade: string | null;
    level: number | null;
}

// ── Attendance ──
export interface AttendanceRecord {
    id: string;
    emp_id: string;
    date: string;
    check_in: string | null;
    check_out: string | null;
    method: 'face' | 'manual' | 'geo' | 'biometric';
    status: 'present' | 'absent' | 'half_day' | 'late' | 'wfh' | 'on_leave';
    duration_mins: number | null;
    overtime_mins: number;
}

export interface Shift {
    id: string;
    name: string;
    start_time: string;
    end_time: string;
    grace_mins: number;
}

// ── Leave ──
export interface LeaveType {
    id: string;
    name: string;
    code: string;
    max_days_per_year: number | null;
    is_paid: boolean;
    requires_doc: boolean;
}

export interface LeaveApplication {
    id: string;
    emp_id: string;
    leave_type_id: string;
    from_date: string;
    to_date: string;
    days: number;
    session: 'full' | 'first_half' | 'second_half';
    reason: string | null;
    status: 'pending' | 'approved' | 'rejected' | 'cancelled' | 'withdrawn';
    created_at: string;
    leave_type?: string;
    employee_name?: string;
}

export interface LeaveBalance {
    id: string;
    emp_id: string;
    leave_type_id: string;
    year: number;
    total: number;
    used: number;
    pending: number;
    available: number;
    leave_type?: string;
}

// ── Payroll ──
export interface PayrollRun {
    id: string;
    org_id: string;
    month: number;
    year: number;
    status: 'draft' | 'processing' | 'pending_approval' | 'approved' | 'disbursed';
    total_gross: number | null;
    total_deductions: number | null;
    total_net: number | null;
}

export interface Payslip {
    id: string;
    emp_id: string;
    payroll_run_id: string;
    gross: number;
    total_deductions: number;
    net_pay: number;
    working_days: number | null;
    lop_days: number;
    breakdown: Record<string, number> | null;
}

export interface Reimbursement {
    id: string;
    emp_id: string;
    category: string;
    amount: number;
    description: string | null;
    status: 'pending' | 'approved' | 'rejected' | 'paid';
    created_at: string;
}

// ── Recruitment ──
export interface JobPosting {
    id: string;
    title: string;
    description: string | null;
    type: 'full_time' | 'part_time' | 'contract' | 'internship';
    openings: number;
    location: string | null;
    is_remote: boolean;
    status: 'draft' | 'open' | 'paused' | 'closed' | 'cancelled';
    created_at: string;
}

export interface Candidate {
    id: string;
    first_name: string;
    last_name: string | null;
    email: string;
    phone: string | null;
    source: string | null;
    created_at: string;
}

// ── Performance ──
export interface ReviewCycle {
    id: string;
    name: string;
    type: 'annual' | 'semi_annual' | 'quarterly';
    year: number;
    start_date: string;
    end_date: string;
    status: 'upcoming' | 'active' | 'completed' | 'archived';
}

export interface Goal {
    id: string;
    emp_id: string;
    title: string;
    description: string | null;
    target_value: number | null;
    current_value: number;
    weight: number;
    type: 'kpi' | 'okr' | 'learning';
    status: 'active' | 'completed' | 'missed' | 'cancelled';
    due_date: string | null;
}

// ── Training ──
export interface Course {
    id: string;
    title: string;
    description: string | null;
    category: string | null;
    type: 'video' | 'pdf' | 'scorm' | 'classroom';
    duration_hrs: number | null;
    is_mandatory: boolean;
}

// ── Notifications ──
export interface Notification {
    id: string;
    type: string;
    title: string;
    message: string | null;
    action_url: string | null;
    is_read: boolean;
    created_at: string;
}

export interface Announcement {
    id: string;
    title: string;
    content: string;
    published_at: string;
    expires_at: string | null;
}

// ── Workflow ──
export interface WorkflowTemplate {
    id: string;
    name: string;
    module: string;
    steps: unknown[];
    is_active: boolean;
}

// ── API Response ──
export interface ApiResponse<T> {
    success: boolean;
    data: T;
    message?: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    data: T[];
    total: number;
    page: number;
    limit: number;
}
