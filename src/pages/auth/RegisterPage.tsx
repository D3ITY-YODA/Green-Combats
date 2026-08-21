import { RegisterForm } from '../../components/auth/RegisterForm';
import { AuthLayout } from '../../components/layout/MainLayout';

export default function RegisterPage() {
  return (
    <AuthLayout>
      <RegisterForm />
    </AuthLayout>
  );
}