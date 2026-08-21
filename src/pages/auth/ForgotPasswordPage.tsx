import { ForgotPasswordForm } from '../../components/auth/ForgotPasswordForm';
import { AuthLayout } from '../../components/layout/MainLayout';

export default function ForgotPasswordPage() {
  return (
    <AuthLayout>
      <ForgotPasswordForm />
    </AuthLayout>
  );
}