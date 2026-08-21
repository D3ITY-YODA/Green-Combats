import { useState } from 'react';
import { qrcode } from 'qrcode-generator';
import { CheckIcon, CopyIcon, AlertTriangleIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { toast } from 'react-hot-toast';

interface TwoFactorSetupProps {
  secret: string;
  qrCode: string;
  backupCodes: string[];
  onVerify: (code: string) => Promise<void>;
  onComplete: () => void;
}

export function TwoFactorSetup({ secret, qrCode, backupCodes, onVerify, onComplete }: TwoFactorSetupProps) {
  const [step, setStep] = useState<'qr' | 'verify' | 'backup'>('qr');
  const [code, setCode] = useState('');
  const [copiedCodes, setCopiedCodes] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(false);

  const handleVerify = async () => {
    if (code.length !== 6) {
      toast.error('Please enter a 6-digit code');
      return;
    }
    setLoading(true);
    try {
      await onVerify(code);
      setStep('backup');
    } catch (error) {
      toast.error('Invalid code. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const copyCode = (index: number) => {
    navigator.clipboard.writeText(backupCodes[index]);
    setCopiedCodes((prev) => new Set([...prev, index]));
    toast.success('Backup code copied!');
    setTimeout(() => {
      setCopiedCodes((prev) => {
        const next = new Set(prev);
        next.delete(index);
        return next;
      });
    }, 2000);
  };

  const downloadCodes = () => {
    const content = `Green Combats - 2FA Backup Codes\nGenerated: ${new Date().toISOString()}\n\n${backupCodes.map((code, i) => `${i + 1}. ${code}`).join('\n')}`;
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'green-combats-2fa-backup-codes.txt';
    a.click();
    URL.revokeObjectURL(url);
    toast.success('Backup codes downloaded!');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-secondary-50 dark:bg-secondary-900 px-4 py-12">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="inline-flex items-center gap-2 mb-6">
            <div className="h-10 w-10 rounded-xl bg-primary-600 flex items-center justify-center">
              <span className="text-white font-bold text-xl">GC</span>
            </div>
            <span className="text-2xl font-bold text-secondary-900 dark:text-white">Green Combats</span>
          </div>
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Two-Factor Authentication</h1>
          <p className="mt-2 text-secondary-600 dark:text-secondary-400">
            {step === 'qr' && 'Set up 2FA to secure your account'}
            {step === 'verify' && 'Verify your authenticator app'}
            {step === 'backup' && 'Save your backup codes'}
          </p>
        </div>

        <Card className="shadow-xl">
          <CardHeader className="pb-4">
            <div className="flex items-center gap-2">
              <div className={cn(
                'h-8 w-8 rounded-full flex items-center justify-center',
                step === 'qr' ? 'bg-primary-100 text-primary-600' : 'bg-green-100 text-green-600'
              )}>
                {step === 'qr' ? <span className="text-sm font-bold">1</span> : <CheckIcon className="h-5 w-5" />}
              </div>
              <div className="flex-1 h-1 bg-secondary-200 dark:bg-secondary-700" />
              <div className={cn(
                'h-8 w-8 rounded-full flex items-center justify-center',
                step === 'verify' ? 'bg-primary-100 text-primary-600' : step === 'qr' ? 'bg-secondary-200 text-secondary-400' : 'bg-green-100 text-green-600'
              )}>
                {step === 'verify' ? <span className="text-sm font-bold">2</span> : step === 'qr' ? <span className="text-sm font-bold">2</span> : <CheckIcon className="h-5 w-5" />}
              </div>
              <div className="flex-1 h-1 bg-secondary-200 dark:bg-secondary-700" />
              <div className={cn(
                'h-8 w-8 rounded-full flex items-center justify-center',
                step === 'backup' ? 'bg-primary-100 text-primary-600' : step !== 'backup' ? 'bg-secondary-200 text-secondary-400' : 'bg-green-100 text-green-600'
              )}>
                {step === 'backup' ? <span className="text-sm font-bold">3</span> : <span className="text-sm font-bold">3</span>}
              </div>
            </div>
          </CardHeader>
          <CardBody className="pt-0">
            {step === 'qr' && (
              <div className="space-y-6">
                <div className="text-center">
                  <p className="text-secondary-600 dark:text-secondary-400 mb-4">
                    Scan this QR code with your authenticator app
                  </p>
                  <div className="inline-block p-4 bg-white dark:bg-secondary-800 rounded-lg">
                    <img src={qrCode} alt="QR Code for 2FA setup" className="h-48 w-48" />
                  </div>
                  <p className="mt-4 text-sm text-secondary-500 dark:text-secondary-400">
                    Or enter this key manually: <code className="font-mono text-primary-600 dark:text-primary-400">{secret}</code>
                  </p>
                </div>
                <Button variant="primary" fullWidth onClick={() => setStep('verify')}>
                  I've scanned the code
                </Button>
              </div>
            )}

            {step === 'verify' && (
              <div className="space-y-6">
                <div className="text-center">
                  <p className="text-secondary-600 dark:text-secondary-400 mb-4">
                    Enter the 6-digit code from your authenticator app
                  </p>
                  <Input
                    label="Authentication Code"
                    type="text"
                    placeholder="123456"
                    value={code}
                    onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                    inputMode="numeric"
                    maxLength={6}
                    autoComplete="one-time-code"
                    className="text-center text-2xl tracking-widest"
                  />
                </div>
                <Button variant="primary" fullWidth onClick={handleVerify} loading={loading}>
                  Verify & Continue
                </Button>
              </div>
            )}

            {step === 'backup' && (
              <div className="space-y-6">
                <div className="flex items-start gap-3 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
                  <AlertTriangleIcon className="h-5 w-5 text-yellow-600 dark:text-yellow-400 mt-0.5 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium text-yellow-800 dark:text-yellow-300">Important!</h4>
                    <p className="mt-1 text-sm text-yellow-700 dark:text-yellow-400">
                      Save these backup codes in a secure place. Each code can only be used once.
                      If you lose access to your authenticator app, these codes are the only way to recover your account.
                    </p>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-2">
                  {backupCodes.map((code, index) => (
                    <div key={index} className="flex items-center gap-2 p-3 bg-secondary-50 dark:bg-secondary-800 rounded-lg">
                      <code className="flex-1 font-mono text-sm text-secondary-900 dark:text-white bg-white dark:bg-secondary-700 px-2 py-1 rounded">
                        {code}
                      </code>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => copyCode(index)}
                        aria-label={copiedCodes.has(index) ? 'Copied' : 'Copy code'}
                      >
                        {copiedCodes.has(index) ? (
                          <CheckIcon className="h-4 w-4 text-green-600" />
                        ) : (
                          <CopyIcon className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  ))}
                </div>
                <div className="flex gap-3">
                  <Button variant="secondary" fullWidth onClick={downloadCodes}>
                    Download Codes
                  </Button>
                  <Button variant="primary" fullWidth onClick={onComplete}>
                    Done
                  </Button>
                </div>
              </div>
            )}
          </CardBody>
        </Card>
      </div>
    </div>
  );
}