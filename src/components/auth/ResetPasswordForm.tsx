import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { EyeIcon, EyeOffIcon, LoaderIcon, LockIcon, CheckIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter } from '../ui/Card';
import { toast } from 'react-hot-toast';

const resetPasswordSchema = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

export function ResetPasswordForm() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [validToken, setValidToken] = useState(true);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
  });

  const password = watch('password', '');

  const calculateStrength = (pwd: string) => {
    let strength = 0;
    if (pwd.length >= 8) strength++;
    if (/[A-Z]/.test(pwd)) strength++;
    if (/[a-z]/.test(pwd)) strength++;
    if (/[0-9]/.test(pwd)) strength++;
    if (/[^A-Za-z0-9]/.test(pwd)) strength++;
    return Math.min(strength, 4);
  };

  const onSubmit = async (data: ResetPasswordFormData) => {
    if (!token) {
      toast.error('Invalid reset token');
      return;
    }

    setLoading(true);
    try {
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1500));
      toast.success('Password has been reset!');
      navigate('/auth/login');
    } catch (error) {
      toast.error('Failed to reset password');
    } finally {
      setLoading(false);
    }
  };

  const strengthLabels = ['Very Weak', 'Weak', 'Fair', 'Strong', 'Very Strong'];
  const strengthColors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-lime-500', 'bg-green-500'];

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-secondary-50 dark:bg-secondary-900 px-4 py-12">
        <div className="w-full max-w-md text-center">
          <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
            <svg className="h-8 w-8 text-red-600 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Invalid reset link</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            This password reset link is invalid or has expired.
          </p>
          <Button variant="primary" className="mt-6" onClick={() => navigate('/auth/forgot-password')}>
            Request new link
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-secondary-50 dark:bg-secondary-900 px-4 py-12">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Link to="/" className="inline-flex items-center gap-2 mb-6">
            <div className="h-10 w-10 rounded-xl bg-primary-600 flex items-center justify-center">
              <span className="text-white font-bold text-xl">GC</span>
            </div>
            <span className="text-2xl font-bold text-secondary-900 dark:text-white">Green Combats</span>
          </Link>
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Set new password</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            Your new password must be different from previously used passwords
          </p>
        </div>

        <Card className="shadow-xl">
          <CardHeader className="pb-4">
            <h2 className="text-lg font-semibold text-secondary-900 dark:text-white">Create new password</h2>
          </CardHeader>
          <CardBody className="pt-0">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
              <div className="relative">
                <Input
                  label="New Password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  iconLeft={<LockIcon className="h-5 w-5" />}
                  iconRight={
                    <button
                      type="button"
                      className="text-secondary-400 hover:text-secondary-600 dark:hover:text-secondary-300"
                      onClick={() => setShowPassword(!showPassword)}
                      aria-label={showPassword ? 'Hide password' : 'Show password'}
                    >
                      {showPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                    </button>
                  }
                  error={errors.password?.message}
                  {...register('password')}
                  autoComplete="new-password"
                  disabled={loading}
                  onChange={(e) => setPasswordStrength(calculateStrength(e.target.value))}
                />
              </div>
              {password && (
                <div className="space-y-1.5">
                  <div className="h-1.5 bg-secondary-200 dark:bg-secondary-700 rounded-full overflow-hidden">
                    <div
                      className={cn(
                        'h-full transition-all duration-300 rounded-full',
                        strengthColors[passwordStrength]
                      )}
                      style={{ width: `${(passwordStrength / 4) * 100}%` }}
                    />
                  </div>
                  <p className="text-xs text-secondary-500 dark:text-secondary-400">
                    Password strength: {strengthLabels[passwordStrength]}
                  </p>
                </div>
              )}
              <div className="relative">
                <Input
                  label="Confirm New Password"
                  type={showConfirmPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  iconLeft={<LockIcon className="h-5 w-5" />}
                  iconRight={
                    <button
                      type="button"
                      className="text-secondary-400 hover:text-secondary-600 dark:hover:text-secondary-300"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
                    >
                      {showConfirmPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                    </button>
                  }
                  error={errors.confirmPassword?.message}
                  {...register('confirmPassword')}
                  autoComplete="new-password"
                  disabled={loading}
                />
              </div>

              <Button
                type="submit"
                variant="primary"
                fullWidth
                loading={loading}
              >
                Reset password
              </Button>
            </form>
          </CardBody>
          <CardFooter className="pt-4">
            <p className="text-center text-sm text-secondary-600 dark:text-secondary-400">
              <Link to="/auth/login" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                Back to Sign In
              </Link>
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}