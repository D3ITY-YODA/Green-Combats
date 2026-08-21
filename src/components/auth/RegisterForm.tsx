import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { EyeIcon, EyeOffIcon, LoaderIcon, MailIcon, LockIcon, UserIcon, ShieldIcon, GithubIcon, GitlabIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter } from '../ui/Card';
import { useAuthStore } from '../../store';
import { toast } from 'react-hot-toast';

const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
  organizationName: z.string().optional(),
  acceptTerms: z.boolean().refine((val) => val === true, 'You must accept the terms and conditions'),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});

type RegisterFormData = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const navigate = useNavigate();
  const { register: registerUser } = useAuthStore();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);

  const {
    register,
    handleSubmit,
    formState: { errors },
    watch,
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
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

  const onSubmit = async (data: RegisterFormData) => {
    setLoading(true);
    try {
      await registerUser(data);
      toast.success('Account created successfully!');
      navigate('/dashboard');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  const strengthLabels = ['Very Weak', 'Weak', 'Fair', 'Strong', 'Very Strong'];
  const strengthColors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-lime-500', 'bg-green-500'];

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
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Create your account</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            Start your climate action intelligence journey
          </p>
        </div>

        <Card className="shadow-xl">
          <CardHeader className="pb-4">
            <h2 className="text-lg font-semibold text-secondary-900 dark:text-white">Sign up</h2>
          </CardHeader>
          <CardBody className="pt-0">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
              <Input
                label="Full Name"
                type="text"
                placeholder="Alex Morgan"
                iconLeft={<UserIcon className="h-5 w-5" />}
                error={errors.name?.message}
                {...register('name')}
                autoComplete="name"
                disabled={loading}
              />
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
              <Input
                label="Organization (Optional)"
                type="text"
                placeholder="Green Combats Global"
                iconLeft={<ShieldIcon className="h-5 w-5" />}
                {...register('organizationName')}
                autoComplete="organization"
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
                  label="Confirm Password"
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
              <div className="flex items-start gap-2">
                <input
                  type="checkbox"
                  id="acceptTerms"
                  className="mt-0.5 h-4 w-4 rounded border-secondary-300 text-primary-600 focus:ring-2 focus:ring-primary-500/20"
                  {...register('acceptTerms')}
                />
                <label htmlFor="acceptTerms" className="text-sm text-secondary-600 dark:text-secondary-400">
                  I agree to the{' '}
                  <Link to="/terms" className="text-primary-600 hover:text-primary-700 dark:text-primary-400">
                    Terms of Service
                  </Link>{' '}
                  and{' '}
                  <Link to="/privacy" className="text-primary-600 hover:text-primary-700 dark:text-primary-400">
                    Privacy Policy
                  </Link>
                </label>
              </div>

              <Button
                type="submit"
                variant="primary"
                fullWidth
                loading={loading}
                className="mt-2"
              >
                Create account
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
              Already have an account?{' '}
              <Link to="/auth/login" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                Sign in
              </Link>
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}