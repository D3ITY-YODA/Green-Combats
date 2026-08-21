import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { MonitorIcon, PaletteIcon, BellIcon, DatabaseIcon, RefreshCwIcon, ZapIcon, ChevronDownIcon, ChevronUpIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Select } from '../ui/Input';
import { Card, CardHeader, CardBody, CardFooter, CardTitle, CardDescription } from '../ui/Card';
import { Switch } from '../ui/Input';
import { Badge } from '../ui/Badge';
import { Tabs, TabItem } from '../ui/Tabs';
import { useAuthStore } from '../../store';
import { useUIStore } from '../../store';
import { toast } from 'react-hot-toast';

const preferencesSchema = z.object({
  theme: z.enum(['light', 'dark', 'system']),
  language: z.string(),
  timezone: z.string(),
  dateFormat: z.string(),
  compactMode: z.boolean(),
  showTooltips: z.boolean(),
  animations: z.boolean(),
  autoRefresh: z.boolean(),
  refreshInterval: z.number(),
  defaultDashboard: z.string(),
  defaultMapBasemap: z.string(),
  defaultMapZoom: z.number(),
  measureUnit: z.enum(['metric', 'imperial']),
  showCoordinates: z.boolean(),
});

type PreferencesFormData = z.infer<typeof preferencesSchema>;

const timezones = ['UTC', 'America/New_York', 'America/Los_Angeles', 'Europe/London', 'Europe/Paris', 'Asia/Tokyo'];
const languages = [
  { value: 'en', label: 'English' },
  { value: 'es', label: 'Spanish' },
  { value: 'fr', label: 'French' },
  { value: 'de', label: 'German' },
  { value: 'pt', label: 'Portuguese' },
  { value: 'ja', label: 'Japanese' },
  { value: 'zh', label: 'Chinese' },
];
const dateFormats = [
  { value: 'YYYY-MM-DD', label: '2024-01-15' },
  { value: 'MM/DD/YYYY', label: '01/15/2024' },
  { value: 'DD/MM/YYYY', label: '15/01/2024' },
  { value: 'MMM D, YYYY', label: 'Jan 15, 2024' },
];
const dashboards = [
  { value: 'overview', label: 'Overview' },
  { value: 'emissions', label: 'Emissions' },
  { value: 'energy', label: 'Energy' },
  { value: 'maps', label: 'Maps' },
];
const basemaps = [
  { value: 'satellite', label: 'Satellite' },
  { value: 'streets', label: 'Streets' },
  { value: 'topographic', label: 'Topographic' },
  { value: 'dark', label: 'Dark' },
  { value: 'light', label: 'Light' },
];
const refreshIntervals = [
  { value: 60000, label: '1 minute' },
  { value: 300000, label: '5 minutes' },
  { value: 600000, label: '10 minutes' },
  { value: 1800000, label: '30 minutes' },
  { value: 3600000, label: '1 hour' },
];

const tabs: TabItem[] = [
  { id: 'appearance', label: 'Appearance', icon: <PaletteIcon className="h-4 w-4" /> },
  { id: 'dashboard', label: 'Dashboard', icon: <MonitorIcon className="h-4 w-4" /> },
  { id: 'maps', label: 'Maps', icon: <GlobeIcon className="h-4 w-4" /> },
  { id: 'data', label: 'Data & Sync', icon: <DatabaseIcon className="h-4 w-4" /> },
  { id: 'accessibility', label: 'Accessibility', icon: <ZapIcon className="h-4 w-4" /> },
];

