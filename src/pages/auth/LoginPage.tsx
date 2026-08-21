import { LoginForm } from '../../components/auth/LoginForm';
import { AuthLayout } from '../../components/layout/MainLayout';

export default function LoginPage() {
  return (
    <AuthLayout>
      <LoginForm />
    </AuthLayout>
  );
}