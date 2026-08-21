import { LineChart, Line, BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, AreaChart, Area } from 'recharts';
import { ChevronDownIcon, DownloadIcon, FilterIcon, CalendarIcon, ArrowUpRightIcon } from 'lucide-react';
import { formatNumber, formatPercentage } from '../../utils/helpers';
import { cn } from '../../utils/helpers';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Select } from '../../components/ui/Input';
import { Tabs, TabItem } from '../../components/ui/Tabs';
import { Breadcrumb } from '../../components/ui/Breadcrumb';
import { Dropdown, DropdownItem } from '../../components/ui/Dropdown';

const tabs: TabItem[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'trends', label: 'Trends' },
  { id: 'comparisons', label: 'Comparisons' },
  { id: 'forecasting', label: 'Forecasting' },
  { id: 'custom', label: 'Custom Reports' },
];

const kpiData = [
  { label: 'Total Emissions', value: '12,450 tCO₂e', trend: -2.3, good: true, benchmark: '11,200 target' },
  { label: 'Carbon Intensity', value: '45.2 tCO₂e/$M', trend: -5.1, good: true, benchmark: '40.0 target' },
  { label: 'Energy Efficiency', value: '0.45 MWh/$M', trend: -3.8, good: true, benchmark: '0.40 target' },
  { label: 'Renewable %', value: '67%', trend: 12.5, good: true, benchmark: '80% target' },
  { label: 'Water Intensity', value: '2.4 m³/$M', trend: 1.2, good: false, benchmark: '2.0 target' },
  { label: 'Waste Diversion', value: '78%', trend: 4.5, good: true, benchmark: '90% target' },
];

const trendData = [
  { period: 'Jan', emissions: 4200, intensity: 48.5, energy: 1250, water: 210000, waste: 850, targetEmissions: 4000 },
  { period: 'Feb', emissions: 4100, intensity: 47.8, energy: 1220, water: 205000, waste: 830, targetEmissions: 3950 },
  { period: 'Mar', emissions: 4000, intensity: 46.9, energy: 1180, water: 198000, waste: 810, targetEmissions: 3900 },
  { period: 'Apr', emissions: 3900, intensity: 46.2, energy: 1150, water: 192000, waste: 790, targetEmissions: 3850 },
  { period: 'May', emissions: 3850, intensity: 45.8, energy: 1130, water: 188000, waste: 770, targetEmissions: 3800 },
  { period: 'Jun', emissions: 3780, intensity: 45.2, energy: 1100, water: 185000, waste: 750, targetEmissions: 3750 },
  { period: 'Jul', emissions: 3720, intensity: 44.9, energy: 1080, water: 182000, waste: 730, targetEmissions: 3700 },
  { period: 'Aug', emissions: 3650, intensity: 44.3, energy: 1050, water: 178000, waste: 710, targetEmissions: 3650 },
  { period: 'Sep', emissions: 3580, intensity: 43.8, energy: 1030, water: 175000, waste: 690, targetEmissions: 3600 },
  { period: 'Oct', emissions: 3520, intensity: 43.2, energy: 1010, water: 172000, waste: 670, targetEmissions: 3550 },
  { period: 'Nov', emissions: 3450, intensity: 42.8, energy: 990, water: 168000, waste: 650, targetEmissions: 3500 },
  { period: 'Dec', emissions: 3400, intensity: 42.3, energy: 970, water: 165000, waste: 630, targetEmissions: 3450 },
];

const facilityData = [
  { facility: 'HQ - San Francisco', emissions: 1200, intensity: 35.2, energy: 3200, water: 45000, waste: 85, status: 'On Track' },
  { facility: 'Plant - Texas', emissions: 3200, intensity: 52.8, energy: 8500, water: 120000, waste: 210, status: 'At Risk' },
  { facility: 'Plant - Arizona', emissions: 2800, intensity: 48.5, energy: 7200, water: 98000, waste: 180, status: 'On Track' },
  { facility: 'DC - Virginia', emissions: 1500, intensity: 42.1, energy: 4100, water: 55000, waste: 95, status: 'On Track' },
  { facility: 'Office - London', emissions: 800, intensity: 28.4, energy: 2200, water: 30000, waste: 45, status: 'Ahead' },
  { facility: 'Office - Singapore', emissions: 600, intensity: 31.2, energy: 1800, water: 25000, waste: 38, status: 'On Track' },
  { facility: 'Plant - Germany', emissions: 1800, intensity: 41.7, energy: 4800, water: 65000, waste: 120, status: 'At Risk' },
  { facility: 'Warehouse - Japan', emissions: 550, intensity: 33.9, energy: 1500, water: 22000, waste: 32, status: 'Ahead' },
];

