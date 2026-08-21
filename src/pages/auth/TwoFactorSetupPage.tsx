import { TwoFactorSetup } from '../../components/auth/TwoFactorSetup';
import { AuthLayout } from '../../components/layout/MainLayout';
import { useState } from 'react';

export default function TwoFactorSetupPage() {
  const [secret] = useState('JBSWY3DPEHPK3PXP');
  const [qrCode] = useState('https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=otpauth%3A%2F%2Ftotp%2FGreenCombats%253Aalex.morgan%40greencombats.io%3Fsecret%3DJBSWY3DPEHPK3PXP%26issuer%3DGreenCombats');
  const [backupCodes] = useState([
    'BC-1234-5678', 'BC-8765-4321', 'BC-1111-2222',
    'BC-3333-4444', 'BC-5555-6666', 'BC-7777-8888',
    'BC-9999-0000', 'BC-AAAA-BBBB', 'BC-CCCC-DDDD', 'BC-EEEE-FFFF',
  ]);

  return (
    <AuthLayout>
      <TwoFactorSetup
        secret={secret}
        qrCode={qrCode}
        backupCodes={backupCodes}
        onVerify={async (code) => {
          if (code !== '123456') throw new Error('Invalid code');
        }}
        onComplete={() => {}}
      />
    </AuthLayout>
  );
}