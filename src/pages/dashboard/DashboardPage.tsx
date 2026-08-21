import { useState } from 'react';
import { LineChart, Line, AreaChart, Area, BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, ComposedChart } from 'recharts';
import { TrendingUpIcon, TrendingDownIcon, MinusIcon, TargetIcon, LeafIcon, ZapIcon, DropletIcon, CloudIcon, ArrowUpRightIcon, DownloadIcon, RefreshCwIcon, CalendarIcon, ChevronDownIcon } from 'lucide-react';
import { formatNumber, formatPercentage, formatRelativeTime } from '../../utils/helpers';
import { cn } from '../../utils/helpers';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Select } from '../../components/ui/Input';
import { Tabs, TabItem } from '../../components/ui/Tabs';
import { Dropdown, DropdownItem } from '../../components/ui/Dropdown';
import { Breadcrumb } from '../../components/ui/Breadcrumb';
import { Skeleton } from '../../components/ui/Skeleton';

const metricCards = [
  { label: 'Total Emissions', value: '12,450', unit: 'tCO₂e', trend: 2.3, trendUp: false, icon: CloudIcon, color: 'bg-red-500', subtitle: 'Scope 1+2+3' },
  { label: 'Carbon Intensity', value: '45.2', unit: 'tCO₂e/$M', trend: 5.1, trendUp: false, icon: TargetIcon, color: 'bg-orange-500', subtitle: 'YoY change' },
  { label: 'Renewable Energy', value: '67%', unit: '', trend: 12.5, trendUp: true, icon: ZapIcon, color: 'bg-green-500', subtitle: 'Of total consumption' },
  { label: 'Water Usage', value: '2.4M', unit: 'm³', trend: 1.8, trendUp: false, icon: DropletIcon, color: 'bg-blue-500', subtitle: 'Annual withdrawal' },
  { label: 'Active Projects', value: '23', unit: '', trend: 3, trendUp: true, icon: LeafIcon, color: 'bg-emerald-500', subtitle: 'Mitigation & adaptation' },
  { label: 'Compliance Score', value: '94%', unit: '', trend: 2, trendUp: true, icon: ShieldCheckIcon, color: 'bg-primary-500', subtitle: 'Regulatory adherence' },
];

const emissionsData = [
  { month: 'Jan', scope1: 1200, scope2: 800, scope3: 2100, target: 3800 },
  { month: 'Feb', scope1: 1150, scope2: 780, scope3: 2050, target: 3700 },
  { month: 'Mar', scope1: 1100, scope2: 750, scope3: 2000, target: 3600 },
  { month: 'Apr', scope1: 1050, scope2: 720, scope3: 1950, target: 3500 },
  { month: 'May', scope1: 1000, scope2: 700, scope3: 1900, target: 3400 },
  { month: 'Jun', scope1: 980, scope2: 680, scope3: 1850, target: 3300 },
  { month: 'Jul', scope1: 950, scope2: 650, scope3: 1800, target: 3200 },
  { month: 'Aug', scope1: 920, scope2: 630, scope3: 1750, target: 3100 },
  { month: 'Sep', scope1: 900, scope2: 620, scope3: 1700, target: 3000 },
  { month: 'Oct', scope1: 880, scope2: 600, scope3: 1650, target: 2900 },
  { month: 'Nov', scope1: 850, scope2: 580, scope3: 1600, target: 2800 },
  { month: 'Dec', scope1: 820, scope2: 550, scope3: 1550, target: 2700 },
];

const energyData = [
  { source: 'Solar', value: 35, color: '#f59e0b' },
  { source: 'Wind', value: 25, color: '#3b82f6' },
  { source: 'Hydro', value: 15, color: '#06b6d4' },
  { source: 'Grid (Green)', value: 12, color: '#22c55e' },
  { source: 'Grid (Standard)', value: 8, color: '#ef4444' },
  { source: 'Other', value: 5, color: '#9ca3af' },
];

const projectsData = [
  { id: 'PRJ-001', name: 'Solar Farm Phase 2', type: 'Mitigation', status: 'Active', progress: 65, emissionsReduction: 1200, location: 'Arizona, USA' },
  { id: 'PRJ-002', name: 'Coastal Restoration', type: 'Adaptation', status: 'Planning', progress: 15, emissionsReduction: 0, location: 'Louisiana, USA' },
  { id: 'PRJ-003', name: 'Energy Efficiency Retrofit', type: 'Mitigation', status: 'Active', progress: 80, emissionsReduction: 850, location: 'Texas, USA' },
  { id: 'PRJ-004', name: 'Reforestation Program', type: 'Mitigation', status: 'Active', progress: 45, emissionsReduction: 2100, location: 'Oregon, USA' },
  { id: 'PRJ-005', name: 'Flood Defense System', type: 'Adaptation', status: 'Completed', progress: 100, emissionsReduction: 0, location: 'Florida, USA' },
];

