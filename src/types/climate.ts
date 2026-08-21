export interface ClimateDataPoint {
  timestamp: string;
  value: number;
  unit: string;
  location?: GeoLocation;
  metadata?: Record<string, unknown>;
}

export interface GeoLocation {
  lat: number;
  lng: number;
  name?: string;
  region?: string;
  country?: string;
}

export interface EmissionsData {
  scope1: ClimateDataPoint[];
  scope2: ClimateDataPoint[];
  scope3: ClimateDataPoint[];
  total: ClimateDataPoint[];
  byCategory: Record<string, ClimateDataPoint[]>;
  byRegion: Record<string, ClimateDataPoint[]>;
  byFacility: Record<string, ClimateDataPoint[]>;
}

export interface CarbonFootprint {
  organizationId: string;
  period: { start: string; end: string };
  totalEmissions: number;
  unit: 'tCO2e' | 'kgCO2e';
  breakdown: {
    scope1: number;
    scope2: number;
    scope3: number;
  };
  intensity: number;
  intensityUnit: string;
  targets: ReductionTarget[];
}

export interface ReductionTarget {
  id: string;
  name: string;
  type: 'absolute' | 'intensity' | 'net-zero';
  baseYear: number;
  targetYear: number;
  baseValue: number;
  targetValue: number;
  currentValue: number;
  unit: string;
  status: 'on-track' | 'at-risk' | 'off-track' | 'achieved';
  scope: ('scope1' | 'scope2' | 'scope3')[];
  milestones: Milestone[];
}

export interface Milestone {
  year: number;
  targetValue: number;
  actualValue?: number;
  status: 'pending' | 'achieved' | 'missed';
}

export interface EnergyData {
  consumption: ClimateDataPoint[];
  production: ClimateDataPoint[];
  renewablePercentage: number;
  bySource: Record<string, ClimateDataPoint[]>;
  byFacility: Record<string, ClimateDataPoint[]>;
  costs: ClimateDataPoint[];
}

export interface WaterData {
  withdrawal: ClimateDataPoint[];
  consumption: ClimateDataPoint[];
  discharge: ClimateDataPoint[];
  recycled: ClimateDataPoint[];
  stressLevel: 'low' | 'medium' | 'high' | 'extreme';
  bySource: Record<string, ClimateDataPoint[]>;
  byFacility: Record<string, ClimateDataPoint[]>;
}

export interface BiodiversityData {
  speciesCount: number;
  habitats: HabitatData[];
  protectedAreas: ProtectedArea[];
  threats: ThreatData[];
  index: number;
  trend: 'improving' | 'stable' | 'declining';
}

export interface HabitatData {
  id: string;
  name: string;
  type: string;
  area: number;
  condition: 'good' | 'fair' | 'poor';
  speciesCount: number;
  location: GeoLocation;
}

export interface ProtectedArea {
  id: string;
  name: string;
  designation: string;
  area: number;
  location: GeoLocation;
  established: string;
}

export interface ThreatData {
  id: string;
  name: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  affectedSpecies: number;
  affectedArea: number;
  trend: 'increasing' | 'stable' | 'decreasing';
}

export interface AirQualityData {
  aqi: number;
  pollutants: Record<string, ClimateDataPoint>;
  byStation: Record<string, ClimateDataPoint[]>;
  healthImpact: 'good' | 'moderate' | 'unhealthy-sensitive' | 'unhealthy' | 'very-unhealthy' | 'hazardous';
}

export interface ClimateScenario {
  id: string;
  name: string;
  description: string;
  model: string;
  timeframe: { start: number; end: number };
  temperatureRise: number;
  emissionsPathway: ClimateDataPoint[];
  impacts: ScenarioImpact[];
  probability: number;
}

export interface ScenarioImpact {
  category: string;
  region: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  description: string;
  estimatedCost: number;
  adaptationOptions: string[];
}

export interface ClimateProject {
  id: string;
  name: string;
  description: string;
  type: 'mitigation' | 'adaptation' | 'resilience' | 'research';
  status: 'planning' | 'active' | 'completed' | 'on-hold' | 'cancelled';
  startDate: string;
  endDate?: string;
  budget: number;
  spent: number;
  location: GeoLocation;
  emissionsReduction: number;
  coBenefits: string[];
  stakeholders: string[];
  progress: number;
  kpis: ProjectKPI[];
}

