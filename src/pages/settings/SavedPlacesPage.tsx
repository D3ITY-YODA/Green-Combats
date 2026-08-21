import { useState } from 'react';
import { MapPinIcon, PlusIcon, EditIcon, Trash2Icon, EyeIcon, GlobeIcon, Building2Icon, HomeIcon, SearchIcon, FilterIcon, ChevronDownIcon, DownloadIcon, Share2Icon, MapIcon } from 'lucide-react';
import { cn } from '../../utils/helpers';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { Select } from '../../components/ui/Input';
import { Card, CardHeader, CardBody, CardTitle, CardDescription } from '../../components/ui/Card';
import { Badge } from '../../components/ui/Badge';
import { Table } from '../../components/ui/Table';
import { Modal, ConfirmModal } from '../../components/ui/Modal';
import { Tabs, TabItem } from '../../components/ui/Tabs';
import { Breadcrumb } from '../../components/ui/Breadcrumb';
import { Dropdown, DropdownItem } from '../../components/ui/Dropdown';
import { Pagination } from '../../components/ui/Pagination';

const savedPlaces = [
  { id: 'loc-001', name: 'HQ - San Francisco', type: 'facility', address: '100 Market St, San Francisco, CA 94102', coordinates: [37.7946, -122.3999], emissions: 1200, energy: 3200, water: 45000, status: 'active', tags: ['headquarters', 'office'], created: '2023-01-15' },
  { id: 'loc-002', name: 'Manufacturing Plant - Texas', type: 'facility', address: '500 Industrial Blvd, Houston, TX 77002', coordinates: [29.7604, -95.3698], emissions: 3200, energy: 8500, water: 120000, status: 'active', tags: ['manufacturing', 'scope1'], created: '2023-02-20' },
  { id: 'loc-003', name: 'Solar Farm - Arizona', type: 'project', address: 'Maricopa County, AZ', coordinates: [33.4484, -112.0740], emissions: -850, energy: -2100, water: 5000, status: 'active', tags: ['renewable', 'solar'], created: '2023-06-10' },
  { id: 'loc-004', name: 'Data Center - Virginia', type: 'facility', address: '21730 Red Rum Dr, Ashburn, VA 20147', coordinates: [39.0438, -77.4874], emissions: 1500, energy: 4100, water: 55000, status: 'active', tags: ['data-center', 'scope2'], created: '2023-03-05' },
  { id: 'loc-005', name: 'London Office', type: 'facility', address: '10 Downing St, London SW1A 2AA, UK', coordinates: [51.5074, -0.1278], emissions: 800, energy: 2200, water: 30000, status: 'active', tags: ['office', 'international'], created: '2023-04-12' },
  { id: 'loc-006', name: 'Wind Farm - North Sea', type: 'project', address: 'North Sea, UK Waters', coordinates: [55.0, 2.0], emissions: -1200, energy: -3500, water: 0, status: 'planning', tags: ['renewable', 'wind', 'offshore'], created: '2023-08-22' },
  { id: 'loc-007', name: 'Singapore Office', type: 'facility', address: '8 Marina Blvd, Singapore 018981', coordinates: [1.2966, 103.8517], emissions: 600, energy: 1800, water: 25000, status: 'active', tags: ['office', 'apac'], created: '2023-05-18' },
  { id: 'loc-008', name: 'Reforestation - Oregon', type: 'project', address: 'Willamette National Forest, OR', coordinates: [44.5, -122.0], emissions: -2100, energy: 0, water: 15000, status: 'active', tags: ['nature-based', 'forestry'], created: '2023-07-30' },
];

const placeTypes = [
  { value: 'all', label: 'All Types' },
  { value: 'facility', label: 'Facilities' },
  { value: 'project', label: 'Projects' },
];

const statusOptions = [
  { value: 'all', label: 'All Status' },
  { value: 'active', label: 'Active' },
  { value: 'planning', label: 'Planning' },
  { value: 'inactive', label: 'Inactive' },
];

const tabs: TabItem[] = [
  { id: 'list', label: 'List View', icon: <MapIcon className="h-4 w-4" /> },
  { id: 'map', label: 'Map View', icon: <GlobeIcon className="h-4 w-4" /> },
];