const alertsData = [
  { id: 'ALT-001', severity: 'critical', title: 'Scope 1 Emissions Spike', message: 'Facility TX-03 exceeded threshold by 15%', time: '15 min ago' },
  { id: 'ALT-002', severity: 'warning', title: 'Water Stress Alert', message: 'Regional water stress index reached 4.2/5', time: '1 hour ago' },
  { id: 'ALT-003', severity: 'info', title: 'New Regulation Published', message: 'SEC Climate Disclosure Rules finalized', time: '3 hours ago' },
  { id: 'ALT-004', severity: 'warning', title: 'Project Delay Risk', message: 'Solar Farm Phase 2 behind schedule by 2 weeks', time: '5 hours ago' },
];

const tabs: TabItem[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'emissions', label: 'Emissions' },
  { id: 'energy', label: 'Energy' },
  { id: 'projects', label: 'Projects' },
];

const timeRanges = [
  { value: '7d', label: 'Last 7 Days' },
  { value: '30d', label: 'Last 30 Days' },
  { value: '90d', label: 'Last 90 Days' },
  { value: '1y', label: 'Last Year' },
];

export default function DashboardPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [timeRange, setTimeRange] = useState('30d');

  const renderMetricCard = (metric: typeof metricCards[0]) => (
    <Card className="relative overflow-hidden">
      <CardBody className="p-6">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-sm text-secondary-500 dark:text-secondary-400">{metric.label}</p>
            <p className="mt-1 text-3xl font-bold text-secondary-900 dark:text-white">{metric.value}</p>
            <p className="text-sm text-secondary-500 dark:text-secondary-400">{metric.subtitle}</p>
          </div>
          <div className={cn('p-3 rounded-xl', metric.color + '/10')}>
            <metric.icon className={cn('h-6 w-6', metric.color.replace('bg-', 'text-'))} />
          </div>
        </div>
        <div className="mt-4 flex items-center gap-2">
          <span className={cn('text-sm font-medium', metric.trendUp ? 'text-green-600' : 'text-red-600')}>
            {metric.trendUp ? <TrendingUpIcon className="h-4 w-4 inline" /> : <TrendingDownIcon className="h-4 w-4 inline" />}
            {formatPercentage(metric.trend)}
          </span>
          <span className="text-xs text-secondary-500 dark:text-secondary-400">vs last period</span>
        </div>
      </CardBody>
    </Card>
  );

  const renderEmissionsChart = () => (
    <ResponsiveContainer width="100%" height={300}>
      <ComposedChart data={emissionsData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey="month" stroke="#9ca3af" fontSize={12} tickLine={false} axisLine={false} />
        <YAxis stroke="#9ca3af" fontSize={12} tickLine={false} axisLine={false} tickFormatter={(v) => `${v/1000}k`} />
        <Tooltip
          formatter={(value: number) => [formatNumber(value), 'tCO₂e']}
          labelFormatter={(label) => `Month: ${label}`}
          contentStyle={{ backgroundColor: '#fff', border: '1px solid #e5e7eb', borderRadius: '8px' }}
        />
        <Legend />
        <Area type="monotone" dataKey="scope1" name="Scope 1" stroke="#ef4444" fill="#ef4444" fillOpacity={0.1} />
        <Area type="monotone" dataKey="scope2" name="Scope 2" stroke="#f59e0b" fill="#f59e0b" fillOpacity={0.1} />
        <Area type="monotone" dataKey="scope3" name="Scope 3" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.1} />
        <Line type="monotone" dataKey="target" name="Target" stroke="#22c55e" strokeDasharray="5 5" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ResponsiveContainer>
  );

  const renderEnergyChart = () => (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={energyData}
          cx="50%"
          cy="50%"
          innerRadius={60}
          outerRadius={100}
          paddingAngle={2}
          dataKey="value"
          nameKey="source"
          label={({ source, percent }) => `${source} ${(percent * 100).toFixed(0)}%`}
          labelLine={false}
        >
          {energyData.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={entry.color} />
          ))}
        </Pie>
        <Tooltip formatter={(value: number) => [`${value}%`, 'Share']} />
      </PieChart>
    </ResponsiveContainer>
  );

  const renderAlerts = () => (
    <div className="space-y-3">
      {alertsData.map((alert) => (
        <div
          key={alert.id}
          className={cn(
            'flex items-start gap-3 p-4 rounded-lg border',
            alert.severity === 'critical' && 'border-red-200 bg-red-50 dark:bg-red-900/20',
            alert.severity === 'warning' && 'border-yellow-200 bg-yellow-50 dark:bg-yellow-900/20',
            alert.severity === 'info' && 'border-blue-200 bg-blue-50 dark:bg-blue-900/20'
          )}
        >
          <div className={cn('flex-shrink-0 h-2 w-2 rounded-full mt-2', alert.severity === 'critical' && 'bg-red-500', alert.severity === 'warning' && 'bg-yellow-500', alert.severity === 'info' && 'bg-blue-500')} />
          <div className="flex-1 min-w-0">
            <p className="font-medium text-secondary-900 dark:text-white">{alert.title}</p>
            <p className="text-sm text-secondary-600 dark:text-secondary-400 mt-0.5">{alert.message}</p>
            <p className="text-xs text-secondary-400 dark:text-secondary-500 mt-1">{alert.time}</p>
          </div>
          <Badge variant={alert.severity === 'critical' ? 'danger' : alert.severity === 'warning' ? 'warning' : 'info'} size="sm">
            {alert.severity}
          </Badge>
        </div>
      ))}
    </div>
  );

  const renderProjectsTable = () => (
    <div className="space-y-3">
      {projectsData.map((project) => (
        <div key={project.id} className="flex items-center gap-4 p-4 rounded-lg border border-secondary-200 dark:border-secondary-700 hover:bg-secondary-50 dark:hover:bg-secondary-800/50">
          <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
            <LeafIcon className="h-5 w-5 text-primary-600 dark:text-primary-400" />
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-medium text-secondary-900 dark:text-white truncate">{project.name}</p>
            <p className="text-sm text-secondary-500 dark:text-secondary-400">{project.location} • {project.type}</p>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-24">
              <div className="h-2 bg-secondary-200 dark:bg-secondary-700 rounded-full overflow-hidden">
                <div className={cn('h-full rounded-full transition-all', project.status === 'Completed' ? 'bg-green-500' : project.status === 'Active' ? 'bg-blue-500' : 'bg-yellow-500')} style={{ width: `${project.progress}%` }} />
              </div>
            </div>
            <Badge variant={project.status === 'Completed' ? 'success' : project.status === 'Active' ? 'primary' : project.status === 'Planning' ? 'warning' : 'secondary'} size="sm">
              {project.status}
            </Badge>
            <span className="text-sm font-mono text-secondary-600 dark:text-secondary-400">
              -{formatNumber(project.emissionsReduction)} tCO₂e
            </span>
          </div>
        </div>
      ))}
    </div>
  );

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <Breadcrumb items={[
            { label: 'Dashboard', href: '/dashboard', current: true },
          ]} />
          <h1 className="mt-2 text-2xl font-bold text-secondary-900 dark:text-white">Dashboard</h1>
          <p className="text-secondary-600 dark:text-secondary-400">Overview of your climate action intelligence</p>
        </div>
        <div className="flex items-center gap-3">
          <Select options={timeRanges} value={timeRange} onChange={(e) => setTimeRange(e.target.value)} className="w-40" />
          <Button variant="outline" size="sm">
            <RefreshCwIcon className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" size="sm">
            <DownloadIcon className="h-4 w-4 mr-2" />
            Export
          </Button>
        </div>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      {activeTab === 'overview' && (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
            {metricCards.map((metric, index) => (
              <div key={index}>{renderMetricCard(metric)}</div>
            ))}
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Emissions Trend</CardTitle>
                  <CardDescription>Monthly Scope 1, 2 & 3 emissions vs target</CardDescription>
                </div>
                <Badge variant="success" size="sm">On Track</Badge>
              </CardHeader>
              <CardBody>{renderEmissionsChart()}</CardBody>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Energy Mix</CardTitle>
                  <CardDescription>Renewable vs non-renewable sources</CardDescription>
                </div>
                <Badge variant="info" size="sm">67% Renewable</Badge>
              </CardHeader>
              <CardBody>{renderEnergyChart()}</CardBody>
            </Card>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Active Alerts</CardTitle>
                  <CardDescription>Critical and warning notifications</CardDescription>
                </div>
                <Badge variant="danger" size="sm">4 Active</Badge>
              </CardHeader>
              <CardBody>{renderAlerts()}</CardBody>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Recent Projects</CardTitle>
                  <CardDescription>Track progress on climate initiatives</CardDescription>
                </div>
                <Button variant="ghost" size="sm" asChild>
                  <a href="/projects">View All <ArrowUpRightIcon className="h-3 w-3 ml-1" /></a>
                </Button>
              </CardHeader>
              <CardBody>{renderProjectsTable()}</CardBody>
            </Card>
          </div>
        </>
      )}

      {activeTab === 'emissions' && (
        <Card>
          <CardHeader>
            <CardTitle>Emissions Deep Dive</CardTitle>
            <CardDescription>Detailed breakdown of emissions by scope, category, and facility</CardDescription>
          </CardHeader>
          <CardBody>
            <div className="grid gap-6 lg:grid-cols-3">
              <Card>
                <CardBody className="text-center">
                  <p className="text-4xl font-bold text-red-600 dark:text-red-400">4,100</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Scope 1 (Direct)</p>
                </CardBody>
              </Card>
              <Card>
                <CardBody className="text-center">
                  <p className="text-4xl font-bold text-orange-600 dark:text-orange-400">2,800</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Scope 2 (Indirect)</p>
                </CardBody>
              </Card>
              <Card>
                <CardBody className="text-center">
                  <p className="text-4xl font-bold text-blue-600 dark:text-blue-400">5,550</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">Scope 3 (Value Chain)</p>
                </CardBody>
              </Card>
            </div>
            <div className="mt-6">
              <ResponsiveContainer width="100%" height={400}>
                <BarChart data={emissionsData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Bar dataKey="scope1" name="Scope 1" fill="#ef4444" radius={[4, 4, 0, 0]} />
                  <Bar dataKey="scope2" name="Scope 2" fill="#f59e0b" radius={[4, 4, 0, 0]} />
                  <Bar dataKey="scope3" name="Scope 3" fill="#3b82f6" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardBody>
        </Card>
      )}

      {activeTab === 'energy' && (
        <div className="grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Consumption Trend</CardTitle>
              <CardDescription>Monthly energy consumption by source</CardDescription>
            </CardHeader>
            <CardBody>
              <ResponsiveContainer width="100%" height={400}>
                <AreaChart data={emissionsData.map((d, i) => ({ ...d, total: d.scope1 + d.scope2 + d.scope3, solar: 350 + i * 10, wind: 250 + i * 5 }))} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Area type="monotone" dataKey="solar" name="Solar" stroke="#f59e0b" fill="#f59e0b" fillOpacity={0.3} />
                  <Area type="monotone" dataKey="wind" name="Wind" stroke="#3b82f6" fill="#3b82f6" fillOpacity={0.3} />
                </AreaChart>
              </ResponsiveContainer>
            </CardBody>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Efficiency Metrics</CardTitle>
              <CardDescription>Key performance indicators</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              {[
                { label: 'Energy Intensity', value: '0.45', unit: 'MWh/$M revenue', trend: -5.2, good: true },
                { label: 'Renewable Ratio', value: '67%', unit: 'of total consumption', trend: 12.5, good: true },
                { label: 'Peak Demand', value: '12.4', unit: 'MW', trend: 3.1, good: false },
                { label: 'Cost per MWh', value: '$42', unit: 'weighted average', trend: -2.8, good: true },
              ].map((metric) => (
                <div key={metric.label} className="flex items-center justify-between p-4 bg-secondary-50 dark:bg-secondary-800/50 rounded-lg">
                  <div>
                    <p className="font-medium text-secondary-900 dark:text-white">{metric.label}</p>
                    <p className="text-sm text-secondary-500 dark:text-secondary-400">{metric.unit}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-bold text-secondary-900 dark:text-white">{metric.value}</p>
                    <span className={cn('text-sm font-medium', metric.good ? 'text-green-600' : 'text-red-600')}>
                      {metric.trend >= 0 ? '+' : ''}{metric.trend}%
                    </span>
                  </div>
                </div>
              ))}
            </CardBody>
          </Card>
        </div>
      )}

      {activeTab === 'projects' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>All Projects</CardTitle>
              <CardDescription>Track progress on all climate initiatives</CardDescription>
            </div>
            <Button variant="primary" asChild>
              <a href="/projects/new">New Project</a>
            </Button>
          </CardHeader>
          <CardBody>{renderProjectsTable()}</CardBody>
        </Card>
      )}
    </div>
  );
}

import { ShieldCheckIcon } from 'lucide-react';
import { Skeleton } from '../../components/ui/Skeleton';