export interface ProjectKPI {
  id: string;
  name: string;
  target: number;
  current: number;
  unit: string;
  frequency: 'monthly' | 'quarterly' | 'annually';
  trend: 'up' | 'down' | 'stable';
}

export interface ClimateInitiative {
  id: string;
  name: string;
  description: string;
  framework: string;
  commitments: Commitment[];
  progress: number;
  status: 'committed' | 'in-progress' | 'completed' | 'expired';
  startDate: string;
  endDate: string;
  reportingFrequency: 'annual' | 'biannual' | 'quarterly';
}

export interface Commitment {
  id: string;
  description: string;
  target: number;
  current: number;
  unit: string;
  deadline: string;
  status: 'on-track' | 'at-risk' | 'off-track' | 'achieved';
}

export interface Regulation {
  id: string;
  name: string;
  jurisdiction: string;
  type: 'mandatory' | 'voluntary' | 'guidance';
  status: 'active' | 'proposed' | 'expired' | 'superseded';
  effectiveDate: string;
  expiryDate?: string;
  requirements: RegulationRequirement[];
  applicability: string[];
  penalties: string;
  url: string;
}

export interface RegulationRequirement {
  id: string;
  description: string;
  metric: string;
  threshold: number;
  unit: string;
  frequency: string;
  verified: boolean;
}

export interface Disclosure {
  id: string;
  framework: string;
  reportingPeriod: { start: string; end: string };
  status: 'draft' | 'submitted' | 'approved' | 'rejected';
  submittedAt?: string;
  approvedAt?: string;
  sections: DisclosureSection[];
  score?: number;
}

export interface DisclosureSection {
  id: string;
  title: string;
  required: boolean;
  completed: boolean;
  dataPoints: DataPoint[];
}

export interface DataPoint {
  id: string;
  label: string;
  value: string | number;
  unit?: string;
  source: string;
  verified: boolean;
  lastUpdated: string;
}

export interface AuditLog {
  id: string;
  timestamp: string;
  userId: string;
  userName: string;
  action: string;
  resource: string;
  resourceId: string;
  changes: AuditChange[];
  ipAddress: string;
  userAgent: string;
}

export interface AuditChange {
  field: string;
  oldValue: unknown;
  newValue: unknown;
}

export interface Certification {
  id: string;
  name: string;
  issuer: string;
  standard: string;
  status: 'valid' | 'expired' | 'suspended' | 'pending';
  issuedDate: string;
  expiryDate: string;
  scope: string[];
  certificateUrl: string;
  verified: boolean;
}

export interface Integration {
  id: string;
  name: string;
  type: 'data-source' | 'api' | 'webhook' | 'file-upload' | 'iot-device';
  provider: string;
  status: 'connected' | 'disconnected' | 'error' | 'pending';
  config: Record<string, unknown>;
  lastSync?: string;
  nextSync?: string;
  syncFrequency: string;
  dataTypes: string[];
}

export interface ApiKey {
  id: string;
  name: string;
  key: string;
  prefix: string;
  permissions: string[];
  createdAt: string;
  lastUsedAt?: string;
  expiresAt?: string;
  status: 'active' | 'revoked' | 'expired';
}

export interface DashboardWidget {
  id: string;
  type: 'metric' | 'chart' | 'map' | 'table' | 'list' | 'progress' | 'gauge';
  title: string;
  dataSource: string;
  config: Record<string, unknown>;
  position: { x: number; y: number; w: number; h: number };
  refreshInterval?: number;
}

export interface Report {
  id: string;
  name: string;
  description: string;
  type: 'scheduled' | 'ad-hoc' | 'template';
  format: 'pdf' | 'excel' | 'csv' | 'html';
  frequency?: 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'annual';
  recipients: string[];
  filters: Record<string, unknown>;
  lastGenerated?: string;
  nextScheduled?: string;
  status: 'active' | 'paused' | 'error';
}