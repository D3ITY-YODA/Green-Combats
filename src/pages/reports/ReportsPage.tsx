import { useState } from 'react';
import { FileTextIcon, DownloadIcon, CalendarIcon, ClockIcon, Trash2Icon, EditIcon, EyeIcon, PlusIcon, FilterIcon, SearchIcon } from 'lucide-react';
import { formatDate, formatRelativeTime } from '../../utils/helpers';
import { cn } from '../../utils/helpers';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Badge } from '../../components/ui/Badge';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Input';
import { Table } from '../../components/ui/Table';
import { Pagination } from '../../components/ui/Pagination';
import { Breadcrumb } from '../../components/ui/Breadcrumb';
import { Modal, ConfirmModal } from '../../components/ui/Modal';
import { Dropdown, DropdownItem } from '../../components/ui/Dropdown';

const reports = [
  { id: 'RPT-001', name: 'Monthly Emissions Report - Dec 2024', type: 'Scheduled', format: 'PDF', status: 'Completed', created: '2024-12-01T08:00:00Z', size: '2.4 MB', pages: 15 },
  { id: 'RPT-002', name: 'Quarterly Carbon Footprint - Q4 2024', type: 'Scheduled', format: 'PDF', status: 'Completed', created: '2024-12-15T09:00:00Z', size: '5.1 MB', pages: 32 },
  { id: 'RPT-003', name: 'Annual Sustainability Report 2024', type: 'Ad-hoc', format: 'PDF', status: 'Generating', created: '2024-12-20T10:30:00Z', size: '-', pages: '-' },
  { id: 'RPT-004', name: 'Facility Energy Audit - Texas Plant', type: 'Ad-hoc', format: 'Excel', status: 'Completed', created: '2024-12-18T14:00:00Z', size: '1.8 MB', pages: 8 },
  { id: 'RPT-005', name: 'Water Risk Assessment', type: 'Template', format: 'PDF', status: 'Draft', created: '2024-12-10T11:00:00Z', size: '-', pages: '-' },
  { id: 'RPT-006', name: 'Scope 3 Category Analysis', type: 'Ad-hoc', format: 'CSV', status: 'Failed', created: '2024-12-12T16:00:00Z', size: '-', pages: '-' },
  { id: 'RPT-007', name: 'Renewable Energy Progress', type: 'Scheduled', format: 'PDF', status: 'Completed', created: '2024-12-01T08:00:00Z', size: '3.2 MB', pages: 12 },
  { id: 'RPT-008', name: 'Compliance Gap Analysis', type: 'Ad-hoc', format: 'PDF', status: 'Completed', created: '2024-11-28T10:00:00Z', size: '4.5 MB', pages: 24 },
  { id: 'RPT-009', name: 'Biodiversity Impact Assessment', type: 'Template', format: 'PDF', status: 'Draft', created: '2024-11-25T09:00:00Z', size: '-', pages: '-' },
  { id: 'RPT-010', name: 'Climate Scenario Modeling Results', type: 'Ad-hoc', format: 'Excel', status: 'Completed', created: '2024-11-20T13:00:00Z', size: '2.7 MB', pages: 18 },
];

const scheduledReports = [
  { id: 'SCH-001', name: 'Monthly Emissions Summary', frequency: 'Monthly', day: '1st', time: '08:00', format: 'PDF', recipients: 5, status: 'Active', nextRun: '2025-01-01' },
  { id: 'SCH-002', name: 'Weekly KPI Dashboard', frequency: 'Weekly', day: 'Monday', time: '07:00', format: 'PDF', recipients: 12, status: 'Active', nextRun: '2024-12-23' },
  { id: 'SCH-003', name: 'Quarterly Board Report', frequency: 'Quarterly', day: '1st of Quarter', time: '09:00', format: 'PDF', recipients: 8, status: 'Active', nextRun: '2025-01-01' },
  { id: 'SCH-004', name: 'Annual CDP Disclosure', frequency: 'Annual', day: 'March 31', time: '10:00', format: 'PDF', recipients: 3, status: 'Paused', nextRun: '2025-03-31' },
];