export function PreferencesSettings() {
  const { user, updateUser } = useAuthStore();
  const { theme: uiTheme, setTheme } = useUIStore();
  const [activeTab, setActiveTab] = useState('appearance');
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
    watch,
  } = useForm<PreferencesFormData>({
    resolver: zodResolver(preferencesSchema),
    defaultValues: {
      theme: user?.preferences.theme || 'system',
      language: user?.preferences.language || 'en',
      timezone: user?.preferences.timezone || 'UTC',
      dateFormat: user?.preferences.dateFormat || 'YYYY-MM-DD',
      compactMode: user?.preferences.dashboard?.compactMode || false,
      showTooltips: user?.preferences.dashboard?.showTooltips ?? true,
      animations: true,
      autoRefresh: user?.preferences.dashboard?.autoRefresh ?? true,
      refreshInterval: user?.preferences.dashboard?.refreshInterval || 300000,
      defaultDashboard: user?.preferences.dashboard?.defaultView || 'overview',
      defaultMapBasemap: user?.preferences.maps?.defaultBasemap || 'satellite',
      defaultMapZoom: user?.preferences.maps?.defaultZoom || 3,
      measureUnit: user?.preferences.maps?.measureUnit || 'metric',
      showCoordinates: user?.preferences.maps?.showCoordinates ?? true,
    },
  });

  const onSubmit = async (data: PreferencesFormData) => {
    setLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      updateUser({
        preferences: {
          theme: data.theme,
          language: data.language,
          timezone: data.timezone,
          dateFormat: data.dateFormat,
          notifications: user!.preferences.notifications,
          dashboard: {
            defaultView: data.defaultDashboard,
            autoRefresh: data.autoRefresh,
            refreshInterval: data.refreshInterval,
            compactMode: data.compactMode,
            showTooltips: data.showTooltips,
          },
          maps: {
            defaultBasemap: data.defaultMapBasemap,
            defaultZoom: data.defaultMapZoom,
            defaultCenter: user!.preferences.maps.defaultCenter,
            showCoordinates: data.showCoordinates,
            measureUnit: data.measureUnit,
          },
        },
      });
      setTheme(data.theme);
      toast.success('Preferences saved!');
      reset(data);
    } catch (error) {
      toast.error('Failed to save preferences');
    } finally {
      setLoading(false);
    }
  };

  const handleThemeChange = (newTheme: 'light' | 'dark' | 'system') => {
    setTheme(newTheme);
    // Update form
    // This would need setValue from react-hook-form
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-secondary-900 dark:text-white">Preferences</h1>
        <p className="mt-1 text-secondary-600 dark:text-secondary-400">
          Customize your experience and default settings
        </p>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Appearance Tab */}
        <div hidden={activeTab !== 'appearance'} id="appearance-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Theme</CardTitle>
              <CardDescription>Choose your preferred color scheme</CardDescription>
            </CardHeader>
            <CardBody>
              <div className="grid gap-4 sm:grid-cols-3">
                {(['light', 'dark', 'system'] as const).map((theme) => (
                  <button
                    key={theme}
                    type="button"
                    onClick={() => handleThemeChange(theme)}
                    className={cn(
                      'relative p-4 rounded-xl border-2 transition-all',
                      'text-left group',
                      uiTheme === theme
                        ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                        : 'border-secondary-200 hover:border-secondary-300 dark:border-secondary-700 dark:hover:border-secondary-600'
                    )}
                  >
                    <div className={cn(
                      'h-24 w-full rounded-lg mb-3 flex items-center justify-center',
                      theme === 'light' && 'bg-white border border-secondary-200',
                      theme === 'dark' && 'bg-secondary-900 border border-secondary-700',
                      theme === 'system' && 'bg-gradient-to-r from-white to-secondary-900 border border-secondary-200'
                    )}>
                      {theme === 'system' && (
                        <div className="flex items-center justify-center gap-2">
                          <span className="text-2xl">☀️</span>
                          <span className="text-2xl">🌙</span>
                        </div>
                      )}
                    </div>
                    <div className="font-medium text-secondary-900 dark:text-white capitalize">{theme}</div>
                    <div className="text-sm text-secondary-500 dark:text-secondary-400">
                      {theme === 'light' && 'Always use light mode'}
                      {theme === 'dark' && 'Always use dark mode'}
                      {theme === 'system' && 'Match system setting'}
                    </div>
                    {uiTheme === theme && (
                      <div className="absolute top-2 right-2 bg-primary-500 text-white rounded-full h-5 w-5 flex items-center justify-center">
                        <svg className="h-3 w-3" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" /></svg>
                      </div>
                    )}
                  </button>
                ))}
              </div>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Language & Region</CardTitle>
              <CardDescription>Set your preferred language and regional settings</CardDescription>
            </CardHeader>
            <CardBody>
              <div className="grid gap-4 sm:grid-cols-3">
                <div>
                  <label className="label">Language</label>
                  <Select options={languages} {...register('language')} />
                </div>
                <div>
                  <label className="label">Timezone</label>
                  <Select options={timezones.map((tz) => ({ value: tz, label: tz }))} {...register('timezone')} />
                </div>
                <div>
                  <label className="label">Date Format</label>
                  <Select options={dateFormats} {...register('dateFormat')} />
                </div>
              </div>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Interface</CardTitle>
              <CardDescription>Adjust UI density and behavior</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Compact Mode</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Reduce spacing for more content</p>
                </div>
                <Switch {...register('compactMode')} />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Show Tooltips</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Display helpful hints on hover</p>
                </div>
                <Switch {...register('showTooltips')} />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Animations</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Enable UI transitions and animations</p>
                </div>
                <Switch {...register('animations')} />
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Dashboard Tab */}
        <div hidden={activeTab !== 'dashboard'} id="dashboard-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Dashboard Behavior</CardTitle>
              <CardDescription>Configure how your dashboard loads and refreshes</CardDescription>
            </CardHeader>
            <CardBody className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="label">Default Dashboard View</label>
                  <Select options={dashboards} {...register('defaultDashboard')} />
                </div>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-secondary-900 dark:text-white">Auto Refresh</p>
                    <p className="text-sm text-secondary-500 dark:text-secondary-400">Automatically refresh data</p>
                  </div>
                  <Switch {...register('autoRefresh')} />
                </div>
              </div>
              <div>
                <label className="label">Refresh Interval</label>
                <Select
                  options={refreshIntervals}
                  {...register('refreshInterval', { valueAsNumber: true })}
                  disabled={!watch('autoRefresh')}
                />
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Maps Tab */}
        <div hidden={activeTab !== 'maps'} id="maps-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Map Settings</CardTitle>
              <CardDescription>Configure default map view and measurements</CardDescription>
            </CardHeader>
            <CardBody className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="label">Default Basemap</label>
                  <Select options={basemaps} {...register('defaultMapBasemap')} />
                </div>
                <div>
                  <label className="label">Default Zoom Level</label>
                  <Select
                    options={Array.from({ length: 18 }, (_, i) => ({ value: i + 1, label: `Level ${i + 1}` }))}
                    {...register('defaultMapZoom', { valueAsNumber: true })}
                  />
                </div>
              </div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <label className="label">Measurement Unit</label>
                  <Select options={[
                    { value: 'metric', label: 'Metric (km, m²)' },
                    { value: 'imperial', label: 'Imperial (mi, ft²)' },
                  ]} {...register('measureUnit')} />
                </div>
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Show Coordinates</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Display lat/lng on map</p>
                </div>
                <Switch {...register('showCoordinates')} />
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Data & Sync Tab */}
        <div hidden={activeTab !== 'data'} id="data-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Data & Synchronization</CardTitle>
              <CardDescription>Manage data refresh and storage preferences</CardDescription>
            </CardHeader>
            <CardBody className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Auto-sync Data</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Automatically sync with data sources</p>
                </div>
                <Switch defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Offline Support</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Cache data for offline viewing</p>
                </div>
                <Switch defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Data Compression</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Compress large datasets for faster loading</p>
                </div>
                <Switch defaultChecked />
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Accessibility Tab */}
        <div hidden={activeTab !== 'accessibility'} id="accessibility-panel" role="tabpanel">
          <Card>
            <CardHeader>
              <CardTitle>Accessibility</CardTitle>
              <CardDescription>Adjust settings for better accessibility</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">High Contrast</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Increase color contrast for better visibility</p>
                </div>
                <Switch />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Reduced Motion</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Minimize animations and transitions</p>
                </div>
                <Switch />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Focus Indicators</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Enhanced focus outlines for keyboard navigation</p>
                </div>
                <Switch defaultChecked />
              </div>
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">Screen Reader Optimizations</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Additional ARIA labels and live regions</p>
                </div>
                <Switch defaultChecked />
              </div>
            </CardBody>
          </Card>
        </div>

        <div className="flex justify-end gap-3">
          <Button type="button" variant="secondary" onClick={() => reset()} disabled={!isDirty}>
            Reset to Defaults
          </Button>
          <Button type="submit" variant="primary" loading={loading}>
            Save Preferences
          </Button>
        </div>
      </form>
    </div>
  );
}