export default function SavedPlacesPage() {
  const [activeTab, setActiveTab] = useState('list');
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');
  const [statusFilter, setStatusFilter] = useState('all');
  const [currentPage, setCurrentPage] = useState(1);
  const [deleteModal, setDeleteModal] = useState<string | null>(null);
  const [addPlaceModal, setAddPlaceModal] = useState(false);
  const [editPlace, setEditPlace] = useState<typeof savedPlaces[0] | null>(null);

  const filteredPlaces = savedPlaces.filter((place) => {
    const matchesSearch = place.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                         place.address.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesType = typeFilter === 'all' || place.type === typeFilter;
    const matchesStatus = statusFilter === 'all' || place.status === statusFilter;
    return matchesSearch && matchesType && matchesStatus;
  });

  const handleDelete = (id: string) => setDeleteModal(id);
  const handleEdit = (place: typeof savedPlaces[0]) => { setEditPlace(place); setAddPlaceModal(true); };
  const handleAdd = () => { setEditPlace(null); setAddPlaceModal(true); };

  const columns = [
    { key: 'name', header: 'Name', accessor: 'name', render: (value: string, row: any) => (
      <div>
        <p className="font-medium text-secondary-900 dark:text-white">{value}</p>
        <p className="text-sm text-secondary-500 dark:text-secondary-400">{row.address}</p>
      </div>
    )},
    { key: 'type', header: 'Type', accessor: 'type', render: (value: string) => (
      <Badge variant={value === 'facility' ? 'primary' : 'success'} size="sm">
        {value === 'facility' ? <Building2Icon className="h-3 w-3 mr-1" /> : <LeafIcon className="h-3 w-3 mr-1" />}
        {value.charAt(0).toUpperCase() + value.slice(1)}
      </Badge>
    )},
    { key: 'coordinates', header: 'Coordinates', accessor: (row: any) => `${row.coordinates[0].toFixed(4)}, ${row.coordinates[1].toFixed(4)}` },
    { key: 'emissions', header: 'Emissions (tCO₂e)', accessor: 'emissions', render: (value: number) => (
      <span className={cn('font-mono font-medium', value < 0 ? 'text-green-600' : 'text-red-600')}>
        {value < 0 ? '' : '+'}{value.toLocaleString()}
      </span>
    )},
    { key: 'energy', header: 'Energy (MWh)', accessor: 'energy', render: (value: number) => (
      <span className={cn('font-mono font-medium', value < 0 ? 'text-green-600' : 'text-secondary-900')}>
        {value < 0 ? '' : '+'}{value.toLocaleString()}
      </span>
    )},
    { key: 'water', header: 'Water (m³)', accessor: 'water', render: (value: number) => <span className="font-mono">{value.toLocaleString()}</span> },
    { key: 'status', header: 'Status', accessor: 'status', render: (value: string) => (
      <Badge variant={value === 'active' ? 'success' : value === 'planning' ? 'warning' : 'secondary'} size="sm">
        {value.charAt(0).toUpperCase() + value.slice(1)}
      </Badge>
    )},
    { key: 'tags', header: 'Tags', accessor: 'tags', render: (value: string[]) => (
      <div className="flex flex-wrap gap-1">
        {value.map((tag, i) => (
          <Badge key={i} variant="outline" size="sm">{tag}</Badge>
        ))}
      </div>
    )},
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <Breadcrumb items={[
            { label: 'Settings', href: '/settings' },
            { label: 'Saved Places', href: '/settings/saved-places', current: true },
          ]} />
          <h1 className="mt-2 text-2xl font-bold text-secondary-900 dark:text-white">Saved Places</h1>
          <p className="text-secondary-600 dark:text-secondary-400">Manage your monitored locations and project sites</p>
        </div>
        <Button variant="primary" onClick={handleAdd}>
          <PlusIcon className="h-4 w-4 mr-2" />
          Add Place
        </Button>
      </div>

      <Tabs tabs={tabs} activeTab={activeTab} onChange={setActiveTab} variant="line" />

      {activeTab === 'list' && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between flex-wrap gap-4">
            <div>
              <CardTitle>Saved Locations</CardTitle>
              <CardDescription>{filteredPlaces.length} places saved</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <div className="relative">
                <SearchIcon className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-secondary-400" />
                <Input
                  placeholder="Search places..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-64 pl-10"
                />
              </div>
              <Select value={typeFilter} onChange={(e) => setTypeFilter(e.target.value)} options={placeTypes} className="w-36" />
              <Select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} options={statusOptions} className="w-36" />
              <Button variant="outline" size="sm">
                <DownloadIcon className="h-4 w-4 mr-2" />
                Export
              </Button>
            </div>
          </CardHeader>
          <CardBody className="p-0">
            <Table
              columns={columns}
              data={filteredPlaces}
              keyExtractor={(row) => row.id}
              pagination={{
                page: currentPage,
                pageSize: 10,
                total: filteredPlaces.length,
                onPageChange: setCurrentPage,
              }}
              onRowClick={(row) => handleEdit(row)}
            />
          </CardBody>
        </Card>
      )}

      {activeTab === 'map' && (
        <Card>
          <CardHeader>
            <CardTitle>Map View</CardTitle>
            <CardDescription>Visualize all saved places on an interactive map</CardDescription>
          </CardHeader>
          <CardBody className="p-0">
            <div className="h-96 bg-secondary-100 dark:bg-secondary-800 relative flex items-center justify-center">
              <div className="text-center text-secondary-500">
                <MapIcon className="h-12 w-12 mx-auto mb-4 opacity-50" />
                <p className="text-lg font-medium text-secondary-900 dark:text-white">Interactive Map</p>
                <p className="text-sm mt-2">Map visualization with Leaflet/Mapbox integration</p>
                <div className="mt-4 flex justify-center gap-2">
                  {filteredPlaces.map((place) => (
                    <Badge key={place.id} variant={place.type === 'facility' ? 'primary' : 'success'} size="sm" className="cursor-pointer hover:scale-105 transition-transform">
                      {place.name}
                    </Badge>
                  ))}
                </div>
              </div>
            </div>
          </CardBody>
        </Card>
      )}

      <ConfirmModal
        isOpen={!!deleteModal}
        onClose={() => setDeleteModal(null)}
        onConfirm={() => { /* delete logic */ setDeleteModal(null); }}
        title="Delete Place"
        message="Are you sure you want to delete this saved place? This action cannot be undone."
        confirmText="Delete"
        variant="danger"
      />

      <Modal
        isOpen={addPlaceModal}
        onClose={() => { setAddPlaceModal(false); setEditPlace(null); }}
        title={editPlace ? 'Edit Place' : 'Add New Place'}
        size="lg"
      >
        <form className="space-y-6">
          <div className="grid gap-4 sm:grid-cols-2">
            <Input label="Place Name" placeholder="e.g. New Manufacturing Plant" defaultValue={editPlace?.name || ''} />
            <Select label="Type" options={[
              { value: 'facility', label: 'Facility' },
              { value: 'project', label: 'Project' },
            ]} defaultValue={editPlace?.type || 'facility'} />
          </div>
          <Input label="Address" placeholder="Full address" defaultValue={editPlace?.address || ''} className="sm:col-span-2" />
          <div className="grid gap-4 sm:grid-cols-2">
            <Input label="Latitude" type="number" step="any" placeholder="37.7749" defaultValue={editPlace?.coordinates[0]?.toString() || ''} />
            <Input label="Longitude" type="number" step="any" placeholder="-122.4194" defaultValue={editPlace?.coordinates[1]?.toString() || ''} />
          </div>
          <div className="grid gap-4 sm:grid-cols-3">
            <Input label="Emissions (tCO₂e)" type="number" placeholder="0" defaultValue={editPlace?.emissions?.toString() || ''} />
            <Input label="Energy (MWh)" type="number" placeholder="0" defaultValue={editPlace?.energy?.toString() || ''} />
            <Input label="Water (m³)" type="number" placeholder="0" defaultValue={editPlace?.water?.toString() || ''} />
          </div>
          <Select label="Status" options={statusOptions.filter(s => s.value !== 'all')} defaultValue={editPlace?.status || 'active'} />
          <Input label="Tags (comma separated)" placeholder="manufacturing, scope1, priority" defaultValue={editPlace?.tags?.join(', ') || ''} />
          <div className="flex justify-end gap-3 pt-4 border-t border-secondary-200 dark:border-secondary-700">
            <Button variant="secondary" onClick={() => { setAddPlaceModal(false); setEditPlace(null); }}>Cancel</Button>
            <Button variant="primary" type="submit">{editPlace ? 'Save Changes' : 'Add Place'}</Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}