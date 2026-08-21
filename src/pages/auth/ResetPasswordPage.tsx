import { ResetPasswordForm } from '../../components/auth/ResetPasswordForm';
import { AuthLayout } from '../../components/layout/MainLayout';

export default function ResetPasswordPage() {
  return (
    <AuthLayout>
      <ResetPasswordForm />
    </AuthLayout>
  );
}