import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { BellIcon, MailIcon, SmartphoneIcon, CalendarIcon, ShieldIcon, LoaderIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../ui/Card';
import { Switch } from '../ui/Input';
import { Badge } from '../ui/Badge';
import { Tabs, TabItem } from '../ui/Tabs';
import { toast } from 'react-hot-toast';

const notificationSchema = z.object({
  email: z.boolean(),
  push: z.boolean(),
  inApp: z.boolean(),
  frequency: z.enum(['immediate', 'hourly', 'daily', 'weekly']),
  alerts: z.boolean(),
  reports: z.boolean(),
  updates: z.boolean(),
  mentions: z.boolean(),
  assignments: z.boolean(),
});

type NotificationFormData = z.infer<typeof notificationSchema>;

const tabs: TabItem[] = [
  { id: 'channels', label: 'Channels', icon: <BellIcon className="h-4 w-4" /> },
  { id: 'types', label: 'Notification Types', icon: <MailIcon className="h-4 w-4" /> },
  { id: 'schedule', label: 'Schedule', icon: <CalendarIcon className="h-4 w-4" /> },
];

const notificationTypes = [
  { key: 'alerts', label: 'Alerts & Warnings', description: 'Critical thresholds, anomalies, compliance issues', icon: <ShieldIcon className="h-4 w-4 text-red-500" /> },
  { key: 'reports', label: 'Reports & Analytics', description: 'Scheduled reports, data exports, insights', icon: <CalendarIcon className="h-4 w-4 text-blue-500" /> },
  { key: 'updates', label: 'Platform Updates', description: 'New features, maintenance, version releases', icon: <SmartphoneIcon className="h-4 w-4 text-green-500" /> },
  { key: 'mentions', label: 'Mentions & Comments', description: 'Team mentions, comments on shared items', icon: <MailIcon className="h-4 w-4 text-purple-500" /> },
  { key: 'assignments', label: 'Task Assignments', description: 'New tasks, deadline changes, approvals', icon: <CalendarIcon className="h-4 w-4 text-orange-500" /> },
];

const frequencies = [
  { value: 'immediate', label: 'Immediate', description: 'Receive notifications as they happen' },
  { value: 'hourly', label: 'Hourly Digest', description: 'Batch notifications every hour' },
  { value: 'daily', label: 'Daily Digest', description: 'Single summary email each morning' },
  { value: 'weekly', label: 'Weekly Digest', description: 'Weekly summary every Monday' },
];

export function NotificationSettings() {
  const [activeTab, setActiveTab] = useState('channels');
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
    watch,
    setValue,
  } = useForm<NotificationFormData>({
    resolver: zodResolver(notificationSchema),
    defaultValues: {
      email: true,
      push: true,
      inApp: true,
      frequency: 'immediate',
      alerts: true,
      reports: true,
      updates: true,
      mentions: true,
      assignments: true,
    },
  });

  const onSubmit = async (data: NotificationFormData) => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      toast.success('Notification preferences saved!');
      reset(data);
    } catch (error) {
      toast.error('Failed to save preferences');
    } finally {
      setLoading(false);
    }
  };

  const handleTestNotification = async (channel: string) => {
    toast.success(`Test ${channel} notification sent!`);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Notifications</h1>
        <p className="mt-1 text-secondary-600 dark:text-secondary-400">
          Configure how and when you receive notifications
        </p>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Channels Tab */}
        <div hidden={activeTab !== 'channels'} id="channels-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Delivery Channels</CardTitle>
              <CardDescription>Choose where you want to receive notifications</CardDescription>
            </CardHeader>
            <CardBody className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
                        <MailIcon className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                      </div>
                      <div>
                        <p className="font-medium text-secondary-900 dark:text-white">Email</p>
                        <p className="text-sm text-secondary-500 dark:text-secondary-400">Receive notifications via email</p>
                      </div>
                    </div>
                    <Switch {...register('email')} />
                  </div>
                  <div className="mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700 flex gap-2">
                    <Button variant="outline" size="sm" onClick={() => handleTestNotification('email')}>Test</Button>
                    <Badge variant="outline" size="sm">Primary</Badge>
                  </div>
                </div>

                <div className="p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
                        <SmartphoneIcon className="h-5 w-5 text-green-600 dark:text-green-400" />
                      </div>
                      <div>
                        <p className="font-medium text-secondary-900 dark:text-white">Push</p>
                        <p className="text-sm text-secondary-500 dark:text-secondary-400">Mobile & desktop push notifications</p>
                      </div>
                    </div>
                    <Switch {...register('push')} />
                  </div>
                  <div className="mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700 flex gap-2">
                    <Button variant="outline" size="sm" onClick={() => handleTestNotification('push')}>Test</Button>
                    <Badge variant="outline" size="sm">Secondary</Badge>
                  </div>
                </div>

                <div className="p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
                        <BellIcon className="h-5 w-5 text-purple-600 dark:text-purple-400" />
                      </div>
                      <div>
                        <p className="font-medium text-secondary-900 dark:text-white">In-App</p>
                        <p className="text-sm text-secondary-500 dark:text-secondary-400">Notifications within the platform</p>
                      </div>
                    </div>
                    <Switch {...register('inApp')} />
                  </div>
                  <div className="mt-4 pt-4 border-t border-secondary-200 dark:border-secondary-700 flex gap-2">
                    <Button variant="outline" size="sm" onClick={() => handleTestNotification('in-app')}>Test</Button>
                    <Badge variant="outline" size="sm">Always On</Badge>
                  </div>
                </div>
              </div>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Email Preferences</CardTitle>
              <CardDescription>Fine-tune your email notification settings</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Marketing Emails</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Product updates, tips, and announcements</p>
                </div>
                <Switch defaultChecked={false} />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Security Emails</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Login alerts, password changes, 2FA</p>
                </div>
                <Switch defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Digest Emails</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Summary of activity and insights</p>
                </div>
                <Switch defaultChecked />
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Types Tab */}
        <div hidden={activeTab !== 'types'} id="types-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Notification Types</CardTitle>
              <CardDescription>Enable or disable specific notification categories</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              {notificationTypes.map((type) => (
                <div
                  key={type.key}
                  className="flex items-center justify-between p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl"
                >
                  <div className="flex items-center gap-4">
                    <span className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                      {type.icon}
                    </span>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">{type.label}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">{type.description}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <Switch
                      {...register(type.key as keyof NotificationFormData)}
                      aria-label={type.label}
                    />
                    <Button variant="outline" size="sm" onClick={() => handleTestNotification(type.key)}>
                      Test
                    </Button>
                  </div>
                </div>
              ))}
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Channel Preferences by Type</CardTitle>
              <CardDescription>Choose which channels deliver each notification type</CardDescription>
            </CardHeader>
            <CardBody>
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-secondary-200 dark:border-secondary-700">
                      <th className="text-left py-3 px-4 font-medium text-secondary-500">Notification Type</th>
                      <th className="text-center py-3 px-4 font-medium text-secondary-500">Email</th>
                      <th className="text-center py-3 px-4 font-medium text-secondary-500">Push</th>
                      <th className="text-center py-3 px-4 font-medium text-secondary-500">In-App</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-secondary-200 dark:divide-secondary-700">
                    {notificationTypes.map((type) => (
                      <tr key={type.key}>
                        <td className="py-3 px-4 font-medium text-secondary-900 dark:text-white">{type.label}</td>
                        <td className="py-3 px-4 text-center">
                          <input type="checkbox" className="h-4 w-4 rounded border-secondary-300 text-primary-600" defaultChecked />
                        </td>
                        <td className="py-3 px-4 text-center">
                          <input type="checkbox" className="h-4 w-4 rounded border-secondary-300 text-primary-600" defaultChecked={type.key !== 'updates'} />
                        </td>
                        <td className="py-3 px-4 text-center">
                          <input type="checkbox" className="h-4 w-4 rounded border-secondary-300 text-primary-600" defaultChecked />
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Schedule Tab */}
        <div hidden={activeTab !== 'schedule'} id="schedule-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Delivery Schedule</CardTitle>
              <CardDescription>Control when notifications are delivered</CardDescription>
            </CardHeader>
            <CardBody className="space-y-6">
              <div>
                <label className="label">Default Frequency</label>
                <div className="grid gap-3 sm:grid-cols-4">
                  {frequencies.map((freq) => (
                    <button
                      key={freq.value}
                      type="button"
                      onClick={() => setValue('frequency', freq.value as 'immediate' | 'hourly' | 'daily' | 'weekly')}
                      className={cn(
                        'p-4 rounded-xl border-2 transition-all text-left',
                        watch('frequency') === freq.value
                          ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                          : 'border-secondary-200 hover:border-secondary-300 dark:border-secondary-700 dark:hover:border-secondary-600'
                      )}
                    >
                      <p className="font-medium text-secondary-900 dark:text-white">{freq.label}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400 mt-1">{freq.description}</p>
                    </button>
                  ))}
                </div>
              </div>

              <div className="border-t border-secondary-200 dark:border-secondary-700 pt-6">
                <h4 className="font-medium text-secondary-900 dark:text-white mb-4">Quiet Hours</h4>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-secondary-900 dark:text-white">Enable Quiet Hours</p>
                    <p className="text-sm text-secondary-500 dark:text-secondary-400">Suppress non-critical notifications during set hours</p>
                  </div>
                  <Switch defaultChecked={false} />
                </div>
                <div className="mt-4 grid gap-4 sm:grid-cols-2">
                  <div>
                    <label className="label">Start Time</label>
                    <Input type="time" defaultValue="22:00" />
                  </div>
                  <div>
                    <label className="label">End Time</label>
                    <Input type="time" defaultValue="07:00" />
                  </div>
                </div>
                <div className="mt-4 flex items-center justify-between">
                  <div>
                    <p className="font-medium text-secondary-900 dark:text-white">Override for Critical Alerts</p>
                    <p className="text-sm text-secondary-500 dark:text-secondary-400">Always deliver critical alerts during quiet hours</p>
                  </div>
                  <Switch defaultChecked />
                </div>
              </div>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Scheduled Reports</CardTitle>
              <CardDescription>Configure automatic report delivery</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              {[
                { name: 'Weekly Emissions Summary', frequency: 'Weekly', day: 'Monday', time: '08:00', enabled: true },
                { name: 'Monthly Carbon Footprint', frequency: 'Monthly', day: '1st', time: '09:00', enabled: true },
                { name: 'Quarterly Compliance Report', frequency: 'Quarterly', day: '1st of Quarter', time: '10:00', enabled: false },
              ].map((report, i) => (
                <div key={i} className="flex items-center justify-between p-4 border border-secondary-200 dark:border-secondary-700 rounded-xl">
                  <div className="flex items-center gap-4">
                    <div className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                      <CalendarIcon className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
                    </div>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">{report.name}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">
                        {report.frequency} on {report.day} at {report.time}
                      </p>
                    </div>
                  </div>
                  <Switch defaultChecked={report.enabled} />
                </div>
              ))}
              <Button variant="outline" size="sm" className="w-full">Add Scheduled Report</Button>
            </CardBody>
          </Card>
        </div>

        <div className="flex justify-end gap-3">
          <Button type="button" variant="secondary" onClick={() => reset()} disabled={!isDirty}>
            Reset to Defaults
          </Button>
          <Button type="submit" variant="primary" loading={loading}>
            Save Notification Settings
          </Button>
        </div>
      </form>
    </div>
  );
}