const templates = [
  { id: 'TMP-001', name: 'GHG Inventory Report', description: 'Standard GHG Protocol compliant inventory', category: 'Emissions', lastUpdated: '2024-11-15' },
  { id: 'TMP-002', name: 'TCFD Disclosure Template', description: 'Task Force on Climate-related Financial Disclosures', category: 'Compliance', lastUpdated: '2024-10-20' },
  { id: 'TMP-003', name: 'CDP Climate Change Response', description: 'CDP questionnaire response template', category: 'Compliance', lastUpdated: '2024-11-01' },
  { id: 'TMP-004', name: 'Energy Management Report', description: 'ISO 50001 aligned energy report', category: 'Energy', lastUpdated: '2024-10-10' },
  { id: 'TMP-005', name: 'Water Stewardship Report', description: 'AWS Standard aligned water report', category: 'Water', lastUpdated: '2024-09-20' },
];

const columns = [
  { key: 'name', header: 'Report Name', accessor: 'name' },
  { key: 'type', header: 'Type', accessor: 'type' },
  { key: 'format', header: 'Format', accessor: 'format' },
  { key: 'status', header: 'Status', accessor: 'status', render: (value: string) => <Badge variant={value === 'Completed' ? 'success' : value === 'Generating' ? 'info' : value === 'Failed' ? 'danger' : value === 'Draft' ? 'secondary' : 'primary'} size="sm">{value}</Badge> },
  { key: 'created', header: 'Created', accessor: (row: any) => formatDate(row.created) },
  { key: 'size', header: 'Size', accessor: 'size' },
];

