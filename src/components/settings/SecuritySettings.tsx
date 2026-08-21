import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { LockIcon, ShieldIcon, SmartphoneIcon, MailIcon, KeyIcon, EyeIcon, EyeOffIcon, LoaderIcon, RotateCcwIcon, Trash2Icon, DownloadIcon, CopyIcon, AlertTriangleIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardTitle, CardDescription, CardFooter } from '../ui/Card';
import { Switch } from '../ui/Input';
import { Badge } from '../ui/Badge';
import { Tabs, TabItem } from '../ui/Tabs';
import { Modal, ConfirmModal } from '../ui/Modal';
import { useAuthStore } from '../../store';
import { toast } from 'react-hot-toast';

const passwordSchema = z.object({
  currentPassword: z.string().min(1, 'Current password is required'),
  newPassword: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});

type PasswordFormData = z.infer<typeof passwordSchema>;

const tabs: TabItem[] = [
  { id: 'password', label: 'Password', icon: <LockIcon className="h-4 w-4" /> },
  { id: '2fa', label: 'Two-Factor Auth', icon: <ShieldIcon className="h-4 w-4" /> },
  { id: 'sessions', label: 'Active Sessions', icon: <SmartphoneIcon className="h-4 w-4" /> },
  { id: 'api', label: 'API Keys', icon: <KeyIcon className="h-4 w-4" /> },
  { id: 'login-history', label: 'Login History', icon: <KeyIcon className="h-4 w-4" /> },
];

