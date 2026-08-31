import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/auth';
import { ThemeProvider } from './context/theme';
import ProtectedRoute from './components/ProtectedRoute';

// Layouts
import DashboardLayout from './layouts/DashboardLayout';
import AuthLayout from './layouts/AuthLayout';

// Auth pages
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import ForgotPasswordPage from './pages/auth/ForgotPasswordPage';

// Module pages
import DashboardPage from './pages/dashboard/DashboardPage';
import EmployeeListPage from './pages/employees/EmployeeListPage';
import EmployeeDetailPage from './pages/employees/EmployeeDetailPage';
import DepartmentsPage from './pages/departments/DepartmentsPage';
import DesignationsPage from './pages/designations/DesignationsPage';
import AttendancePage from './pages/attendance/AttendancePage';
import LeaveListPage from './pages/leaves/LeaveListPage';
import LeaveApplyPage from './pages/leaves/LeaveApplyPage';
import PayrollRunsPage from './pages/payroll/PayrollRunsPage';
import PayslipsPage from './pages/payroll/PayslipsPage';
import ReimbursementsPage from './pages/payroll/ReimbursementsPage';
import JobsPage from './pages/recruitment/JobsPage';
import CandidatesPage from './pages/recruitment/CandidatesPage';
import ReviewsPage from './pages/performance/ReviewsPage';
import GoalsPage from './pages/performance/GoalsPage';
import CoursesPage from './pages/training/CoursesPage';
import NotificationsPage from './pages/notifications/NotificationsPage';
import AnnouncementsPage from './pages/notifications/AnnouncementsPage';
import RolesPage from './pages/settings/RolesPage';
import WorkflowsPage from './pages/settings/WorkflowsPage';

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            {/* Auth routes */}
            <Route element={<AuthLayout />}>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              <Route path="/forgot-password" element={<ForgotPasswordPage />} />
            </Route>

            {/* Protected dashboard routes */}
            <Route
              element={
                <ProtectedRoute>
                  <DashboardLayout />
                </ProtectedRoute>
              }
            >
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/employees" element={<EmployeeListPage />} />
              <Route path="/employees/:id" element={<EmployeeDetailPage />} />
              <Route path="/departments" element={<DepartmentsPage />} />
              <Route path="/designations" element={<DesignationsPage />} />
              <Route path="/attendance" element={<AttendancePage />} />
              <Route path="/leaves" element={<LeaveListPage />} />
              <Route path="/leaves/apply" element={<LeaveApplyPage />} />
              <Route path="/payroll" element={<PayrollRunsPage />} />
              <Route path="/payroll/payslips" element={<PayslipsPage />} />
              <Route path="/payroll/reimbursements" element={<ReimbursementsPage />} />
              <Route path="/recruitment" element={<JobsPage />} />
              <Route path="/recruitment/candidates" element={<CandidatesPage />} />
              <Route path="/performance" element={<ReviewsPage />} />
              <Route path="/performance/goals" element={<GoalsPage />} />
              <Route path="/training" element={<CoursesPage />} />
              <Route path="/notifications" element={<NotificationsPage />} />
              <Route path="/announcements" element={<AnnouncementsPage />} />
              <Route path="/settings/roles" element={<RolesPage />} />
              <Route path="/settings/workflows" element={<WorkflowsPage />} />
            </Route>

            {/* Redirects */}
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}