const forecastData = [
  { period: '2024 Q1', actual: 3400, forecast: 3450, lower: 3300, upper: 3600 },
  { period: '2024 Q2', actual: null, forecast: 3350, lower: 3150, upper: 3550 },
  { period: '2024 Q3', actual: null, forecast: 3280, lower: 3050, upper: 3500 },
  { period: '2024 Q4', actual: null, forecast: 3200, lower: 2950, upper: 3450 },
  { period: '2025 Q1', actual: null, forecast: 3120, lower: 2850, upper: 3400 },
  { period: '2025 Q2', actual: null, forecast: 3050, lower: 2750, upper: 3350 },
];

export default function AnalyticsPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [timeRange, setTimeRange] = useState('12m');

  const renderKPIs = () => (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6">
      {kpiData.map((kpi, index) => (
        <Card key={index}>
          <CardBody className="p-6">
            <p className="text-sm text-secondary-500 dark:text-secondary-400">{kpi.label}</p>
            <p className="mt-1 text-3xl font-bold text-secondary-900 dark:text-white">{kpi.value}</p>
            <div className="mt-3 flex items-center justify-between">
              <span className={cn('text-sm font-medium', kpi.good ? 'text-green-600' : 'text-red-600')}>
                {kpi.trend >= 0 ? '+' : ''}{formatPercentage(kpi.trend)}
              </span>
              <span className="text-xs text-secondary-400 dark:text-secondary-500">{kpi.benchmark}</span>
            </div>
          </CardBody>
        </Card>
      ))}
    </div>
  );

  const renderTrendChart = () => (
    <ResponsiveContainer width="100%" height={350}>
      <ComposedChart data={trendData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="period" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="emissions" name="Emissions (tCO₂e)" stroke="#ef4444" strokeWidth={2} dot={false} />
        <Line type="monotone" dataKey="targetEmissions" name="Target" stroke="#22c55e" strokeDasharray="5 5" strokeWidth={2} dot={false} />
        <Line type="monotone" dataKey="energy" name="Energy (MWh)" stroke="#f59e0b" strokeWidth={2} dot={false} yAxisId="right" />
      </ComposedChart>
    </ResponsiveContainer>
  );

  const renderFacilityComparison = () => (
    <div className="space-y-3">
      {facilityData.map((facility, index) => (
        <Card key={index}>
          <CardBody className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className={cn('h-8 w-8 rounded-lg flex items-center justify-center', facility.status === 'Ahead' ? 'bg-green-100' : facility.status === 'On Track' ? 'bg-blue-100' : 'bg-yellow-100')}>
                  <span className="text-xs font-bold text-secondary-700">{facility.facility.split(' ')[0].charAt(0)}</span>
                </div>
                <div>
                  <p className="font-medium text-secondary-900 dark:text-white">{facility.facility}</p>
                  <p className="text-sm text-secondary-500 dark:text-secondary-400">
                    {formatNumber(facility.emissions)} tCO₂e • {formatNumber(facility.energy)} MWh
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <Badge variant={facility.status === 'Ahead' ? 'success' : facility.status === 'On Track' ? 'primary' : 'warning'} size="sm">
                  {facility.status}
                </Badge>
                <div className="text-right">
                  <p className="text-sm font-medium text-secondary-900 dark:text-white">{formatPercentage(facility.intensity)} tCO₂e/$M</p>
                  <p className="text-xs text-secondary-500 dark:text-secondary-400">Intensity</p>
                </div>
              </div>
            </div>
          </CardBody>
        </Card>
      ))}
    </div>
  );

  const renderForecastChart = () => (
    <ResponsiveContainer width="100%" height={350}>
      <ComposedChart data={forecastData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="period" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Area type="monotone" dataKey="upper" name="Upper Bound" stroke="#22c55e" fill="#22c55e" fillOpacity={0.1} />
        <Area type="monotone" dataKey="lower" name="Lower Bound" stroke="#22c55e" fill="#fff" fillOpacity={1} />
        <Line type="monotone" dataKey="forecast" name="Forecast" stroke="#3b82f6" strokeWidth={2} strokeDasharray="5 5" dot={false} />
        <Line type="monotone" dataKey="actual" name="Actual" stroke="#ef4444" strokeWidth={2} dot={false} />
      </ComposedChart>
    </ResponsiveContainer>
  );

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <Breadcrumb items={[
            { label: 'Analytics', href: '/analytics', current: true },
          ]} />
          <h1 className="mt-2 text-2xl font-bold text-secondary-900 dark:text-white">Analytics</h1>
          <p className="text-secondary-600 dark:text-secondary-400">Deep insights into your climate performance</p>
        </div>
        <div className="flex items-center gap-3">
          <Select
            options={[
              { value: '1m', label: '1 Month' },
              { value: '3m', label: '3 Months' },
              { value: '6m', label: '6 Months' },
              { value: '12m', label: '12 Months' },
              { value: 'ytd', label: 'YTD' },
            ]}
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="w-36"
          />
          <Button variant="outline">
            <DownloadIcon className="h-4 w-4 mr-2" />
            Export
          </Button>
        </div>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      {activeTab === 'overview' && (
        <>
          {renderKPIs()}
          <div className="grid gap-6 lg:grid-cols-2 mt-6">
            <Card>
              <CardHeader>
                <CardTitle>Emissions Trend</CardTitle>
                <CardDescription>12-month rolling emissions with target</CardDescription>
              </CardHeader>
              <CardBody>{renderTrendChart()}</CardBody>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Key Metrics Trend</CardTitle>
                <CardDescription>Multiple indicators over time</CardDescription>
              </CardHeader>
              <CardBody>
                <ResponsiveContainer width="100%" height={350}>
                  <LineChart data={trendData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" />
                    <YAxis yAxisId="left" />
                    <YAxis yAxisId="right" orientation="right" />
                    <Tooltip />
                    <Legend />
                    <Line yAxisId="left" type="monotone" dataKey="emissions" name="Emissions" stroke="#ef4444" />
                    <Line yAxisId="left" type="monotone" dataKey="intensity" name="Intensity" stroke="#f59e0b" />
                    <Line yAxisId="right" type="monotone" dataKey="energy" name="Energy" stroke="#3b82f6" />
                  </LineChart>
                </ResponsiveContainer>
              </CardBody>
            </Card>
          </div>
        </>
      )}

      {activeTab === 'trends' && (
        <>
          <div className="grid gap-6 lg:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Emissions Decomposition</CardTitle>
                <CardDescription>Trend, seasonal, and residual components</CardDescription>
              </CardHeader>
              <CardBody>
                <ResponsiveContainer width="100%" height={400}>
                  <AreaChart data={trendData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Area type="monotone" dataKey="emissions" name="Actual" stroke="#ef4444" fill="#ef4444" fillOpacity={0.3} />
                    <Line type="monotone" dataKey="targetEmissions" name="Target" stroke="#22c55e" strokeDasharray="5 5" />
                  </AreaChart>
                </ResponsiveContainer>
              </CardBody>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>Year-over-Year Comparison</CardTitle>
                <CardDescription>Current vs previous year performance</CardDescription>
              </CardHeader>
              <CardBody>
                <ResponsiveContainer width="100%" height={400}>
                  <BarChart data={trendData.slice(0, 12)} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="period" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="emissions" name="2024" fill="#ef4444" radius={[4, 4, 0, 0]} />
                    <Bar dataKey="targetEmissions" name="2023" fill="#9ca3af" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </CardBody>
            </Card>
          </div>
        </>
      )}

      {activeTab === 'comparisons' && (
        <Card>
          <CardHeader>
            <CardTitle>Facility Benchmarking</CardTitle>
            <CardDescription>Compare performance across all facilities</CardDescription>
          </CardHeader>
          <CardBody>{renderFacilityComparison()}</CardBody>
        </Card>
      )}

      {activeTab === 'forecasting' && (
        <div className="grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Emissions Forecast</CardTitle>
              <CardDescription>Predicted trajectory with confidence intervals</CardDescription>
            </CardHeader>
            <CardBody>{renderForecastChart()}</CardBody>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>Scenario Analysis</CardTitle>
              <CardDescription>Projected outcomes under different scenarios</CardDescription>
            </CardHeader>
            <CardBody className="space-y-4">
              {[
                { name: 'Business as Usual', endValue: 3800, color: '#ef4444' },
                { name: 'Current Policies', endValue: 3200, color: '#f59e0b' },
                { name: 'Net Zero 2050', endValue: 1800, color: '#22c55e' },
                { name: '1.5°C Aligned', endValue: 1200, color: '#06b6d4' },
              ].map((scenario) => (
                <div key={scenario.name} className="p-4 bg-secondary-50 dark:bg-secondary-800/50 rounded-lg flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="h-3 w-3 rounded-full" style={{ backgroundColor: scenario.color }} />
                    <span className="font-medium text-secondary-900 dark:text-white">{scenario.name}</span>
                  </div>
                  <div className="text-right">
                    <p className="text-2xl font-bold text-secondary-900 dark:text-white">{formatNumber(scenario.endValue)}</p>
                    <p className="text-sm text-secondary-500 dark:text-secondary-400">tCO₂e by 2030</p>
                  </div>
                </div>
              ))}
            </CardBody>
          </Card>
        </div>
      )}

      {activeTab === 'custom' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Custom Report Builder</CardTitle>
              <CardDescription>Create tailored analytics reports</CardDescription>
            </div>
            <Button variant="primary">Create Report</Button>
          </CardHeader>
          <CardBody className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-2">
              <Card>
                <CardHeader>
                  <CardTitle>Available Metrics</CardTitle>
                </CardHeader>
                <CardBody className="space-y-2 max-h-96 overflow-y-auto">
                  {[
                    'Total Emissions (Scope 1+2+3)',
                    'Scope 1 Emissions',
                    'Scope 2 Emissions (Location-based)',
                    'Scope 2 Emissions (Market-based)',
                    'Scope 3 Emissions (15 categories)',
                    'Carbon Intensity (tCO₂e/$M)',
                    'Energy Consumption (MWh)',
                    'Renewable Energy %',
                    'Water Withdrawal (m³)',
                    'Water Consumption (m³)',
                    'Waste Generated (tonnes)',
                    'Waste Diversion Rate (%)',
                    'Air Quality Index',
                    'Biodiversity Index',
                    'Compliance Score',
                  ].map((metric) => (
                    <label key={metric} className="flex items-center gap-2 p-2 hover:bg-secondary-100 dark:hover:bg-secondary-800 rounded cursor-pointer">
                      <input type="checkbox" className="h-4 w-4 rounded border-secondary-300 text-primary-600" />
                      <span className="text-sm text-secondary-700 dark:text-secondary-300">{metric}</span>
                    </label>
                  ))}
                </CardBody>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Report Configuration</CardTitle>
                </CardHeader>
                <CardBody className="space-y-4">
                  <div>
                    <label className="label">Report Name</label>
                    <Input placeholder="My Custom Report" />
                  </div>
                  <div>
                    <label className="label">Frequency</label>
                    <Select options={[
                      { value: 'manual', label: 'Manual' },
                      { value: 'daily', label: 'Daily' },
                      { value: 'weekly', label: 'Weekly' },
                      { value: 'monthly', label: 'Monthly' },
                      { value: 'quarterly', label: 'Quarterly' },
                    ]} />
                  </div>
                  <div>
                    <label className="label">Format</label>
                    <Select options={[
                      { value: 'pdf', label: 'PDF' },
                      { value: 'excel', label: 'Excel' },
                      { value: 'csv', label: 'CSV' },
                      { value: 'dashboard', label: 'Interactive Dashboard' },
                    ]} />
                  </div>
                  <div>
                    <label className="label">Recipients</label>
                    <Input placeholder="email@domain.com" />
                  </div>
                  <Button variant="primary" className="w-full">Save & Generate</Button>
                </CardBody>
              </Card>
            </div>
          </CardBody>
        </Card>
      )}
    </div>
  );
}

import { useState } from 'react';