export function SecuritySettings() {
  const { user, updateUser } = useAuthStore();
  const [activeTab, setActiveTab] = useState('password');
  const [loading, setLoading] = useState(false);
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [changePasswordModal, setChangePasswordModal] = useState(false);
  const [revokeSessionModal, setRevokeSessionModal] = useState<{ id: string; current: boolean } | null>(null);
  const [deleteApiKeyModal, setDeleteApiKeyModal] = useState<string | null>(null);
  const [showBackupCodes, setShowBackupCodes] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
    setValue,
  } = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema),
  });

  const newPassword = watch('newPassword', '');

  const calculateStrength = (pwd: string) => {
    let strength = 0;
    if (pwd.length >= 8) strength++;
    if (/[A-Z]/.test(pwd)) strength++;
    if (/[a-z]/.test(pwd)) strength++;
    if (/[0-9]/.test(pwd)) strength++;
    if (/[^A-Za-z0-9]/.test(pwd)) strength++;
    return Math.min(strength, 4);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue('newPassword', e.target.value);
    setPasswordStrength(calculateStrength(e.target.value));
  };

  const onPasswordSubmit = async (data: PasswordFormData) => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      toast.success('Password changed successfully!');
      setChangePasswordModal(false);
      reset();
    } catch (error) {
      toast.error('Failed to change password');
    } finally {
      setLoading(false);
    }
  };

  const revokeSession = async (sessionId: string, current: boolean) => {
    if (current) {
      toast.error('Cannot revoke current session');
      return;
    }
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 500));
      toast.success('Session revoked');
      setRevokeSessionModal(null);
    } catch (error) {
      toast.error('Failed to revoke session');
    } finally {
      setLoading(false);
    }
  };

  const revokeAllSessions = async () => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      toast.success('All other sessions revoked');
    } catch (error) {
      toast.error('Failed to revoke sessions');
    } finally {
      setLoading(false);
    }
  };

  const deleteApiKey = async (keyId: string) => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 500));
      toast.success('API key deleted');
      setDeleteApiKeyModal(null);
    } catch (error) {
      toast.error('Failed to delete API key');
    } finally {
      setLoading(false);
    }
  };

  const strengthLabels = ['Very Weak', 'Weak', 'Fair', 'Strong', 'Very Strong'];
  const strengthColors = ['bg-red-500', 'bg-orange-500', 'bg-yellow-500', 'bg-lime-500', 'bg-green-500'];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Security</h1>
        <p className="mt-1 text-secondary-600 dark:text-secondary-400">
          Manage your account security settings
        </p>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      {/* Password Tab */}
      <div hidden={activeTab !== 'password'} id="password-panel" role="tabpanel">
        <Card>
          <CardHeader>
            <CardTitle>Change Password</CardTitle>
            <CardDescription>Update your password regularly for security</CardDescription>
          </CardHeader>
          <CardBody>
            <Button variant="outline" onClick={() => setChangePasswordModal(true)}>
              <LockIcon className="h-4 w-4 mr-2" />
              Change Password
            </Button>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Password Requirements</CardTitle>
            <CardDescription>Your password must meet the following criteria</CardDescription>
          </CardHeader>
          <CardBody className="space-y-2">
            {[
              { label: 'At least 8 characters', met: true },
              { label: 'One uppercase letter', met: true },
              { label: 'One lowercase letter', met: true },
              { label: 'One number', met: true },
              { label: 'One special character', met: false },
            ].map((req, i) => (
              <div key={i} className="flex items-center gap-2">
                <span className={cn('h-5 w-5 rounded-full flex items-center justify-center', req.met ? 'bg-green-100 text-green-600' : 'bg-secondary-100 text-secondary-400')}>
                  {req.met ? <svg className="h-3 w-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg> : <span className="text-xs">•</span>}
                </span>
                <span className="text-sm text-secondary-700 dark:text-secondary-300">{req.label}</span>
              </div>
            ))}
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Password History</CardTitle>
            <CardDescription>Last changed: {user?.security.passwordLastChanged ? new Date(user.security.passwordLastChanged).toLocaleDateString() : 'Unknown'}</CardDescription>
          </CardHeader>
          <CardBody>
            <Button variant="outline" onClick={() => setChangePasswordModal(true)}>
              <RotateCcwIcon className="h-4 w-4 mr-2" />
              Change Password Now
            </Button>
          </CardBody>
        </Card>
      </div>

      {/* 2FA Tab */}
      <div hidden={activeTab !== '2fa'} id="2fa-panel" role="tabpanel">
        <Card>
          <CardHeader>
            <CardTitle>Two-Factor Authentication</CardTitle>
            <CardDescription>Add an extra layer of security to your account</CardDescription>
          </CardHeader>
          <CardBody className="space-y-6">
            <div className="flex items-center justify-between p-4 bg-secondary-50 dark:bg-secondary-800 rounded-xl">
              <div className="flex items-center gap-4">
                <div className="p-2 bg-primary-100 dark:bg-primary-900/30 rounded-lg">
                  <ShieldIcon className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                </div>
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Authenticator App</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">
                    {user?.security.twoFactorEnabled ? 'Enabled' : 'Not configured'}
                  </p>
                </div>
              </div>
              {user?.security.twoFactorEnabled ? (
                <div className="flex gap-2">
                  <Button variant="outline" size="sm" onClick={() => setShowBackupCodes(true)}>
                    <EyeIcon className="h-4 w-4 mr-2" />
                    View Backup Codes
                  </Button>
                  <Button variant="danger" size="sm" onClick={() => toast.success('2FA disabled')}>
                    <Trash2Icon className="h-4 w-4 mr-2" />
                    Disable
                  </Button>
                </div>
              ) : (
                <Button variant="primary" onClick={() => toast.success('2FA setup initiated')}>
                  <ShieldIcon className="h-4 w-4 mr-2" />
                  Enable 2FA
                </Button>
              )}
            </div>

            {user?.security.twoFactorEnabled && (
              <div className="space-y-4">
                <div className="p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
                  <h4 className="font-medium text-secondary-900 dark:text-white mb-2">Backup Codes</h4>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400 mb-4">
                    Store these codes in a safe place. Each code can only be used once.
                  </p>
                  <div className="grid gap-2 sm:grid-cols-2">
                    {user?.security.backupCodes?.map((code, index) => (
                      <div key={index} className="flex items-center justify-between p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                        <code className="font-mono text-sm text-secondary-900 dark:text-white">{code}</code>
                        <Button variant="ghost" size="sm" onClick={() => toast.success('Copied!')}>
                          <CopyIcon className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                  <div className="flex gap-2 mt-4">
                    <Button variant="outline" size="sm" onClick={() => toast.success('Codes regenerated')}>
                      <RotateCcwIcon className="h-4 w-4 mr-2" />
                      Regenerate Codes
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => toast.success('Codes downloaded')}>
                      <DownloadIcon className="h-4 w-4 mr-2" />
                      Download
                    </Button>
                  </div>
                </div>
              </div>
            )}
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Backup Methods</CardTitle>
            <CardDescription>Additional verification methods for account recovery</CardDescription>
          </CardHeader>
          <CardBody className="space-y-4">
            <div className="flex items-center justify-between p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
              <div className="flex items-center gap-4">
                <div className="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
                  <MailIcon className="h-5 w-5 text-green-600 dark:text-green-400" />
                </div>
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Email Verification</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Backup codes sent to your email</p>
                </div>
              </div>
              <Badge variant="success" size="sm">Verified</Badge>
            </div>
            <div className="flex items-center justify-between p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
              <div className="flex items-center gap-4">
                <div className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                  <SmartphoneIcon className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
                </div>
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">SMS Verification</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Backup codes sent via SMS</p>
                </div>
              </div>
              <Button variant="outline" size="sm">Set Up</Button>
            </div>
          </CardBody>
        </Card>
      </div>

      {/* Sessions Tab */}
      <div hidden={activeTab !== 'sessions'} id="sessions-panel" role="tabpanel">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Active Sessions</CardTitle>
              <CardDescription>Manage your active login sessions</CardDescription>
            </div>
            <Button variant="outline" size="sm" onClick={revokeAllSessions} disabled={loading}>
              <RotateCcwIcon className="h-4 w-4 mr-2" />
              Revoke All Others
            </Button>
          </CardHeader>
          <CardBody className="p-0">
            <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
              {user?.security.sessions.map((session) => (
                <div key={session.id} className="p-4 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                      <SmartphoneIcon className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
                    </div>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">{session.device}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">
                        {session.browser} on {session.os}
                      </p>
                      <p className="text-xs text-secondary-400 dark:text-secondary-500">
                        {session.location} • {session.ip}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    {session.current && (
                      <Badge variant="success" size="sm">Current</Badge>
                    )}
                    <span className="text-sm text-secondary-500 dark:text-secondary-400">
                      Last active: {new Date(session.lastActiveAt).toLocaleString()}
                    </span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setRevokeSessionModal({ id: session.id, current: session.current })}
                      disabled={session.current || loading}
                    >
                      <Trash2Icon className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>

      {/* API Keys Tab */}
      <div hidden={activeTab !== 'api'} id="api-panel" role="tabpanel">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>API Keys</CardTitle>
              <CardDescription>Manage your API keys for programmatic access</CardDescription>
            </div>
            <Button variant="primary" onClick={() => toast.success('New API key created')}>
              <KeyIcon className="h-4 w-4 mr-2" />
              Create API Key
            </Button>
          </CardHeader>
          <CardBody className="p-0">
            <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
              {[
                { id: '1', name: 'Production API', prefix: 'gc_live_', permissions: ['read', 'write', 'reports'], lastUsed: '2 hours ago', status: 'active' as const },
                { id: '2', name: 'Development Key', prefix: 'gc_test_', permissions: ['read'], lastUsed: '3 days ago', status: 'active' as const },
                { id: '3', name: 'Analytics Export', prefix: 'gc_live_', permissions: ['read', 'export'], lastUsed: 'Never', status: 'revoked' as const },
              ].map((key) => (
                <div key={key.id} className="p-4 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                      <KeyIcon className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
                    </div>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">{key.name}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">
                        {key.prefix}•••••••• • Last used: {key.lastUsed}
                      </p>
                      <div className="flex gap-1 mt-1">
                        {key.permissions.map((p) => (
                          <Badge key={p} variant="outline" size="sm">{p}</Badge>
                        ))}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant={key.status === 'active' ? 'success' : 'secondary'} size="sm">
                      {key.status}
                    </Badge>
                    <Button variant="ghost" size="sm" onClick={() => setDeleteApiKeyModal(key.id)}>
                      <Trash2Icon className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>

      {/* Login History Tab */}
      <div hidden={activeTab !== 'login-history'} id="login-history-panel" role="tabpanel">
        <Card>
          <CardHeader>
            <CardTitle>Login History</CardTitle>
            <CardDescription>Recent sign-in activity for your account</CardDescription>
          </CardHeader>
          <CardBody className="p-0">
            <div className="divide-y divide-secondary-200 dark:divide-secondary-700">
              {user?.security.loginHistory.map((login) => (
                <div key={login.id} className="p-4 flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className={cn('p-2 rounded-lg', login.success ? 'bg-green-100 dark:bg-green-900/30' : 'bg-red-100 dark:bg-red-900/30')}>
                      <LockIcon className={cn('h-5 w-5', login.success ? 'text-green-600' : 'text-red-600')} />
                    </div>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">
                        {login.success ? 'Successful sign in' : 'Failed sign in attempt'}
                      </p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">
                        {login.device} • {login.browser} • {login.location}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-secondary-900 dark:text-white">{new Date(login.timestamp).toLocaleString()}</p>
                    <p className="text-xs text-secondary-400 dark:text-secondary-500">{login.ip}</p>
                    {!login.success && login.failureReason && (
                      <p className="text-xs text-red-600 dark:text-red-400 mt-1">{login.failureReason}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      </div>

      {/* Modals */}
      <Modal
        isOpen={changePasswordModal}
        onClose={() => setChangePasswordModal(false)}
        title="Change Password"
        size="md"
        footer={
          <div className="flex justify-end gap-3">
            <Button variant="secondary" onClick={() => setChangePasswordModal(false)}>Cancel</Button>
            <Button variant="primary" onClick={() => {}} loading={loading}>Change Password</Button>
          </div>
        }
      >
        <form onSubmit={handleSubmit(onPasswordSubmit)} className="space-y-4">
          <div className="relative">
            <Input
              label="Current Password"
              type={showCurrentPassword ? 'text' : 'password'}
              placeholder="••••••••"
              iconLeft={<LockIcon className="h-5 w-5" />}
              iconRight={
                <button type="button" className="text-secondary-400 hover:text-secondary-600" onClick={() => setShowCurrentPassword(!showCurrentPassword)}>
                  {showCurrentPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                </button>
              }
              error={errors.currentPassword?.message}
              {...register('currentPassword')}
              autoComplete="current-password"
            />
          </div>
          <div className="relative">
            <Input
              label="New Password"
              type={showNewPassword ? 'text' : 'password'}
              placeholder="••••••••"
              iconLeft={<LockIcon className="h-5 w-5" />}
              iconRight={
                <button type="button" className="text-secondary-400 hover:text-secondary-600" onClick={() => setShowNewPassword(!showNewPassword)}>
                  {showNewPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                </button>
              }
              error={errors.newPassword?.message}
              {...register('newPassword')}
              autoComplete="new-password"
              onChange={handlePasswordChange}
            />
          </div>
          {newPassword && (
            <div className="space-y-1.5">
              <div className="h-1.5 bg-secondary-200 dark:bg-secondary-700 rounded-full overflow-hidden">
                <div className={cn('h-full transition-all duration-300 rounded-full', strengthColors[passwordStrength])} style={{ width: `${(passwordStrength / 4) * 100}%` }} />
              </div>
              <p className="text-xs text-secondary-500 dark:text-secondary-400">Password strength: {strengthLabels[passwordStrength]}</p>
            </div>
          )}
          <div className="relative">
            <Input
              label="Confirm New Password"
              type={showConfirmPassword ? 'text' : 'password'}
              placeholder="••••••••"
              iconLeft={<LockIcon className="h-5 w-5" />}
              iconRight={
                <button type="button" className="text-secondary-400 hover:text-secondary-600" onClick={() => setShowConfirmPassword(!showConfirmPassword)}>
                  {showConfirmPassword ? <EyeOffIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
                </button>
              }
              error={errors.confirmPassword?.message}
              {...register('confirmPassword')}
              autoComplete="new-password"
            />
          </div>
        </form>
      </Modal>

      <ConfirmModal
        isOpen={!!revokeSessionModal}
        onClose={() => setRevokeSessionModal(null)}
        onConfirm={() => revokeSessionModal && revokeSession(revokeSessionModal.id, revokeSessionModal.current)}
        title="Revoke Session"
        message={revokeSessionModal?.current ? 'Cannot revoke current session' : 'Are you sure you want to revoke this session? You will be logged out on that device.'}
        confirmText="Revoke"
        variant="danger"
        loading={loading}
      />

      <ConfirmModal
        isOpen={!!deleteApiKeyModal}
        onClose={() => setDeleteApiKeyModal(null)}
        onConfirm={() => deleteApiKeyModal && deleteApiKey(deleteApiKeyModal)}
        title="Delete API Key"
        message="Are you sure you want to delete this API key? This action cannot be undone."
        confirmText="Delete"
        variant="danger"
        loading={loading}
      />

      <Modal
        isOpen={showBackupCodes}
        onClose={() => setShowBackupCodes(false)}
        title="Backup Codes"
        size="lg"
      >
        <div className="space-y-4">
          <p className="text-secondary-600 dark:text-secondary-400">
            Store these codes in a secure place. Each code can only be used once.
          </p>
          <div className="grid gap-2 sm:grid-cols-2">
            {user?.security.backupCodes?.map((code, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                <code className="font-mono text-sm text-secondary-900 dark:text-white">{code}</code>
                <Button variant="ghost" size="sm" onClick={() => toast.success('Copied!')}>
                  <CopyIcon className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
          <div className="flex gap-2 pt-4 border-t border-secondary-200 dark:border-secondary-700">
            <Button variant="outline" onClick={() => toast.success('Codes downloaded')}>
              <DownloadIcon className="h-4 w-4 mr-2" />
              Download
            </Button>
            <Button variant="outline" onClick={() => toast.success('Codes regenerated')}>
              <RotateCcwIcon className="h-4 w-4 mr-2" />
              Regenerate
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}