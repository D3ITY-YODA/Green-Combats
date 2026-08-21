import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { EyeIcon, EyeOffIcon, LoaderIcon, MailIcon, LockIcon, UserIcon, GithubIcon, GitlabIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter } from '../ui/Card';
import { useAuthStore } from '../../store';
import { toast } from 'react-hot-toast';

const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  rememberMe: z.boolean().optional(),
  twoFactorCode: z.string().optional(),
});

type LoginFormData = z.infer<typeof loginSchema>;

export function LoginForm() {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  const [showPassword, setShowPassword] = useState(false);
  const [is2FA, setIs2FA] = useState(false);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      rememberMe: true,
    },
  });

  const onSubmit = async (data: LoginFormData) => {
    setLoading(true);
    try {
      await login(data);
      toast.success('Welcome back!');
      navigate('/dashboard');
    } catch (error) {
      if (error instanceof Error && error.message.includes('2FA')) {
        setIs2FA(true);
      } else {
        toast.error(error instanceof Error ? error.message : 'Login failed');
      }
    } finally {
      setLoading(false);
    }
  };

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
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">
            {is2FA ? 'Two-Factor Authentication' : 'Welcome back'}
          </h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            {is2FA
              ? 'Enter the code from your authenticator app'
              : 'Sign in to your account to continue'}
          </p>
        </div>

        <Card className="shadow-xl">
          <CardHeader className="pb-4">
            <h2 className="text-lg font-semibold text-secondary-900 dark:text-white">
              {is2FA ? 'Verify your identity' : 'Sign in'}
            </h2>
          </CardHeader>
          <CardBody className="pt-0">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
              {!is2FA && (
                <>
                  <Input
                    label="Email"
                    type="email"
                    placeholder="you@example.com"
                    iconLeft={<MailIcon className="h-5 w-5" />}
                    error={errors.email?.message}
                    {...register('email')}
                    autoComplete="email"
                    disabled={loading}
                  />
                  <div className="relative">
                    <Input
                      label="Password"
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
                      autoComplete="current-password"
                      disabled={loading}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <label className="flex items-center gap-2 cursor-pointer">
                      <input
                        type="checkbox"
                        className="h-4 w-4 rounded border-secondary-300 text-primary-600 focus:ring-2 focus:ring-primary-500/20"
                        {...register('rememberMe')}
                      />
                      <span className="text-sm text-secondary-600 dark:text-secondary-400">Remember me</span>
                    </label>
                    <Link
                      to="/auth/forgot-password"
                      className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                    >
                      Forgot password?
                    </Link>
                  </div>
                </>
              )}

              {is2FA && (
                <Input
                  label="Authentication Code"
                  type="text"
                  placeholder="123456"
                  iconLeft={<LockIcon className="h-5 w-5" />}
                  error={errors.twoFactorCode?.message}
                  {...register('twoFactorCode')}
                  autoComplete="one-time-code"
                  inputMode="numeric"
                  maxLength={6}
                  disabled={loading}
                />
              )}

              <Button
                type="submit"
                variant="primary"
                fullWidth
                loading={loading}
                className="mt-2"
              >
                {is2FA ? 'Verify' : 'Sign in'}
              </Button>
            </form>

            <div className="mt-6">
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-secondary-200 dark:border-secondary-700" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-4 bg-white text-secondary-500 dark:bg-secondary-800 dark:text-secondary-400">
                    Or continue with
                  </span>
                </div>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-3">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => toast.success('GitHub login coming soon')}
                  disabled={loading}
                >
                  <GithubIcon className="h-5 w-5" />
                  <span>GitHub</span>
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => toast.success('GitLab login coming soon')}
                  disabled={loading}
                >
                  <GitlabIcon className="h-5 w-5" />
                  <span>GitLab</span>
                </Button>
              </div>
            </div>
          </CardBody>
          <CardFooter className="pt-4">
            <p className="text-center text-sm text-secondary-600 dark:text-secondary-400">
              Don't have an account?{' '}
              <Link to="/auth/register" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                Sign up
              </Link>
            </p>
          </CardFooter>
        </Card>

        <p className="mt-6 text-center text-xs text-secondary-500 dark:text-secondary-400">
          Demo credentials: demo@greencombats.io / demo123
        </p>
      </div>
    </div>
  );
}