export default function ReportsPage() {
  const [activeTab, setActiveTab] = useState('generated');
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [currentPage, setCurrentPage] = useState(1);
  const [deleteModal, setDeleteModal] = useState<string | null>(null);

  const filteredReports = reports.filter((r) => {
    const matchesSearch = r.name.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesStatus = statusFilter === 'all' || r.status.toLowerCase() === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const handleDelete = (id: string) => {
    setDeleteModal(id);
  };

  const confirmDelete = () => {
    if (deleteModal) {
      // Delete logic here
      toast.success('Report deleted');
      setDeleteModal(null);
    }
  };

  const tabs = [
    { id: 'generated', label: 'Generated Reports', count: reports.length },
    { id: 'scheduled', label: 'Scheduled Reports', count: scheduledReports.length },
    { id: 'templates', label: 'Templates', count: templates.length },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <Breadcrumb items={[
            { label: 'Reports', href: '/reports', current: true },
          ]} />
          <h1 className="mt-2 text-2xl font-bold text-secondary-900 dark:text-white">Reports</h1>
          <p className="text-secondary-600 dark:text-secondary-400">Manage and generate climate reports</p>
        </div>
        <Button variant="primary" asChild>
          <a href="/reports/new">
            <PlusIcon className="h-4 w-4 mr-2" />
            Create Report
          </a>
        </Button>
      </div>

      <Tabs
        tabs={tabs.map(t => ({ id: t.id, label: t.label, count: t.count }))}
        activeTab={activeTab}
        onChange={setActiveTab}
        variant="line"
      />

      {activeTab === 'generated' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between flex-wrap gap-4">
            <div>
              <CardTitle>Generated Reports</CardTitle>
              <CardDescription>View and download your generated reports</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Input
                placeholder="Search reports..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-64"
                iconLeft={<SearchIcon className="h-4 w-4" />}
              />
              <Select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                options={[
                  { value: 'all', label: 'All Status' },
                  { value: 'completed', label: 'Completed' },
                  { value: 'generating', label: 'Generating' },
                  { value: 'failed', label: 'Failed' },
                  { value: 'draft', label: 'Draft' },
                ]}
                className="w-40"
              />
            </div>
          </CardHeader>
          <CardBody className="p-0">
            <Table
              columns={columns}
              data={filteredReports}
              keyExtractor={(row) => row.id}
              pagination={{
                page: currentPage,
                pageSize: 10,
                total: filteredReports.length,
                onPageChange: setCurrentPage,
              }}
              onRowClick={(row) => {}}
            />
          </CardBody>
        </Card>
      )}

      {activeTab === 'scheduled' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Scheduled Reports</CardTitle>
              <CardDescription>Automated report generation schedules</CardDescription>
            </div>
            <Button variant="primary" asChild>
              <a href="/reports/schedule">New Schedule</a>
            </Button>
          </CardHeader>
          <CardBody>
            <div className="space-y-3">
              {scheduledReports.map((report) => (
                <div key={report.id} className="flex items-center justify-between p-4 border border-secondary-200 dark:border-secondary-700 rounded-lg">
                  <div className="flex items-center gap-4">
                    <div className="p-2 bg-secondary-100 dark:bg-secondary-800 rounded-lg">
                      <FileTextIcon className="h-5 w-5 text-secondary-600 dark:text-secondary-400" />
                    </div>
                    <div>
                      <p className="font-medium text-secondary-900 dark:text-white">{report.name}</p>
                      <p className="text-sm text-secondary-500 dark:text-secondary-400">
                        {report.frequency} • {report.day} at {report.time} • {report.format}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <Badge variant={report.status === 'Active' ? 'success' : 'secondary'} size="sm">
                      {report.status}
                    </Badge>
                    <span className="text-sm text-secondary-500 dark:text-secondary-400">
                      Next: {formatDate(report.nextRun)}
                    </span>
                    <div className="flex items-center gap-2">
                      <Button variant="ghost" size="sm" asChild>
                        <a href={`/reports/schedule/${report.id}`}><EditIcon className="h-4 w-4" /></a>
                      </Button>
                      <Button variant="ghost" size="sm" onClick={() => handleDelete(report.id)}>
                        <Trash2Icon className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardBody>
        </Card>
      )}

      {activeTab === 'templates' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Report Templates</CardTitle>
              <CardDescription>Reusable templates for consistent reporting</CardDescription>
            </div>
            <Button variant="primary" asChild>
              <a href="/reports/templates/new">Create Template</a>
            </Button>
          </CardHeader>
          <CardBody>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {templates.map((template) => (
                <Card key={template.id} className="hover:shadow-md transition-shadow">
                  <CardBody className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="p-2 bg-primary-100 dark:bg-primary-900/30 rounded-lg">
                        <FileTextIcon className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                      </div>
                      <Badge variant="outline" size="sm">{template.category}</Badge>
                    </div>
                    <p className="mt-4 font-medium text-secondary-900 dark:text-white">{template.name}</p>
                    <p className="mt-1 text-sm text-secondary-500 dark:text-secondary-400">{template.description}</p>
                    <div className="mt-4 flex items-center justify-between">
                      <span className="text-xs text-secondary-400 dark:text-secondary-500">
                        Updated {formatRelativeTime(template.lastUpdated)}
                      </span>
                      <div className="flex gap-2">
                        <Button variant="ghost" size="sm" asChild>
                          <a href={`/reports/templates/${template.id}`}><EyeIcon className="h-4 w-4" /></a>
                        </Button>
                        <Button variant="ghost" size="sm" asChild>
                          <a href={`/reports/new?template=${template.id}`}><PlusIcon className="h-4 w-4" /></a>
                        </Button>
                      </div>
                    </div>
                  </CardBody>
                </Card>
              ))}
            </div>
          </CardBody>
        </Card>
      )}

      <ConfirmModal
        isOpen={!!deleteModal}
        onClose={() => setDeleteModal(null)}
        onConfirm={confirmDelete}
        title="Delete Report"
        message="Are you sure you want to delete this report? This action cannot be undone."
        confirmText="Delete"
        variant="danger"
      />
    </div>
  );
}

import { toast } from 'react-hot-toast';