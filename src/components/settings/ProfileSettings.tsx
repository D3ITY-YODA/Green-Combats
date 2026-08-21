import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { CameraIcon, LoaderIcon, MapPinIcon, GlobeIcon, BellIcon, ShieldIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Select } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter, CardTitle, CardDescription } from '../ui/Card';
import { Avatar } from '../ui/Avatar';
import { Badge } from '../ui/Badge';
import { useAuthStore } from '../../store';
import { toast } from 'react-hot-toast';

const profileSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  jobTitle: z.string().optional(),
  department: z.string().optional(),
  phone: z.string().optional(),
  location: z.string().optional(),
  timezone: z.string(),
  language: z.string(),
  dateFormat: z.string(),
  bio: z.string().max(500, 'Bio must be 500 characters or less').optional(),
});

type ProfileFormData = z.infer<typeof profileSchema>;

const timezones = [
  'UTC', 'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
  'Europe/London', 'Europe/Paris', 'Europe/Berlin', 'Europe/Rome',
  'Asia/Tokyo', 'Asia/Shanghai', 'Asia/Singapore', 'Asia/Dubai',
  'Australia/Sydney', 'Pacific/Auckland',
];

const languages = [
  { value: 'en', label: 'English' },
  { value: 'es', label: 'Spanish' },
  { value: 'fr', label: 'French' },
  { value: 'de', label: 'German' },
  { value: 'it', label: 'Italian' },
  { value: 'pt', label: 'Portuguese' },
  { value: 'ja', label: 'Japanese' },
  { value: 'zh', label: 'Chinese' },
  { value: 'ko', label: 'Korean' },
];

const dateFormats = [
  { value: 'YYYY-MM-DD', label: '2024-01-15 (ISO)' },
  { value: 'MM/DD/YYYY', label: '01/15/2024 (US)' },
  { value: 'DD/MM/YYYY', label: '15/01/2024 (UK/EU)' },
  { value: 'MMM D, YYYY', label: 'Jan 15, 2024' },
];

export function ProfileSettings() {
  const { user, updateUser } = useAuthStore();
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
    watch,
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name || '',
      email: user?.email || '',
      jobTitle: '',
      department: '',
      phone: '',
      location: '',
      timezone: user?.preferences.timezone || 'UTC',
      language: user?.preferences.language || 'en',
      dateFormat: user?.preferences.dateFormat || 'YYYY-MM-DD',
      bio: '',
    },
  });

  const onSubmit = async (data: ProfileFormData) => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      updateUser({
        name: data.name,
        email: data.email,
        preferences: {
          ...user!.preferences,
          timezone: data.timezone,
          language: data.language,
          dateFormat: data.dateFormat,
        },
      });
      toast.success('Profile updated successfully!');
      reset(data);
    } catch (error) {
      toast.error('Failed to update profile');
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setAvatarPreview(e.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Profile Settings</h1>
        <p className="mt-1 text-secondary-600 dark:text-secondary-400">
          Manage your personal information and preferences
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>Profile Picture</CardTitle>
            <CardDescription>Update your avatar and basic info</CardDescription>
          </CardHeader>
          <CardBody className="space-y-6">
            <div className="flex flex-col items-center text-center">
              <div className="relative">
                <Avatar
                  src={avatarPreview || user?.avatar}
                  name={user?.name}
                  size="2xl"
                />
                <label className="absolute bottom-0 right-0 p-2 bg-primary-600 text-white rounded-full cursor-pointer hover:bg-primary-700 transition-colors">
                  <CameraIcon className="h-4 w-4" />
                  <input type="file" accept="image/*" className="sr-only" onChange={handleAvatarChange} />
                </label>
              </div>
              <div className="mt-4 space-y-2">
                <p className="font-medium text-secondary-900 dark:text-white">{user?.name}</p>
                <p className="text-sm text-secondary-500 dark:text-secondary-400">{user?.email}</p>
                <Badge variant="primary" size="sm">{user?.role}</Badge>
              </div>
            </div>
            <div className="border-t border-secondary-200 dark:border-secondary-700 pt-4 space-y-4">
              <div>
                <label className="label">Job Title</label>
                <Input placeholder="e.g. Climate Analyst" {...register('jobTitle')} />
              </div>
              <div>
                <label className="label">Department</label>
                <Input placeholder="e.g. Sustainability" {...register('department')} />
              </div>
              <div>
                <label className="label">Phone</label>
                <Input type="tel" placeholder="+1 (555) 000-0000" {...register('phone')} />
              </div>
              <div>
                <label className="label">Location</label>
                <div className="relative">
                  <MapPinIcon className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-secondary-400" />
                  <Input placeholder="San Francisco, CA" {...register('location')} className="pl-10" />
                </div>
              </div>
              <div>
                <label className="label">Bio</label>
                <textarea
                  {...register('bio')}
                  className="input min-h-[100px] resize-y"
                  placeholder="Tell us about yourself..."
                />
              </div>
            </div>
          </CardBody>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
            <CardDescription>Update your account details and preferences</CardDescription>
          </CardHeader>
          <CardBody>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="label">Full Name</label>
                  <Input error={errors.name?.message} {...register('name')} />
                </div>
                <div>
                  <label className="label">Email</label>
                  <Input type="email" error={errors.email?.message} {...register('email')} />
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div>
                  <label className="label">Timezone</label>
                  <Select
                    options={timezones.map((tz) => ({ value: tz, label: tz }))}
                    {...register('timezone')}
                  />
                </div>
                <div>
                  <label className="label">Language</label>
                  <Select options={languages} {...register('language')} />
                </div>
                <div>
                  <label className="label">Date Format</label>
                  <Select options={dateFormats} {...register('dateFormat')} />
                </div>
              </div>

              <CardFooter className="border-t border-secondary-200 dark:border-secondary-700">
                <div className="flex justify-end gap-3">
                  <Button type="button" variant="secondary" onClick={() => reset()} disabled={!isDirty}>
                    Reset
                  </Button>
                  <Button type="submit" variant="primary" loading={loading}>
                    Save Changes
                  </Button>
                </div>
              </CardFooter>
            </form>
          </CardBody>
        </Card>
      </div>
    </div>
  );
}