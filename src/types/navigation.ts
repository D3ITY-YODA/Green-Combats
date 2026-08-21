export interface NavItem {
  id: string;
  label: string;
  icon?: React.ComponentType<{ className?: string }>;
  href?: string;
  children?: NavItem[];
  badge?: string | number;
  badgeVariant?: 'primary' | 'success' | 'warning' | 'danger' | 'info';
  disabled?: boolean;
  roles?: UserRole[];
  external?: boolean;
  matchPaths?: string[];
}

export interface NavSection {
  id: string;
  label?: string;
  items: NavItem[];
  collapsible?: boolean;
  defaultOpen?: boolean;
}

export interface BreadcrumbItem {
  label: string;
  href?: string;
  icon?: React.ComponentType<{ className?: string }>;
}

export interface TabItem {
  id: string;
  label: string;
  icon?: React.ComponentType<{ className?: string }>;
  count?: number;
  disabled?: boolean;
  badge?: string | number;
}

export interface MenuItem {
  id: string;
  label: string;
  icon?: React.ComponentType<{ className?: string }>;
  shortcut?: string;
  disabled?: boolean;
  divider?: boolean;
  submenu?: MenuItem[];
  action?: () => void;
  danger?: boolean;
}

export type UserRole = 'admin' | 'analyst' | 'viewer' | 'contributor';

export const NAV_SECTIONS: NavSection[] = [
  {
    id: 'main',
    label: 'Main',
    items: [
      { id: 'dashboard', label: 'Dashboard', href: '/dashboard', icon: HomeIcon, matchPaths: ['/dashboard'] },
      { id: 'analytics', label: 'Analytics', href: '/analytics', icon: ChartBarIcon, matchPaths: ['/analytics'] },
      { id: 'maps', label: 'Maps', href: '/maps', icon: MapIcon, matchPaths: ['/maps'] },
      { id: 'reports', label: 'Reports', href: '/reports', icon: DocumentTextIcon, matchPaths: ['/reports'] },
    ],
  },
  {
    id: 'data',
    label: 'Data & Intelligence',
    items: [
      { id: 'emissions', label: 'Emissions Tracking', href: '/emissions', icon: CloudIcon, matchPaths: ['/emissions'] },
      { id: 'carbon', label: 'Carbon Footprint', href: '/carbon-footprint', icon: LeafIcon, matchPaths: ['/carbon-footprint'] },
      { id: 'energy', label: 'Energy Consumption', href: '/energy', icon: BoltIcon, matchPaths: ['/energy'] },
      { id: 'water', label: 'Water Resources', href: '/water', icon: DropIcon, matchPaths: ['/water'] },
      { id: 'biodiversity', label: 'Biodiversity Index', href: '/biodiversity', icon: BugIcon, matchPaths: ['/biodiversity'] },
      { id: 'air-quality', label: 'Air Quality', href: '/air-quality', icon: WindIcon, matchPaths: ['/air-quality'] },
    ],
  },
  {
    id: 'climate',
    label: 'Climate Action',
    items: [
      { id: 'targets', label: 'Reduction Targets', href: '/targets', icon: TargetIcon, matchPaths: ['/targets'] },
      { id: 'initiatives', label: 'Initiatives', href: '/initiatives', icon: RocketLaunchIcon, matchPaths: ['/initiatives'] },
      { id: 'projects', label: 'Projects', href: '/projects', icon: FolderIcon, matchPaths: ['/projects'] },
      { id: 'scenarios', label: 'Scenario Modeling', href: '/scenarios', icon: ChartBarIcon, matchPaths: ['/scenarios'] },
      { id: 'offsets', label: 'Carbon Offsets', href: '/offsets', icon: ArrowPathIcon, matchPaths: ['/offsets'] },
    ],
  },
  {
    id: 'compliance',
    label: 'Compliance & Reporting',
    items: [
      { id: 'regulations', label: 'Regulations', href: '/regulations', icon: ShieldCheckIcon, matchPaths: ['/regulations'] },
      { id: 'disclosures', label: 'Disclosures', href: '/disclosures', icon: DocumentTextIcon, matchPaths: ['/disclosures'] },
      { id: 'audits', label: 'Audit Trail', href: '/audits', icon: MagnifyingGlassIcon, matchPaths: ['/audits'] },
      { id: 'certifications', label: 'Certifications', href: '/certifications', icon: TrophyIcon, matchPaths: ['/certifications'] },
    ],
  },
  {
    id: 'admin',
    label: 'Administration',
    items: [
      { id: 'users', label: 'User Management', href: '/admin/users', icon: UsersIcon, roles: ['admin'], matchPaths: ['/admin/users'] },
      { id: 'organizations', label: 'Organizations', href: '/admin/organizations', icon: BuildingOfficeIcon, roles: ['admin'], matchPaths: ['/admin/organizations'] },
      { id: 'integrations', label: 'Integrations', href: '/admin/integrations', icon: PuzzlePieceIcon, roles: ['admin'], matchPaths: ['/admin/integrations'] },
      { id: 'api', label: 'API Management', href: '/admin/api', icon: CodeBracketIcon, roles: ['admin'], matchPaths: ['/admin/api'] },
      { id: 'settings', label: 'Platform Settings', href: '/admin/settings', icon: CogIcon, roles: ['admin'], matchPaths: ['/admin/settings'] },
    ],
    collapsible: true,
    defaultOpen: false,
  },
];

import {
  HomeIcon,
  ChartBarIcon,
  MapIcon,
  DocumentTextIcon,
  CloudIcon,
  LeafIcon,
  BoltIcon,
  DropIcon,
  BugIcon,
  WindIcon,
  TargetIcon,
  RocketLaunchIcon,
  FolderIcon,
  ArrowPathIcon,
  ShieldCheckIcon,
  MagnifyingGlassIcon,
  TrophyIcon,
  UsersIcon,
  BuildingOfficeIcon,
  PuzzlePieceIcon,
  CodeBracketIcon,
  CogIcon,
} from 'lucide-react';