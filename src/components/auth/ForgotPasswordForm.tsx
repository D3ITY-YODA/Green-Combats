import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { MailIcon, LoaderIcon, ArrowLeftIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter } from '../ui/Card';
import { useAuthStore } from '../../store';
import { toast } from 'react-hot-toast';

const forgotPasswordSchema = z.object({
  email: z.string().email('Invalid email address'),
});

type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export function ForgotPasswordForm() {
  const navigate = useNavigate();
  const [step, setStep] = useState<'request' | 'sent'>('request');
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordFormData) => {
    setLoading(true);
    try {
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1500));
      toast.success('Password reset email sent!');
      setStep('sent');
      reset();
    } catch (error) {
      toast.error('Failed to send reset email');
    } finally {
      setLoading(false);
    }
  };

  if (step === 'sent') {
    return (
      <div className="min-h-screen flex items-center justify-center bg-secondary-50 dark:bg-secondary-900 px-4 py-12">
        <div className="w-full max-w-md text-center">
          <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
            <svg className="h-8 w-8 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Check your email</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            We've sent password reset instructions to your email address.
          </p>
          <div className="mt-6 p-4 bg-secondary-100 dark:bg-secondary-800 rounded-lg text-sm text-secondary-600 dark:text-secondary-400">
            <p className="font-medium">Didn't receive the email?</p>
            <ul className="mt-2 space-y-1 text-left">
              <li>• Check your spam or junk folder</li>
              <li>• Make sure you entered the correct email</li>
              <li>• Wait a few minutes for delivery</li>
            </ul>
          </div>
          <div className="mt-6 flex flex-col gap-3">
            <Button variant="primary" onClick={() => navigate('/auth/login')}>
              Back to Sign In
            </Button>
            <Button variant="ghost" onClick={() => setStep('request')}>
              Resend Email
            </Button>
          </div>
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
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Forgot password?</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            Enter your email and we'll send you a link to reset your password
          </p>
        </div>

        <Card className="shadow-xl">
          <CardHeader className="pb-4">
            <h2 className="text-lg font-semibold text-secondary-900 dark:text-white">Reset password</h2>
          </CardHeader>
          <CardBody className="pt-0">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
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
              <Button
                type="submit"
                variant="primary"
                fullWidth
                loading={loading}
              >
                Send reset link
              </Button>
            </form>
          </CardBody>
          <CardFooter className="pt-4">
            <p className="text-center text-sm text-secondary-600 dark:text-secondary-400">
              <Link to="/auth/login" className="flex items-center justify-center gap-1 text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
                <ArrowLeftIcon className="h-4 w-4" />
                Back to Sign In
              </Link>
            </p>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}